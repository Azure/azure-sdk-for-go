// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"fmt"

	"github.com/google/uuid"
)

// -----------------------------------------------------------------------------
// HTTP Status Codes (matching Azure Cosmos DB)
// -----------------------------------------------------------------------------

// StatusCode represents HTTP-compatible status codes used by RNTBD.
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusNoContent           = 204
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusMethodNotAllowed    = 405
	StatusRequestTimeout      = 408
	StatusConflict            = 409
	StatusGone                = 410
	StatusPreconditionFailed  = 412
	StatusRequestEntityTooBig = 413
	StatusLocked              = 423
	StatusTooManyRequests     = 429
	StatusRetryWith           = 449
	StatusInternalServerError = 500
	StatusServiceUnavailable  = 503
)

// SubStatusCode represents sub-status codes for finer-grained error classification.
const (
	SubStatusUnknown                    = 0
	SubStatusNameCacheIsStale           = 1000 // InvalidPartitionException
	SubStatusPartitionKeyRangeGone      = 1002 // PartitionKeyRangeGoneException
	SubStatusCompletingSplit            = 1007 // PartitionKeyRangeIsSplittingException
	SubStatusCompletingPartitionMigrate = 1008 // PartitionIsMigratingException
	SubStatusLeaseNotFound              = 1022 // LeaseNotFoundException
	SubStatusTimeoutGenerated410        = 20002
)

// -----------------------------------------------------------------------------
// Error Types
// -----------------------------------------------------------------------------

// RntbdError represents an error from an RNTBD response.
// This is the base error type that all specific error types embed.
type RntbdError struct {
	StatusCode    int32
	SubStatusCode int64
	ActivityID    uuid.UUID
	Message       string
	LSN           int64
	PartitionID   string
}

func (e *RntbdError) Error() string {
	if e.SubStatusCode != 0 {
		return fmt.Sprintf("rntbd: status=%d substatus=%d activity=%s: %s",
			e.StatusCode, e.SubStatusCode, e.ActivityID, e.Message)
	}
	return fmt.Sprintf("rntbd: status=%d activity=%s: %s",
		e.StatusCode, e.ActivityID, e.Message)
}

// IsRetriable returns true if the error is potentially retriable.
func (e *RntbdError) IsRetriable() bool {
	switch e.StatusCode {
	case StatusRequestTimeout, StatusServiceUnavailable, StatusTooManyRequests, StatusGone:
		return true
	default:
		return false
	}
}

// -----------------------------------------------------------------------------
// Specific Error Types
// -----------------------------------------------------------------------------

// BadRequestError indicates a malformed request (400).
type BadRequestError struct {
	RntbdError
}

// UnauthorizedError indicates an authentication failure (401).
type UnauthorizedError struct {
	RntbdError
}

// ForbiddenError indicates access denied (403).
type ForbiddenError struct {
	RntbdError
}

// NotFoundError indicates the resource was not found (404).
type NotFoundError struct {
	RntbdError
}

// MethodNotAllowedError indicates an unsupported HTTP method (405).
type MethodNotAllowedError struct {
	RntbdError
}

// RequestTimeoutError indicates the request timed out (408).
type RequestTimeoutError struct {
	RntbdError
}

// ConflictError indicates a conflict with the current state (409).
type ConflictError struct {
	RntbdError
}

// GoneError indicates the resource is no longer available (410).
// This is the base type for partition-related gone errors.
type GoneError struct {
	RntbdError
}

// InvalidPartitionError indicates the partition is invalid due to stale cache (410, substatus 1000).
type InvalidPartitionError struct {
	RntbdError
}

// PartitionKeyRangeGoneError indicates the partition key range is gone (410, substatus 1002).
type PartitionKeyRangeGoneError struct {
	RntbdError
}

// PartitionKeyRangeIsSplittingError indicates the partition is being split (410, substatus 1007).
type PartitionKeyRangeIsSplittingError struct {
	RntbdError
}

// PartitionIsMigratingError indicates the partition is being migrated (410, substatus 1008).
type PartitionIsMigratingError struct {
	RntbdError
}

// LeaseNotFoundError indicates the lease was not found (410, substatus 1022).
type LeaseNotFoundError struct {
	RntbdError
}

