// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ThroughputProperties selects manual or autoscale throughput. Exactly one must be set on writes.
type ThroughputProperties struct {
	// ManualThroughput is provisioned RU/s. Nil selects autoscale.
	ManualThroughput *uint64 `json:"offerThroughput,omitempty"`
	// Autoscale configures autoscale throughput. Nil selects manual throughput.
	Autoscale *AutoscaleSettings `json:"offerAutopilotSettings,omitempty"`
}

// AutoscaleSettings configures autoscale throughput.
type AutoscaleSettings struct {
	// MaxThroughput is the maximum RU/s.
	MaxThroughput uint64 `json:"maxThroughput"`
	// AutoUpgradePolicy optionally increases the autoscale maximum as storage grows.
	AutoUpgradePolicy *AutoscaleAutoUpgradePolicy `json:"autoUpgradePolicy,omitempty"`
}

// AutoscaleAutoUpgradePolicy configures automatic maximum-throughput increases.
type AutoscaleAutoUpgradePolicy struct {
	ThroughputPolicy *AutoscaleThroughputPolicy `json:"throughputPolicy,omitempty"`
}

// AutoscaleThroughputPolicy configures the percentage increase.
type AutoscaleThroughputPolicy struct {
	IncrementPercent uint32 `json:"incrementPercent"`
}

func (p ThroughputProperties) validate() error {
	if (p.ManualThroughput == nil) == (p.Autoscale == nil) {
		return errors.New("azcosmos: set exactly one of ManualThroughput or Autoscale")
	}
	if p.ManualThroughput != nil && *p.ManualThroughput == 0 || p.Autoscale != nil && p.Autoscale.MaxThroughput == 0 {
		return errors.New("azcosmos: throughput must be positive")
	}
	return nil
}

func managementCreateOptions(operation OperationOptions, throughput *ThroughputProperties) (OperationOptions, error) {
	if throughput == nil {
		return operation, nil
	}
	if err := throughput.validate(); err != nil {
		return OperationOptions{}, err
	}
	operation = operation.clone()
	for name := range operation.CustomHeaders {
		if strings.EqualFold(name, "x-ms-offer-throughput") || strings.EqualFold(name, "x-ms-cosmos-offer-autopilot-settings") {
			return OperationOptions{}, errors.New("azcosmos: throughput headers conflict with Throughput")
		}
	}
	return operation, nil
}

// ReadThroughputOptions configures a throughput read. Offer operations require account keys.
type ReadThroughputOptions struct {
	Operation OperationOptions
}

// BeginReplaceThroughputOptions configures throughput replacement.
type BeginReplaceThroughputOptions struct {
	Operation OperationOptions
	// IfMatchETag makes the offer replacement conditional. Nil is unconditional.
	IfMatchETag *azcore.ETag
}

// ThroughputResponse describes a resource's own offer, never its parent's offer.
type ThroughputResponse struct {
	Response
	// ThroughputProperties is nil when no throughput offer exists for this resource.
	// A shared-throughput container and a serverless resource have no dedicated offer.
	ThroughputProperties *ThroughputProperties
	ETag                 azcore.ETag
	SessionToken         SessionToken
}

type throughputOffer struct {
	ID      string               `json:"id"`
	ETag    azcore.ETag          `json:"_etag"`
	Content ThroughputProperties `json:"content"`
	fields  map[string]json.RawMessage
}

func decodeThroughputOffer(body []byte) (throughputOffer, error) {
	var offer throughputOffer
	if err := json.Unmarshal(body, &offer); err != nil {
		return offer, err
	}
	err := json.Unmarshal(body, &offer.fields)
	if err == nil && (offer.Content.ManualThroughput == nil && offer.Content.Autoscale == nil) {
		err = errors.New("offer response has no throughput content")
	}
	return offer, err
}

