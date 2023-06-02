//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/template/aztemplate/internal/tests"
	"github.com/stretchr/testify/require"
)

func TestPublishingAndReceivingCloudEvents(t *testing.T) {
	env := tests.LoadEnv()

	c, err := NewClientFromSharedKey(env.Key, nil)
	require.NoError(t, err)
	require.NotNil(t, c)

	topicName := env.Topic
	subscriptionName := env.Subscription

	newCTX := runtime.WithHTTPHeader(context.Background(), http.Header{
		"Content-type": []string{"application/cloudevents-batch+json; charset=utf-8"},
	})

	_, err = c.PublishCloudEvents(newCTX, env.Endpoint, topicName, []*CloudEvent{

		{
			Data:   "Hello World",
			Source: to.Ptr("hello-source"),
			Type:   to.Ptr("world"),

			// TODO: ID should be auto-assigned?
			ID:          to.Ptr("hello"),
			Specversion: to.Ptr("1.0")},
	}, nil)
	require.NoError(t, err)

	resp, err := c.ReceiveCloudEvents(context.Background(), env.Endpoint, topicName, subscriptionName, nil)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Value)

	ackArgs := AcknowledgeOptions{}

	for _, e := range resp.Value {
		require.NotNil(t, e.BrokerProperties.LockToken)
		ackArgs.LockTokens = append(ackArgs.LockTokens, e.BrokerProperties.LockToken)
	}

	// TODO: it's weird that these are marked optional.
	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), env.Endpoint, topicName, subscriptionName, ackArgs, nil)
	require.NoError(t, err)

	for _, flt := range ackResp.FailedLockTokens {
		require.Nil(t, flt.ErrorCode)
		// switch *flt.ErrorCode {
		// case "BadToken":
		// case "TokenLost":
		// case "InternalServerError":
		// default:
		// }
	}

	require.Empty(t, ackResp.FailedLockTokens)
	require.NotEmpty(t, ackResp.SucceededLockTokens)
}
