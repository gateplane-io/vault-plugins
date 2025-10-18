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

type ConfigAccessPolicyGate struct {
	Policies []string `json:"policies"`
}

type ConfigAccessOktaGroupGate struct {
	GroupID   string `json:"okta_group_id"`
	GroupName string `json:"okta_group_name"`
}
