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

import (
// "time"
)

type ConfigResponse struct {
	RequireJustification bool `json:"require_justification"`
	RequiredApprovals    int  `json:"required_approvals"`
	// AllowRejection       bool  `json:"allow_rejection"`

	// Unix Time
	RequestTTL  float64 `json:"request_ttl"`
	DeleteAfter float64 `json:"delete_after"`
}
