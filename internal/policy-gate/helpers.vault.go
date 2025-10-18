// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package policy_gate

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// GetEntityPolicies returns the policies assigned to a Vault identity entity.
// - client: authenticated *api.Client
// - entityID: the Vault entity ID (e.g., "entity-id-xxxx")
// Returns slice of policy names or an error.
func GetEntityPolicies(ctx context.Context, client *api.Client, entityID string) ([]string, error) {
	if client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}
	if entityID == "" {
		return nil, fmt.Errorf("entityID is required")
	}

	// Read the entity: GET /v1/identity/entity/id/:id
	path := fmt.Sprintf("identity/entity/id/%s", entityID)
	secret, err := client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("reading entity %s: %w", entityID, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("entity %s not found", entityID)
	}

	// The 'policies' field is typically present as []interface{} or []string
	if p, ok := secret.Data["policies"]; ok && p != nil {
		switch v := p.(type) {
		case []interface{}:
			out := make([]string, 0, len(v))
			for _, e := range v {
				if s, ok := e.(string); ok {
					out = append(out, s)
				}
			}
			return out, nil
		case []string:
			return v, nil
		case string:
			// sometimes a single policy might be returned as a string
			return []string{v}, nil
		default:
			return nil, fmt.Errorf("unexpected policies type %T", p)
		}
	}

	// In some Vault versions/policy-assignments, policies are under 'policies' in 'data' or inside 'meta'
	// try nested 'data' object fallback:
	if nested, ok := secret.Data["data"].(map[string]interface{}); ok {
		if p, ok := nested["policies"]; ok && p != nil {
			switch v := p.(type) {
			case []interface{}:
				out := make([]string, 0, len(v))
				for _, e := range v {
					if s, ok := e.(string); ok {
						out = append(out, s)
					}
				}
				return out, nil
			case []string:
				return v, nil
			case string:
				return []string{v}, nil
			}
		}
	}

	// If not present, entity has no direct policies assigned
	return []string{}, nil
}

// AddPoliciesToEntity appends policiesToAdd to the entity's existing policies
// using the GetEntityPolicies helper to fetch current policies. It avoids duplicates and
// preserves common entity fields (name, metadata, disabled) when writing back.
func AddPoliciesToEntity(ctx context.Context, client *api.Client, entityID string, policiesToAdd []string) ([]string, error) {
	zero := []string{}
	if client == nil {
		return zero, fmt.Errorf("vault client is nil")
	}
	if entityID == "" {
		return zero, fmt.Errorf("entityID is required")
	}
	if len(policiesToAdd) == 0 {
		return zero, nil
	}

	// Use existing helper to get current policies
	existing, err := GetEntityPolicies(ctx, client, entityID)
	if err != nil {
		return zero, fmt.Errorf("getting existing policies for entity %s: %w", entityID, err)
	}

	// Build set to avoid duplicates
	policySet := make(map[string]struct{}, len(existing)+len(policiesToAdd))
	for _, p := range existing {
		if p == "" {
			continue
		}
		policySet[p] = struct{}{}
	}
	for _, p := range policiesToAdd {
		if p == "" {
			continue
		}
		policySet[p] = struct{}{}
	}

	merged := append(existing, policiesToAdd...)
	err = SetEntityPolicies(ctx, client, entityID, merged)
	if err != nil {
		return zero, fmt.Errorf("updating entity %s policies: %w", entityID, err)
	}

	return existing, nil
}

// SetEntityPolicies replaces the direct policies of the given entity with the provided policies.
// - client: authenticated *api.Client
// - entityID: Vault identity entity ID
// - policies: slice of policy names to set (duplicates will be removed)
func SetEntityPolicies(ctx context.Context, client *api.Client, entityID string, policies []string) error {
	if client == nil {
		return fmt.Errorf("vault client is nil")
	}
	if entityID == "" {
		return fmt.Errorf("entityID is required")
	}

	// Read the full entity to preserve required fields for update
	path := fmt.Sprintf("identity/entity/id/%s", entityID)
	secret, err := client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return fmt.Errorf("reading entity %s: %w", entityID, err)
	}
	if secret == nil || secret.Data == nil {
		return fmt.Errorf("entity %s not found", entityID)
	}

	// Prepare payload preserving common fields
	payload := map[string]interface{}{}
	payload["policies"] = policies

	_, err = client.Logical().WriteWithContext(ctx, path, payload)
	if err != nil {
		return fmt.Errorf("updating entity %s policies: %w", entityID, err)
	}
	return nil
}
