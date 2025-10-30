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

	"github.com/gateplane-io/vault-plugins/pkg/models"
)

/* ======================== CRUD Request*/

func (b *BaseBackend) GetRequest(ctx context.Context, req *logical.Request, requestID string) (*AccessRequest, error) {
	entityID := req.EntityID

	accessRequest, err := b.GetRequestFromStorage(ctx, req.Storage, requestID)
	if err != nil {
		b.Logger().Error("[-] Could not retrieve request from storage",
			"EntityID", entityID,
			"RequestorID", requestID,
			"error", err,
		)
		return nil, fmt.Errorf("Could not retrieve Access Request from Backend")
	}
	return accessRequest, nil
}

func (b *BaseBackend) GetRequestFromStorage(ctx context.Context, storage logical.Storage, requestID string) (*AccessRequest, error) {

	entry, err := storage.Get(ctx, storageKeyForRequest(requestID))
	if err != nil {
		b.Logger().Error("[-] Could not retrieve request from storage",
			"RequestorID", requestID,
			"error", err,
		)
		return nil, fmt.Errorf("Could not retrieve request from BaseBackend")
	}
	if entry == nil {
		b.Logger().Warn("[!] Missing request entry in storage",
			"RequestorID", requestID,
		)
		return nil, nil
	}

	var accessRequest AccessRequest
	if err := json.Unmarshal(entry.Value, &accessRequest); err != nil {
		b.Logger().Error("[-] Failed to unmarshal AccessRequest",
			"RequestorID", requestID,
			"error", err,
		)
		return nil, fmt.Errorf("Request could not be retrieved")
	}

	requestDirty := false
	// Set terminal states in case of expiration
	if (accessRequest.Status != models.Expired &&
		accessRequest.Status != models.Revoked &&
		accessRequest.Status != models.Abandoned) && requestHasExpired(accessRequest) {
		if accessRequest.Status == models.Pending {
			accessRequest.Status = models.Abandoned
		} else {
			accessRequest.Status = models.Expired
		}
		requestDirty = true
		b.Logger().Info("[*] Request status set",
			"RequestorID", requestID,
			"Status", accessRequest.Status,
			"Expiration", accessRequest.Expiration,
		)
	}

	if requestIsDeletable(accessRequest) {
		err := b.DeleteRequestFromStorage(ctx, storage, accessRequest)
		if err != nil {
			b.Logger().Error("[-] Failed to Deleted AccessRequest in Storage",
				"RequestorID", requestID,
				"Expiration", accessRequest.Expiration,
				"DeleteAfter", accessRequest.Deletion,
				"error", err,
			)
			return nil, err
		}
		b.Logger().Info("[*] Deleted AccessRequest after EOL",
			"RequestorID", requestID,
			"Expiration", accessRequest.Expiration,
			"DeleteAfter", accessRequest.Deletion,
		)
		// Return nothing as the AccessRequest was to be deleted
		return nil, nil
	}

	// Handle Approval Status
	if accessRequest.Status == models.Pending &&
		requestIsApproved(accessRequest) {
		accessRequest.Status = models.Approved
		requestDirty = true
	}
	if requestDirty {
		err := b.StoreRequestToStorage(ctx, storage, &accessRequest)
		if err != nil {
			b.Logger().Error("[-] Could not store changed AccessRequest",
				"RequestorID", requestID,
				"error", err,
			)
			return nil, err
		}
	}
	return &accessRequest, nil
}
