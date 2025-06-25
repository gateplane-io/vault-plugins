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

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// Path for plugin configuration
func PathOktaApiTest(b *Backend) *framework.Path {
	return &framework.Path{
		Pattern: "test/okta/api",
		Fields: map[string]*framework.FieldSchema{
			"user_id": {
				Type:        framework.TypeString,
				Description: "An Okta User ID, email or name to test functionality against",
				Required:    true,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleOktaApiTest,
		},

		HelpSynopsis: "Tests the Okta API access",
		HelpDescription: `This endpoint tests the Okta API configuration.
		<TODO>
		`,
	}
}

func (b *Backend) handleOktaApiTest(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entityID := req.EntityID

	b.BaseMutex.Lock()
	defer b.BaseMutex.Unlock()

	userId, exists := d.GetOk("user_id")
	if exists == false {
		b.Logger().Info("[-] Testing Okta Configuration, 'user_id' not set",
			"EntityID", entityID,
			"OktaUserID", userId,
		)
		return logical.ErrorResponse("'user_id' must be set"), logical.ErrMissingRequiredState
	}

	config, err := b.GetConfiguration(ctx, req)
	if err != nil {
		b.Logger().Info("[-] Could not retrieve Configuration",
			"EntityID", entityID,
			"OktaUserID", userId,
			"err", err,
		)
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	oktaConfig, err := b.GetOktaApiConfiguration(ctx, req)
	if err != nil {
		b.Logger().Info("[-] Could not retrieve Configuration",
			"EntityID", entityID,
			"OktaUserID", userId,
			"err", err,
		)
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	err = oktaConfig.Test(userId.(string), config.OktaGroupID)
	if err != nil {
		b.Logger().Info("[-] Testing Okta Configuration failed",
			"EntityID", entityID,
			"OktaUserID", userId,
			"OrgUrl", oktaConfig.OrgUrl,
			"err", err,
		)
		return &logical.Response{Data: map[string]interface{}{
			"test": "fail",
		}}, nil
	}

	return &logical.Response{Data: map[string]interface{}{
		"test": "success",
	}}, nil
}
