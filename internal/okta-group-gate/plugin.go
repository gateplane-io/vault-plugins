// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package okta_group_gate

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/okta/okta-sdk-golang/v5/okta"

	"github.com/gateplane-io/vault-plugins/internal/base"
	clients "github.com/gateplane-io/vault-plugins/internal/clients"
	clientConfig "github.com/gateplane-io/vault-plugins/internal/clients/config"
	"github.com/gateplane-io/vault-plugins/internal/utils"
)

/* Storage Keys */
const ConfigAPIOktaKey = "config/api/okta"
const ConfigAccessKey = "config/access"

type Backend struct {
	*base.BaseBackend
	Mutex      sync.Mutex
	OktaClient *okta.APIClient
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

	config := clientConfig.NewConfigApiOkta()
	oktaExists, err := base.StoreConfigurationToStorageIfNotPresent[*clientConfig.ConfigApiOkta](ctx,
		b.BaseBackend, req.Storage, &config, ConfigAPIOktaKey,
	)
	if err != nil {
		b.Logger().Error("[-] Could not initialize Okta API Configuration")
		return err
	}

	b.BaseBackend.ClaimArray = utils.NewCallbackArray(
		(func(ctx context.Context, requ *logical.Request, ownerID string) (map[string]interface{}, error) { // Append
			err := b.EnsureOktaAPI(ctx, req.Storage)
			if err != nil {
				return nil, err
			}

			cfg, err := base.GetConfigurationFromStorage[*ConfigAccess](ctx,
				b.BaseBackend, req.Storage, ConfigAccessKey,
			)
			if err != nil {
				return nil, err
			}
			oktaConfig, err := base.GetConfigurationFromStorage[*clientConfig.ConfigApiOkta](ctx,
				b.BaseBackend, req.Storage, ConfigAPIOktaKey,
			)
			if err != nil {
				return nil, err
			}

			groupName := cfg.GroupName
			if groupName == "" {
				return nil, fmt.Errorf("Could not retrieve Okta Group details from API")
			}

			groupID := cfg.GroupID
			oktaUserID, err := b.GetOktaUserID(ctx, oktaConfig, ownerID)
			if err != nil && oktaUserID != "" {
				return nil, err
			}

			err = oktaAddToGroup(ctx, b.OktaClient, groupID, oktaUserID)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"okta_group_id": groupID,
				"okta_user_id":  oktaUserID,
			}, nil
		}),
		(func(ctx context.Context, requ *logical.Request, ownerID string, internalData map[string]interface{}) error { // Append
			err := b.EnsureOktaAPI(ctx, req.Storage)
			if err != nil {
				return err
			}

			groupID := internalData["okta_group_id"].(string)
			oktaUserID := internalData["okta_user_id"].(string)

			err = oktaRemoveFromGroup(ctx, b.OktaClient, groupID, oktaUserID)
			if err != nil {
				return err
			}
			return err
		}),
	)
	b.Logger().Info("GatePlane Okta Group Gate initialized with default configuration",
		"path", ConfigAPIOktaKey,
		"OktaURL", config.OrgUrl,
		"OktaConfigExists", oktaExists,
	)
	return nil
}

func (b *Backend) EnsureOktaAPI(ctx context.Context, storage logical.Storage) error {

	cfg, err := base.GetConfigurationFromStorage[*clientConfig.ConfigApiOkta](ctx,
		b.BaseBackend, storage, ConfigAPIOktaKey,
	)
	if err != nil {
		return err
	}

	if cfg.ApiToken != "" {
		oktaClient, err := clients.NewOktaClient(cfg.OrgUrl, cfg.ApiToken)
		if err != nil {
			b.Logger().Error("[-] Could not create Okta Client")
			return err
		}
		b.OktaClient = oktaClient
	}
	return nil
}

func (b *Backend) GetOktaUserID(ctx context.Context, config *clientConfig.ConfigApiOkta, entityID string) (string, error) {

	view := b.System()
	entity, err := view.EntityInfo(entityID)
	if err != nil {
		return "", err
	}

	entityMeta := entity.GetMetadata()
	oktaUserId, exists := entityMeta[config.OktaEntityKey]
	if exists {
		return oktaUserId, nil
	}

	for _, alias := range entity.Aliases {
		if alias.MountAccessor != config.OktaOIDCMountAccessor {
			continue
		}
		oktaUserId = alias.Name
	}
	if oktaUserId != "" {
		return oktaUserId, nil
	}
	b.Logger().Error("[-] Could not retrieve Okta User ID from Vault/OpenBao Entity",
		"EntityID", entityID,
		"EntityMetadata", entityMeta,
		"OktaOIDCMountAccessor", config.OktaOIDCMountAccessor,
		"OktaEntityKey", config.OktaEntityKey,
	)

	return "", fmt.Errorf("Could not retrieve Okta User ID from Vault/OpenBao Entity")
}
