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
	elements []string
	onAppend func(context.Context, *logical.Request, string) (map[string]interface{}, error)
	onRemove func(context.Context, *logical.Request, string, map[string]interface{}) error
}

func NewCallbackArray(
	onAppend func(context.Context, *logical.Request, string) (map[string]interface{}, error),
	onRemove func(context.Context, *logical.Request, string, map[string]interface{}) error,
) *CallbackArray {
	return &CallbackArray{
		elements: []string{},
		onAppend: onAppend,
		onRemove: onRemove,
	}
}

// Append adds an element to the array and runs the callback
func (m *CallbackArray) Append(ctx context.Context, req *logical.Request, element string) (map[string]interface{}, error) {
	if m.onAppend == nil {
		return nil, fmt.Errorf("No Append Callback defined")
	}
	m.elements = append(m.elements, element)
	return m.onAppend(ctx, req, element) // Run the function when an element is added
}

// Remove removes an element from the array and runs the onRemove callback
func (m *CallbackArray) Remove(ctx context.Context, req *logical.Request, element string, internalData map[string]interface{}) (bool, error) {
	if m.onRemove == nil {
		return false, fmt.Errorf("No Remove Callback defined")
	}
	// Find the index of the element to remove
	for i, e := range m.elements {
		if e == element {
			// Remove the element by slicing around it
			m.elements = append(m.elements[:i], m.elements[i+1:]...)
			return true, m.onRemove(ctx, req, element, internalData) // Run the function when an element is removed
		}
	}
	return false, nil // Element not found
}
