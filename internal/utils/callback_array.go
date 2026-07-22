// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"
)

type CallbackArray struct {
	onAppend func(context.Context, *logical.Request, string) (map[string]interface{}, error)
	onRemove func(context.Context, *logical.Request, string, map[string]interface{}) error
}

func NewCallbackArray(
	onAppend func(context.Context, *logical.Request, string) (map[string]interface{}, error),
	onRemove func(context.Context, *logical.Request, string, map[string]interface{}) error,
) *CallbackArray {
	return &CallbackArray{
		onAppend: onAppend,
		onRemove: onRemove,
	}
}

// Append runs the callback that grants access for an element.
func (m *CallbackArray) Append(ctx context.Context, req *logical.Request, element string) (map[string]interface{}, error) {
	if m.onAppend == nil {
		return nil, fmt.Errorf("No Append Callback defined")
	}
	return m.onAppend(ctx, req, element)
}

// Remove runs the idempotent callback that revokes access for an element.
// Revocation must not depend on process-local state because Vault may invoke it
// after a plugin restart or route it to another plugin process.
func (m *CallbackArray) Remove(ctx context.Context, req *logical.Request, element string, internalData map[string]interface{}) (bool, error) {
	if m.onRemove == nil {
		return false, fmt.Errorf("No Remove Callback defined")
	}
	if err := m.onRemove(ctx, req, element, internalData); err != nil {
		return false, err
	}
	return true, nil
}
