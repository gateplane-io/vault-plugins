// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package config

import "fmt"

type ConfigApiVaultPeriodicToken struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

func NewConfigApiVaultPeriodicToken() ConfigApiVaultPeriodicToken {
	return ConfigApiVaultPeriodicToken{Url: "https://localhost:8200"}
}

func (c *ConfigApiVaultPeriodicToken) SetConfigurationKey(key string, value interface{}) error {
	switch key {
	case "url":
		url, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid type for 'url', expected string")
		}
		c.Url = url
	case "token":
		token, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid type for 'token', expected string")
		}
		c.Token = token
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
	return nil
}
