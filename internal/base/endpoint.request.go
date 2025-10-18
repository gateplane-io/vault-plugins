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
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/gateplane-io/vault-plugins/pkg/responses"
)

// Path for user to request access
func PathRequest(b *BaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: "request/?",
		Fields: map[string]*framework.FieldSchema{
			"justification": {
				Type:        framework.TypeString,
				Description: "Reason/Objective for requesting access",
				Required:    false,
			},
			"ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Duration of the requested access",
				Required:    false,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleRequestUpdate,
			logical.ListOperation:   b.handleRequestList,
			logical.ReadOperation:   b.handleRequestRead,
		},

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: "gateplane-base",
			OperationSuffix: "AccessRequest",
		},

		HelpSynopsis: "Creates AccessRequests and checks its status",
		HelpDescription: `This endpoint can create AccessRequests for a requestor (using 'update'),
		check the AccessRequest created by the requestor (using 'read')
		and list all AccessRequests created by this backend (using 'list').

		The 'justification' parameter can be mandatory if 'require_justification' is set under '/config' endpoint.

		The 'ttl' parameter is the duration that the requested access will be in effect (must be below 'lease_max')
		`,
	}
}

func (b *BaseBackend) handleRequestUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID
	if entityID == "" {
		return logical.ErrorResponse("Token has no EntityID assigned"), logical.ErrPermissionDenied
	}

	justification := d.Get("justification").(string)
	ttl := time.Duration(d.Get("ttl").(int))

	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	exists, _ := b.GetRequest(ctx, req, entityID)
	overwrite := exists != nil

	// IF exists and status: active, Revoke before re-creating
	config, err := GetConfiguration[*Config](ctx, b, req, ConfigKey)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	configLease, err := GetConfiguration[*ConfigLease](ctx, b, req, ConfigLeaseKey)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	b.Logger().Info("[+] Access Requested",
		"EntityID", entityID,
		"PreviousRequestExistence", overwrite,
		"Justification", justification,
		"JustificationRequired", config.RequireJustification,
		"JustificationNoWhitspaceLength", len(strings.TrimSpace(justification)),
		"ClaimTTL", ttl,
	)

	accessRequest, err := NewAccessRequest(config, configLease, entityID, ttl, justification)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrPermissionDenied
	}

	b.Logger().Info("[+] Access Request Created",
		"Object", accessRequest,
	)

	err = b.StoreRequest(ctx, req, accessRequest)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	responseObj := responses.AccessRequestCreationResponse{
		Overwrite:     overwrite,
		Justification: accessRequest.Justification,
		OwnerID:       accessRequest.OwnerID,

		CreatedAt:  accessRequest.CreatedAt.Unix(),
		Expiration: accessRequest.Expiration.Unix(),
		Deletion:   accessRequest.Deletion.Unix(),

		RequiredApprovals: accessRequest.RequiredApprovals,
		Status:            UncapitalizeFirstLetter(accessRequest.Status.String()),
		NumOfApprovals:    len(accessRequest.Approvals),

		ClaimTTL: accessRequest.ClaimTTL,
	}

	responseData, err := StructToMap(responseObj)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	return &logical.Response{Data: responseData}, nil
}

func (b *BaseBackend) handleRequestRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID
	requestID := entityID

	if entityID == "" {
		return logical.ErrorResponse("Token has no EntityID assigned"), logical.ErrPermissionDenied
	}
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	accessRequest, err := b.GetRequest(ctx, req, requestID)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}
	if accessRequest == nil {
		return &logical.Response{Warnings: []string{"Request does not exist"}}, nil
	}

	responseObj := responses.AccessRequestResponse{
		Justification: accessRequest.Justification,
		OwnerID:       accessRequest.OwnerID,

		CreatedAt:  accessRequest.CreatedAt.Unix(),
		Expiration: accessRequest.Expiration.Unix(),
		Deletion:   accessRequest.Deletion.Unix(),

		RequiredApprovals: accessRequest.RequiredApprovals,
		Status:            UncapitalizeFirstLetter(accessRequest.Status.String()),
		NumOfApprovals:    len(accessRequest.Approvals),

		ClaimTTL:       accessRequest.ClaimTTL,
		ClaimCreatedAt: accessRequest.ClaimCreatedAt.Unix(),

		HaveApproved: accessRequest.isApprovedBy(entityID),
	}

	responseData, err := StructToMap(responseObj)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}

	return &logical.Response{Data: responseData}, nil
}

func (b *BaseBackend) handleRequestList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID

	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	resultsFull := map[string]interface{}{}
	results := []string{}

	accessRequests, err := b.ListRequests(ctx, req)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	for _, accessRequest := range accessRequests {
		results = append(results, accessRequest.OwnerID)

		responseObj := responses.AccessRequestResponse{
			Justification: accessRequest.Justification,
			OwnerID:       accessRequest.OwnerID,

			CreatedAt:  accessRequest.CreatedAt.Unix(),
			Expiration: accessRequest.Expiration.Unix(),
			Deletion:   accessRequest.Deletion.Unix(),

			RequiredApprovals: accessRequest.RequiredApprovals,
			Status:            UncapitalizeFirstLetter(accessRequest.Status.String()),
			NumOfApprovals:    len(accessRequest.Approvals),

			ClaimTTL:       accessRequest.ClaimTTL,
			ClaimCreatedAt: accessRequest.ClaimCreatedAt.Unix(),

			HaveApproved: accessRequest.isApprovedBy(entityID),
		}

		responseData, err := StructToMap(responseObj)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprint(err)), nil
		}
		resultsFull[accessRequest.OwnerID] = responseData
	}

	return logical.ListResponseWithInfo(
		results,
		resultsFull, // for the 'vault list -detailed path/' command
	), nil
}
