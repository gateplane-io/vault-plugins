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
	"sync"
	// "fmt"

	// "github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	vault "github.com/hashicorp/vault/api"

	"github.com/gateplane-io/vault-plugins/internal/base"
	clients "github.com/gateplane-io/vault-plugins/internal/clients"
	clientConfig "github.com/gateplane-io/vault-plugins/internal/clients/config"
	"github.com/gateplane-io/vault-plugins/internal/utils"
)

/* Storage Keys */
const ConfigAPIVaultKey = "config/api/vault"
const ConfigAccessKey = "config/access"

type Backend struct {
	*base.BaseBackend
	Mutex       sync.Mutex
	VaultClient *vault.Client
}

func (b *Backend) Initialize(ctx context.Context, req *logical.InitializationRequest) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	b.Logger().Info("Initializing plugin configuration")

	// Initialize the base
	err := b.BaseBackend.Initialize(ctx, req)
	if err != nil {
		b.Logger().Error("Could not initialize Plugin Base")
		return err
	}

	configAccess := NewConfigAccess()
	_, err = base.StoreConfigurationToStorageIfNotPresent[*ConfigAccess](ctx,
		b.BaseBackend, req.Storage, &configAccess, ConfigAccessKey,
	)

	config := clientConfig.NewConfigApiVaultAppRole()
	exists, err := base.StoreConfigurationToStorageIfNotPresent[*clientConfig.ConfigApiVaultAppRole](ctx,
		b.BaseBackend, req.Storage, &config, ConfigAPIVaultKey,
	)
	if err != nil {
		b.Logger().Error("[-] Could not initialize Vault API Configuration")
		return err
	}

	b.BaseBackend.ClaimArray = utils.NewCallbackArray(
		(func(ctx context.Context, requ *logical.Request, ownerID string) (map[string]interface{}, error) { // Append
			err := b.EnsureVaultAPI(ctx, req.Storage)
			if err != nil {
				return nil, err
			}

			cfg, err := base.GetConfigurationFromStorage[*ConfigAccess](ctx,
				b.BaseBackend, req.Storage, ConfigAccessKey,
			)
			if err != nil {
				return nil, err
			}

			newPolicies := cfg.Policies //[]string{"test-pgate"} // to be fetched from config
			existingPolicies, err := AddPoliciesToEntity(ctx, b.VaultClient, ownerID, newPolicies)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"previous_policies": existingPolicies,
				"new_policies":      newPolicies,
			}, nil
		}),
		(func(ctx context.Context, requ *logical.Request, ownerID string, internalData map[string]interface{}) error { // Append
			err := b.EnsureVaultAPI(ctx, req.Storage)
			if err != nil {
				return err
			}

			policies := Subtract[string](
				InterfaceSliceToStringsStrict(internalData["previous_policies"].([]interface{})),
				InterfaceSliceToStringsStrict(internalData["new_policies"].([]interface{})),
			)
			err = SetEntityPolicies(ctx, b.VaultClient, ownerID, policies)

			return err
		}),
	)
	b.Logger().Info("GatePlane Policy Gate initialized with default configuration",
		"path", ConfigAPIVaultKey,
		"AppRoleID", config.RoleID,
		"AppRoleMount", config.AppRoleMount,
		"already_set", exists,
	)
	return nil
}

func (b *Backend) EnsureVaultAPI(ctx context.Context, storage logical.Storage) error {

	cfg, err := base.GetConfigurationFromStorage[*clientConfig.ConfigApiVaultAppRole](ctx,
		b.BaseBackend, storage, ConfigAPIVaultKey,
	)
	if err != nil {
		return err
	}

	err = clients.EnsureAuthenticationVault(ctx,
		b.VaultClient, cfg.RoleID, cfg.RoleSecret, cfg.AppRoleMount,
	)
	if err != nil {
		b.Logger().Error("[-] Could not ensure authentication to the Vault API",
			"configuration", cfg,
			"error", err,
		)
		return err
	}

	return nil
}
