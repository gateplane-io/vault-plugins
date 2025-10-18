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

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/gateplane-io/vault-plugins/internal/base"
	clients "github.com/gateplane-io/vault-plugins/internal/clients"
	clientConfig "github.com/gateplane-io/vault-plugins/internal/clients/config"

	"github.com/gateplane-io/vault-plugins/pkg/responses"
)

// Path for plugin configuration
func PathConfigApiOkta(b *Backend) *framework.Path {
	return &framework.Path{
		Pattern: ConfigAPIOktaKey,
		Fields: map[string]*framework.FieldSchema{
			"org_url": {
				Type:        framework.TypeString,
				Description: "The Okta Org Url (e.g: 'https://<org>.okta.com')",
				Required:    true,
			},
			"api_token": {
				Type:        framework.TypeString,
				Description: "The Okta SSWS API Token",
				Required:    true,
			},
			"okta_entity_key": {
				Type:        framework.TypeString,
				Description: "The key that containing Okta UserID in Vault/OpenBao Entity",
				Required:    false,
			},
			"auth_mount_accessor": {
				Type:        framework.TypeString,
				Description: "The Mount Accessor that authenticates users to Vault through Okta",
				Required:    false,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleConfigApiOktaUpdate,
			logical.ReadOperation:   b.handleConfigApiOktaRead,
		},

		HelpSynopsis: "Configures access and authentication to Okta API",
		HelpDescription: `This endpoint sets the Okta API configuration for this backend.
		'okta_entity_key': The Vault/OpenBao Entity Metadata key that holds the Okta ID of the user.

		'auth_mount_accessor': In case the user logs into Vault/OpenBao through Okta,
		and the token's 'sub' value, containing the Okta User ID, is stored as Alias Name,
		this field fetches the Okta ID.
		`,
	}
}

func (b *Backend) handleConfigApiOktaUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := base.GetConfiguration[*clientConfig.ConfigApiOkta](ctx, b.BaseBackend, req, "")
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

	err = base.StoreConfiguration[*clientConfig.ConfigApiOkta](ctx, b.BaseBackend, req, config, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	oktaClient, err := clients.NewOktaClient(config.OrgUrl, config.ApiToken)
	if err != nil {
		return &logical.Response{Warnings: []string{
			fmt.Sprintf("%s", err),
		}}, nil
	}
	b.OktaClient = oktaClient

	return &logical.Response{}, nil
}

func (b *Backend) handleConfigApiOktaRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := base.GetConfiguration[*clientConfig.ConfigApiOkta](ctx, b.BaseBackend, req, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	responseObj := responses.ConfigApiOktaResponse{
		OrgUrl:                config.OrgUrl,
		OktaOIDCMountAccessor: config.OktaOIDCMountAccessor,
		OktaEntityKey:         config.OktaEntityKey,
		ApiTokenSet:           config.ApiToken != "",
	}

	responseData, err := base.StructToMap(responseObj)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	return &logical.Response{Data: responseData}, nil
}
