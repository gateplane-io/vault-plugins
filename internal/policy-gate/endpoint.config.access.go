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
	"github.com/gateplane-io/vault-plugins/pkg/responses"
)

// Path for plugin configuration
func PathConfigAccess(b *Backend) *framework.Path {
	return &framework.Path{
		Pattern: ConfigAccessKey,
		Fields: map[string]*framework.FieldSchema{
			"policies": {
				Type:        framework.TypeStringSlice,
				Description: "The Vault/OpenBao policies to be assigned to approved requestors",
				Required:    true,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleConfigAccessUpdate,
			logical.ReadOperation:   b.handleConfigAccessRead,
		},

		HelpSynopsis: "Configures the access claimed by this backend",
		HelpDescription: `This endpoint sets the Vault/OpenBao policies to be conditionally assigned to the requesting Entity.

		'policies' contains a list of Vault/OpenBao Policies that will be assigned (and retracted) when an AccessRequest is claimed.
		`,
	}
}

func (b *Backend) handleConfigAccessUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := base.GetConfiguration[*ConfigAccess](ctx, b.BaseBackend, req, "")
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

	err = base.StoreConfiguration[*ConfigAccess](ctx, b.BaseBackend, req, config, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	return &logical.Response{}, nil
}

func (b *Backend) handleConfigAccessRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := base.GetConfiguration[*ConfigAccess](ctx, b.BaseBackend, req, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	responseObj := responses.ConfigAccessPolicyGate{
		Policies: config.Policies,
	}

	responseData, err := base.StructToMap(responseObj)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	return &logical.Response{Data: responseData}, nil
}
