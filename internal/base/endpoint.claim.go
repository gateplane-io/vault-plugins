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
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/gateplane-io/vault-plugins/pkg/models"
)

func PathClaim(b *BaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: "claim",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleClaim,
		},
		HelpSynopsis: "Provides access for an approved AccessRequest.",
		HelpDescription: `This endpoint reads a user's AccessRequest
		issued by the '/request' endpoint and grants access if it is in 'approved' state,

		The successful response of this endpoint provides the access supported by the plugin.
		`,
	}
}

func (b *BaseBackend) handleClaim(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID

	if entityID == "" {
		return logical.ErrorResponse("Token has no EntityID assigned"), logical.ErrPermissionDenied
	}

	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	requestorID := req.EntityID

	accessRequest, err := b.GetRequest(ctx, req, requestorID)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}
	if accessRequest == nil {
		return &logical.Response{Warnings: []string{"Request does not exist"}}, nil
	}

	if accessRequest.Status != models.Approved {
		return logical.ErrorResponse(
			fmt.Sprintf(
				"Cannot Claim an AccessRequest that is not in 'approved' state (state: %s, approvals %d/%d)",
				accessRequest.Status,
				len(accessRequest.Approvals),
				accessRequest.RequiredApprovals,
			),
		), nil
	}

	b.Logger().Info("[+] Claiming access through the Lease Append hook",
		"RequestorID", accessRequest.OwnerID,
	)
	internalData, err := b.ClaimArray.Append(ctx, req, accessRequest.OwnerID)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}
	internalData["requestor_id"] = accessRequest.OwnerID
	b.Logger().Info("[+] Updating AccessRequest status to 'Active'",
		"RequestorID", accessRequest.OwnerID,
		"InternalData", internalData,
	)

	accessRequest.Claim()
	err = b.StoreRequest(ctx, req, accessRequest)
	if err != nil {
		_, err2 := b.ClaimArray.Remove(ctx, req, accessRequest.OwnerID, internalData)
		if err2 != nil {
			return logical.ErrorResponse(fmt.Sprint(err2)), nil
		}
		return logical.ErrorResponse(fmt.Sprint(err)), nil
	}
	// data := NormalizeMapStrings(internalData)
	b.Logger().Info("[+] Creating Response Secret",
		"RequestorID", accessRequest.OwnerID,
		// "Data", data,
		"InternalData", internalData,
		"AccessRequest", accessRequest,
	)

	respSecret := b.Secret(SecretType)
	respSecret.Type = SecretType
	resp := respSecret.Response(internalData, internalData)
	resp.Secret.TTL = accessRequest.ClaimTTL * time.Second
	b.Logger().Warn("[+] Claimed AccessRequest",
		"RequestorID", accessRequest.OwnerID,
		"ClaimTime", accessRequest.ClaimCreatedAt,
		"LeaseExpiration", resp.Secret.LeaseOptions.ExpirationTime(),
		"TTL", resp.Secret.TTL,
	)
	return resp, nil
}

func (b *BaseBackend) handleClaimRevocation(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	entityID := req.EntityID
	requestorID := req.Secret.InternalData["requestor_id"].(string)

	b.Logger().Info("[*] Revoke AccessRequest",
		"RequestorID", requestorID,
		"EntityID", entityID,
		"LeaseExpiration", req.Secret.LeaseOptions.ExpirationTime(),
		"TTL", req.Secret.TTL,
		"Revoker", req.DisplayName,
	)

	accessRequest, err := b.GetRequest(ctx, req, requestorID)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrNotFound
	}
	now := time.Now()

	if accessRequest.Status == models.Active {
		removed, err := b.ClaimArray.Remove(ctx, req, accessRequest.OwnerID, req.Secret.InternalData)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprint(err)), nil
		}
		if !removed {
			return logical.ErrorResponse(fmt.Sprint(err)), nil
		}
		accessRequest.Status = models.Revoked
		/*
			If the revocation happended by a nameless token
			and close to the real expiration (2s) we assume that
			the request is made by Vault/OpenBao core.
			This sets the request as Expired
			(Revoked is when it is manually revoked)
		*/
		if entityID == "" &&
			req.DisplayName == "" &&
			areDatesClose(
				now,
				req.Secret.LeaseOptions.ExpirationTime(),
				time.Second*2) {
			accessRequest.Status = models.Expired
		}
		err = b.StoreRequest(ctx, req, accessRequest)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprint(err)), nil
		}
	}
	b.Logger().Info("[+] AccessRequest Revoked",
		"RequestorID", requestorID,
		"EntityID", entityID,
		"LeaseExpiration", req.Secret.LeaseOptions.ExpirationTime(),
		"Status", accessRequest.Status,
		"TTL", req.Secret.TTL,
	)

	// Empty response
	return nil, nil //&logical.Response{}
}
