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
)

// Path for plugin configuration
func PathConfig(b *Backend) *framework.Path {
	basePaths := base.PathConfig(b.BaseBackend)
	basePaths.Fields["policies"] = &framework.FieldSchema{
		Type:        framework.TypeCommaStringSlice,
		Description: "Policies to be attached to the issued token.",
		Required:    false,
	}
	basePaths.Callbacks = map[logical.Operation]framework.OperationFunc{
		logical.UpdateOperation: b.handleConfigUpdate,
		logical.ReadOperation:   b.handleConfigRead,
	}
	basePaths.HelpDescription = fmt.Sprintf(`%s

	policies: The list of Vault/OpenBao policies to attach to the issued token
	`, basePaths.HelpDescription)
	return basePaths
}

func (b *Backend) handleConfigUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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

func (b *Backend) handleConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
		"policies":           config.Policies,
	}

	return &logical.Response{Data: responseData}, nil
}
