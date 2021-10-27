// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/devigned/tab"
)

func (ac *AdminClient) createOrUpdateQueueImpl(ctx context.Context, props *QueueProperties, creating bool) (*QueueProperties, error) {
	if props == nil {
		return nil, errors.New("properties are required and cannot be nil")
	}

	reqBytes, mw, err := serializeQueueProperties(props, ac.em.TokenProvider())

	if err != nil {
		return nil, err
	}

	if !creating {
		// an update requires the entity to already exist.
		mw = append(mw, func(next atom.RestHandler) atom.RestHandler {
			return func(ctx context.Context, req *http.Request) (*http.Response, error) {
				req.Header.Set("If-Match", "*")
				return next(ctx, req)
			}
		})
	}

	resp, err := ac.em.Put(ctx, "/"+props.Name, reqBytes, mw...)
	defer atom.CloseRes(ctx, resp)

	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	atomResp, err := deserializeQueueEnvelope(resp.Body)

	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	return newQueueProperties(props.Name, &atomResp.Content.QueueDescription)
}

// queuePropertiesPager provides iteration over QueueProperties pages.
type queuePropertiesPager struct {
	adminClient *AdminClient

	pageSize int
	skip     int

	lastErr      error
	lastResponse []*QueueProperties
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *queuePropertiesPager) NextPage(ctx context.Context) bool {
	p.lastResponse, p.lastErr = p.getQueuePage(ctx)
	p.skip += len(p.lastResponse)
	return p.lastResponse != nil
}

// PageResponse returns the current page.
func (p *queuePropertiesPager) PageResponse() []*QueueProperties {
	return p.lastResponse
}

// Err returns the last error encountered while paging.
func (p *queuePropertiesPager) Err() error {
	return p.lastErr
}

func (p *queuePropertiesPager) getQueuePage(ctx context.Context) ([]*QueueProperties, error) {
	var envelopes []atom.QueueEnvelope
	envelopes, err := p.adminClient.getQueuePage(ctx, p.pageSize, p.skip)

	if err != nil {
		return nil, err
	}

	var all []*QueueProperties

	for _, env := range envelopes {
		props, err := newQueueProperties(env.Title, &env.Content.QueueDescription)

		if err != nil {
			return nil, err
		}

		all = append(all, props)
	}

	return all, nil
}

// queueRuntimePropertiesPager provides iteration over QueueRuntimeProperties pages.
type queueRuntimePropertiesPager struct {
	adminClient *AdminClient

	pageSize int
	skip     int

	lastErr      error
	lastResponse []*QueueRuntimeProperties
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *queueRuntimePropertiesPager) NextPage(ctx context.Context) bool {
	p.lastResponse, p.lastErr = p.getQueuePage(ctx)
	p.skip += len(p.lastResponse)
	return p.lastResponse != nil
}

// PageResponse returns the current page.
func (p *queueRuntimePropertiesPager) PageResponse() []*QueueRuntimeProperties {
	return p.lastResponse
}

// Err returns the last error encountered while paging.
func (p *queueRuntimePropertiesPager) Err() error {
	return p.lastErr
}

func (p *queueRuntimePropertiesPager) getQueuePage(ctx context.Context) ([]*QueueRuntimeProperties, error) {
	var envelopes []atom.QueueEnvelope
	envelopes, err := p.adminClient.getQueuePage(ctx, p.pageSize, p.skip)

	if err != nil {
		return nil, err
	}

	var all []*QueueRuntimeProperties

	for _, env := range envelopes {
		runtimeProps := newQueueRuntimeProperties(env.Title, &env.Content.QueueDescription)

		if err != nil {
			return nil, err
		}

		all = append(all, runtimeProps)
	}

	return all, nil
}

func (ac *AdminClient) getQueuePage(ctx context.Context, top int, skip int) ([]atom.QueueEnvelope, error) {
	fragment := "/$Resources/Queues?"

	if top > 0 {
		fragment += fmt.Sprintf("&$top=%d", top)
	}

	if skip > 0 {
		fragment += fmt.Sprintf("&$skip=%d", skip)
	}

	resp, err := ac.em.Get(ctx, fragment)

	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var feed *atom.QueueFeed

	if err := xml.Unmarshal(bytes, &feed); err != nil {
		return nil, err
	}

	return feed.Entries, nil
}

func (ac *AdminClient) getQueueImpl(ctx context.Context, queueName string) (string, *atom.QueueDescription, error) {
	ctx, span := tab.StartSpan(ctx, tracing.SpanGetEntity)
	defer span.End()

	resp, err := ac.em.Get(ctx, "/"+queueName)

	if err != nil {
		return "", nil, err
	}

	atomResp, err := deserializeQueueEnvelope(resp.Body)

	if err != nil {
		return "", nil, err
	}

	return atomResp.Title, &atomResp.Content.QueueDescription, nil
}

func deserializeQueueEnvelope(body io.Reader) (*atom.QueueEnvelope, error) {
	b, err := ioutil.ReadAll(body)

	if err != nil {
		return nil, err
	}

	var atomResp atom.QueueEnvelope

	if err := xml.Unmarshal(b, &atomResp); err != nil {
		return nil, atom.FormatManagementError(b)
	}

	return &atomResp, nil
}

