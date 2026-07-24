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
	"github.com/gateplane-io/vault-plugins/internal/clients"
	clientConfig "github.com/gateplane-io/vault-plugins/internal/clients/config"
	"github.com/gateplane-io/vault-plugins/pkg/responses"
)

func PathConfigApiVault(b *Backend) *framework.Path {
	return &framework.Path{
		Pattern: ConfigAPIVaultKey,
		Fields: map[string]*framework.FieldSchema{
			"url": {
				Type:        framework.TypeString,
				Description: "The URL of the Vault/OpenBao instance",
				Required:    false,
			},
			"wrapped_token": {
				Type:        framework.TypeString,
				Description: "A single-use response-wrapped periodic orphan token",
				Required:    true,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleConfigApiVaultUpdate,
			logical.ReadOperation:   b.handleConfigApiVaultRead,
		},
		HelpSynopsis: "Configures access and authentication to Vault/OpenBao API",
		HelpDescription: `This endpoint configures the Vault/OpenBao endpoint and the token used
			to assign and remove policies from Vault/OpenBao entities.

			'url' configures the Vault/OpenBao network endpoint. When omitted, the existing URL is retained.

			'wrapped_token' must contain a single-use response-wrapped renewable periodic orphan token.
			The wrapping token is consumed once and is never stored. Only the unwrapped token is persisted
			in the plugin's encrypted storage and renewed through 'auth/token/renew-self'.

			The token must be able to 'read' and 'update' 'identity/entity/id/*', 'read'
			'auth/token/lookup-self', and 'update' 'auth/token/renew-self'.`,
	}
}

func (b *Backend) handleConfigApiVaultUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := base.GetConfiguration[*clientConfig.ConfigApiVaultPeriodicToken](ctx, b.BaseBackend, req, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	if rawURL, ok := d.GetOk("url"); ok {
		url, ok := rawURL.(string)
		if !ok {
			return logical.ErrorResponse("invalid type for 'url', expected string"), logical.ErrInvalidRequest
		}
		config.Url = url
	}

	rawWrappedToken, ok := d.GetOk("wrapped_token")
	if !ok {
		return logical.ErrorResponse("wrapped_token is required"), logical.ErrInvalidRequest
	}
	wrappedToken, ok := rawWrappedToken.(string)
	if !ok || wrappedToken == "" {
		return logical.ErrorResponse("wrapped_token is required"), logical.ErrInvalidRequest
	}

	vaultClient, tokenInfo, err := clients.NewVaultWrappedTokenClient(ctx, config.Url, wrappedToken, nil)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrInvalidRequest
	}
	config.Token = vaultClient.Token()

	if err := base.StoreConfiguration[*clientConfig.ConfigApiVaultPeriodicToken](ctx, b.BaseBackend, req, config, ""); err != nil {
		_ = vaultClient.Auth().Token().RevokeSelfWithContext(ctx, config.Token)
		vaultClient.ClearToken()
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	if err := b.ReplaceVaultClient(vaultClient, tokenInfo); err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	return &logical.Response{}, nil
}

func (b *Backend) handleConfigApiVaultRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := base.GetConfiguration[*clientConfig.ConfigApiVaultPeriodicToken](ctx, b.BaseBackend, req, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	tokenSet, tokenPeriod := b.VaultTokenStatus()
	responseData, err := base.StructToMap(responses.ConfigApiVaultResponse{
		Url:         config.Url,
		TokenSet:    tokenSet || config.Token != "",
		TokenPeriod: tokenPeriod,
	})
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	return &logical.Response{Data: responseData}, nil
}
