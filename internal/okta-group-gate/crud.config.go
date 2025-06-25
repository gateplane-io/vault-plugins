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

/* ======================== CRUD Config*/
const ConfigKey = "config"

func (b *Backend) GetConfiguration(ctx context.Context, req *logical.Request) (*Config, error) {

	entry, err := req.Storage.Get(ctx, ConfigKey)
	if err != nil {
		b.Logger().Error("[-] Could not retrieve configuration from storage",
			"error", err,
		)
		return nil, fmt.Errorf("Could not retrieve configuration from PluginBackend")
	}

	var config Config
	if err := json.Unmarshal(entry.Value, &config); err != nil {
		b.Logger().Error("[-] Failed to unmarshal Config",
			"error", err,
		)
		return nil, fmt.Errorf("Configuration could not be retrieved")
	}
	return &config, nil
}

func (b *Backend) StoreConfigurationToStorage(ctx context.Context, storage logical.Storage, config *Config) error {
	configJSON, err := json.Marshal(*config)
	if err != nil {
		b.Logger().Error("[-] Could not marshal configuration to JSON",
			"config", config,
			"error", err,
		)
		return err
	}

	err = storage.Put(ctx, &logical.StorageEntry{
		Key:   ConfigKey,
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
	return nil
}

func (b *Backend) StoreConfiguration(ctx context.Context, req *logical.Request, config *Config) error {
	return b.StoreConfigurationToStorage(ctx, req.Storage, config)
}
