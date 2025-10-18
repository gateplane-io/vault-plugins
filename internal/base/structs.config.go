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
	"fmt"
	"time"
)

type Config struct {
	RequireJustification bool `json:"require_justification"`
	RequiredApprovals    int  `json:"required_approvals"`
	// AllowRejection       bool  `json:"allow_rejection"`

	RequestTTL  time.Duration `json:"request_ttl"`
	DeleteAfter time.Duration `json:"delete_after"`
}

func NewConfig() Config {
	return Config{
		RequireJustification: false, // Default: Do not require justification
		RequiredApprovals:    1,     // Default: Require at least 1 approval
		// AllowRejection:       true,                 // Default: Allow rejection
		RequestTTL:  1 * time.Hour,  // Default: 1 hour for request TTL
		DeleteAfter: 24 * time.Hour, // Default: 24 hours for deletion
	}
}

type PluginConfig interface {
	// For struct to be a Config it needs to properly update its fields independently
	SetConfigurationKey(key string, value interface{}) error
}

func (c *Config) SetConfigurationKey(key string, value interface{}) error {
	switch key {
	case "required_approvals":
		if v, ok := value.(int); ok {
			c.RequiredApprovals = v
		} else {
			return fmt.Errorf("invalid type for required_approvals, expected int")
		}
	case "require_justification":
		if v, ok := value.(bool); ok {
			c.RequireJustification = v
		} else {
			return fmt.Errorf("invalid type for require_reason, expected bool")
		}
	case "request_ttl":
		c.RequestTTL = time.Duration(value.(int)) * time.Second
	case "delete_after":
		c.DeleteAfter = time.Duration(value.(int)) * time.Second
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
	return nil
}

type ConfigLease struct {
	Lease    time.Duration `json:"lease"`
	LeaseMax time.Duration `json:"lease_max"`
}

func NewConfigLease() ConfigLease {
	return ConfigLease{
		Lease:    30 * time.Minute,
		LeaseMax: 1 * time.Hour,
	}
}

func (c *ConfigLease) SetConfigurationKey(key string, value interface{}) error {
	switch key {
	case "lease":
		c.Lease = time.Duration(value.(int)) * time.Second
	case "lease_max":
		c.LeaseMax = time.Duration(value.(int)) * time.Second
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
	return nil
}
