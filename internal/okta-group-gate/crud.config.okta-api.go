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
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"
)

const OktaConfigKey = "config/okta/api"

// =================== Okta API Config

func (b *Backend) GetOktaApiConfiguration(ctx context.Context, req *logical.Request) (*OktaApiConfig, error) {

	entry, err := req.Storage.Get(ctx, OktaConfigKey)
	if err != nil {
		b.Logger().Error("[-] Could not retrieve configuration from storage",
			"error", err,
		)
		return nil, fmt.Errorf("Could not retrieve configuration from PluginBackend")
	}

	var config OktaApiConfig
	if err := json.Unmarshal(entry.Value, &config); err != nil {
		b.Logger().Error("[-] Failed to unmarshal OktaApiConfig",
			"error", err,
		)
		return nil, fmt.Errorf("Configuration could not be retrieved")
	}
	return &config, nil
}

func (b *Backend) StoreOktaApiConfigurationToStorage(ctx context.Context, storage logical.Storage, config *OktaApiConfig) error {
	configJSON, err := json.Marshal(*config)
	if err != nil {
		b.Logger().Error("[-] Could not marshal configuration to JSON",
			"config", config,
			"error", err,
		)
		return err
	}

	err = storage.Put(ctx, &logical.StorageEntry{
		Key:   OktaConfigKey,
		Value: configJSON,
	})
	if err != nil {
		b.Logger().Error("[-] Could not store configuration to storage",
			"config", config,
			"configJSON", configJSON,
			"error", err,
		)
		return err
	}

	b.oktaClient, err = createOktaClient(config)
	if err != nil {
		b.Logger().Error("[-] Could not store Okta Client",
			"config", config,
			"error", err,
		)
		return err
	}

	return nil
}

func (b *Backend) StoreOktaApiConfiguration(ctx context.Context, req *logical.Request, config *OktaApiConfig) error {
	return b.StoreOktaApiConfigurationToStorage(ctx, req.Storage, config)
}
