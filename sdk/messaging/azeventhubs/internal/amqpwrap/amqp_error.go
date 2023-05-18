// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package amqpwrap

import (
	"errors"
)

// Error is a wrapper that has the context of which connection and
// link the error happened with.
type Error struct {
	ConnID   uint64
	LinkName string
	Err      error
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e Error) As(target any) bool {
	return errors.As(e.Err, target)
}

func (e Error) Is(target error) bool {
	return errors.Is(e.Err, target)
}

func NewError(err error, connID uint64, linkName string) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(Error); ok {
		return err
	}

	return Error{
		ConnID:   connID,
		LinkName: linkName,
		Err:      err,
	}
}
