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

func storageKeyForRequest(requestID string) string {
	return fmt.Sprintf("%s/%s", RequestKey, requestID)
}

func requestHasExpired(accessRequest AccessRequest) bool {
	return accessRequest.Expiration.Before(time.Now())
}

func requestIsDeletable(accessRequest AccessRequest) bool {
	return accessRequest.Deletion.Before(time.Now())
}

func validApprovalsNum(accessRequest AccessRequest) int {
	return len(accessRequest.Approvals)
}

func requestIsApproved(accessRequest AccessRequest) bool {
	numOfvalidApprovals := validApprovalsNum(accessRequest)
	return numOfvalidApprovals >= accessRequest.RequiredApprovals
}
