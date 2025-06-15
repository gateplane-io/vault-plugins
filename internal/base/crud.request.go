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
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"
)

/* ======================== CRUD Request*/

func (b *BaseBackend) getRequest(ctx context.Context, req *logical.Request, requestID string) (*AccessRequest, error) {
	entityID := req.EntityID
	// Not Relevant in the /claim endpoint (it's unauth'd in PolicyGate)
	/* Disallow getting an AccessRequest object to non-Entities
	if entityID == "" {
		return nil, fmt.Errorf("user entity ID is missing")
	}
	*/

	entry, err := req.Storage.Get(ctx, storageKeyForRequest(requestID))
	if err != nil {
		b.Logger().Error("[-] Could not retrieve request from storage",
			"EntityID", entityID,
			"RequestID", requestID,
			"error", err,
		)
		return nil, fmt.Errorf("Could not retrieve request from BaseBackend")
	}
	if entry == nil {
		b.Logger().Warn("[!] Missing request entry in storage",
			"EntityID", entityID,
			"RequestID", requestID,
		)
		return nil, fmt.Errorf("Request Not Found")
	}

	var accessRequest AccessRequest
	if err := json.Unmarshal(entry.Value, &accessRequest); err != nil {
		b.Logger().Error("[-] Failed to unmarshal AccessRequest",
			"EntityID", entityID,
			"RequestID", requestID,
			"error", err,
		)
		return nil, fmt.Errorf("Request could not be retrieved")
	}

	// Check ALL expirations
	if requestHasExpired(accessRequest) {
		accessRequest.Status = AccessRequestStatus(Expired)
		b.Logger().Info("[*] Request set to Expired",
			"EntityID", entityID,
			"RequestID", requestID,
			"Expiration", accessRequest.Expiration,
		)
	}
	for _, approval := range accessRequest.Approvals {
		if approvalHasExpired(approval) {
			approval.Status = ApprovalStatus(ApprovalExpired)
			b.Logger().Info("[*] Request set to Expired",
				"EntityID", entityID,
				"RequestID", requestID,
				"ApprovalID", approval.ID,
				"Expiration", approval.Expiration,
			)
		}
	}

	config, err := b.GetConfiguration(ctx, req)
	if err != nil {
		b.Logger().Error("[-] Failed to retrieve Configuration. AccessRequest Deletion is not possible",
			"EntityID", entityID,
			"RequestID", requestID,
			"error", err,
		)
		return nil, fmt.Errorf("Could not retrieve Plugin Configuration")
	}

	// Handle Access Request Deletion after EOL
	deleteAfter := config.DeleteAfter
	if requestIsDeletable(accessRequest, deleteAfter) {
		err := b.deleteRequest(ctx, req, accessRequest)
		if err != nil {
			b.Logger().Error("[-] Failed to Deleted AccessRequest after EOL",
				"EntityID", entityID,
				"RequestID", requestID,
				"Expiration", accessRequest.Expiration,
				"DeleteAfter", deleteAfter,
				"error", err,
			)
		}
		b.Logger().Info("[*] Deleted AccessRequest after EOL",
			"EntityID", entityID,
			"RequestID", requestID,
			"Expiration", accessRequest.Expiration,
			"DeleteAfter", deleteAfter,
		)
		// Return nothing as the AccessRequest was to be deleted
		return nil, nil
	}

	// Handle Approval Status
	requiredApprovals := config.RequiredApprovals
	if accessRequest.Status == Pending && requestIsApproved(accessRequest, requiredApprovals) {
		accessRequest.Status = Approved
		err = accessRequest.createGrantCode()
		err2 := b.storeGrantCodeMap(ctx, req, &accessRequest)
		err3 := b.storeRequest(ctx, req, &accessRequest)
		if err != nil {
			b.Logger().Error("[-] GrantCode could not be created",
				"EntityID", entityID,
				"RequestID", requestID,
				"error", err,
			)
			return nil, err
		} else if err2 != nil {
			b.Logger().Error("[-] GrantCode could not be mapped to AccessRequest",
				"EntityID", entityID,
				"RequestID", requestID,
				"error", err2,
			)
			return nil, err2
		} else if err3 != nil {
			b.Logger().Error("[-] Could not store Approved AccessRequest",
				"EntityID", entityID,
				"RequestID", requestID,
				"error", err3,
			)
			return nil, err3
		} else {
			b.Logger().Info("[+] GrantCode created",
				"EntityID", entityID,
				"RequestID", requestID,
			)
		}
	}
	return &accessRequest, nil
}

func (b *BaseBackend) listRequests(ctx context.Context, req *logical.Request) ([]AccessRequest, error) {
	entityID := req.EntityID
	// This is not sensitive to Entities only
	/* Disallow listing AccessRequest object to non-Entities
	if entityID == "" {
		return nil, fmt.Errorf("user entity ID is missing")
	}
	*/

	accessRequests := []AccessRequest{}

	entries, err := req.Storage.List(ctx, storageKeyForRequest(""))
	if err != nil {
		b.Logger().Error("[-] Could not list request entries",
			"EntityID", entityID,
			"error", err,
		)
		return nil, fmt.Errorf("unable to list requests: %w", err)
	}

	for _, requestID := range entries {
		// b.Logger().Info("[*] Iterating key",
		// 	"key", key,
		// )

		// getRequest refreshes the AccessRequest state (Approvals, )
		accessRequest, err := b.getRequest(ctx, req, requestID)
		if err != nil {
			b.Logger().Error("[-] Could not retrieve AccessRequest",
				"EntityID", entityID,
				"RequestID", requestID,
				"error", err,
			)
			continue
		}
		// It is possible that the fetched AccessRequest has been deleted as EOL is passed.
		// in that case 'getRequest' yields 'nil'
		if accessRequest != nil {
			accessRequests = append(accessRequests, *accessRequest)
		}
	}
	return accessRequests, nil
}

func (b *BaseBackend) deleteRequest(ctx context.Context, req *logical.Request, accessRequest AccessRequest) error {
	entityID := req.EntityID
	// This is not sensitive to Entities only
	/* Disallow listing AccessRequest object to non-Entities
	if entityID == "" {
		return nil, fmt.Errorf("user entity ID is missing")
	}
	*/
	requestID := accessRequest.ID
	err := req.Storage.Delete(ctx, storageKeyForRequest(requestID))
	if err != nil {
		b.Logger().Error("[-] Could not delete AccessRequest",
			"EntityID", entityID,
			"RequestID", requestID,
			"error", err,
		)
		return err
	}
	return nil
}

func (b *BaseBackend) storeRequest(ctx context.Context, req *logical.Request, accessRequest *AccessRequest) error {
	entityID := req.EntityID
	// This is not sensitive to Entities only
	/* Disallow listing AccessRequest object to non-Entities
	if entityID == "" {
		return nil, fmt.Errorf("user entity ID is missing")
	}
	*/

	requestJSON, err := json.Marshal(*accessRequest)
	if err != nil {
		b.Logger().Error("[-] Could not marshal AccessRequest to JSON",
			"AccessRequest", accessRequest,
			"error", err,
		)
		return err
	}

	requestID := accessRequest.ID
	err = req.Storage.Put(ctx, &logical.StorageEntry{
		Key:   storageKeyForRequest(requestID),
		Value: requestJSON,
	})
	if err != nil {
		b.Logger().Error("[-] Could not store AccessRequest",
			"EntityID", entityID,
			"RequestID", requestID,
			"RequestJSON", requestJSON,
			"error", err,
		)
		return err
	}
	return nil
}
