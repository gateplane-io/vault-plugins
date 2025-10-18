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
	"fmt"
)

type ConfigAccess struct {
	GroupID   string `json:"okta_group_id"`
	GroupName string `json:"okta_group_name"`
}

func (c *ConfigAccess) SetConfigurationKey(key string, value interface{}) error {
	switch key {
	case "okta_group_id":
		if v, ok := value.(string); ok {
			c.GroupID = v
		} else {
			return fmt.Errorf("invalid type for 'okta_group_id', expected string")
		}
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
	return nil
}

func NewConfigAccess() ConfigAccess {
	return ConfigAccess{}
}
