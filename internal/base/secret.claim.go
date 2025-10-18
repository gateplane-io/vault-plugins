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
	"github.com/hashicorp/vault/sdk/framework"
)

const SecretType = "claim"

func ClaimSecret(b *BaseBackend) *framework.Secret {
	return &framework.Secret{
		Type: SecretType,
		Fields: map[string]*framework.FieldSchema{
			"requestor_id": {
				Type:        framework.TypeString,
				Description: "The Entity ID of the AccessRequest owner",
			},
			// More fields are set adhoc in each plugin
		},

		Renew:  nil,
		Revoke: b.handleClaimRevocation,
	}
}
