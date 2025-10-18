// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package clients

import (
	"github.com/okta/okta-sdk-golang/v5/okta"
)

func NewOktaClient(orgUrl, apiToken string) (*okta.APIClient, error) {

	config, err := okta.NewConfiguration(
		okta.WithOrgUrl(orgUrl),
		okta.WithToken(apiToken),
	)
	if err != nil {
		return nil, err
	}

	return okta.NewAPIClient(config), nil
}
