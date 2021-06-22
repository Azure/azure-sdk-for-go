// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package utils

import "time"

func StringPtr(s string) *string {
	return &s
}

func TimePtr(t time.Time) *time.Time {
	return &t
}