func serializeQueueProperties(props *QueueProperties, tokenProvider auth.TokenProvider) ([]byte, []atom.MiddlewareFunc, error) {
	qpr := &atom.QueueDescription{
		LockDuration:                        utils.DurationToStringPtr(props.LockDuration),
		MaxSizeInMegabytes:                  props.MaxSizeInMegabytes,
		RequiresDuplicateDetection:          props.RequiresDuplicateDetection,
		RequiresSession:                     props.RequiresSession,
		DefaultMessageTimeToLive:            utils.DurationToStringPtr(props.DefaultMessageTimeToLive),
		DeadLetteringOnMessageExpiration:    props.DeadLetteringOnMessageExpiration,
		DuplicateDetectionHistoryTimeWindow: utils.DurationToStringPtr(props.DuplicateDetectionHistoryTimeWindow),
		MaxDeliveryCount:                    props.MaxDeliveryCount,
		EnableBatchedOperations:             props.EnableBatchedOperations,
		Status:                              (*atom.EntityStatus)(props.Status),
		AutoDeleteOnIdle:                    utils.DurationToStringPtr(props.AutoDeleteOnIdle),
		EnablePartitioning:                  props.EnablePartitioning,
		ForwardTo:                           props.ForwardTo,
		ForwardDeadLetteredMessagesTo:       props.ForwardDeadLetteredMessagesTo,
	}

	atomRequest, mw := atom.ConvertToQueueRequest(qpr, tokenProvider)

	bytes, err := xml.MarshalIndent(atomRequest, "", "  ")

	if err != nil {
		return nil, nil, err
	}

	return bytes, mw, nil
}

func newQueueProperties(name string, desc *atom.QueueDescription) (*QueueProperties, error) {
	lockDuration, err := utils.ISO8601StringToDuration(desc.LockDuration)

	if err != nil {
		return nil, err
	}

	defaultMessageTimeToLive, err := utils.ISO8601StringToDuration(desc.DefaultMessageTimeToLive)

	if err != nil {
		return nil, err
	}

	duplicateDetectionHistoryTimeWindow, err := utils.ISO8601StringToDuration(desc.DuplicateDetectionHistoryTimeWindow)

	if err != nil {
		return nil, err
	}

	autoDeleteOnIdle, err := utils.ISO8601StringToDuration(desc.AutoDeleteOnIdle)

	if err != nil {
		return nil, err
	}

	queuePropsResult := &QueueProperties{
		Name:                                name,
		LockDuration:                        lockDuration,
		MaxSizeInMegabytes:                  desc.MaxSizeInMegabytes,
		RequiresDuplicateDetection:          desc.RequiresDuplicateDetection,
		RequiresSession:                     desc.RequiresSession,
		DefaultMessageTimeToLive:            defaultMessageTimeToLive,
		DeadLetteringOnMessageExpiration:    desc.DeadLetteringOnMessageExpiration,
		DuplicateDetectionHistoryTimeWindow: duplicateDetectionHistoryTimeWindow,
		MaxDeliveryCount:                    desc.MaxDeliveryCount,
		EnableBatchedOperations:             desc.EnableBatchedOperations,
		Status:                              (*EntityStatus)(desc.Status),
		AutoDeleteOnIdle:                    autoDeleteOnIdle,
		EnablePartitioning:                  desc.EnablePartitioning,
		ForwardTo:                           desc.ForwardTo,
		ForwardDeadLetteredMessagesTo:       desc.ForwardDeadLetteredMessagesTo,
	}

	return queuePropsResult, nil
}

func newQueueRuntimeProperties(name string, desc *atom.QueueDescription) *QueueRuntimeProperties {
	return &QueueRuntimeProperties{
		Name:                           name,
		SizeInBytes:                    int64OrZero(desc.SizeInBytes),
		CreatedAt:                      dateTimeToTime(desc.CreatedAt),
		UpdatedAt:                      dateTimeToTime(desc.UpdatedAt),
		AccessedAt:                     dateTimeToTime(desc.AccessedAt),
		TotalMessageCount:              int64OrZero(desc.MessageCount),
		ActiveMessageCount:             int32OrZero(desc.CountDetails.ActiveMessageCount),
		DeadLetterMessageCount:         int32OrZero(desc.CountDetails.DeadLetterMessageCount),
		ScheduledMessageCount:          int32OrZero(desc.CountDetails.ScheduledMessageCount),
		TransferDeadLetterMessageCount: int32OrZero(desc.CountDetails.TransferDeadLetterMessageCount),
		TransferMessageCount:           int32OrZero(desc.CountDetails.TransferMessageCount),
	}
}

func dateTimeToTime(t *date.Time) time.Time {
	if t == nil {
		return time.Time{}
	}

	return t.Time
}

func int32OrZero(i *int32) int32 {
	if i == nil {
		return 0
	}

	return *i
}

func int64OrZero(i *int64) int64 {
	if i == nil {
		return 0
	}

	return *i
}
