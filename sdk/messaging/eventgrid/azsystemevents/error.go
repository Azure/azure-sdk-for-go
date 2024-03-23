//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents

// Error is an error that is included as part of a system event.
type Error struct {
	// Code is an error code from the service producing the system event.
	Code string

	// InnerError can provide more information about the origin of this error.
	InnerError *Error

	// Details contains related errors for this error.
	Details []*Error

	message string
}

// Error returns human-readable text describing the error.
func (e *Error) Error() string {
	return e.message
}
