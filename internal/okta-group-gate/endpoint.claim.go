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

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// Path for user to claim token
func PathClaim(b *Backend) *framework.Path {
	return &framework.Path{
		Pattern: "claim",
		Fields: map[string]*framework.FieldSchema{
			"grant_code": {
				Type:        framework.TypeString,
				Description: "The GrantCode found in the approved AccessRequest",
				Required:    true,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleClaim,
		},
	}
}

func (b *Backend) handleClaim(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	grantCode := d.Get("grant_code").(string)

	config, err := b.GetConfiguration(ctx, req)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	accessRequest, err := b.GetAccessRequestByGrantCode(ctx, req, grantCode)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}
	entityID := req.EntityID

	reqEntityID := accessRequest.ID

	if reqEntityID != entityID {
		return logical.ErrorResponse("This 'grant_code' has been issued for another Entity"), logical.ErrPermissionDenied
	}

	view := b.System()
	entity, err := view.EntityInfo(entityID)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	oktaConfig, err := b.GetOktaApiConfiguration(ctx, req)
	if err != nil {
		b.Logger().Info("[-] Could not retrieve Configuration",
			"EntityID", entityID,
			"err", err,
		)
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	oktaUserId := ""
	for _, alias := range entity.Aliases {
		if alias.MountAccessor != oktaConfig.OktaOIDCMountAccessor {
			continue
		}
		// alias.Metadata
		oktaUserId = alias.Name // TODO: This will be de-hardcoded later
	}

	ttl := view.DefaultLeaseTTL() / time.Second

	b.Logger().Info("[+] Claimed Access",
		"entityID", entityID,
		"OktaUserID", oktaUserId,
		"ttl", ttl,
	)

	groupName, err := getGroupNameById(ctx, b.oktaClient, config.OktaGroupID)
	if err != nil {
		b.Logger().Info("[-] Could not retrieve Okta Group details from API",
			"EntityID", entityID,
			"OktaGroupID", config.OktaGroupID,
			"OrgUrl", oktaConfig.OrgUrl,
			"err", err,
		)
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	err = oktaConfig.Test(oktaUserId, config.OktaGroupID)
	if err != nil {
		b.Logger().Info("[-] Testing Okta add/remove to Group failed",
			"EntityID", entityID,
			"OrgUrl", oktaConfig.OrgUrl,
			"OktaUserID", oktaUserId,
			"OktaGroupName", groupName,
			"err", err,
		)
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	expiration, err := b.StoreMembership(ctx, req, oktaUserId, ttl)
	if err != nil {
		b.Logger().Info("[-] Group Membership addition failed",
			"EntityID", entityID,
			"OktaUserID", oktaUserId,
			"ExpirationDate", expiration,
			"err", err,
		)
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	err = oktaAddToGroup(ctx, b.oktaClient, config.OktaGroupID, oktaUserId)
	if err != nil {
		b.Logger().Info("[-] Failed to add Okta User to Group",
			"EntityID", entityID,
			"OktaUserID", oktaUserId,
			"OktaGroupID", config.OktaGroupID,
			"OktaGroupName", groupName,
			"err", err,
		)
		return logical.ErrorResponse(fmt.Sprint(err)), logical.ErrMissingRequiredState
	}

	return &logical.Response{Data: map[string]interface{}{
		"exp":             expiration,
		"okta_user_id":    oktaUserId,
		"okta_group_name": groupName,
	}}, nil

}
