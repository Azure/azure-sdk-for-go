// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// TransactionalBatchResponse contains the received transaction outcome and ordered operation results.
type TransactionalBatchResponse struct {
	Response

	// SessionToken is the session token reported for the batch.
	SessionToken SessionToken

	// ETag is the batch-level ETag, when supplied. Use OperationResults for item ETags.
	ETag azcore.ETag

	// RetryAfter is the batch-level retry delay, when supplied.
	RetryAfter time.Duration

	// Body is an owned copy of the native JSON result envelope.
	Body []byte

	// OperationResults has exactly one result per submitted operation, in the same order.
	OperationResults []TransactionalBatchResult

	// Success reports a committed transaction. A rolled-back transaction has Success=false
	// even though execution returned a response without an error, such as HTTP 207.
	// Native diagnostics describe request completion, not transaction commitment.
	Success bool
}

// FailedOperationIndex returns the zero-based index of the operation that caused rollback.
// The second return value is false for a successful batch or a zero response.
// Dependency failures (HTTP 424) are not attributed as the cause.
func (r TransactionalBatchResponse) FailedOperationIndex() (int, bool) {
	if !r.Success {
		for i, result := range r.OperationResults {
			if (result.StatusCode < 200 || result.StatusCode >= 300) && result.StatusCode != http.StatusFailedDependency {
				return i, true
			}
		}
	}
	return -1, false
}

// TransactionalBatchResult describes one operation's outcome.
type TransactionalBatchResult struct {
	// StatusCode is the operation's raw HTTP status. HTTP 424 means the transaction
	// was rolled back because another operation failed, not that this operation committed.
	StatusCode int

	// SubStatus is the operation's substatus, or zero when absent.
	SubStatus int

	// RequestCharge is the operation's request-unit charge, or zero when absent.
	RequestCharge float64

	// ResourceBody is owned raw JSON returned for this operation. It is nil when absent.
	// Reads return bodies regardless of the write-content-response option.
	ResourceBody []byte

	// ETag is the item ETag reported for this operation, when supplied.
	ETag azcore.ETag

	// RetryAfter is the operation's reported retry delay, or zero when absent.
	RetryAfter time.Duration
}

type transactionalBatchResultJSON struct {
	StatusCode             *int            `json:"statusCode"`
	SubStatus              *int            `json:"substatusCode"`
	RequestCharge          float64         `json:"requestCharge"`
	ResourceBody           json.RawMessage `json:"resourceBody"`
	ETag                   azcore.ETag     `json:"eTag"`
	RetryAfterMilliseconds uint64          `json:"retryAfterMilliseconds"`
}

func decodeTransactionalBatchResponse(response ItemResponse, body []byte, retryAfter time.Duration, fromWire bool, operationCount int) (TransactionalBatchResponse, error) {
	fail := func(cause error) (TransactionalBatchResponse, error) {
		return TransactionalBatchResponse{}, &Error{
			Code: CodeSerializationFailed, Message: "decoding transactional batch response",
			StatusCode: response.StatusCode, SubStatus: response.SubStatus,
			Diagnostics: response.Diagnostics, AttemptCount: response.AttemptCount,
			RequestCharge: response.RequestCharge, ActivityID: response.ActivityID,
			SessionToken: response.SessionToken, ETag: response.ETag, RetryAfter: retryAfter,
			Body: append([]byte(nil), body...), FromWire: fromWire, cause: cause,
		}
	}
	if !utf8.Valid(body) {
		return fail(errors.New("batch result envelope must contain valid UTF-8"))
	}
	var encoded []transactionalBatchResultJSON
	if err := json.Unmarshal(body, &encoded); err != nil {
		return fail(err)
	}
	if operationCount == 0 || len(encoded) != operationCount {
		return fail(fmt.Errorf("expected %d operation results; received %d", operationCount, len(encoded)))
	}
	result := TransactionalBatchResponse{
		Response: response.Response, SessionToken: response.SessionToken, ETag: response.ETag,
		RetryAfter: retryAfter, Body: append([]byte(nil), body...),
		OperationResults: make([]TransactionalBatchResult, len(encoded)),
		Success:          response.StatusCode != http.StatusMultiStatus,
	}
	failureCount := 0
	for i, operation := range encoded {
		if operation.StatusCode == nil || *operation.StatusCode <= 0 || *operation.StatusCode > math.MaxUint16 {
			return fail(fmt.Errorf("operation %d has no valid statusCode", i))
		}
		subStatus := 0
		if operation.SubStatus != nil {
			subStatus = *operation.SubStatus
			if subStatus < 0 || uint64(subStatus) > math.MaxUint32 {
				return fail(fmt.Errorf("operation %d has an invalid substatusCode", i))
			}
		}
		if operation.RequestCharge < 0 || operation.RetryAfterMilliseconds > uint64(math.MaxInt64/int64(time.Millisecond)) {
			return fail(fmt.Errorf("operation %d has invalid charge or retry metadata", i))
		}
		status := *operation.StatusCode
		succeeded := status >= 200 && status < 300
		if result.Success && !succeeded {
			return fail(fmt.Errorf("successful batch contains failing operation %d", i))
		}
		if !result.Success && succeeded {
			// The released in-memory emulator can return 207 after committing writes around a
			// conditional 304 read. Do not turn that contradictory envelope into a rollback claim.
			return fail(fmt.Errorf("rollback response contains successful operation %d", i))
		}
		if !succeeded && status != http.StatusFailedDependency {
			failureCount++
		}
		result.OperationResults[i] = TransactionalBatchResult{
			StatusCode: status, SubStatus: subStatus, RequestCharge: operation.RequestCharge,
			ResourceBody: append([]byte(nil), operation.ResourceBody...), ETag: operation.ETag,
			RetryAfter: time.Duration(operation.RetryAfterMilliseconds) * time.Millisecond,
		}
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusMultiStatus {
		return fail(fmt.Errorf("unexpected batch status %d", response.StatusCode))
	}
	if !result.Success && failureCount != 1 {
		return fail(errors.New("rollback response must identify exactly one failed operation"))
	}
	return result, nil
}
