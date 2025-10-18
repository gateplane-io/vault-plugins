// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package policy_gate

import (
	"fmt"
)

type ConfigAccess struct {
	Policies []string `json:"policies"`
}

func (c *ConfigAccess) SetConfigurationKey(key string, value interface{}) error {
	switch key {
	case "policies":
		if v, ok := value.([]string); ok {
			c.Policies = v
		} else {
			return fmt.Errorf("invalid type for policies, expected []string")
		}
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
	return nil
}

func NewConfigAccess() ConfigAccess {
	return ConfigAccess{
		Policies: []string{},
	}
}
