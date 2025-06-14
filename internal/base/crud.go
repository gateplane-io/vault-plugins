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
	"fmt"
	"time"
)

/* Storage Keys */
const RequestKey = "request"
const ConfigKey = "config"
const GrantCodeKey = "grant_code"

/* ======================== Helpers */

func storageKeyForRequest(requestID string) string {
	return fmt.Sprintf("%s/%s", RequestKey, requestID)
}

func requestHasExpired(accessRequest AccessRequest) bool {
	return accessRequest.Expiration.Before(time.Now())
}

func requestIsDeletable(accessRequest AccessRequest, deleteAfter time.Duration) bool {
	if accessRequest.Status != Expired {
		return false
	}
	return accessRequest.Expiration.Before(time.Now().Add(-deleteAfter))
}

func approvalHasExpired(approval Approval) bool {
	return approval.Expiration.Before(time.Now())
}

func validApprovalsNum(accessRequest AccessRequest) int {
	numOfvalidApprovals := 0
	for _, approval := range accessRequest.Approvals {
		// If someone changed their mind, the request can't be granted
		// if approval.Status == ApprovalRetracted{
		// 	return false
		// }
		if approval.Status == ApprovalActive {
			numOfvalidApprovals++
		}
	}
	return numOfvalidApprovals
}

func requestIsApproved(accessRequest AccessRequest, requiredApprovals int) bool {
	numOfvalidApprovals := validApprovalsNum(accessRequest)
	return numOfvalidApprovals >= requiredApprovals
}
