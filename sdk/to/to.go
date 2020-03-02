// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package to

// BoolPtr returns a pointer to the provided bool.
func BoolPtr(b bool) *bool {
	return &b
}

// Float32Ptr returns a pointer to the provided float32.
func Float32Ptr(i float32) *float32 {
	return &i
}

// Float64Ptr returns a pointer to the provided float64.
func Float64Ptr(i float64) *float64 {
	return &i
}

// Int32Ptr returns a pointer to the provided int32.
func Int32Ptr(i int32) *int32 {
	return &i
}

// Int64Ptr returns a pointer to the provided int64.
func Int64Ptr(i int64) *int64 {
	return &i
}

// StringPtr returns a pointer to the provided string.
func StringPtr(s string) *string {
	return &s
}
