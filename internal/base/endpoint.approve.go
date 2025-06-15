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

// Path for gatekeeper to grant access
func PathApprove(b *BaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: "approve",
		Fields: map[string]*framework.FieldSchema{
			"request_id": {
				Type:        framework.TypeString,
				Description: "The entity ID of the user being granted access.",
				Required:    true,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleApprove,
		},

		HelpSynopsis:    "Adds an Approval to an AccessRequest",
		HelpDescription: `This endpoint approves AccessRequests designated by 'request_id' parameter.`,
	}
}

func (b *BaseBackend) handleApprove(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID
	requestID := d.Get("request_id").(string)

	if entityID == "" {
		return logical.ErrorResponse("Token has no EntityID assigned"), logical.ErrPermissionDenied
	}

	if requestID == "" {
		return nil, fmt.Errorf("No requestor EntityID set")
	}

	if entityID == requestID {
		return logical.ErrorResponse("Entities cannot approve their own requests"), logical.ErrPermissionDenied
	}

	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := b.GetConfiguration(ctx, req)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	accessRequest, err := b.getRequest(ctx, req, requestID)
	if err != nil {
		return logical.ErrorResponse("Corresponding request does not exist"), logical.ErrNotFound
	}

	if accessRequest.Status != Pending {
		return logical.ErrorResponse("Request cannot be approved as it is not in pending state"), logical.ErrPermissionDenied
	}

	_, haveApproved := accessRequest.Approvals[entityID]
	if haveApproved {
		return logical.ErrorResponse("Approval by this entity already exists"), logical.ErrPermissionDenied
	}

	approval := NewApproval(config, entityID, requestID)
	accessRequest.Approvals[entityID] = approval
	approved := requestIsApproved(*accessRequest, config.RequiredApprovals)

	err = b.storeRequest(ctx, req, accessRequest)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	b.Logger().Info("[+] Approval created",
		"EntityID", entityID,
		"RequestID", requestID,
	)

	return &logical.Response{
		Data: map[string]interface{}{
			"message":         "access approved",
			"iat":             approval.CreatedAt,
			"exp":             approval.Expiration,
			"approval_id":     approval.ID,
			"access_approved": approved,
		},
	}, nil
}