func (o throughputOffer) replaceBody(properties ThroughputProperties) ([]byte, error) {
	fields := maps.Clone(o.fields)
	content, err := json.Marshal(properties)
	if err != nil {
		return nil, err
	}
	fields["content"] = content
	fields["offerVersion"] = json.RawMessage(`"V2"`)
	return json.Marshal(fields)
}

// ReadThroughput reads shared database throughput. A nil properties pointer means no offer.
// Offer operations require account-key authentication, not Microsoft Entra ID.
func (d *DatabaseClient) ReadThroughput(ctx context.Context, options *ReadThroughputOptions) (ThroughputResponse, error) {
	return d.client.readThroughput(ctx, d.id, "", options)
}

// ReadThroughput reads dedicated container throughput. It never falls back to database throughput.
// Offer operations require account-key authentication, not Microsoft Entra ID.
func (c *ContainerClient) ReadThroughput(ctx context.Context, options *ReadThroughputOptions) (ThroughputResponse, error) {
	return c.database.client.readThroughput(ctx, c.database.id, c.id, options)
}

func (c *Client) readThroughput(ctx context.Context, databaseID, containerID string, options *ReadThroughputOptions) (ThroughputResponse, error) {
	operation := OperationOptions{}
	if options != nil {
		operation = options.Operation
	}
	if err := validateManagementContext(ctx, operation); err != nil {
		return ThroughputResponse{}, err
	}
	if !c.usesKeyAuth {
		return ThroughputResponse{}, errKeyCredentialRequired()
	}
	release, err := c.acquire()
	if err != nil {
		return ThroughputResponse{}, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return ThroughputResponse{}, err
	}
	response, _, err := c.findThroughputOffer(ctx, databaseID, containerID, operation)
	return response, err
}

// findThroughputOffer runs under the caller's lifetime guard. Every page contributes to RU totals.
func (c *Client) findThroughputOffer(ctx context.Context, databaseID, containerID string, operation OperationOptions) (ThroughputResponse, throughputOffer, error) {
	kind := operationKindReadDatabase
	if containerID != "" {
		kind = operationKindReadContainer
	}
	metadata, err := c.executeManagementRequest(ctx, managementRequest{kind: kind, databaseID: databaseID, containerID: containerID, operation: operation})
	if err != nil {
		return ThroughputResponse{}, throughputOffer{}, err
	}
	var resource struct {
		ResourceID string `json:"_rid"`
	}
	if err := json.Unmarshal(metadata.body, &resource); err != nil || resource.ResourceID == "" {
		if err == nil {
			err = errors.New("resource response has no _rid")
		}
		return ThroughputResponse{}, throughputOffer{}, newManagementSerializationError(metadata.response, metadata.body, "throughput resource", err)
	}
	query, err := NewQuery("SELECT * FROM c WHERE c.offerResourceId = @rid").WithParameter("@rid", resource.ResourceID)
	if err != nil {
		return ThroughputResponse{}, throughputOffer{}, addQuerySetupCharge(err, metadata.response.Response)
	}
	body, err := query.body()
	if err != nil {
		return ThroughputResponse{}, throughputOffer{}, addQuerySetupCharge(err, metadata.response.Response)
	}
	cursor, err := c.openManagementQuery(ctx, managementRequest{kind: operationKindQueryOffers, body: body, operation: operation}, ManagementQueryOptions{Operation: operation})
	if err != nil {
		return ThroughputResponse{}, throughputOffer{}, addQuerySetupCharge(err, metadata.response.Response)
	}
	defer cursor.close()
	total := metadata.response.Response
	var lastSessionToken SessionToken
	for {
		page, end, err := cursor.next(ctx)
		if err != nil {
			return ThroughputResponse{}, throughputOffer{}, addQuerySetupCharge(err, total)
		}
		charge := total.RequestCharge + page.RequestCharge
		if !end {
			total = page.Response
			lastSessionToken = page.SessionToken
		}
		total.RequestCharge = charge
		if len(page.Items) > 0 {
			offer, err := decodeThroughputOffer(page.Items[0])
			if err != nil {
				return ThroughputResponse{}, throughputOffer{}, newManagementSerializationError(ItemResponse{Response: total}, page.Items[0], "offer", err)
			}
			return ThroughputResponse{Response: total, ThroughputProperties: &offer.Content, ETag: offer.ETag, SessionToken: page.SessionToken}, offer, nil
		}
		if end {
			return ThroughputResponse{Response: total, SessionToken: lastSessionToken}, throughputOffer{}, nil
		}
	}
}

