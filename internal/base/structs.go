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
	"strings"
	"time"

	"github.com/gateplane-io/vault-plugins/pkg/models"
)

// AccessRequest
type AccessRequest struct {
	OwnerID    string    `json:"owner_id"`
	CreatedAt  time.Time `json:"iat"`
	Expiration time.Time `json:"exp"`
	Deletion   time.Time `json:"deleted_after"`

	Justification     string `json:"justification"` // provided by the requestor
	RequiredApprovals int    `json:"required_approvals"`

	ClaimCreatedAt time.Time `json:"claim_iat"`
	// ClaimTTL is stored as a number of seconds for API and storage compatibility.
	ClaimTTL time.Duration `json:"claim_ttl"`

	Status    models.AccessRequestStatus `json:"status"`
	Approvals map[string]*Approval       `json:"approvals"`
}

func NewAccessRequest(config *Config, configLease *ConfigLease, ownerID string, ttlSeconds time.Duration, justification string) (*AccessRequest, error) {
	leaseSeconds := configLease.Lease / time.Second
	leaseMaxSeconds := configLease.LeaseMax / time.Second

	if ttlSeconds == 0 {
		ttlSeconds = leaseSeconds
	}

	if ttlSeconds < leaseSeconds {
		return nil, fmt.Errorf("the requested TTL (%s) is lower than the minimum lease of the backend (%s)", ttlSeconds*time.Second, configLease.Lease)
	}
	if ttlSeconds > leaseMaxSeconds {
		return nil, fmt.Errorf("the requested TTL (%s) is higher than the maximum lease of the backend (%s)", ttlSeconds*time.Second, configLease.LeaseMax)
	}

	if config.RequireJustification && strings.TrimSpace(justification) == "" {
		return nil, fmt.Errorf("The backend requires a justification, but it is not provided")
	}

	now := time.Now()

	return &AccessRequest{
		Status: models.Pending,

		OwnerID:    ownerID,
		CreatedAt:  now,
		Expiration: now.Add(config.RequestTTL),
		Deletion:   now.Add(config.DeleteAfter),

		Justification:     justification,
		RequiredApprovals: config.RequiredApprovals,

		ClaimTTL:       ttlSeconds,
		ClaimCreatedAt: time.Unix(0, 0),

		Approvals: map[string]*Approval{},
	}, nil
}

func (a AccessRequest) Equals(b AccessRequest) bool {
	return a.OwnerID == b.OwnerID
}

// Approval
type Approval struct {
	OwnerID   string    `json:"requestor_id"`
	CreatedAt time.Time `json:"iat"`
}

func (req *AccessRequest) Approve(approverID string) (*Approval, bool, error) {
	if req.Status != models.Pending {
		return nil, false, fmt.Errorf(
			"The AccessRequest cannot be approved, as it is in '%s' state",
			req.Status,
		)
	}

	lastApproval := len(req.Approvals) == req.RequiredApprovals
	if lastApproval {
		req.Status = models.Approved
	}

	now := time.Now()
	approval := &Approval{
		OwnerID:   approverID,
		CreatedAt: now,
	}

	req.Approvals[approverID] = approval
	return approval, lastApproval, nil
}

func (req *AccessRequest) isApprovedBy(approverID string) bool {
	_, ok := req.Approvals[approverID]
	return ok
}

func (req *AccessRequest) Claim() error {
	if req.Status != models.Approved {
		return fmt.Errorf(
			"The AccessRequest cannot be claimed, as it is in '%s' state",
			req.Status,
		)
	}

	now := time.Now()
	req.Status = models.Active
	req.ClaimCreatedAt = now

	return nil
}
