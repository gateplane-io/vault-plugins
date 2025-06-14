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

	"github.com/gateplane-io/vault-plugins/internal/base"
)

/*
Embedding new fields to existing Struct
and re-implementing Methods the CRUD
*/

type Config struct {
	*base.Config
	Policies []string `json:"policies"`
}

func NewDefaultConfig() Config {
	bConfig := base.NewDefaultConfig()
	return Config{
		Config:   &bConfig,
		Policies: []string{},
	}
}

func (c *Config) SetConfigurationKey(key string, value interface{}) error {
	switch key {
	// add the case of the new key
	case "policies":
		if v, ok := value.([]string); ok {
			c.Policies = v
		} else {
			return fmt.Errorf("invalid type for policies, expected []string")
		}
	default:
		return c.Config.SetConfigurationKey(key, value)
	}
	return nil
}