// PreconditionFailedError indicates an ETag/If-Match precondition failed (412).
type PreconditionFailedError struct {
	RntbdError
}

// RequestEntityTooLargeError indicates the request payload is too large (413).
type RequestEntityTooLargeError struct {
	RntbdError
}

// LockedError indicates the resource is locked (423).
type LockedError struct {
	RntbdError
}

// RequestRateTooLargeError indicates rate limiting (429).
type RequestRateTooLargeError struct {
	RntbdError
	RetryAfterMs int64 // Retry-after duration in milliseconds
}

// RetryWithError indicates the client should retry with a different endpoint (449).
type RetryWithError struct {
	RntbdError
}

// InternalServerError indicates a server-side error (500).
type InternalServerError struct {
	RntbdError
}

// ServiceUnavailableError indicates the service is temporarily unavailable (503).
type ServiceUnavailableError struct {
	RntbdError
}

// -----------------------------------------------------------------------------
// Error Factory
// -----------------------------------------------------------------------------

// ErrorFromResponse creates an appropriate error from an RNTBD response.
// Returns nil if the response indicates success (2xx status).
func ErrorFromResponse(response *ResponseMessage) error {
	if response == nil {
		return fmt.Errorf("rntbd: nil response")
	}

	statusCode := response.Frame.Status
	activityID := response.Frame.ActivityID

	// Success - no error
	if statusCode >= StatusOK && statusCode < 300 {
		return nil
	}

	// Extract substatus if available
	var subStatusCode int64
	if token := response.Headers.Get(uint16(ResponseHeaderSubStatus)); token != nil && token.IsPresent() {
		if val, err := token.GetValue(); err == nil {
			switch v := val.(type) {
			case uint64:
				subStatusCode = int64(v)
			case int64:
				subStatusCode = v
			case uint32:
				subStatusCode = int64(v)
			case int32:
				subStatusCode = int64(v)
			}
		}
	}

	// Extract LSN if available
	var lsn int64
	if token := response.Headers.Get(uint16(ResponseHeaderLSN)); token != nil && token.IsPresent() {
		if val, err := token.GetValue(); err == nil {
			switch v := val.(type) {
			case int64:
				lsn = v
			case uint64:
				lsn = int64(v)
			}
		}
	}

	// Extract PartitionKeyRangeId if available
	var partitionID string
	if token := response.Headers.Get(uint16(ResponseHeaderPartitionKeyRangeId)); token != nil && token.IsPresent() {
		if val, err := token.GetValue(); err == nil {
			if s, ok := val.(string); ok {
				partitionID = s
			}
		}
	}

	// Build message from payload if present
	message := fmt.Sprintf("request failed with status %d", statusCode)
	if len(response.Payload) > 0 {
		message = string(response.Payload)
	}

	// Create base error
	baseError := RntbdError{
		StatusCode:    statusCode,
		SubStatusCode: subStatusCode,
		ActivityID:    activityID,
		Message:       message,
		LSN:           lsn,
		PartitionID:   partitionID,
	}

	// Map to specific error type
	switch statusCode {
	case StatusBadRequest:
		return &BadRequestError{baseError}

	case StatusUnauthorized:
		return &UnauthorizedError{baseError}

	case StatusForbidden:
		return &ForbiddenError{baseError}

	case StatusNotFound:
		return &NotFoundError{baseError}

	case StatusMethodNotAllowed:
		return &MethodNotAllowedError{baseError}

	case StatusRequestTimeout:
		return &RequestTimeoutError{baseError}

	case StatusConflict:
		return &ConflictError{baseError}

	case StatusGone:
		// Sub-status determines specific Gone error type
		switch subStatusCode {
		case SubStatusNameCacheIsStale:
			return &InvalidPartitionError{baseError}
		case SubStatusPartitionKeyRangeGone:
			return &PartitionKeyRangeGoneError{baseError}
		case SubStatusCompletingSplit:
			return &PartitionKeyRangeIsSplittingError{baseError}
		case SubStatusCompletingPartitionMigrate:
			return &PartitionIsMigratingError{baseError}
		case SubStatusLeaseNotFound:
			return &LeaseNotFoundError{baseError}
		default:
			return &GoneError{baseError}
		}

	case StatusPreconditionFailed:
		return &PreconditionFailedError{baseError}

	case StatusRequestEntityTooBig:
		return &RequestEntityTooLargeError{baseError}

	case StatusLocked:
		return &LockedError{baseError}

	case StatusTooManyRequests:
		retryErr := &RequestRateTooLargeError{RntbdError: baseError}
		// Extract retry-after if available
		if token := response.Headers.Get(uint16(ResponseHeaderRetryAfterMilliseconds)); token != nil && token.IsPresent() {
			if val, err := token.GetValue(); err == nil {
				switch v := val.(type) {
				case uint32:
					retryErr.RetryAfterMs = int64(v)
				case int32:
					retryErr.RetryAfterMs = int64(v)
				}
			}
		}
		return retryErr

	case StatusRetryWith:
		return &RetryWithError{baseError}

	case StatusInternalServerError:
		return &InternalServerError{baseError}

	case StatusServiceUnavailable:
		return &ServiceUnavailableError{baseError}

	default:
		// Unknown status code - return base error
		return &baseError
	}
}

