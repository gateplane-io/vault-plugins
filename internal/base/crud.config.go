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

/* ======================== CRUD Config*/

func GetConfiguration[T PluginConfig](ctx context.Context, b *BaseBackend, req *logical.Request, path string) (T, error) {
	// If called without 'path', assume it from the request
	// (means it comes from a /config* endpoint)
	configStorageKey := path
	if configStorageKey == "" {
		configStorageKey = req.Path
	}

	return GetConfigurationFromStorage[T](ctx, b, req.Storage, configStorageKey)
}

func GetConfigurationFromStorage[T PluginConfig](ctx context.Context, b *BaseBackend, storage logical.Storage, path string) (T, error) {
	var zero T

	entry, err := storage.Get(ctx, path)
	if err != nil || entry == nil {
		b.Logger().Error("[-] Could not retrieve configuration from storage",
			"error", err,
		)
		return zero, fmt.Errorf("Could not retrieve configuration from BaseBackend")
	}

	var config T
	if err := json.Unmarshal(entry.Value, &config); err != nil {
		b.Logger().Error("[-] Failed to unmarshal Config",
			"error", err,
		)
		return zero, fmt.Errorf("Configuration could not be retrieved")
	}
	return config, nil
}

func StoreConfigurationToStorage[T PluginConfig](ctx context.Context, b *BaseBackend, storage logical.Storage, config T, path string) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		b.Logger().Error("[-] Could not marshal configuration to JSON",
			"path", path,
			"config", config,
			"error", err,
		)
		return err
	}
	err = storage.Put(ctx, &logical.StorageEntry{
		Key:   path,
		Value: configJSON,
	})
	if err != nil {
		b.Logger().Error("[-] Could not store configuration to storage",
			"path", path,
			"config", config,
			"configJSON", configJSON,
			"error", err,
		)
		return err
	}
	return nil
}

func StoreConfiguration[T PluginConfig](ctx context.Context, b *BaseBackend, req *logical.Request, config T, path string) error {
	// If called without 'path', assume it from the request
	// (means it comes from a /config* endpoint)
	configStorageKey := path
	if configStorageKey == "" {
		configStorageKey = req.Path
	}

	return StoreConfigurationToStorage(ctx, b, req.Storage, config, configStorageKey)
}

func StoreConfigurationToStorageIfNotPresent[T PluginConfig](ctx context.Context, b *BaseBackend, storage logical.Storage, config T, path string) (bool, error) {
	// var zero T

	_, err := GetConfigurationFromStorage[T](ctx, b, storage, path)
	if err != nil {
		return false, StoreConfigurationToStorage[T](ctx, b, storage, config, path)
	}
	return true, nil
}
