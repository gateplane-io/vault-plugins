// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package responses

type ConfigApiOktaResponse struct {
	OrgUrl                string `json:"org_url"`
	OktaOIDCMountAccessor string `json:"auth_mount_accessor"`
	OktaEntityKey         string `json:"okta_entity_key"`
	ApiTokenSet           bool   `json:"ApiTokenSet"`
}
