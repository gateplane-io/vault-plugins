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
)

// Path for plugin configuration
func PathOktaApiConfig(b *Backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/okta/api",
		Fields: map[string]*framework.FieldSchema{
			"org_url": {
				Type:        framework.TypeString,
				Description: "The Okta Org Url (e.g: '<org>.okta.com')",
				Required:    false,
			},
			"api_token": {
				Type:        framework.TypeString,
				Description: "The Okta SSWS API Token",
				Required:    false,
			},
			"okta_entity_key": {
				Type:        framework.TypeString,
				Description: "The key that stores in Okta UserID in Vault Entity",
				Required:    false,
			},
			"auth_mount_accessor": {
				Type:        framework.TypeString,
				Description: "The Mount Accessor that authenticates users to Vault through Okta",
				Required:    false,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleOktaApiConfigUpdate,
			logical.ReadOperation:   b.handleOktaApiConfigRead,
		},

		HelpSynopsis: "Configures the Okta API access",
		HelpDescription: `This endpoint sets the Okta API configuration for this backend.
		<TODO>
		`,
	}
}

func (b *Backend) handleOktaApiConfigUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := b.GetOktaApiConfiguration(ctx, req)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	for key := range d.Raw {
		value := d.Get(key)
		logValue := value
		if key == "api_token" {
			logValue = "<redacted>"
		}
		b.Logger().Info("[*] Replacing configuration value",
			"EntityID", entityID,
			"ConfigKey", key,
			// "OldValue", config
			"NewValue", logValue,
		)

		err := config.SetConfigurationKey(key, value)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrPermissionDenied
		}
	}

	err = b.StoreOktaApiConfiguration(ctx, req, config)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	return &logical.Response{Data: map[string]interface{}{}}, nil
}

func (b *Backend) handleOktaApiConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := b.GetOktaApiConfiguration(ctx, req)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}
	token := "<not set>"
	if config.ApiToken != "" {
		token = "<set>"
	}

	responseData := map[string]interface{}{ // Just rewrite here for control
		"org_url":             config.OrgUrl,
		"api_token":           token,
		"okta_entity_key":     config.OktaEntityKey,
		"auth_mount_accessor": config.OktaOIDCMountAccessor,
	}

	return &logical.Response{Data: responseData}, nil
}
