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

func PathApprove(b *BaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: "approve",
		Fields: map[string]*framework.FieldSchema{
			"requestor_id": {
				Type:        framework.TypeString,
				Description: "The RequestorID of the AccessRequest to approve",
				Required:    false,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleApprove,
		},
		HelpSynopsis: "Approves the AccessRequest created by the provided RequestorID",
		HelpDescription: `This endpoint approves AccessRequests.

		'requestor_id' designates the owner of the AccessRequest to be approved.
		`,
	}
}

func (b *BaseBackend) handleApprove(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID

	if entityID == "" {
		return logical.ErrorResponse("Token has no EntityID assigned"), logical.ErrPermissionDenied
	}

	requestorID := d.Get("requestor_id").(string)
	// if !ok {
	// 	return logical.ErrorResponse(fmt.Sprint(ok)), nil
	// }

	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	approverID := entityID

	if requestorID == approverID {
		return logical.ErrorResponse("Entities cannot approve their own requests"), logical.ErrPermissionDenied
	}

	accessRequest, err := b.GetRequest(ctx, req, requestorID)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	if accessRequest == nil {
		return &logical.Response{Warnings: []string{"Request does not exist"}}, nil
	}

	if accessRequest.Status != Pending {
		return logical.ErrorResponse(
			fmt.Sprintf(
				"Cannot Approve an AccessRequest that is not in 'pending' state (state: %s)",
				accessRequest.Status.String(),
			),
		), nil
	}

	alreadyApproved := accessRequest.isApprovedBy(approverID)
	if alreadyApproved {
		return &logical.Response{Warnings: []string{"Request already approved by this user"}}, nil
	}

	_, _, err = accessRequest.Approve(approverID) // lastApproval
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	err = b.StoreRequest(ctx, req, accessRequest)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"status": accessRequest.Status,
		},
	}, nil
}
