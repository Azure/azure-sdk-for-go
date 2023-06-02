//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/template/aztemplate/internal/tests"
	"github.com/stretchr/testify/require"
)

func TestPublishingAndReceivingCloudEvents(t *testing.T) {
	env := tests.LoadEnv()

	log.SetListener(func(e log.Event, s string) {
		fmt.Printf("[%s]: %s\n", e, s)
	})

	c, err := NewClientFromSharedKey(env.Endpoint, env.Key, nil)
	require.NoError(t, err)
	require.NotNil(t, c)

	topicName := env.Topic
	subscriptionName := env.Subscription

	_, err = c.PublishCloudEvents(context.Background(), topicName, []*CloudEvent{
		{
			Data:   []byte("hello World"),
			Source: to.Ptr("hello-source"),
			Type:   to.Ptr("world"),
		},
	}, nil)
	require.NoError(t, err)

	resp, err := c.ReceiveCloudEvents(context.Background(), topicName, subscriptionName, nil)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Value)

	ackArgs := AcknowledgeOptions{}

	for _, e := range resp.Value {
		require.NotNil(t, e.BrokerProperties.LockToken)
		ackArgs.LockTokens = append(ackArgs.LockTokens, e.BrokerProperties.LockToken)
	}

	// TODO: it's weird that these are marked optional.
	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), topicName, subscriptionName, ackArgs, nil)
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
