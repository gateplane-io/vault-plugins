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
	"time"

	"github.com/hashicorp/vault/api"

	clientConfig "github.com/gateplane-io/vault-plugins/internal/clients/config"
)

func NewVaultAppRoleClient(ctx context.Context, cfg clientConfig.ConfigApiVaultAppRole, httpClient *http.Client) (*api.Client, error) {
	if cfg.Url == "" {
		return nil, fmt.Errorf("vault URL is required")
	}

	// create Vault config
	vaultCfg := api.DefaultConfig()
	vaultCfg.Address = cfg.Url
	if httpClient != nil {
		vaultCfg.HttpClient = httpClient
	}
	// reduce long blocking by setting a reasonable timeout if none set on httpClient
	if vaultCfg.HttpClient == nil {
		vaultCfg.HttpClient = &http.Client{Timeout: 10 * time.Second}
	}

	client, err := api.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	token, err := LoginWithAppRole(ctx, client, cfg.RoleID, cfg.RoleSecret, cfg.AppRoleMount)
	if err != nil {
		return nil, fmt.Errorf("approle login failed: %w", err)
	}

	// set token on client
	client.SetToken(token)
	return client, nil
}

func LoginWithAppRole(ctx context.Context, client *api.Client, roleID, secretID, mount string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("vault client is nil")
	}
	if roleID == "" {
		return "", fmt.Errorf("roleID is required")
	}
	if secretID == "" {
		return "", fmt.Errorf("secretID is required")
	}
	if mount == "" {
		mount = "approle"
	}

	// Normalize mount: remove leading "auth/" if provided by caller
	if len(mount) >= 5 && mount[:5] == "auth/" {
		mount = mount[5:]
	}

	payload := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := client.Logical().WriteWithContext(ctx, path, payload)
	if err != nil {
		return "", fmt.Errorf("approle login failed: %w", err)
	}
	if secret == nil || secret.Auth == nil || secret.Auth.ClientToken == "" {
		return "", fmt.Errorf("approle login returned no token")
	}

	client.SetToken(secret.Auth.ClientToken)
	return secret.Auth.ClientToken, nil
}

// IsAuthenticated checks whether the provided Vault client currently has a valid token.
// It returns true if the client has a non-empty token and the token lookup succeeds.
func IsAuthenticatedVault(ctx context.Context, client *api.Client) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("vault client is nil")
	}
	token := client.Token()
	if token == "" {
		return false, nil
	}

	// token/lookup-self works even for accessor-based checks; use ReadWithContext for context support
	secret, err := client.Logical().ReadWithContext(ctx, "auth/token/lookup-self")
	if err != nil {
		return false, fmt.Errorf("token lookup failed: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return false, nil
	}

	// Optionally ensure token is not expired/renewable info present; presence of data means valid
	return true, nil
}

// EnsureAuthentication checks if the client is authenticated and, if not, logs in using AppRole.
// - client must be preconfigured (address, http client, TLS).
// - roleID and secretID are used for AppRole login when needed.
// - mount is the approle mount path (e.g., "approle"); empty defaults to "approle".
// Returns an error on failure.
func EnsureAuthenticationVault(ctx context.Context, client *api.Client, roleID, secretID, mount string) error {
	if client == nil {
		return fmt.Errorf("vault client is nil")
	}

	ok, err := IsAuthenticatedVault(ctx, client)
	if err != nil {
		return fmt.Errorf("authentication check failed: %w", err)
	}
	if ok {
		return nil
	}

	_, err = LoginWithAppRole(ctx, client, roleID, secretID, mount)
	if err != nil {
		return fmt.Errorf("approle login failed: %w", err)
	}
	return nil
}
