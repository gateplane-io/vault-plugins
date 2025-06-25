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
	"fmt"
	"time"

	"github.com/gateplane-io/vault-plugins/internal/base"
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

/*
Configuration of Direct Okta API Access
*/

type OktaApiConfig struct {
	OrgUrl   string `json:"org_url"`
	ApiToken string `json:"api_token"`

	// OktaEntityKey: "name", // Default when 'user_claim=sub' in OIDC
	OktaEntityKey string `json:"okta_entity_key"`

	// used to check Entity Aliases
	// https://pkg.go.dev/github.com/hashicorp/vault/sdk@v0.18.0/logical#Alias
	// Mount Accessor: mount_accessor must be of type OIDC
	OktaOIDCMountAccessor string `json:"auth_mount_accessor"`
}

func NewDefaultOktaApiConfig() OktaApiConfig {
	return OktaApiConfig{
		// Default when 'user_claim=sub' in OIDC
		OktaEntityKey: "name",
	}
}

func (c *OktaApiConfig) SetConfigurationKey(key string, value interface{}) error {
	switch key {
	case "org_url":
		if v, ok := value.(string); ok {
			c.OrgUrl = v
		} else {
			return fmt.Errorf("invalid type for 'org_url', expected string")
		}
	case "api_token":
		if v, ok := value.(string); ok {
			c.ApiToken = v
		} else {
			return fmt.Errorf("invalid type for 'api_token', expected string")
		}
	case "okta_entity_key":
		if v, ok := value.(string); ok {
			c.OktaEntityKey = v
		} else {
			return fmt.Errorf("invalid type for 'okta_entity_key', expected string")
		}
	case "auth_mount_accessor":
		if v, ok := value.(string); ok {
			c.OktaOIDCMountAccessor = v
		} else {
			return fmt.Errorf("invalid type for 'auth_mount_accessor', expected string")
		}
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
	return nil
}

func (config *OktaApiConfig) Test(userID string, groupID string) error {

	client, err := createOktaClient(config)
	if err != nil {
		return err
	}

	ttl := 1 * time.Second

	err = oktaAddRemoveFromGroup(context.TODO(), client, ttl, groupID, userID)

	return err
}
