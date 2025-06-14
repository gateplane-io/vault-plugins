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

	"github.com/hashicorp/vault/sdk/logical"
)

func (b *BaseBackend) GetAccessRequestByGrantCode(ctx context.Context, req *logical.Request, grantCode string) (*AccessRequest, error) {

	requestIDEntry, err := req.Storage.Get(ctx, fmt.Sprintf("%s/%s", GrantCodeKey, grantCode))
	if err != nil {
		b.Logger().Error("[-] Could not retrieve AccessRequest ID from GrantCode",
			"FailingGrantCode", grantCode,
			"error", err,
		)
		return nil, err
	}

	if requestIDEntry == nil {
		err = fmt.Errorf("Provided GrantCode does not map to any AccessRequests")
		b.Logger().Error("[-] Provided GrantCode does not map to any AccessRequests",
			"FailingGrantCode", grantCode,
			"error", err,
		)
		return nil, err
	}

	requestID := string(requestIDEntry.Value)

	accessRequest, err := b.getRequest(ctx, req, requestID)
	if err != nil {
		b.Logger().Error("[-] Could not retrieve AccessRequest from GrantCode",
			"RequestID", requestID,
			// Might need to remove from logs
			"FailingGrantCode", grantCode,
			"error", err,
		)
		return nil, err
	}

	if accessRequest.GrantCode != grantCode {
		err = fmt.Errorf("GrantCode does not match the stored one")
		b.Logger().Error("[-] GrantCode stored in AccessRequest does not match the provided one",
			"RequestID", requestID,
			"FailingGrantCode", grantCode,
			"StoredFailingGrantCode", accessRequest.GrantCode,
			"error", err,
		)
		return nil, err
	}

	err = req.Storage.Delete(ctx, fmt.Sprintf("%s/%s", GrantCodeKey, grantCode))
	if err != nil {
		b.Logger().Error("[-] Could not delete GrantCode mapping to AccessRequest",
			"RequestID", requestID,
			"FailingGrantCode", grantCode,
			"error", err,
		)
		return nil, err
	}

	accessRequest.Status = AccessRequestStatus(Active)
	accessRequest.GrantCode = ""

	err = b.storeRequest(ctx, req, accessRequest)
	if err != nil {
		b.Logger().Error("[-] Could not store Active AccessRequest",
			"RequestID", requestID,
			"FailingGrantCode", grantCode,
			"error", err,
		)
		return nil, err
	}
	return accessRequest, nil
}

func (b *BaseBackend) storeGrantCodeMap(ctx context.Context, req *logical.Request, accessRequest *AccessRequest) error {

	requestID := accessRequest.ID
	requestIDEntry := []byte(requestID)

	err := req.Storage.Put(ctx, &logical.StorageEntry{
		Key:   fmt.Sprintf("%s/%s", GrantCodeKey, accessRequest.GrantCode),
		Value: requestIDEntry,
	})
	if err != nil {
		b.Logger().Error("[-] Could not store GrantCode Map",
			"EntityID", req.EntityID,
			"RequestID", requestID,
			"error", err,
		)
		return err
	}
	return nil
}
