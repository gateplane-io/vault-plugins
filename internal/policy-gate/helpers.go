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

func Subtract[T comparable](a, b []T) []T {
	if len(b) == 0 {
		return append([]T(nil), a...)
	}
	m := make(map[T]struct{}, len(b))
	for _, v := range b {
		m[v] = struct{}{}
	}
	out := make([]T, 0, len(a))
	for _, v := range a {
		if _, found := m[v]; !found {
			out = append(out, v)
		}
	}
	return out
}

func InterfaceSliceToStringsStrict(in []interface{}) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
