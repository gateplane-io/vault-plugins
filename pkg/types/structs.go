// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode"

	uuid "github.com/hashicorp/go-uuid"
)

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s // Return the empty string as is
	}
	runes := []rune(s)                   // Convert to runes to handle Unicode
	runes[0] = unicode.ToUpper(runes[0]) // Capitalize the first rune
	return string(runes)
}

/* ======================== Statuses */

type AccessRequestStatus int
type ApprovalStatus int

const (
	Pending AccessRequestStatus = iota
	Approved
	Active
	Expired
	Denied
)

var AccessRequestStatusStrings = []string{"Pending", "Approved", "Active", "Expired", "Denied"}
var ApprovalStatusStrings = []string{"Active", "Expired", "Retracted"}

const (
	ApprovalActive ApprovalStatus = iota
	ApprovalExpired
	ApprovalRetracted
)

func (s AccessRequestStatus) String() string {
	return strings.ToLower(AccessRequestStatusStrings[s])
}

func (s AccessRequestStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *AccessRequestStatus) UnmarshalJSON(data []byte) error {
	var statusString string
	if err := json.Unmarshal(data, &statusString); err != nil {
		return err
	}

	statusString = CapitalizeFirstLetter(statusString)
	for i, validStatus := range AccessRequestStatusStrings {
		if statusString == validStatus {
			*s = AccessRequestStatus(i)
			return nil
		}
	}
	return fmt.Errorf("invalid AccessRequestStatus: %s", statusString)
}

func (s ApprovalStatus) String() string {
	return strings.ToLower(ApprovalStatusStrings[s])
}

func (s ApprovalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *ApprovalStatus) UnmarshalJSON(data []byte) error {
	var statusString string
	if err := json.Unmarshal(data, &statusString); err != nil {
		return err
	}

	statusString = CapitalizeFirstLetter(statusString)
	for i, validStatus := range ApprovalStatusStrings {
		if statusString == validStatus {
			*s = ApprovalStatus(i)
			return nil
		}
	}
	return fmt.Errorf("invalid ApprovalStatus: %s", statusString)
}

/* ======================== AccessRequest */

type RequestID string
type ResponseID string

func newResponseID(requestID string, responderID string, denial bool) ResponseID {
	symbol := "+"
	if denial {
		symbol = "!"
	}
	return ResponseID(fmt.Sprintf("%s:%s:%s", requestID, responderID, symbol))
}

type AccessRequest struct {
	CreatedAt  time.Time           `json:"at"`        // Timestamp when the request was made
	Expiration time.Time           `json:"exp"`       // Expiration timestamp of the request
	Reason     string              `json:"reason"`    // Reason for requesting access
	Approvals  map[string]Approval `json:"approvals"` // Map of Gatekeeper ID -> Approval details
	Status     AccessRequestStatus `json:"status"`
	ID         string              `json:"requestor"`
	GrantCode  string              `json:"grant_code"`
}

func NewAccessRequest(config *Config, entityID string, reason string) AccessRequest {
	now := time.Now()
	return AccessRequest{
		CreatedAt:  now,
		Expiration: now.Add(config.RequestTTL),
		Reason:     reason,
		Approvals:  make(map[string]Approval),
		Status:     Pending,
		ID:         entityID,
	}
}

func (accessRequest *AccessRequest) CreateGrantCode() error {
	if accessRequest.GrantCode != "" {
		return fmt.Errorf("GrantCode has already been created for AccessRequest %s", accessRequest.ID)
	}
	code, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	accessRequest.GrantCode = code
	return nil
}

type Approval struct {
	CreatedAt  time.Time      `json:"at"`  // Timestamp when the approval was granted
	Expiration time.Time      `json:"exp"` // Expiration timestamp of the approval
	Approver   string         `json:"approver"`
	Status     ApprovalStatus `json:"status"`
	ID         ResponseID     `json:"response_id"`
}

func NewApproval(config *Config, entityID string, requestId string) Approval {
	now := time.Now()
	return Approval{
		CreatedAt:  now,
		Expiration: now.Add(config.ApprovalTTL),
		Approver:   entityID,
		Status:     ApprovalActive,
		ID:         newResponseID(requestId, entityID, false),
	}
}

/* ======================== Config */

type Config struct {
	ApprovalTTL       time.Duration `json:"approval_ttl"`       // Approval time to live
	RequestTTL        time.Duration `json:"request_ttl"`        // Request time to live
	RequiredApprovals int           `json:"required_approvals"` // Number of required approvals
	RequireReason     bool          `json:"require_reason"`     // Whether a reason is required
	DeleteAfter       time.Duration `json:"delete_after"`
}

func NewDefaultConfig() Config {
	return Config{
		ApprovalTTL:       1 * time.Hour,
		RequestTTL:        1 * time.Hour,
		RequiredApprovals: 1,
		RequireReason:     false,
		DeleteAfter:       8 * time.Hour,
	}
}

func (c *Config) SetConfigurationKey(key string, value interface{}) error {
	switch key {
	case "required_approvals":
		if v, ok := value.(int); ok {
			c.RequiredApprovals = v
		} else {
			return fmt.Errorf("invalid type for required_approvals, expected int")
		}
	case "require_reason":
		if v, ok := value.(bool); ok {
			c.RequireReason = v
		} else {
			return fmt.Errorf("invalid type for require_reason, expected bool")
		}
	case "approval_ttl":
		c.ApprovalTTL = time.Duration(value.(int)) * time.Second
	case "request_ttl":
		c.RequestTTL = time.Duration(value.(int)) * time.Second
	case "delete_after":
		c.DeleteAfter = time.Duration(value.(int)) * time.Second
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
	return nil
}