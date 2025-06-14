// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package mock

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
		HelpSynopsis: "[Mock] Accepts a GrantCode to Claim an approved AccessRequest.",
		HelpDescription: `This endpoint accepts a GrantCode (grant_code="<uuid>"),
		issued by the '/request' endpoint to the AccessRequest creator,
		when sufficient approvals have been submitted to this AccessRequest.

		The successful response of this endpoint provides the access supported by the plugin.

		This plugin is Mock, only returning a successful response for testing purposes.
		`,
	}
}

func (b *Backend) handleClaim(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	grantCode := d.Get("grant_code").(string)

	b.BaseBackend.BaseMutex.Lock()
	defer b.BaseBackend.BaseMutex.Unlock()

	accessRequest, err := b.GetAccessRequestByGrantCode(ctx, req, grantCode)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrNotFound
	}

	b.Logger().Info("[+] Claimed Token",
		"entityID", accessRequest.ID,
		"approvals", accessRequest.Approvals,
	)

	return &logical.Response{
		Data: map[string]interface{}{
			"claimed": true,
		},
	}, nil
}