// BeginReplaceThroughput begins replacing shared database throughput.
// Poll the returned poller to observe terminal success or failure.
func (d *DatabaseClient) BeginReplaceThroughput(ctx context.Context, properties ThroughputProperties, options *BeginReplaceThroughputOptions) (*ThroughputPoller, error) {
	return d.client.beginReplaceThroughput(ctx, d.id, "", properties, options)
}

// BeginReplaceThroughput begins replacing dedicated container throughput.
// Shared-throughput containers return an error instead of targeting the database's offer.
func (c *ContainerClient) BeginReplaceThroughput(ctx context.Context, properties ThroughputProperties, options *BeginReplaceThroughputOptions) (*ThroughputPoller, error) {
	return c.database.client.beginReplaceThroughput(ctx, c.database.id, c.id, properties, options)
}

func (c *Client) beginReplaceThroughput(ctx context.Context, databaseID, containerID string, properties ThroughputProperties, options *BeginReplaceThroughputOptions) (*ThroughputPoller, error) {
	if err := properties.validate(); err != nil {
		return nil, err
	}
	operation := OperationOptions{}
	if options != nil {
		operation = options.Operation.clone()
	}
	replaceOperation := operation.clone()
	var precondition Precondition
	if options != nil && options.IfMatchETag != nil {
		var err error
		precondition, err = IfMatch(*options.IfMatchETag)
		if err != nil {
			return nil, err
		}
	}
	if err := validateManagementContext(ctx, operation); err != nil {
		return nil, err
	}
	if !c.usesKeyAuth {
		return nil, errKeyCredentialRequired()
	}
	release, err := c.acquire()
	if err != nil {
		return nil, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	current, offer, err := c.findThroughputOffer(ctx, databaseID, containerID, operation)
	if err != nil {
		return nil, err
	}
	if current.ThroughputProperties == nil {
		return nil, &Error{Code: CodeClientError, StatusCode: 400, SubStatus: 20301, RequestCharge: current.RequestCharge, ActivityID: current.ActivityID, Diagnostics: current.Diagnostics, AttemptCount: current.AttemptCount, SessionToken: current.SessionToken, Message: "azcosmos: no throughput offer exists for this resource"}
	}
	if offer.ID == "" || strings.ContainsRune(offer.ID, 0) {
		return nil, newManagementSerializationError(ItemResponse{Response: current.Response}, nil, "offer ID", errors.New("service returned an invalid offer ID"))
	}
	body, err := offer.replaceBody(properties)
	if err != nil {
		return nil, addQuerySetupCharge(err, current.Response)
	}
	// Conditional headers belong only to the replace, not metadata reads or later polls.
	forceContentResponseOnWrite(&replaceOperation)
	req := managementRequest{kind: operationKindReplaceOffer, resourceLink: offer.ID, body: body, operation: replaceOperation, precondition: precondition}
	initial, err := c.executeManagementRequest(ctx, req)
	if err != nil {
		return nil, addQuerySetupCharge(err, current.Response)
	}
	result, err := throughputResult(initial)
	if err != nil {
		return nil, addQuerySetupCharge(err, current.Response)
	}
	result.RequestCharge += current.RequestCharge
	readOptions := OperationOptions{}
	if options != nil {
		readOptions = options.Operation.clone()
	}
	return &ThroughputPoller{
		result: result, done: !initial.pending,
		poll: func(ctx context.Context) (managementResult, error) {
			if err := validateManagementContext(ctx, readOptions); err != nil {
				return managementResult{}, err
			}
			release, err := c.acquire()
			if err != nil {
				return managementResult{}, err
			}
			defer release()
			if err := ctx.Err(); err != nil {
				return managementResult{}, err
			}
			return c.executeManagementRequest(ctx, managementRequest{kind: operationKindReadOffer, resourceLink: offer.ID, operation: readOptions})
		},
	}, nil
}

func throughputResult(result managementResult) (ThroughputResponse, error) {
	out := ThroughputResponse{Response: result.response.Response, ETag: result.response.ETag, SessionToken: result.response.SessionToken}
	if len(result.body) == 0 && result.pending {
		return out, nil
	}
	offer, err := decodeThroughputOffer(result.body)
	if err != nil {
		return ThroughputResponse{}, newManagementSerializationError(result.response, result.body, "offer", err)
	}
	out.ThroughputProperties = &offer.Content
	return out, nil
}

// ThroughputPoller observes a throughput replacement. It retains no native handles.
// Stopping polling does not cancel the service-side replacement.
type ThroughputPoller struct {
	mu     sync.Mutex
	result ThroughputResponse
	done   bool
	err    error
	poll   func(context.Context) (managementResult, error)
}

// Done reports whether terminal success or failure has been observed.
func (p *ThroughputPoller) Done() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.done
}

