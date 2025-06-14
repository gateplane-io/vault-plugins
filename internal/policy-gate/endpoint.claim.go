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
)

// Path for user to claim token
func PathClaim(b *Backend) *framework.Path {
	return &framework.Path{
		Pattern: "claim",
		Fields: map[string]*framework.FieldSchema{
			"grant_code": {
				Type:        framework.TypeString,
				Description: "The GrantCode found in the approved AccessRequest",
				Required:    true,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleClaim,
		},
	}
}

func (b *Backend) handleClaim(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	grantCode := d.Get("grant_code").(string)

	config, err := b.GetConfiguration(ctx, req)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	accessRequest, err := b.GetAccessRequestByGrantCode(ctx, req, grantCode)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	b.Logger().Info("[+] Claimed Token",
		"policies", config.Policies,
		"entityID", accessRequest.ID,
	)

	return &logical.Response{
		Auth: &logical.Auth{
			Policies:        config.Policies,
			NoDefaultPolicy: true,
			Metadata: map[string]string{
				"grantedBy": "gateplane",
			},
			LeaseOptions: logical.LeaseOptions{
				Renewable: false,
			},
			EntityID: accessRequest.ID,
		},
	}, nil
}