// -----------------------------------------------------------------------------
// Error Type Assertions
// -----------------------------------------------------------------------------

// IsBadRequest returns true if the error is a BadRequestError.
func IsBadRequest(err error) bool {
	_, ok := err.(*BadRequestError)
	return ok
}

// IsUnauthorized returns true if the error is an UnauthorizedError.
func IsUnauthorized(err error) bool {
	_, ok := err.(*UnauthorizedError)
	return ok
}

// IsForbidden returns true if the error is a ForbiddenError.
func IsForbidden(err error) bool {
	_, ok := err.(*ForbiddenError)
	return ok
}

// IsNotFound returns true if the error is a NotFoundError.
func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

// IsConflict returns true if the error is a ConflictError.
func IsConflict(err error) bool {
	_, ok := err.(*ConflictError)
	return ok
}

// IsGone returns true if the error indicates the resource is gone (any Gone variant).
func IsGone(err error) bool {
	switch err.(type) {
	case *GoneError, *InvalidPartitionError, *PartitionKeyRangeGoneError,
		*PartitionKeyRangeIsSplittingError, *PartitionIsMigratingError:
		return true
	default:
		return false
	}
}

// IsPreconditionFailed returns true if the error is a PreconditionFailedError.
func IsPreconditionFailed(err error) bool {
	_, ok := err.(*PreconditionFailedError)
	return ok
}

// IsRequestRateTooLarge returns true if the error is a RequestRateTooLargeError.
func IsRequestRateTooLarge(err error) bool {
	_, ok := err.(*RequestRateTooLargeError)
	return ok
}

// IsServiceUnavailable returns true if the error is a ServiceUnavailableError.
func IsServiceUnavailable(err error) bool {
	_, ok := err.(*ServiceUnavailableError)
	return ok
}

// IsRetriable returns true if the error is potentially retriable.
func IsRetriable(err error) bool {
	if rntbdErr, ok := err.(*RntbdError); ok {
		return rntbdErr.IsRetriable()
	}
	// Check embedded RntbdError types
	switch e := err.(type) {
	case *BadRequestError:
		return e.IsRetriable()
	case *UnauthorizedError:
		return e.IsRetriable()
	case *ForbiddenError:
		return e.IsRetriable()
	case *NotFoundError:
		return e.IsRetriable()
	case *MethodNotAllowedError:
		return e.IsRetriable()
	case *RequestTimeoutError:
		return e.IsRetriable()
	case *ConflictError:
		return e.IsRetriable()
	case *GoneError:
		return e.IsRetriable()
	case *InvalidPartitionError:
		return e.IsRetriable()
	case *PartitionKeyRangeGoneError:
		return e.IsRetriable()
	case *PartitionKeyRangeIsSplittingError:
		return e.IsRetriable()
	case *PartitionIsMigratingError:
		return e.IsRetriable()
	case *PreconditionFailedError:
		return e.IsRetriable()
	case *RequestEntityTooLargeError:
		return e.IsRetriable()
	case *LockedError:
		return e.IsRetriable()
	case *RequestRateTooLargeError:
		return e.IsRetriable()
	case *RetryWithError:
		return e.IsRetriable()
	case *InternalServerError:
		return e.IsRetriable()
	case *ServiceUnavailableError:
		return e.IsRetriable()
	default:
		return false
	}
}
