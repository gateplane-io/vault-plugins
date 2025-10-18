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

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/gateplane-io/vault-plugins/internal/base"
	clients "github.com/gateplane-io/vault-plugins/internal/clients"
	clientConfig "github.com/gateplane-io/vault-plugins/internal/clients/config"

	"github.com/gateplane-io/vault-plugins/pkg/responses"
)

// Path for plugin configuration
func PathConfigApiVault(b *Backend) *framework.Path {
	return &framework.Path{
		Pattern: ConfigAPIVaultKey,
		Fields: map[string]*framework.FieldSchema{
			"url": {
				Type:        framework.TypeString,
				Description: "The URL of the Vault/OpenBao instance",
				Required:    true,
			},
			"role_id": {
				Type:        framework.TypeString,
				Description: "The AppRole ID to authenticate against",
				Required:    true,
			},
			"role_secret": {
				Type:        framework.TypeString,
				Description: "The AppRole Secret to authenticate with",
				Required:    false,
			},
			"approle_mount": {
				Type:        framework.TypeString,
				Description: "The Vault Auth Mount for AppRole",
				Required:    false,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleConfigApiVaultUpdate,
			logical.ReadOperation:   b.handleConfigApiVaultRead,
		},

		HelpSynopsis: "Configures access and authentication to Vault/OpenBao API",
		HelpDescription: `This endpoint sets the Vault/OpenBao endpoints and AppRole configuration
		needed to assign and remove Policies to Vault/OpenBao Entities.

		'url' configures the Vault/OpenBao network endpoint to be used to contact the API

		'role_id' and 'role_secret' configure the AppRole credentials to authenticate to and interact with the API

		'approle_mount' configures the mountpoint of the AppRole Auth Method, where the AppRole should authenticate against.

		The provided AppRole must be able to
		'read' and 'update' the 'identity/entity/id/*',
		and 'read' the 'auth/token/lookup-self' paths.
		`,
	}
}

func (b *Backend) handleConfigApiVaultUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := base.GetConfiguration[*clientConfig.ConfigApiVaultAppRole](ctx, b.BaseBackend, req, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	for key := range d.Raw {
		value, ok := d.GetOk(key)
		if !ok {
			continue
		}
		b.Logger().Info("[*] Replacing configuration value",
			"EntityID", entityID,
			"ConfigKey", key,
			// "OldValue", config
			"NewValue", value,
		)

		err := config.SetConfigurationKey(key, value)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrPermissionDenied
		}
	}

	err = base.StoreConfiguration[*clientConfig.ConfigApiVaultAppRole](ctx, b.BaseBackend, req, config, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	vaultClient, err := clients.NewVaultAppRoleClient(ctx, *config, nil)
	if err != nil {
		return &logical.Response{Warnings: []string{
			fmt.Sprintf("%s", err),
		}}, nil
	}
	b.VaultClient = vaultClient

	return &logical.Response{}, nil
}

func (b *Backend) handleConfigApiVaultRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := base.GetConfiguration[*clientConfig.ConfigApiVaultAppRole](ctx, b.BaseBackend, req, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	responseObj := responses.ConfigApiVaultResponse{
		Url:           config.Url,
		RoleID:        config.RoleID,
		AppRoleMount:  config.AppRoleMount,
		RoleSecretSet: config.RoleSecret != "",
	}

	responseData, err := base.StructToMap(responseObj)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	return &logical.Response{Data: responseData}, nil
}
