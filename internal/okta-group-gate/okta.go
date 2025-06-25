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
	"github.com/okta/okta-sdk-golang/v5/okta"
	"time"
)

func oktaAddToGroup(ctx context.Context, client *okta.APIClient, groupId string, userId string) error {
	_, err := client.GroupAPI.AssignUserToGroup(ctx, groupId, userId).Execute()
	return err
}

func oktaRemoveFromGroup(ctx context.Context, client *okta.APIClient, groupId string, userId string) error {
	_, err := client.GroupAPI.UnassignUserFromGroup(ctx, groupId, userId).Execute()
	return err
}

func oktaAddRemoveFromGroup(ctx context.Context, client *okta.APIClient, ttl time.Duration, groupId string, userId string) error {

	_, err := client.GroupAPI.AssignUserToGroup(ctx, groupId, userId).Execute()
	if err != nil {
		return err
	}

	timer := time.After(ttl)
	<-timer

	_, err = client.GroupAPI.UnassignUserFromGroup(ctx, groupId, userId).Execute()
	if err != nil {
		return err
	}
	return nil
}

func getGroupNameById(ctx context.Context, client *okta.APIClient, groupId string) (string, error) {

	group, _, err := client.GroupAPI.GetGroup(ctx, groupId).Execute()
	if err != nil {
		return "", err
	}

	return *group.Profile.Name, nil
}

func createOktaClient(oktaApiConfig *OktaApiConfig) (*okta.APIClient, error) {

	config, err := okta.NewConfiguration(
		okta.WithOrgUrl(oktaApiConfig.OrgUrl),
		okta.WithToken(oktaApiConfig.ApiToken),
	)
	if err != nil {
		return nil, err
	}

	return okta.NewAPIClient(config), nil
}
