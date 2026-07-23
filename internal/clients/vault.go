// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package clients

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/api"
)

type PeriodicTokenInfo struct {
	Secret        *api.Secret
	PeriodSeconds int
}

func NewVaultClient(url string, httpClient *http.Client) (*api.Client, error) {
	if strings.TrimSpace(url) == "" {
		return nil, fmt.Errorf("vault URL is required")
	}

	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = url
	if httpClient != nil {
		vaultConfig.HttpClient = httpClient
	} else if vaultConfig.HttpClient == nil {
		vaultConfig.HttpClient = &http.Client{Timeout: 10 * time.Second}
	}

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}
	return client, nil
}

func NewVaultWrappedTokenClient(ctx context.Context, url, wrappedToken string, httpClient *http.Client) (*api.Client, *PeriodicTokenInfo, error) {
	if strings.TrimSpace(wrappedToken) == "" {
		return nil, nil, fmt.Errorf("wrapped token is required")
	}

	client, err := NewVaultClient(url, httpClient)
	if err != nil {
		return nil, nil, err
	}

	secret, err := client.Logical().UnwrapWithContext(ctx, wrappedToken)
	client.ClearToken()
	if err != nil {
		return nil, nil, tokenSafeError("unwrapping token", err, wrappedToken)
	}

	token, err := secret.TokenID()
	if err != nil {
		return nil, nil, tokenSafeError("reading unwrapped token", err, wrappedToken)
	}
	if strings.TrimSpace(token) == "" {
		return nil, nil, fmt.Errorf("wrapped response did not contain a token")
	}

	client.SetToken(token)
	info, err := ValidatePeriodicOrphanToken(ctx, client)
	if err != nil {
		client.ClearToken()
		return nil, nil, tokenSafeError("validating unwrapped token", err, wrappedToken, token)
	}
	return client, info, nil
}

func NewVaultPeriodicTokenClient(ctx context.Context, url, token string, httpClient *http.Client) (*api.Client, *PeriodicTokenInfo, error) {
	if strings.TrimSpace(token) == "" {
		return nil, nil, fmt.Errorf("vault token is required")
	}

	client, err := NewVaultClient(url, httpClient)
	if err != nil {
		return nil, nil, err
	}
	client.SetToken(token)

	info, err := ValidatePeriodicOrphanToken(ctx, client)
	if err != nil {
		client.ClearToken()
		return nil, nil, tokenSafeError("validating configured token", err, token)
	}
	return client, info, nil
}

func ValidatePeriodicOrphanToken(ctx context.Context, client *api.Client) (*PeriodicTokenInfo, error) {
	if client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}
	if client.Token() == "" {
		return nil, fmt.Errorf("vault token is required")
	}

	lookup, err := client.Auth().Token().LookupSelfWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("looking up token: %w", err)
	}
	if lookup == nil || lookup.Data == nil {
		return nil, fmt.Errorf("token lookup returned no data")
	}

	orphan, err := parseutil.ParseBool(lookup.Data["orphan"])
	if err != nil {
		return nil, fmt.Errorf("parsing token orphan status: %w", err)
	}
	if !orphan {
		return nil, fmt.Errorf("token must be orphan")
	}

	renewable, err := lookup.TokenIsRenewable()
	if err != nil {
		return nil, fmt.Errorf("parsing token renewable status: %w", err)
	}
	if !renewable {
		return nil, fmt.Errorf("token must be renewable")
	}

	ttl, err := lookup.TokenTTL()
	if err != nil {
		return nil, fmt.Errorf("parsing token TTL: %w", err)
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("token TTL must be greater than zero")
	}

	period, err := parseutil.ParseDurationSecond(lookup.Data["period"])
	if err != nil {
		return nil, fmt.Errorf("parsing token period: %w", err)
	}
	if period <= 0 {
		return nil, fmt.Errorf("token must be periodic")
	}

	if rawExplicitMaxTTL, ok := lookup.Data["explicit_max_ttl"]; ok && rawExplicitMaxTTL != nil {
		explicitMaxTTL, err := parseutil.ParseDurationSecond(rawExplicitMaxTTL)
		if err != nil {
			return nil, fmt.Errorf("parsing token explicit max TTL: %w", err)
		}
		if explicitMaxTTL > 0 {
			return nil, fmt.Errorf("token must not have an explicit max TTL")
		}
	}

	periodSeconds := int(period / time.Second)
	return &PeriodicTokenInfo{
		PeriodSeconds: periodSeconds,
		Secret: &api.Secret{Auth: &api.SecretAuth{
			ClientToken:   client.Token(),
			Orphan:        true,
			Renewable:     true,
			LeaseDuration: int(ttl / time.Second),
		}},
	}, nil
}

func tokenSafeError(operation string, err error, tokens ...string) error {
	message := err.Error()
	for _, token := range tokens {
		if token != "" {
			message = strings.ReplaceAll(message, token, "[redacted]")
		}
	}
	return fmt.Errorf("%s: %s", operation, message)
}
