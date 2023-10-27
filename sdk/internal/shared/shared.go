//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
)

// NonRetriableError marks the specified error as non-retriable.
func NonRetriableError(err error) error {
	return &nonRetriableError{err}
}

type nonRetriableError struct {
	error
}

func (p *nonRetriableError) Error() string {
	return p.error.Error()
}

func (*nonRetriableError) NonRetriable() {
	// marker method
}

func (p *nonRetriableError) Unwrap() error {
	return p.error
}

var _ errorinfo.NonRetriable = (*nonRetriableError)(nil)
