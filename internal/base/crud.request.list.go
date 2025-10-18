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

func (b *BaseBackend) ListRequests(ctx context.Context, req *logical.Request) ([]AccessRequest, error) {
	entityID := req.EntityID

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
		accessRequest, err := b.GetRequest(ctx, req, requestID)
		if err != nil {
			b.Logger().Error("[-] Could not retrieve AccessRequest",
				"EntityID", entityID,
				"RequestorID", requestID,
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
