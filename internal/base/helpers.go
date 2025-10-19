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
	"encoding/json"
	"time"
)

// Function to convert a struct to a map[string]interface{}
func StructToMap(obj interface{}) (map[string]interface{}, error) {
	// Marshal the struct to JSON ([]byte)
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a map[string]interface{}
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Function to check if two dates are close within a given duration
func areDatesClose(date1, date2 time.Time, threshold time.Duration) bool {
	// Calculate the difference between the two dates
	diff := date1.Sub(date2)
	if diff < 0 {
		diff = -diff // Ensure positive difference
	}

	// Check if the difference is within the threshold
	return diff <= threshold
}