// Poll reads the same offer by ID. Response charge includes submission and all polls.
// Context cancellation leaves the poller resumable; other failures are terminal.
func (p *ThroughputPoller) Poll(ctx context.Context) (ThroughputResponse, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if ctx == nil {
		return ThroughputResponse{}, errors.New("azcosmos: context must not be nil")
	}
	if p.done {
		if p.err != nil {
			return ThroughputResponse{}, p.err
		}
		return p.result, nil
	}
	native, err := p.poll(ctx)
	if err != nil {
		err = addQuerySetupCharge(err, p.result.Response)
		if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			p.done, p.err = true, err
		} else {
			var cosmosErr *Error
			if errors.As(err, &cosmosErr) {
				p.result.RequestCharge = cosmosErr.RequestCharge
			}
		}
		return ThroughputResponse{}, err
	}
	response, err := throughputResult(native)
	if err != nil {
		p.done, p.err = true, addQuerySetupCharge(err, p.result.Response)
		return ThroughputResponse{}, p.err
	}
	response.RequestCharge += p.result.RequestCharge
	p.result, p.done = response, !native.pending
	return response, nil
}

// Result returns the terminal response, or an error if pending or terminally failed.
func (p *ThroughputPoller) Result() (ThroughputResponse, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.done {
		return ThroughputResponse{}, errors.New("azcosmos: throughput replacement is pending")
	}
	if p.err != nil {
		return ThroughputResponse{}, p.err
	}
	return p.result, nil
}

// PollUntilDoneOptions configures the delay between offer reads.
type PollUntilDoneOptions struct {
	// Frequency defaults to five seconds. Negative values are invalid.
	Frequency time.Duration
}

// PollUntilDone waits for terminal success or failure, or ctx cancellation.
func (p *ThroughputPoller) PollUntilDone(ctx context.Context, options *PollUntilDoneOptions) (ThroughputResponse, error) {
	if ctx == nil {
		return ThroughputResponse{}, errors.New("azcosmos: context must not be nil")
	}
	frequency := 5 * time.Second
	if options != nil && options.Frequency != 0 {
		frequency = options.Frequency
	}
	if frequency < 0 {
		return ThroughputResponse{}, fmt.Errorf("azcosmos: polling frequency must not be negative")
	}
	for !p.Done() {
		timer := time.NewTimer(frequency)
		select {
		case <-ctx.Done():
			timer.Stop()
			p.mu.Lock()
			err := addSetupCharge(ctx.Err(), p.result.Response, "polling throughput replacement")
			p.mu.Unlock()
			return ThroughputResponse{}, err
		case <-timer.C:
		}
		if _, err := p.Poll(ctx); err != nil {
			return ThroughputResponse{}, err
		}
	}
	return p.Result()
}
