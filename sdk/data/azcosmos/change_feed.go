// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ChangeFeedMode selects which changes the service returns.
type ChangeFeedMode string

const (
	// ChangeFeedModeLatestVersion returns the latest versions of changed items.
	ChangeFeedModeLatestVersion ChangeFeedMode = "LatestVersion"
	// ChangeFeedModeAllVersionsAndDeletes returns intermediate versions and deletes.
	// The account must support this mode and checkpoints must remain within its retention window.
	ChangeFeedModeAllVersionsAndDeletes ChangeFeedMode = "AllVersionsAndDeletes"
)

// ChangeFeedStartFrom is an explicit initial position. Its zero value is invalid.
// When resuming with a checkpoint, supply the same start position as the original pager.
type ChangeFeedStartFrom struct {
	kind uint32
	time time.Time
}

// NewChangeFeedStartFromBeginning reads available history in LatestVersion mode.
// AllVersionsAndDeletes does not support this fresh start; the service rejects it.
func NewChangeFeedStartFromBeginning() ChangeFeedStartFrom {
	return ChangeFeedStartFrom{kind: 1}
}

// NewChangeFeedStartFromNow begins when the native driver first polls each range.
func NewChangeFeedStartFromNow() ChangeFeedStartFrom {
	return ChangeFeedStartFrom{kind: 2}
}

// NewChangeFeedStartFromPointInTime starts LatestVersion reads at the supplied time.
// The service uses second precision. AllVersionsAndDeletes rejects this fresh start.
func NewChangeFeedStartFromPointInTime(startTime time.Time) ChangeFeedStartFrom {
	return ChangeFeedStartFrom{kind: 3, time: startTime.UTC()}
}

// ChangeFeedOptions configures pull reads. Nil selects LatestVersion and native feed defaults.
type ChangeFeedOptions struct {
	// Operation holds shared settings. Its timeout applies to each native page fetch.
	// Context cancellation bounds the Go wait, not admitted native execution.
	Operation OperationOptions
	// Feed controls page sizing, initial fan-out and opaque checkpoint resumption.
	Feed FeedOptions
	// Mode selects the changes to read. Zero selects LatestVersion.
	Mode ChangeFeedMode
	// SessionToken overrides the client's captured session token. Empty uses the client's token.
	SessionToken SessionToken
}

// ChangeFeedResponse contains one native change-feed page.
// Empty and HTTP 304 pages remain pollable. One idle range does not imply container-wide catch-up.
type ChangeFeedResponse struct {
	Response
	// Items contains independently owned raw JSON feed entries in feed order.
	// Use DecodeChangeFeedItem to read entries that are change envelopes into a typed model.
	Items [][]byte
	// ETag is the page's service position, not a resumable planner checkpoint.
	ETag azcore.ETag
	// SessionToken is the session token returned for this page.
	SessionToken SessionToken
}

// ChangeFeedItem contains the item images and metadata supplied for one change.
// Each field is nil when its value is absent or null.
type ChangeFeedItem[T any] struct {
	Current  *T                  `json:"current"`
	Previous *T                  `json:"previous"`
	Metadata *ChangeFeedMetadata `json:"metadata"`
}

// ChangeFeedMetadata describes a change. The feed mode determines which fields are available.
type ChangeFeedMetadata struct {
	// OperationType is the operation reported by the service, such as "create", "replace",
	// or "delete". It is empty if not supplied. Other values can be returned.
	OperationType string `json:"operationType"`
	// LSN is the logical sequence number, or nil if not supplied.
	LSN *int64 `json:"lsn"`
}

// DecodeChangeFeedItem decodes one change envelope from ChangeFeedResponse.Items.
// It returns the zero item and an error if raw is not an envelope or a field fails to decode.
// If the backend returns a plain document instead, decode it directly into the document model.
// Unknown fields are ignored in the typed result but remain available in raw.
func DecodeChangeFeedItem[T any](raw []byte) (ChangeFeedItem[T], error) {
	var envelope struct {
		Current  json.RawMessage `json:"current"`
		Previous json.RawMessage `json:"previous"`
		Metadata json.RawMessage `json:"metadata"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return ChangeFeedItem[T]{}, err
	}
	if envelope.Current == nil && envelope.Previous == nil && envelope.Metadata == nil {
		return ChangeFeedItem[T]{}, errors.New("azcosmos: expected a change feed envelope")
	}
	var item ChangeFeedItem[T]
	if err := json.Unmarshal(raw, &item); err != nil {
		return ChangeFeedItem[T]{}, err
	}
	return item, nil
}

type changeFeedRequest struct {
	databaseID    string
	containerID   string
	partitionKey  PartitionKey
	fullContainer bool
	startFrom     ChangeFeedStartFrom
	options       ChangeFeedOptions
}

func newChangeFeedRequest(scope FeedScope, startFrom ChangeFeedStartFrom, options *ChangeFeedOptions) (changeFeedRequest, error) {
	req := changeFeedRequest{
		partitionKey: scope.partitionKey, fullContainer: scope.fullContainer, startFrom: startFrom,
	}
	if !scope.fullContainer {
		if err := validateItemArguments(scope.partitionKey); err != nil {
			return req, err
		}
	}
	switch startFrom.kind {
	case 1, 2:
	case 3:
		if startFrom.time.Year() < 0 || startFrom.time.Year() > 9999 {
			return req, errors.New("azcosmos: change feed start time must have an RFC 3339 year")
		}
	default:
		return req, errors.New("azcosmos: change feed requires an explicit start position")
	}
	if options != nil {
		req.options = *options
		req.options.Operation = options.Operation.clone()
	}
	switch req.options.Mode {
	case "":
		req.options.Mode = ChangeFeedModeLatestVersion
	case ChangeFeedModeLatestVersion, ChangeFeedModeAllVersionsAndDeletes:
	default:
		return req, errors.New("azcosmos: invalid change feed mode")
	}
	if err := req.options.Operation.validate(); err != nil {
		return req, err
	}
	if err := req.options.SessionToken.validate(); err != nil {
		return req, err
	}
	if req.options.Feed.PageSizeHint < 0 {
		return req, errors.New("azcosmos: page size hint must not be negative")
	}
	if strings.IndexByte(req.options.Feed.ContinuationToken, 0) >= 0 {
		return req, errors.New("azcosmos: continuation token must not contain a NUL byte")
	}
	return req, nil
}

func decodeChangeFeedPage(body []byte, page ChangeFeedResponse) (ChangeFeedResponse, error) {
	if len(body) == 0 {
		return page, nil
	}
	items, err := decodeFeedItems(body)
	if err != nil {
		return ChangeFeedResponse{}, changeFeedResponseError(page, err)
	}
	page.Items = items
	return page, nil
}

func changeFeedResponseError(page ChangeFeedResponse, cause error) *Error {
	return &Error{
		Code: CodeSerializationFailed, Message: "decoding change feed response",
		Diagnostics: page.Diagnostics, StatusCode: page.StatusCode, SubStatus: page.SubStatus,
		AttemptCount: page.AttemptCount, RequestCharge: page.RequestCharge, ActivityID: page.ActivityID,
		SessionToken: page.SessionToken, ETag: page.ETag, cause: cause,
	}
}
