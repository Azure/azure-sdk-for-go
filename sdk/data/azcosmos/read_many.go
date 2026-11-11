// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
)

// maxPartitionKeyComponents is the deepest hierarchical partition key Cosmos DB supports.
const maxPartitionKeyComponents = 3

// ItemIdentity names one item by partition key and ID for [NewReadManyItemsSelection].
// The zero value is invalid; use [NewItemIdentity].
type ItemIdentity struct {
	partitionKey PartitionKey
	id           string
}

// NewItemIdentity identifies an item. A hierarchical partition key prefix selects every item
// with that ID across the logical partitions sharing the prefix.
func NewItemIdentity(partitionKey PartitionKey, id string) ItemIdentity {
	return ItemIdentity{partitionKey: partitionKey, id: id}
}

type readManySelectionKind uint32

const (
	readManySelectionItems      readManySelectionKind = 1
	readManySelectionPartitions readManySelectionKind = 2
)

// ReadManySelection describes which items [ContainerClient.NewReadManyItemsPager] reads.
// Values are immutable and can be shared. The zero value is invalid; use
// [NewReadManyItemsSelection] or [NewReadManyPartitionsSelection].
type ReadManySelection struct {
	kind       readManySelectionKind
	identities []ItemIdentity
	filter     string
	parameters []queryParameter
}

// NewReadManyItemsSelection selects items by partition key and ID. The items are copied.
// An empty selection reads nothing; it never selects the full container.
func NewReadManyItemsSelection(items ...ItemIdentity) ReadManySelection {
	return ReadManySelection{kind: readManySelectionItems, identities: append([]ItemIdentity(nil), items...)}
}

// NewReadManyPartitionsSelection selects every item in the given logical partitions or
// hierarchical partition key prefixes. The keys are copied. An empty selection reads nothing;
// it never selects the full container.
func NewReadManyPartitionsSelection(partitionKeys ...PartitionKey) ReadManySelection {
	identities := make([]ItemIdentity, len(partitionKeys))
	for i, partitionKey := range partitionKeys {
		identities[i] = ItemIdentity{partitionKey: partitionKey}
	}
	return ReadManySelection{kind: readManySelectionPartitions, identities: identities}
}

// WithFilter returns a new selection that also applies a scalar predicate to each selected item.
// The predicate refers to the item as c and omits WHERE, for example "c.status = @status".
// Empty clears the filter.
func (s ReadManySelection) WithFilter(predicate string) ReadManySelection {
	s.filter = predicate
	return s
}

// WithParameter returns a new selection with a JSON snapshot of value for the filter, replacing
// any parameter with the same name. Names must start with @; values must be JSON-serializable.
func (s ReadManySelection) WithParameter(name string, value any) (ReadManySelection, error) {
	if len(name) < 2 || name[0] != '@' || strings.ContainsAny(name, " \t\r\n\x00") {
		return ReadManySelection{}, errors.New("azcosmos: read many parameter name must start with @ and contain no whitespace or NUL")
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return ReadManySelection{}, &Error{Code: CodeSerializationFailed, Message: "encoding read many parameter", cause: err}
	}
	parameters := append([]queryParameter(nil), s.parameters...)
	for i := range parameters {
		if parameters[i].Name == name {
			parameters[i].Value = encoded
			s.parameters = parameters
			return s, nil
		}
	}
	parameters = append(parameters, queryParameter{Name: name, Value: encoded})
	s.parameters = parameters
	return s, nil
}

func (s ReadManySelection) validate() error {
	if s.kind != readManySelectionItems && s.kind != readManySelectionPartitions {
		return errors.New("azcosmos: read many selection must be created with NewReadManyItemsSelection or NewReadManyPartitionsSelection")
	}
	for i, identity := range s.identities {
		if n := identity.partitionKey.Len(); n == 0 || n > maxPartitionKeyComponents {
			return fmt.Errorf("azcosmos: read many selection %d must have between 1 and %d partition key components; got %d", i, maxPartitionKeyComponents, n)
		}
		if s.kind == readManySelectionItems {
			if identity.id == "" {
				return fmt.Errorf("azcosmos: read many selection %d must have a non-empty item id", i)
			}
			if strings.IndexByte(identity.id, 0) >= 0 {
				return fmt.Errorf("azcosmos: read many selection %d item id must not contain a NUL byte", i)
			}
		}
	}
	if strings.IndexByte(s.filter, 0) >= 0 {
		return errors.New("azcosmos: read many filter must not contain a NUL byte")
	}
	if len(s.parameters) != 0 && strings.TrimSpace(s.filter) == "" {
		return errors.New("azcosmos: read many parameters require a filter")
	}
	return nil
}

