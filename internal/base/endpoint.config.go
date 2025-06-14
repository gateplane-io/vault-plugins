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
)

// Path for plugin configuration
func PathConfig(b *BaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"require_reason": {
				Type:        framework.TypeBool,
				Description: "Whether the 'reason' parameter is required in /request endpoint.",
				Required:    false,
			},
			"request_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Time until a request expires.",
				Required:    false,
			},
			"approval_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Time until a granted request expires.",
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

		HelpSynopsis: "Configures this backend",
		HelpDescription: `This endpoint sets the base configuration for this backend.

		'require_reason' configures whether a non-empty 'reason' parameter is required for AccessRequest creation.

		'request_ttl', 'approval_ttl' and 'delete_after' configure the lifetime of AccessRequests and Approvals.

		'required_approvals' sets the umber of approval for an AccessRequest to get to a claimable state.
		`,
	}
}

func (b *BaseBackend) handleConfigUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := b.GetConfiguration(ctx, req)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	for key := range d.Raw {
		value := d.Get(key)
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

	err = b.StoreConfiguration(ctx, req, config)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	return &logical.Response{Data: map[string]interface{}{
		"config": config,
	}}, nil
}

func (b *BaseBackend) handleConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := b.GetConfiguration(ctx, req)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}
	responseData := map[string]interface{}{ // Just rewrite here for control
		"required_approvals": config.RequiredApprovals,
		"require_reason":     config.RequireReason,
		"approval_ttl":       config.ApprovalTTL.Seconds(),
		"request_ttl":        config.RequestTTL.Seconds(),
		"delete_after":       config.DeleteAfter.Seconds(),
	}

	return &logical.Response{Data: responseData}, nil
}
