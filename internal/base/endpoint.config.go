// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package base

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/gateplane-io/vault-plugins/pkg/responses"
)

// Path for plugin configuration
func PathConfig(b *BaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: ConfigKey,
		Fields: map[string]*framework.FieldSchema{
			"require_justification": {
				Type:        framework.TypeBool,
				Description: "Whether the 'reason' parameter is required in /request endpoint.",
				Required:    false,
			},
			"request_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Time until a request expires.",
				Required:    false,
			},
			"delete_after": {
				Type:        framework.TypeDurationSecond,
				Description: "Time until a granted request expires.",
				Required:    false,
			},
			"required_approvals": {
				Type:        framework.TypeInt,
				Description: "Required number of approvals before claiming.",
				Required:    false,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleConfigUpdate,
			logical.ReadOperation:   b.handleConfigRead,
		},

		HelpSynopsis: "Configure this backend conditional workflows and AccessRequest validity periods",
		HelpDescription: `This endpoint sets the base configuration for this backend.

		'require_justification' configures whether a non-empty 'justification' parameter is required for the creation of an AccessRequest.

		'request_ttl' and 'delete_after' configure the lifetime of AccessRequests.

		'required_approvals' sets the number of approvals for an AccessRequest required to reach the 'approved' state (can be positive integer or 0).
		`,
	}
}

func (b *BaseBackend) handleConfigUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := GetConfiguration[*Config](ctx, b, req, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	for key := range d.Raw {
		value, ok := d.GetOk(key)
		if !ok {
			continue
		}
		b.Logger().Info("[*] Replacing configuration value",
			"Path", req.Path,
			"EntityID", entityID,
			"ConfigKey", key,
			// "OldValue", config,
			"NewValue", value,
		)

		err := config.SetConfigurationKey(key, value)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrPermissionDenied
		}
	}
	err = StoreConfiguration[*Config](ctx, b, req, config, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	return &logical.Response{}, nil
}

func (b *BaseBackend) handleConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := GetConfiguration[*Config](ctx, b, req, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	responseObj := responses.ConfigResponse{
		RequiredApprovals:    config.RequiredApprovals,
		RequireJustification: config.RequireJustification,
		// RequestTTL:           config.RequestTTL,
		// DeleteAfter:          config.DeleteAfter,
		// If I need seconds
		RequestTTL:  config.RequestTTL.Seconds(),
		DeleteAfter: config.DeleteAfter.Seconds(),
	}

	responseData, err := StructToMap(responseObj)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	return &logical.Response{Data: responseData}, nil
}
