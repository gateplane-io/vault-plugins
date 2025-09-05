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
	"github.com/gateplane-io/vault-plugins/pkg/types"
)

// Re-export types from pkg/types for backwards compatibility
type AccessRequestStatus = types.AccessRequestStatus
type ApprovalStatus = types.ApprovalStatus
type RequestID = types.RequestID
type ResponseID = types.ResponseID
type AccessRequest = types.AccessRequest
type Approval = types.Approval
type Config = types.Config

// Re-export constants
const (
	Pending  = types.Pending
	Approved = types.Approved
	Active   = types.Active
	Expired  = types.Expired
	Denied   = types.Denied
)

const (
	ApprovalActive    = types.ApprovalActive
	ApprovalExpired   = types.ApprovalExpired
	ApprovalRetracted = types.ApprovalRetracted
)

// Re-export functions
var NewAccessRequest = types.NewAccessRequest
var NewApproval = types.NewApproval
var NewDefaultConfig = types.NewDefaultConfig
