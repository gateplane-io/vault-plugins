// Copyright (C) 2025 Ioannis Torakis <john.torakis@gmail.com>
// SPDX-License-Identifier: Elastic-2.0
//
// Licensed under the Elastic License 2.0.
// You may obtain a copy of the license at:
// https://www.elastic.co/licensing/elastic-license
//
// Use, modification, and redistribution permitted under the terms of the license,
// except for providing this software as a commercial service or product.

package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s // Return the empty string as is
	}
	runes := []rune(s)                   // Convert to runes to handle Unicode
	runes[0] = unicode.ToUpper(runes[0]) // Capitalize the first rune
	return string(runes)
}

type AccessRequestStatus int

const (
	Pending AccessRequestStatus = iota
	Approved
	Active
	Expired
	Abandoned
	Rejected
	Revoked
)

var AccessRequestStatusStrings = []string{"Pending", "Approved", "Active", "Expired", "Abandoned", "Rejected", "Revoked"}

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

	statusString = capitalizeFirstLetter(statusString)
	for i, validStatus := range AccessRequestStatusStrings {
		if statusString == validStatus {
			*s = AccessRequestStatus(i)
			return nil
		}
	}
	return fmt.Errorf("invalid AccessRequestStatus: %s", statusString)
}
