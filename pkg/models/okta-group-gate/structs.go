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

	"github.com/gateplane-io/vault-plugins/pkg/models/base"
)

/*
Embedding new fields to existing Struct
and re-implementing Methods the CRUD
*/

type Config struct {
	*base.Config
	OktaGroupID string `json:"okta_group_id"`
}

func NewDefaultConfig() Config {
	bConfig := base.NewDefaultConfig()
	return Config{
		Config:      &bConfig,
		OktaGroupID: "",
	}
}

func (c *Config) SetConfigurationKey(key string, value interface{}) error {
	switch key {
	// add the case of the new key
	case "okta_group_id":
		if v, ok := value.(string); ok {
			c.OktaGroupID = v
		} else {
			return fmt.Errorf("invalid type for 'okta_group_id', expected string")
		}

	default:
		return c.Config.SetConfigurationKey(key, value)
	}
	return nil
}
