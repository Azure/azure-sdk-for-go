// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build race
// +build race

package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func (suite *eventHubSuite) TestConcurrency() {
	tests := map[string]func(context.Context, *testing.T, *Hub, string){
		"TestConcurrentSendWithRecover": testConcurrentSendWithRecover,
	}

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			hub, cleanup := suite.RandomHub()
			defer cleanup()
			partitionID := (*hub.PartitionIds)[0]
			client, closer := suite.newClient(t, *hub.Name, HubWithPartitionedSender(partitionID))
			defer closer()
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			testFunc(ctx, t, client, partitionID)
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func testConcurrentSendWithRecover(ctx context.Context, t *testing.T, client *Hub, _ string) {
	var wg sync.WaitGroup
	var err error
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// we don't check for errors here as the call to Recover()
			// can cancel any in-flight calls to Send().  this is only
			// an interesting test when race detection is enabled.
			client.Send(ctx, NewEventFromString("Hello!"))
			if inner := client.sender.Recover(ctx); inner != nil {
				err = inner
			}
		}()
	}
	end, _ := ctx.Deadline()
	waitUntil(t, &wg, time.Until(end))
	assert.NoError(t, err)
}
