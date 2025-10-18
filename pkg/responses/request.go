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
	"time"
)

type AccessRequestCreationResponse struct {
	OwnerID    string `json:"requestor_id"`
	CreatedAt  int64  `json:"iat"`
	Expiration int64  `json:"exp"`
	Deletion   int64  `json:"deleted_after"`

	Justification     string `json:"justification"`
	RequiredApprovals int    `json:"required_approvals"`

	Status string `json:"status"`

	NumOfApprovals int           `json:"num_of_approvals"`
	Overwrite      bool          `json:"overwrite"`
	ClaimTTL       time.Duration `json:"claim_ttl"`
}

type AccessRequestResponse struct {
	OwnerID    string `json:"requestor_id"`
	CreatedAt  int64  `json:"iat"`
	Expiration int64  `json:"exp"`
	Deletion   int64  `json:"deleted_after"`

	Justification     string `json:"justification"`
	RequiredApprovals int    `json:"required_approvals"`

	Status string `json:"status"`

	NumOfApprovals int  `json:"num_of_approvals"`
	HaveApproved   bool `json:"have_approved"`

	ClaimCreatedAt int64         `json:"claim_iat"`
	ClaimTTL       time.Duration `json:"claim_ttl"`
}
