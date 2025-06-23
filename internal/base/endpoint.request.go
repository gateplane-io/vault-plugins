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
	// "encoding/json"
	"fmt"
	// "time"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// Path for user to request access
func PathRequest(b *BaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: "request/?",
		Fields: map[string]*framework.FieldSchema{
			"reason": {
				Type:        framework.TypeString,
				Description: "Reason of requested access",
				Required:    false,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleRequestUpdate,
			logical.ListOperation:   b.handleRequestList,
			logical.ReadOperation:   b.handleRequestRead,
		},

		HelpSynopsis: "Creates AccessRequests and checks its status",
		HelpDescription: `This endpoint can create AccessRequests for a requestor (issuer of 'update'),
		check the AccessRequest created by the requestor (using 'read')
		and list all AccessRequests created by this backend (using 'list').

		The 'reason' parameter can be mandatory if set so in '/config' endpoint.
		`,
	}
}

func (b *BaseBackend) handleRequestUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID
	if entityID == "" {
		return logical.ErrorResponse("Token has no EntityID assigned"), logical.ErrPermissionDenied
	}

	reason := d.Get("reason").(string)

	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	exists, _ := b.getRequest(ctx, req, entityID)
	overwrite := exists != nil

	config, err := b.GetConfiguration(ctx, req)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	requireReason := config.RequireReason

	if requireReason && len(strings.TrimSpace(reason)) == 0 {
		return logical.ErrorResponse("A 'reason' parameter is required by this BaseBackend"), logical.ErrPermissionDenied
	}

	b.Logger().Info("[+] Access Requested",
		"EntityID", entityID,
		"PreviousRequestExistence", overwrite,
		"Reason", reason,
		"ReasonRequired", requireReason,
		"ReasonNoWhitspaceLength", len(strings.TrimSpace(reason)),
	)

	accessRequest := NewAccessRequest(config, entityID, reason)

	err = b.storeRequest(ctx, req, &accessRequest)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"overwrite":  overwrite,
			"iat":        accessRequest.CreatedAt,
			"exp":        accessRequest.Expiration,
			"reason":     accessRequest.Reason,
			"status":     accessRequest.Status.String(),
			"request_id": accessRequest.ID,
		},
	}, nil
}

func (b *BaseBackend) handleRequestRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID
	requestID := entityID

	if entityID == "" {
		return logical.ErrorResponse("Token has no EntityID assigned"), logical.ErrPermissionDenied
	}
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	accessRequest, err := b.getRequest(ctx, req, requestID)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}
	if accessRequest == nil {
		return &logical.Response{Data: map[string]interface{}{}}, nil
	}

	data := map[string]interface{}{
		"iat":        accessRequest.CreatedAt,
		"exp":        accessRequest.Expiration,
		"reason":     accessRequest.Reason,
		"status":     accessRequest.Status.String(),
		"request_id": accessRequest.ID,
	}

	// This check is always true now, yet,
	// there will be the possibility to query AccessRequests not belonging to oneself
	if entityID == requestID {
		// Only reveal the GrantCode if the Requestor is the Caller
		data["grant_code"] = accessRequest.GrantCode
	}

	return &logical.Response{
		Data: data,
	}, nil
}

func (b *BaseBackend) handleRequestList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID

	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	config, err := b.GetConfiguration(ctx, req)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}
	requiredApprovals := config.RequiredApprovals

	resultsFull := map[string]interface{}{}
	results := []string{}

	accessRequests, err := b.listRequests(ctx, req)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	var numOfvalidApprovals int
	for _, accessRequest := range accessRequests {
		numOfvalidApprovals = validApprovalsNum(accessRequest)
		_, haveApproved := accessRequest.Approvals[entityID]

		results = append(results, accessRequest.ID)
		resultsFull[accessRequest.ID] = map[string]interface{}{
			"request_id":          accessRequest.ID,
			"iat":                 accessRequest.CreatedAt,
			"exp":                 accessRequest.Expiration,
			"reason":              accessRequest.Reason,
			"num_approvals":       numOfvalidApprovals,
			"remaining_approvals": requiredApprovals - numOfvalidApprovals,
			"status":              accessRequest.Status.String(),
			"have_approved":       haveApproved,
		}
	}

	return logical.ListResponseWithInfo(
		results,
		resultsFull, // for the 'vault list -detailed path/' command
	), nil
}
