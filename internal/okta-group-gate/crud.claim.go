// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package okta_group_gate

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

/* ======================== CRUD Config*/
const GroupMembershipKey = "claims"

func (b *Backend) GetMembership(ctx context.Context, req *logical.Request, userId string) (time.Time, error) {

	key := fmt.Sprintf("%s/%s", GroupMembershipKey, userId)
	timeError := time.Now()

	entry, err := req.Storage.Get(ctx, key)
	if err != nil {
		b.Logger().Error("[-] Could not retrieve Membership Entry",
			"error", err,
		)
		return timeError, err
	}

	expBytes := entry.Value
	secs, err := strconv.ParseInt(string(expBytes), 10, 64)
	if err != nil {
		b.Logger().Error("[-] Could not decode Membership Entry",
			"error", err,
			"expBytes", expBytes,
		)
		return timeError, err
	}
	exp := time.Unix(secs, 0)

	return exp, nil
}

func (b *Backend) StoreMembership(ctx context.Context, req *logical.Request, userId string, ttl time.Duration) (time.Time, error) {

	exp := time.Now().Add(ttl)
	s := strconv.FormatInt(exp.Unix(), 10) // e.g. "1617181723"
	expBytes := []byte(s)

	key := fmt.Sprintf("%s/%s", GroupMembershipKey, userId)
	err := req.Storage.Put(ctx, &logical.StorageEntry{
		Key:   key,
		Value: expBytes,
	})
	if err != nil {
		b.Logger().Error("[-] Could not store Membership Entry",
			"key", key,
			"exp", exp,
			"expBytes", expBytes,
			"error", err,
		)
		return exp, err
	}
	return exp, nil
}

func (b *Backend) CleanMembershipsPartial(ctx context.Context, req *logical.Request, all bool, continueOnError bool) error {

	entries, err := req.Storage.List(ctx, fmt.Sprintf("%s/", GroupMembershipKey))
	if err != nil {
		b.Logger().Error("[-] Could not List Memberships",
			"error", err,
		)
		return err
	}

	config, err := b.GetConfiguration(ctx, req)
	if err != nil {
		b.Logger().Error("[-] Could retrieve Configuration",
			"error", err,
		)
		return err
	}

	now := time.Now()

	for _, id := range entries {

		b.Logger().Trace("[*] About to check Expiration for key",
			"now", now,
			"UserID", id,
		)

		if strings.Contains(id, "/") {
			// This means it is a subkey, shouldn't happen
			continue
		}

		exp, err := b.GetMembership(ctx, req, id)

		if err != nil {
			if continueOnError {
				continue // retry in the next tick
			}
			return err
		}

		b.Logger().Debug("[*] Checking Expiration",
			"now", now,
			"exp", exp,
			"UserID", id,
		)
		if all || now.After(exp) {
			err = oktaRemoveFromGroup(ctx, b.oktaClient, config.OktaGroupID, id)
			if err != nil {
				b.Logger().Error("[-] Could remove from API Group",
					"GroupID", config.OktaGroupID,
					"UserID", id,
					"error", err,
				)
				if continueOnError {
					continue // retry in the next tick
				}
			}
			key := fmt.Sprintf("%s/%s", GroupMembershipKey, id)
			err := req.Storage.Delete(ctx, key)
			if err != nil {
				b.Logger().Error("[-] Could not delete Membership Entry from Storage",
					"error", err,
					"key", key,
				)
			}
		}
	}
	return nil
}

func (b *Backend) CleanMemberships(ctx context.Context, req *logical.Request) error {
	return b.CleanMembershipsPartial(ctx, req, true, true)
}

func (b *Backend) CleanExpiredMemberships(ctx context.Context, req *logical.Request) error {
	return b.CleanMembershipsPartial(ctx, req, false, true)
}
