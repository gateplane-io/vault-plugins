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
func PathConfigLease(b *BaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: ConfigLeaseKey,
		Fields: map[string]*framework.FieldSchema{
			"lease": {
				Type:        framework.TypeDurationSecond,
				Description: "Default lease for an Access Requests Claim",
				Required:    false,
			},
			"lease_max": {
				Type:        framework.TypeDurationSecond,
				Description: "Maximum lease for an Access Requests Claim",
				Required:    false,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleConfigLeaseUpdate,
			logical.ReadOperation:   b.handleConfigLeaseRead,
		},

		HelpSynopsis: "Configure the default lease durations of AccessRequest claims",
		HelpDescription: `This endpoint sets the lease configuration for this backend.

		'lease' configures duration of the claimed access, if not specific 'ttl' is set at creation time.

		'lease_max' configures the maximum duration that an AccessRequest can be claimed for.

		The format for the 'lease' and 'lease_max' is "1h" or integer and then unit.
		`,
	}
}

func (b *BaseBackend) handleConfigLeaseUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := GetConfiguration[*ConfigLease](ctx, b, req, "")
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

	err = StoreConfiguration[*ConfigLease](ctx, b, req, config, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	return &logical.Response{}, nil
}

func (b *BaseBackend) handleConfigLeaseRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := GetConfiguration[*ConfigLease](ctx, b, req, "")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	responseObj := responses.ConfigLeaseResponse{
		Lease:    config.Lease.Seconds(),
		LeaseMax: config.LeaseMax.Seconds(),
	}

	responseData, err := StructToMap(responseObj)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	return &logical.Response{Data: responseData}, nil
}
