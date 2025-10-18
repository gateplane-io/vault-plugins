// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package mock

import (
	"sync"
	// "time"
	"context"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/gateplane-io/vault-plugins/internal/base"

	"github.com/gateplane-io/vault-plugins/internal/utils"
)

type Backend struct {
	*base.BaseBackend
	Mutex sync.Mutex
}

func (b *Backend) Initialize(ctx context.Context, req *logical.InitializationRequest) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	b.Logger().Info("Initializing plugin configuration")

	// Initialize the base
	err := b.BaseBackend.Initialize(ctx, req)
	if err != nil {
		b.Logger().Error("Could not initialize Plugin Base")
		return err
	}

	b.ClaimArray = utils.NewCallbackArray(
		(func(ctx context.Context, requ *logical.Request, ownerID string) (map[string]interface{}, error) { // Append

			areq, err := b.GetRequest(ctx, requ, ownerID)
			b.Logger().Warn(
				"Function Claim",
				"RequestorID", ownerID,
				"AccessRequest", areq,
				"Error", err,
			)
			ret := map[string]interface{}{
				"claimed": true,
			}
			return ret, err
		}),
		(func(ctx context.Context, requ *logical.Request, ownerID string, internalData map[string]interface{}) error { // Remove

			areq, err := b.GetRequest(ctx, requ, ownerID)
			b.Logger().Warn(
				"Function UnClaim",
				"RequestorID", ownerID,
				"AccessRequest", areq,
				"InternalData", internalData,
				"Error", err,
			)
			return err
		}),
	)

	b.Logger().Info("GatePlane Mock initialized with default configuration")
	return nil
}
