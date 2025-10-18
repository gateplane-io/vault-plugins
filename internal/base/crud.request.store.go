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

func (b *BaseBackend) StoreRequest(ctx context.Context, req *logical.Request, accessRequest *AccessRequest) error {
	entityID := req.EntityID

	err := b.StoreRequestToStorage(ctx, req.Storage, accessRequest)
	if err != nil {
		b.Logger().Error("[-] Could not store request from storage",
			"EntityID", entityID,
			"RequestorID", accessRequest.OwnerID,
			"error", err,
		)
		return fmt.Errorf("Could not store Access Request to Backend")
	}
	return nil
}

func (b *BaseBackend) StoreRequestToStorage(ctx context.Context, storage logical.Storage, accessRequest *AccessRequest) error {

	requestJSON, err := json.Marshal(*accessRequest)
	if err != nil {
		b.Logger().Error("[-] Could not marshal AccessRequest to JSON",
			"AccessRequest", accessRequest,
			"error", err,
		)
		return err
	}

	requestID := accessRequest.OwnerID
	err = storage.Put(ctx, &logical.StorageEntry{
		Key:   storageKeyForRequest(requestID),
		Value: requestJSON,
	})
	if err != nil {
		b.Logger().Error("[-] Could not store AccessRequest",
			"RequestorID", requestID,
			"RequestJSON", requestJSON,
			"error", err,
		)
		return err
	}
	return nil
}
