// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// errNotImplemented is returned by operations whose driver binding has not landed yet. It is
// temporary scaffolding for the incremental v2 rollout and will be removed once every operation
// is wired up. It is an [Error] so that the documented errors.As idiom works during the preview.
//
//nolint:unused // returned by the operations that land in later changes.
var errNotImplemented = &Error{
	Code:    CodeClientError,
	Message: "not implemented yet in this preview of v2",
}

// Code classifies a Cosmos DB failure so that callers can branch on it programmatically instead
// of interpreting status and sub-status codes themselves.
//
// New codes may be added as the service and driver grow, so treat an unrecognized code the same
// as [CodeUnknown] rather than assuming the set is closed.
type Code string

const (
	// CodeUnknown means the failure could not be classified. [Error.StatusCode] and
	// [Error.SubStatus] carry whatever the service reported.
	CodeUnknown Code = ""

	// CodeBadRequest means the request was malformed or rejected as invalid.
	CodeBadRequest Code = "BadRequest"

	// CodeUnauthorized means the credential was missing, malformed or rejected.
	CodeUnauthorized Code = "Unauthorized"

	// CodeForbidden means the credential is valid but not permitted to perform the operation.
	CodeForbidden Code = "Forbidden"

	// CodeNotFound means the requested resource does not exist.
	CodeNotFound Code = "NotFound"

	// CodeRequestTimeout means the service did not complete the request in time.
	CodeRequestTimeout Code = "RequestTimeout"

	// CodeConflict means a resource with the same ID already exists.
	CodeConflict Code = "Conflict"

	// CodeGone means the addressed resource moved. The driver retries these internally, so this
	// surfaces to callers only once those retries are exhausted.
	CodeGone Code = "Gone"

	// CodePartitionKeyRangeGone means the addressed physical partition no longer exists, usually
	// because the container was split or recreated, so cached routing information is stale.
	CodePartitionKeyRangeGone Code = "PartitionKeyRangeGone"

	// CodePreconditionFailed means an ETag precondition did not hold.
	CodePreconditionFailed Code = "PreconditionFailed"

	// CodeThrottled means the request exceeded the provisioned throughput. Callers that disable
	// the built-in throttle retries are expected to honor [Error.RetryAfter].
	CodeThrottled Code = "Throttled"

	// CodeServiceUnavailable means the service was temporarily unable to serve the request.
	CodeServiceUnavailable Code = "ServiceUnavailable"

	// CodeSessionUnavailable means the replica that served the request had not caught up to the
	// supplied session token. The resource itself may well exist, so this is deliberately not
	// [CodeNotFound] even though the service reports it with the same status code.
	CodeSessionUnavailable Code = "SessionUnavailable"

	// CodeClientError means the client produced the failure without a response from the service,
	// for example a transport, serialization or configuration failure.
	CodeClientError Code = "ClientError"
)

// Error describes a failed Cosmos DB operation.
//
// It reports failures produced by the service and failures produced locally by the client, which
// [Error.FromWire] distinguishes. Retrieve it with errors.As and branch on [Error.Code]:
//
//	var cosmosErr *azcosmos.Error
//	if errors.As(err, &cosmosErr) && cosmosErr.Code == azcosmos.CodeNotFound {
//		// the item, container or database does not exist
//	}
//
// Unlike v1, this is not an [azcore.ResponseError] and carries no *http.Response: v2 operations
// are executed by the Cosmos driver rather than by an azcore HTTP pipeline, so there is no HTTP
// response for a caller to inspect. The fields below carry the diagnostic values that were
// previously read off response headers.
//
// [Error.Code] is part of the published API. The string returned by [Error.Error] is not, and is
// subject to change.
type Error struct {
	// Code classifies the failure. Prefer it over StatusCode and SubStatus, which are reported
	// verbatim and are harder to interpret correctly.
	Code Code

	// StatusCode is the HTTP status code reported for the operation. It is zero when the failure
	// occurred before the service responded.
	StatusCode int

	// SubStatus is the Cosmos DB sub-status code, which qualifies StatusCode. It is zero when the
	// service did not report one.
	SubStatus int

	// Message describes the failure.
	Message string

	// ActivityID correlates the operation with server-side telemetry. It is empty when the
	// failure occurred before the service responded.
	ActivityID string

	// SessionToken is the session token reported alongside the failure, if any.
	SessionToken string

	// ETag is the entity tag reported alongside the failure, if any.
	ETag azcore.ETag

	// RetryAfter is how long the service asked the caller to wait before retrying. It is zero
	// when the service did not ask for a delay.
	RetryAfter time.Duration

	// FromWire reports whether the service produced this failure. It is false for failures the
	// client produced without a service response, such as transport and serialization failures.
	FromWire bool
}

// Error implements the error interface. The format is not part of the published API; branch on
// [Error.Code] rather than parsing this string.
func (e *Error) Error() string {
	var b strings.Builder
	b.WriteString("azcosmos: operation failed")
	if e.Code != CodeUnknown {
		fmt.Fprintf(&b, ": %s", e.Code)
	}
	switch {
	case e.StatusCode != 0 && e.SubStatus != 0:
		fmt.Fprintf(&b, ": status %d/%d", e.StatusCode, e.SubStatus)
	case e.StatusCode != 0:
		fmt.Fprintf(&b, ": status %d", e.StatusCode)
	case e.SubStatus != 0:
		fmt.Fprintf(&b, ": sub-status %d", e.SubStatus)
	}
	if e.Message != "" {
		fmt.Fprintf(&b, ": %s", e.Message)
	}
	if e.RetryAfter > 0 {
		fmt.Fprintf(&b, " (retry after %s)", e.RetryAfter)
	}
	if e.ActivityID != "" {
		fmt.Fprintf(&b, " (activity id: %s)", e.ActivityID)
	}
	return b.String()
}

// subStatusReadSessionNotAvailable is reported with 404 when the serving replica has not caught up
// to the supplied session token, and with 410 when the addressed partition key range is gone.
const subStatusReadSessionNotAvailable = 1002

// codeForStatus classifies a status and sub-status pair reported by the service.
//
// Cosmos DB overloads some status codes with sub-status codes that mean something categorically
// different from the status alone, so the sub-status is consulted first. Pairs that are not
// recognized fall back to the status, and statuses that are not recognized yield [CodeUnknown]
// rather than a guess.
//
//nolint:unused // consumed once operations translate driver errors.
func codeForStatus(statusCode, subStatus int) Code {
	if subStatus == subStatusReadSessionNotAvailable {
		switch statusCode {
		case 404:
			return CodeSessionUnavailable
		case 410:
			return CodePartitionKeyRangeGone
		}
	}

	switch statusCode {
	case 400:
		return CodeBadRequest
	case 401:
		return CodeUnauthorized
	case 403:
		return CodeForbidden
	case 404:
		return CodeNotFound
	case 408:
		return CodeRequestTimeout
	case 409:
		return CodeConflict
	case 410:
		return CodeGone
	case 412:
		return CodePreconditionFailed
	case 429:
		return CodeThrottled
	case 503:
		return CodeServiceUnavailable
	default:
		return CodeUnknown
	}
}
