// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package internal

import (
	"time"
)

// Retryable represents an error which should be able to be retried
type Retryable string

// Error implementation for Retryable
func (r Retryable) Error() string {
	return string(r)
}

// Retry will attempt to retry an action a number of times if the action returns a retryable error
func Retry(times int, delay time.Duration, action func() (interface{}, error)) (interface{}, error) {
	var lastErr error
	for i := 0; i < times; i++ {
		item, err := action()
		if err != nil {
			if retryable, ok := err.(Retryable); ok {
				lastErr = retryable
				time.Sleep(delay)
				continue
			} else {
				return nil, err
			}
		}
		return item, nil
	}
	return nil, lastErr
}
