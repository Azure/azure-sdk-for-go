// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

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

	// CodeServiceError means the service responded with a failure that this package does not
	// classify further. [Error.StatusCode] carries what it reported.
	CodeServiceError Code = "ServiceError"

	// CodeSessionUnavailable means the replica that served the request had not caught up to the
	// supplied session token. The resource itself may well exist, so this is deliberately not
	// [CodeNotFound] even though the service reports it with the same status code.
	CodeSessionUnavailable Code = "SessionUnavailable"

	// CodeClientError means the client produced the failure without a response from the service,
	// for example a transport, serialization or configuration failure.
	CodeClientError Code = "ClientError"

	// CodeClientClosed means the operation was attempted on a [Client] that has been closed.
	CodeClientClosed Code = "ClientClosed"

	// CodeTransportFailure means the request could not be delivered to the service, for example
	// because the connection failed or DNS could not be resolved.
	CodeTransportFailure Code = "TransportFailure"

	// CodeSerializationFailed means a payload could not be encoded or decoded.
	CodeSerializationFailed Code = "SerializationFailed"

	// CodeAuthenticationFailed means a token could not be acquired for the credential. It differs
	// from [CodeUnauthorized], which is the service rejecting a token that was acquired.
	CodeAuthenticationFailed Code = "AuthenticationFailed"

	// CodeClientOperationTimeout means the operation exhausted its client-side budget before the
	// service responded. It differs from [CodeRequestTimeout], which the service reports.
	CodeClientOperationTimeout Code = "ClientOperationTimeout"

	// CodeOperationCancelled means an operation already in flight was cancelled before it
	// completed, because the caller's context ended or the native driver cancelled it. An [Error]
	// carrying it unwraps to the context error that caused cancellation,
	// so errors.Is reports it the same way the rest of the standard library does. An operation
	// started after the client was closed is [CodeClientClosed] instead, since it never ran.
	CodeOperationCancelled Code = "OperationCancelled"
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
// A failed operation returns the zero response value, so this type is what carries whatever the
// service sent back with the failure. That includes the request charge: Cosmos DB bills for failed
// requests, so [Error.RequestCharge] is the only place to account for them.
//
// [Error.Code] is part of the published API. The string returned by [Error.Error] is not, and is
// subject to change.
type Error struct {
	// Code classifies the failure. Prefer it over StatusCode and SubStatus, which are reported
	// verbatim and are harder to interpret correctly.
	Code Code

	// StatusCode is the HTTP status code reported for the operation, as defined in
	// https://pkg.go.dev/net/http#pkg-constants. It is zero when the failure occurred before the
	// service responded.
	StatusCode int

	// SubStatus is the Cosmos DB sub-status code, which qualifies StatusCode
	// (`x-ms-substatus`). It is zero when the service did not report one.
	//
	// Its values are only meaningful together with StatusCode, which reuses them: sub-status 1002
	// means the session token could not be satisfied under status 404, and that the partition key
	// range is gone under status 410. Prefer Code, which resolves the pair.
	SubStatus int

	// Message describes the failure.
	Message string

	// RequestCharge is the number of request units the failed operation consumed
	// (`x-ms-request-charge`). Failed requests are still billed, most notably when they are
	// throttled. See https://learn.microsoft.com/azure/cosmos-db/request-units.
	RequestCharge float64

	// ActivityID correlates the operation with server-side telemetry (`x-ms-activity-id`). It is
	// empty when the failure occurred before the service responded.
	ActivityID string

	// SessionToken is the session token reported alongside the failure, if any
	// (`x-ms-session-token`).
	SessionToken SessionToken

	// ETag is the entity tag reported alongside the failure, if any (`etag`).
	ETag azcore.ETag

	// RetryAfter is how long the service asked the caller to wait before retrying
	// (`x-ms-retry-after-ms`). It is zero when the service did not ask for a delay.
	RetryAfter time.Duration

	// Body is the error document the service returned, verbatim and unparsed. It is nil for
	// failures the client produced, and is the payload to quote when reporting a problem.
	Body []byte

	// FromWire reports whether the service produced this failure, which is to say whether the
	// request reached it at all. It is false for failures the client produced without a service
	// response, such as transport and serialization failures.
	//
	// StatusCode does not imply this: the driver synthesizes statuses for client-side failures
	// too, so a non-zero status is not evidence the service replied. The distinction matters most
	// for writes, where a wire failure means the operation definitely reached the service and a
	// client-side one leaves that unknown. The response-derived fields above are only populated
	// when this is true.
	FromWire bool

	cause error
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
	case e.StatusCode != 0 && e.SubStatus > 0:
		fmt.Fprintf(&b, ": status %d/%d", e.StatusCode, e.SubStatus)
	case e.StatusCode != 0:
		fmt.Fprintf(&b, ": status %d", e.StatusCode)
	case e.SubStatus > 0:
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

// Unwrap reports the standard library error a failure corresponds to, so that callers can use
// errors.Is for the conditions Go already has a vocabulary for. A cancelled operation unwraps to
// its caller's context error, or [context.Canceled] when the native driver initiated cancellation.
func (e *Error) Unwrap() error {
	if e.cause != nil {
		return e.cause
	}
	if e.Code == CodeOperationCancelled {
		return context.Canceled
	}
	return nil
}

func newOperationCancelledError(cause error, requestCharge float64, activityID string) *Error {
	return &Error{
		Code:          CodeOperationCancelled,
		Message:       "azcosmos: the operation was cancelled",
		RequestCharge: requestCharge,
		ActivityID:    activityID,
		cause:         cause,
	}
}

func cloneError(err error) error {
	cosmosErr, ok := err.(*Error)
	if !ok {
		return err
	}
	clone := *cosmosErr
	clone.Body = append([]byte(nil), cosmosErr.Body...)
	return &clone
}

// subStatusReadSessionNotAvailable is reported with 404 when the serving replica has not caught up
// to the supplied session token, and with 410 when the addressed partition key range is gone.
const subStatusReadSessionNotAvailable = 1002

// driverStatus mirrors cosmos_error_code_t, the coarse classification the driver reports on every
// completion. The driver populates it whether or not the service responded, which is why it, and
// not the HTTP status, is what [Code] is derived from: failures the client produced carry no HTTP
// status at all.
type driverStatus int32

// The subset of cosmos_error_code_t that maps onto a distinct [Code]. Codes outside this set are
// classified by the band they fall in; see codeForDriverStatus.
const (
	driverStatusSuccess                driverStatus = 0
	driverStatusBadRequest             driverStatus = 2400
	driverStatusUnauthorized           driverStatus = 2401
	driverStatusForbidden              driverStatus = 2403
	driverStatusNotFound               driverStatus = 2404
	driverStatusTimeout                driverStatus = 2408
	driverStatusConflict               driverStatus = 2409
	driverStatusGone                   driverStatus = 2410
	driverStatusPreconditionFailed     driverStatus = 2412
	driverStatusThrottled              driverStatus = 2429
	driverStatusServiceUnavailable     driverStatus = 2503
	driverStatusServiceError           driverStatus = 2999
	driverStatusClientError            driverStatus = 3001
	driverStatusTransportFailure       driverStatus = 3002
	driverStatusSerializationFailed    driverStatus = 3003
	driverStatusAuthenticationFailed   driverStatus = 3004
	driverStatusClientOperationTimeout driverStatus = 3005
	driverStatusOperationCancelled     driverStatus = 4012
)

// The bands cosmos_error_code_t is laid out in. The header requires consumers to classify codes
// they do not recognize by band rather than rejecting them, so that codes added later stay usable:
//
//	1..=999     FFI and argument validation
//	1001..=1999 auth and conversion
//	2001..=2999 mapped from a wire HTTP status
//	3001..=3999 FFI plumbing
//	4001..=4999 driver-wrapper fatal
//	5001..=5999 non-fatal warnings, where the operation still delivered its result
const (
	driverStatusWireMin     driverStatus = 2001
	driverStatusPlumbingMin driverStatus = 3001
	driverStatusWarningMin  driverStatus = 5001

	// driverStatusWireBase is the offset the driver adds to an HTTP status to form its wire code,
	// so that 404 becomes driverStatusNotFound. See driverStatusForHTTP.
	driverStatusWireBase driverStatus = 2000
)

// normalizeSubStatus maps the driver's absent sentinel onto the zero value [Error.SubStatus]
// documents. cosmos_completion_t reports -1 when the service sent no sub-status, which would
// otherwise read as a real sub-status and surface in [Error.Error].
//
// Any negative value is treated as absent, not just -1, so a sentinel added later cannot leak.
//
//nolint:unused // consumed once operations translate driver errors.
func normalizeSubStatus(subStatus int32) int {
	if subStatus < 0 {
		return 0
	}
	return int(subStatus)
}

// Sub-status codes the driver reports on failures it produced itself. It pairs them with a
// synthetic 408 or 503, so the sub-status rather than the status is what says what went wrong.
const (
	subStatusTransportGenerated503  = 20003
	subStatusClientOperationTimeout = 20008
	subStatusTransportMin           = 20010
	subStatusTransportMax           = 20015
	// Serialization covers both directions: 20020 is a response body the client could not read,
	// 20021 a request body it could not write.
	subStatusSerializationMin = 20020
	subStatusSerializationMax = 20021
	// The driver treats a signature it generated a 401 for and a token it could not acquire as
	// the same class of failure, and so does this.
	subStatusClientGenerated401   = 20401
	subStatusAuthenticationFailed = 20402
)

// codeForSyntheticSubStatus classifies a failure the client produced rather than the service.
//
// These carry a synthetic HTTP status, so classifying them by it would report a transport failure
// as [CodeServiceUnavailable] and a client-side timeout as [CodeRequestTimeout] — in both cases
// blaming the service for something local. The driver's own classifier consults the sub-status
// first for exactly this reason, and this mirrors it.
func codeForSyntheticSubStatus(subStatus int) Code {
	switch {
	case subStatus == subStatusClientGenerated401, subStatus == subStatusAuthenticationFailed:
		return CodeAuthenticationFailed
	case subStatus == subStatusTransportGenerated503,
		subStatus >= subStatusTransportMin && subStatus <= subStatusTransportMax:
		return CodeTransportFailure
	case subStatus == subStatusClientOperationTimeout:
		return CodeClientOperationTimeout
	case subStatus >= subStatusSerializationMin && subStatus <= subStatusSerializationMax:
		return CodeSerializationFailed
	default:
		return CodeClientError
	}
}

// codeForRichError classifies a failure reported through a rich cosmos_error_t or an equivalent
// HTTP and sub-status pair, following the driver's own order: a failure the client produced is
// classified by its sub-status, and only a wire response is classified by its HTTP status.
func codeForRichError(fromWire bool, httpStatus, subStatus int) Code {
	if !fromWire {
		return codeForSyntheticSubStatus(subStatus)
	}
	return codeForDriverStatus(driverStatusForHTTP(httpStatus), subStatus)
}

// driverStatusForHTTP maps a wire HTTP status onto the driver status that reports it, so that a
// packed cosmos_status_code_t, which carries an HTTP status, classifies through the same table as
// a completion, which carries a driver status.
//
// The driver derives its wire codes as 2000 plus the HTTP status, so the mapping is that offset
// rather than a second table that could drift from the first. A status with no dedicated code
// still lands in the wire band and classifies as [CodeServiceError], matching what the driver does
// with one.
func driverStatusForHTTP(httpStatus int) driverStatus {
	if httpStatus <= 0 {
		return driverStatusClientError
	}
	return driverStatusWireBase + driverStatus(httpStatus)
}

// codeForDriverStatus classifies a completion from the coarse status the driver reports, refined
// by the sub-status where the sub-status changes the meaning.
//
// Cosmos DB overloads sub-status 1002 across two statuses that mean something categorically
// different from the status alone, and the driver does not model that distinction, so it is
// applied here.
//
// A status the driver adds later still classifies, because codes that are not recognized fall back
// to their band: a wire code becomes [CodeServiceError] and a client-side code [CodeClientError].
// Only a status outside every band, or a success, yields [CodeUnknown].
//
//nolint:unused // consumed once operations translate driver errors.
func codeForDriverStatus(status driverStatus, subStatus int) Code {
	if subStatus == subStatusReadSessionNotAvailable {
		switch status {
		case driverStatusNotFound:
			return CodeSessionUnavailable
		case driverStatusGone:
			return CodePartitionKeyRangeGone
		}
	}

	switch status {
	case driverStatusBadRequest:
		return CodeBadRequest
	case driverStatusUnauthorized:
		return CodeUnauthorized
	case driverStatusForbidden:
		return CodeForbidden
	case driverStatusNotFound:
		return CodeNotFound
	case driverStatusTimeout:
		return CodeRequestTimeout
	case driverStatusConflict:
		return CodeConflict
	case driverStatusGone:
		return CodeGone
	case driverStatusPreconditionFailed:
		return CodePreconditionFailed
	case driverStatusThrottled:
		return CodeThrottled
	case driverStatusServiceUnavailable:
		return CodeServiceUnavailable
	case driverStatusServiceError:
		return CodeServiceError
	case driverStatusClientError:
		return CodeClientError
	case driverStatusTransportFailure:
		return CodeTransportFailure
	case driverStatusSerializationFailed:
		return CodeSerializationFailed
	case driverStatusAuthenticationFailed:
		return CodeAuthenticationFailed
	case driverStatusClientOperationTimeout:
		return CodeClientOperationTimeout
	case driverStatusOperationCancelled:
		return CodeOperationCancelled
	}

	switch {
	case status <= driverStatusSuccess:
		// Success is not a failure, and a negative code is not one the ABI defines.
		return CodeUnknown
	case status >= driverStatusWarningMin:
		// The warning band still delivers its result, so it is not a failure either and no
		// operation should be building an Error from it.
		return CodeUnknown
	case status >= driverStatusWireMin && status < driverStatusPlumbingMin:
		// The service responded, just not with something this package models.
		return CodeServiceError
	default:
		// Everything else is produced locally: argument validation, auth and conversion, FFI
		// plumbing, and the driver's own fatal codes. They are all caller-visible client
		// failures, so they share one code rather than spending a public constant each.
		return CodeClientError
	}
}
