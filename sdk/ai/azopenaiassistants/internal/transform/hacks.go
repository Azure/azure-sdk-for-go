//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

// The changes in here are all workarounds for issues with generation that should
// go away if we fix the generator.

import (
	"fmt"
	"strings"
)

// hackFixTimestamps is a workaround for a bug where the typespec compiler
// doesn't appear to be propagating the date/time format attribute for all
// attributes, resulting in Unix timestamps failing to deserialized as RFC1139.
func (t *transformer) hackFixTimestamps() error {
	return transformFiles(t.fileCache, "fix timestamps", []string{"models_serde.go"}, func(text string) (string, error) {
		fixes := []struct {
			JSONFieldName string
			FieldName     string
			ObjectName    string
		}{
			//
			{"cancelled_at", "CancelledAt", "r"},
			{"completed_at", "CompletedAt", "r"},
			{"expired_at", "ExpiredAt", "r"},
			{"failed_at", "FailedAt", "r"},

			// ThreadRun
			{"cancelled_at", "CancelledAt", "t"},
			{"expires_at", "ExpiresAt", "t"},
			{"failed_at", "FailedAt", "t"},
			{"completed_at", "CompletedAt", "t"},
			{"started_at", "StartedAt", "t"},
		}

		for _, fix := range fixes {
			searchStr := fmt.Sprintf(`populateDateTimeRFC3339(objectMap, "%s", %s.%s)`, fix.JSONFieldName, fix.ObjectName, fix.FieldName)
			newData := strings.Replace(text,
				searchStr,
				fmt.Sprintf(`populateTimeUnix(objectMap, "%s", %s.%s)`, fix.JSONFieldName, fix.ObjectName, fix.FieldName), -1)

			if newData == text {
				return "", fmt.Errorf("No replacement matched: '%s'", searchStr)
			}

			text = newData

			searchStr = fmt.Sprintf(`err = unpopulateDateTimeRFC3339(val, "%s", &%s.%s)`, fix.FieldName, fix.ObjectName, fix.FieldName)

			newData = strings.Replace(text,
				searchStr,
				fmt.Sprintf(`err = unpopulateTimeUnix(val, "%s", &%s.%s)`, fix.FieldName, fix.ObjectName, fix.FieldName), -1)

			if newData == text {
				return "", fmt.Errorf("No replacement matched: %q", searchStr)
			}

			text = newData
		}

		return text, nil
	}, nil)
}
