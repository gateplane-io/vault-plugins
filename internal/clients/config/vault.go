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

import (
	"fmt"
)

type ConfigApiVaultAppRole struct {
	Url          string `json:"url"`
	RoleID       string `json:"role_id"`
	RoleSecret   string `json:"role_secret"`
	AppRoleMount string `json:"approle_mount"`
}

func NewConfigApiVaultAppRole() ConfigApiVaultAppRole {
	return ConfigApiVaultAppRole{
		Url:          "https://localhost:8200",
		AppRoleMount: "approle",
	}
}

func (c *ConfigApiVaultAppRole) SetConfigurationKey(key string, value interface{}) error {
	switch key {
	case "url":
		c.Url = value.(string)
	case "role_id":
		c.RoleID = value.(string)
	case "role_secret":
		c.RoleSecret = value.(string)
	case "approle_mount":
		c.AppRoleMount = value.(string)
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
	return nil
}
