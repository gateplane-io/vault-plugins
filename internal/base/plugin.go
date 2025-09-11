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

	models "github.com/gateplane-io/vault-plugins/pkg/models/base"
)

type BaseBackend struct {
	*framework.Backend
	BaseMutex sync.Mutex
}

func (b *BaseBackend) Initialize(ctx context.Context, req *logical.InitializationRequest) error {
	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()
	b.Logger().Info("Initializing plugin configuration")

	defaultConfig := models.NewDefaultConfig()
	b.StoreConfigurationToStorage(ctx, req.Storage, &defaultConfig)

	b.Logger().Info("Vault auth plugin initialized with default configuration",
		"configuration", defaultConfig,
	)
	return nil
}
