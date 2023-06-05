//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Failed to load .env file, no integration tests will run: %s", err)
	} else {
		fmt.Printf("BEGIN: Purging old events before starting tests...\n")
		// purge the subscription before the test
		c := newClientForTest()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		for {
			resp, err := c.ReceiveCloudEvents(ctx, c.TestVars.Topic, c.TestVars.Subscription, &azeventgrid.ReceiveCloudEventsOptions{
				MaxEvents:   to.Ptr[int32](100),
				MaxWaitTime: to.Ptr[int32](10),
			})

			if errors.Is(err, context.DeadlineExceeded) || len(resp.Value) == 0 {
				break
			}

			if err != nil {
				panic(err)
			}

			var lockTokens []*string

			for _, v := range resp.Value {
				lockTokens = append(lockTokens, v.BrokerProperties.LockToken)
			}

			_, err = c.AcknowledgeCloudEvents(ctx, c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.AcknowledgeOptions{
				LockTokens: lockTokens,
			}, nil)

			if err != nil {
				panic(err)
			}
		}
		fmt.Printf("END: Purging old events before starting tests...\n")
	}

	os.Exit(m.Run())
}
