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
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/gateplane-io/vault-plugins/internal/utils"
)

/* Storage Keys */
const RequestKey = "request"
const ConfigKey = "config"
const ConfigLeaseKey = "config/lease"

type BaseBackend struct {
	*framework.Backend
	BaseMutex  sync.Mutex
	ClaimArray *utils.CallbackArray
}

func (b *BaseBackend) Initialize(ctx context.Context, req *logical.InitializationRequest) error {
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()
	b.Logger().Info("Initializing plugin Base configuration")

	// If configuration NOT THERE - load default
	config := NewConfig()
	// StoreConfigurationToStorageIfNotPresent[*Config](ctx, b, req.Storage, &config, ConfigKey)
	exists, err := StoreConfigurationToStorageIfNotPresent[*Config](ctx, b, req.Storage, &config, ConfigKey)
	b.Logger().Info("GatePlane Base initialized with default configuration",
		"configuration", config,
		"Existing", exists,
		"Error", err,
	)
	configLease := NewConfigLease()
	// StoreConfigurationToStorageIfNotPresent(ctx, b, req.Storage, &configLease, ConfigLeaseKey)
	exists, err = StoreConfigurationToStorageIfNotPresent(ctx, b, req.Storage, &configLease, ConfigLeaseKey)

	b.Logger().Info("GatePlane Base initialized with default configuration",
		"configuration", config,
		"Existing", exists,
		"Error", err,
	)
	return nil
}
