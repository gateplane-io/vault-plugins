// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package okta_group_gate

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/okta/okta-sdk-golang/v5/okta"

	"github.com/gateplane-io/vault-plugins/internal/base"

	models "github.com/gateplane-io/vault-plugins/pkg/models/okta-group-gate"
)

type Backend struct {
	*base.BaseBackend
	oktaClient *okta.APIClient
}

func (b *Backend) Initialize(ctx context.Context, req *logical.InitializationRequest) error {
	b.BaseBackend.BaseMutex.Lock()
	defer b.BaseBackend.BaseMutex.Unlock()
	b.Logger().Info("Initializing plugin configuration")

	// hadle defaults like this
	// https://github.com/jfrog/vault-plugin-secrets-artifactory/blob/master/backend.go#L86
	defaultConfig := models.NewDefaultConfig()
	b.StoreConfigurationToStorage(ctx, req.Storage, &defaultConfig)

	if b.oktaClient == nil {
		// try if there is already an api config in place
		// b.oktaClient = createOktaClient(fetch-config)
		defaultOktaApiConfig := NewDefaultOktaApiConfig()
		b.StoreOktaApiConfigurationToStorage(ctx, req.Storage, &defaultOktaApiConfig)
	}
	b.Logger().Info("Vault auth plugin initialized with default configuration",
		"configuration", defaultConfig,
	)
	return nil
}
