// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

type (
	// ErrNoMessages is returned when an operation returned no messages. It is not indicative that there will not be
	// more messages in the future.
	ErrNoMessages struct{}
)

func (e ErrNoMessages) Error() string {
	return "no messages available"
}