// ReadManyOptions configures [ContainerClient.NewReadManyItemsPager]. A nil *ReadManyOptions
// selects defaults.
type ReadManyOptions struct {
	// Operation holds shared settings. Its timeout applies separately to each page fetch.
	Operation OperationOptions

	// SessionToken overrides the client's captured session token. Empty uses the client's token.
	SessionToken SessionToken

	// PageSizeHint is a positive maximum-item-count hint, not a guaranteed page size.
	// Zero uses the driver default. Negative values are invalid.
	PageSizeHint int32

	// MaxFanOut limits the physical partitions read concurrently. Zero uses the native default.
	MaxFanOut uint32
}

type readManyRequest struct {
	databaseID  string
	containerID string
	selection   ReadManySelection
	options     ReadManyOptions
}

func newReadManyRequest(selection ReadManySelection, options *ReadManyOptions) (readManyRequest, error) {
	req := readManyRequest{selection: selection}
	if err := selection.validate(); err != nil {
		return req, err
	}
	if options != nil {
		req.options = *options
		req.options.Operation = options.Operation.clone()
	}
	if err := req.options.Operation.validate(); err != nil {
		return req, err
	}
	if err := req.options.SessionToken.validate(); err != nil {
		return req, err
	}
	if req.options.PageSizeHint < 0 {
		return req, errors.New("azcosmos: page size hint must not be negative")
	}
	return req, nil
}

// NewReadManyItemsPager reads the selected items, decoding each as T with encoding/json.
// Missing items are omitted rather than reported as errors, and results are not ordered.
// Construction performs no network I/O; argument errors are reported by NextPage.
//
// Defer Close when stopping before exhaustion. Empty pages do not imply exhaustion.
// Read many cannot be resumed with a continuation token; restart with a new pager instead.
func (c *ContainerClient) NewReadManyItemsPager[T any](selection ReadManySelection, options *ReadManyOptions) *ReadManyItemsPager[T] {
	req, err := newReadManyRequest(selection, options)
	req.databaseID = c.database.id
	req.containerID = c.id
	return &ReadManyItemsPager[T]{client: c.database.client, req: req, validationErr: err}
}

// ReadManyItemsPager iterates read many results. Close releases native resources when
// iteration stops early. Iteration calls must not be concurrent. Close waits for native
// fetching, but not model conversion.
type ReadManyItemsPager[T any] struct {
	mu            sync.Mutex
	client        *Client
	req           readManyRequest
	validationErr error
	cursor        queryCursor
	done          bool
	converting    bool
}

// More reports whether another page may be fetched. It does not perform I/O.
func (p *ReadManyItemsPager[T]) More() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return !p.done
}

// NextPage fetches one page. A final empty page can signal exhaustion through More.
// Failures after execution starts terminate this pager; retries belong to the native driver.
// The context bounds the Go wait, not admitted native execution.
// Decoding failures return a zero page and stop iteration.
func (p *ReadManyItemsPager[T]) NextPage(ctx context.Context) (QueryItemsResponse[T], error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.validationErr != nil {
		return QueryItemsResponse[T]{}, p.validationErr
	}
	if ctx == nil {
		return QueryItemsResponse[T]{}, errors.New("azcosmos: context must not be nil")
	}
	release, err := p.client.acquire()
	if err != nil {
		return QueryItemsResponse[T]{}, err
	}
	defer func() {
		if release != nil {
			release()
		}
	}()
	if p.converting {
		return QueryItemsResponse[T]{}, &Error{Code: CodeClientError, Message: "azcosmos: read many page conversion is in progress"}
	}
	if p.done {
		return QueryItemsResponse[T]{}, &Error{Code: CodeClientError, Message: "azcosmos: read many pager is closed or exhausted"}
	}
	if err := ctx.Err(); err != nil {
		return QueryItemsResponse[T]{}, err
	}
	if len(p.req.selection.identities) == 0 {
		p.finish()
		return QueryItemsResponse[T]{}, nil
	}
	ctx, cancel, err := p.client.queryContext(ctx, p.req.options.Operation)
	if err != nil {
		return QueryItemsResponse[T]{}, err
	}
	defer cancel()
	if p.cursor == nil {
		p.cursor, err = p.client.openReadMany(ctx, &p.req)
		if err != nil {
			p.finish()
			return QueryItemsResponse[T]{}, err
		}
	}
	page, end, err := p.cursor.next(ctx)
	if err != nil {
		p.finish()
		return QueryItemsResponse[T]{}, err
	}
	p.converting = true
	defer func() { p.converting = false }()
	release()
	release = nil
	decoded, err := func() (QueryItemsResponse[T], error) {
		p.mu.Unlock()
		defer p.mu.Lock()
		return decodeModels[T](ctx, page)
	}()
	if err != nil {
		p.finish()
		return QueryItemsResponse[T]{}, err
	}
	if end {
		p.finish()
	}
	return decoded, nil
}

// Close releases native resources. It is idempotent and safe alongside Client.Close.
func (p *ReadManyItemsPager[T]) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.finish()
	return nil
}

func (p *ReadManyItemsPager[T]) finish() {
	p.done = true
	if p.cursor != nil {
		p.cursor.close()
		p.cursor = nil
	}
}
