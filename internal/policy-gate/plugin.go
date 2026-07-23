// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package policy_gate

import (
	"context"
	"fmt"
	"sync"

	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/gateplane-io/vault-plugins/internal/base"
	"github.com/gateplane-io/vault-plugins/internal/clients"
	clientConfig "github.com/gateplane-io/vault-plugins/internal/clients/config"
	"github.com/gateplane-io/vault-plugins/internal/utils"
)

const ConfigAPIVaultKey = "config/api/vault"
const ConfigAccessKey = "config/access"

type Backend struct {
	*base.BaseBackend
	Mutex             sync.Mutex
	VaultClient       *vault.Client
	vaultTokenWatcher *vault.LifetimeWatcher
	vaultTokenPeriod  int
}

func (b *Backend) Initialize(ctx context.Context, req *logical.InitializationRequest) error {
	b.Logger().Info("Initializing plugin configuration")

	if err := b.BaseBackend.Initialize(ctx, req); err != nil {
		b.Logger().Error("Could not initialize Plugin Base")
		return err
	}

	configAccess := NewConfigAccess()
	if _, err := base.StoreConfigurationToStorageIfNotPresent[*ConfigAccess](ctx,
		b.BaseBackend, req.Storage, &configAccess, ConfigAccessKey,
	); err != nil {
		return err
	}

	config := clientConfig.NewConfigApiVaultPeriodicToken()
	exists, err := base.StoreConfigurationToStorageIfNotPresent[*clientConfig.ConfigApiVaultPeriodicToken](ctx,
		b.BaseBackend, req.Storage, &config, ConfigAPIVaultKey,
	)
	if err != nil {
		b.Logger().Error("[-] Could not initialize Vault API Configuration")
		return err
	}
	if exists {
		storedConfig, err := base.GetConfigurationFromStorage[*clientConfig.ConfigApiVaultPeriodicToken](ctx,
			b.BaseBackend, req.Storage, ConfigAPIVaultKey,
		)
		if err != nil {
			return err
		}
		if storedConfig.Token != "" {
			if err := b.EnsureVaultAPI(ctx, req.Storage); err != nil {
				b.Logger().Error("[-] Could not initialize Vault API with existing configuration",
					"path", ConfigAPIVaultKey,
					"error", err,
				)
				return err
			}
		}
	}

	b.BaseBackend.ClaimArray = utils.NewCallbackArray(
		func(ctx context.Context, _ *logical.Request, ownerID string) (map[string]interface{}, error) {
			if err := b.EnsureVaultAPI(ctx, req.Storage); err != nil {
				return nil, err
			}
			vaultClient := b.currentVaultClient()
			if vaultClient == nil {
				return nil, fmt.Errorf("vault API is not configured")
			}

			cfg, err := base.GetConfigurationFromStorage[*ConfigAccess](ctx,
				b.BaseBackend, req.Storage, ConfigAccessKey,
			)
			if err != nil {
				return nil, err
			}

			newPolicies := cfg.Policies
			existingPolicies, err := AddPoliciesToEntity(ctx, vaultClient, ownerID, newPolicies)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"previous_policies": existingPolicies,
				"new_policies":      newPolicies,
			}, nil
		},
		func(ctx context.Context, _ *logical.Request, ownerID string, internalData map[string]interface{}) error {
			if err := b.EnsureVaultAPI(ctx, req.Storage); err != nil {
				return err
			}
			vaultClient := b.currentVaultClient()
			if vaultClient == nil {
				return fmt.Errorf("vault API is not configured")
			}

			policies := Subtract[string](
				InterfaceSliceToStringsStrict(internalData["previous_policies"].([]interface{})),
				InterfaceSliceToStringsStrict(internalData["new_policies"].([]interface{})),
			)
			return SetEntityPolicies(ctx, vaultClient, ownerID, policies)
		},
	)

	b.Logger().Info("GatePlane Policy Gate initialized with default configuration",
		"path", ConfigAPIVaultKey,
		"token_set", b.currentVaultClient() != nil,
		"already_set", exists,
	)
	return nil
}

func (b *Backend) EnsureVaultAPI(ctx context.Context, storage logical.Storage) error {
	b.Mutex.Lock()
	clientReady := b.VaultClient != nil && b.vaultTokenWatcher != nil
	client := b.VaultClient
	b.Mutex.Unlock()

	if clientReady {
		if _, err := client.Auth().Token().LookupSelfWithContext(ctx); err != nil {
			return fmt.Errorf("checking configured vault token: %w", err)
		}
		return nil
	}

	config, err := base.GetConfigurationFromStorage[*clientConfig.ConfigApiVaultPeriodicToken](ctx,
		b.BaseBackend, storage, ConfigAPIVaultKey,
	)
	if err != nil {
		return err
	}
	if config.Token == "" {
		return fmt.Errorf("vault API token is not configured")
	}

	client, tokenInfo, err := clients.NewVaultPeriodicTokenClient(ctx, config.Url, config.Token, nil)
	if err != nil {
		return err
	}
	return b.ReplaceVaultClient(client, tokenInfo)
}

func (b *Backend) ReplaceVaultClient(client *vault.Client, tokenInfo *clients.PeriodicTokenInfo) error {
	if client == nil {
		return fmt.Errorf("vault client is nil")
	}
	if tokenInfo == nil || tokenInfo.Secret == nil || tokenInfo.PeriodSeconds <= 0 {
		return fmt.Errorf("periodic token information is invalid")
	}

	watcher, err := client.NewLifetimeWatcher(&vault.LifetimeWatcherInput{
		Secret:    tokenInfo.Secret,
		Increment: tokenInfo.PeriodSeconds,
	})
	if err != nil {
		return fmt.Errorf("creating vault token lifetime watcher: %w", err)
	}

	b.Mutex.Lock()
	oldWatcher := b.vaultTokenWatcher
	b.VaultClient = client
	b.vaultTokenWatcher = watcher
	b.vaultTokenPeriod = tokenInfo.PeriodSeconds
	b.Mutex.Unlock()

	if oldWatcher != nil {
		oldWatcher.Stop()
	}

	go b.runVaultTokenWatcher(watcher)
	return nil
}

func (b *Backend) runVaultTokenWatcher(watcher *vault.LifetimeWatcher) {
	watcher.Start()
	err := <-watcher.DoneCh()

	b.Mutex.Lock()
	if b.vaultTokenWatcher == watcher {
		b.vaultTokenWatcher = nil
	}
	b.Mutex.Unlock()

	if err != nil {
		b.Logger().Error("Vault API periodic token renewal stopped", "error", err)
	}
}

func (b *Backend) CleanupVaultClient(_ context.Context) {
	b.Mutex.Lock()
	watcher := b.vaultTokenWatcher
	b.vaultTokenWatcher = nil
	b.VaultClient = nil
	b.vaultTokenPeriod = 0
	b.Mutex.Unlock()

	if watcher != nil {
		watcher.Stop()
	}
}

func (b *Backend) VaultTokenStatus() (bool, int) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	return b.VaultClient != nil, b.vaultTokenPeriod
}

func (b *Backend) currentVaultClient() *vault.Client {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	return b.VaultClient
}
