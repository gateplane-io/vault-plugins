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

func (b *BaseBackend) DeleteRequest(ctx context.Context, req *logical.Request, accessRequest AccessRequest) error {
	entityID := req.EntityID

	err := b.DeleteRequestFromStorage(ctx, req.Storage, accessRequest)
	if err != nil {
		b.Logger().Error("[-] Could not delete request from storage",
			"EntityID", entityID,
			"RequestorID", accessRequest.OwnerID,
			"error", err,
		)
		return fmt.Errorf("Could not delete Access Request from Backend")
	}
	return nil
}

func (b *BaseBackend) DeleteRequestFromStorage(ctx context.Context, storage logical.Storage, accessRequest AccessRequest) error {

	requestID := accessRequest.OwnerID
	err := storage.Delete(ctx, storageKeyForRequest(requestID))
	if err != nil {
		b.Logger().Error("[-] Could not delete AccessRequest",
			"RequestorID", requestID,
			"error", err,
		)
		return err
	}
	return nil
}
