// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import "time"

// NOTE: this is exposed via type-aliasing in azeventhubs/client.go

// RetryOptions represent the options for retries.
type RetryOptions struct {
	// MaxRetries specifies the maximum number of attempts a failed operation will be retried
	// before producing an error.
	// The default value is three.  A value less than zero means one try and no retries.
	MaxRetries int32

	// RetryDelay specifies the initial amount of delay to use before retrying an operation.
	// The delay increases exponentially with each retry up to the maximum specified by MaxRetryDelay.
	// The default value is four seconds.  A value less than zero means no delay between retries.
	RetryDelay time.Duration

	// MaxRetryDelay specifies the maximum delay allowed before retrying an operation.
	// Typically the value is greater than or equal to the value specified in RetryDelay.
	// The default Value is 120 seconds.  A value less than zero means there is no cap.
	MaxRetryDelay time.Duration

	// LinkRecoveryDelay specifies a fixed delay to use after a link recovery failure, instead of
	// the normal exponential backoff specified by RetryDelay. This only applies when an AMQP link
	// needs to be recovered (e.g., link detached errors). For other types of failures, the normal
	// RetryDelay with exponential backoff is used.
	// The default value is zero, which means link recovery retries will use the normal RetryDelay
	// exponential backoff behavior. A value less than zero means no delay after link recovery.
	// Positive values let you slow down repeated link recovery attempts if, for example, recreating
	// links is putting pressure on your namespace, while negative values let you immediately try
	// again when link recovery failures happen.
	LinkRecoveryDelay time.Duration
}
