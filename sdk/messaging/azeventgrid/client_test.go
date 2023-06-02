//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestPublishingAndReceivingCloudEvents(t *testing.T) {
	keyLogWriter, err := os.OpenFile(keyFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	require.NoError(t, err)

	defer keyLogWriter.Close()

	key := os.Getenv("EVENTGRID_KEY")
	endpoint := os.Getenv("EVENTGRID_ENDPOINT")

	cp := http.DefaultTransport.(*http.Transport).Clone()
	cp.TLSClientConfig = &tls.Config{
		KeyLogWriter: keyLogWriter,
	}

	myClient := &http.Client{
		Transport: cp,
	}

	p := runtime.NewPipeline("azeventgrid", "0.1", runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			&EventGridSharedKeyPolicy{Key: key},
		},
	}, &policy.ClientOptions{
		Transport: myClient,
	})
	c := azeventgrid.NewClient(p)
	require.NotNil(t, c)

	topicName := "testtopic1"
	subscriptionName := "testsubscription1"

	/*
		{
			"error": {
			"code": "BadRequest",
			"message": "Cannot parse CONTENT-TYPE header",
			"timestamp_utc": "2023-06-01T01:49:53.197646029+00:00",
			"tracking_id": "4F15A8D7-59BB-4C3A-A0D2-197CE85D6260"
			}
		}
	*/

	// azlog.SetListener(func(e azlog.Event, s string) {
	// 	log.Printf("%s: %s", e, s)
	// })

	// this gets us a little further. But it seems like the payload isn't getting serialized properly here.
	newCTX := runtime.WithHTTPHeader(context.Background(), http.Header{
		"Content-type": []string{"application/cloudevents-batch+json; charset=utf-8"},
	})
	// newCTX = context.Background()

	_, err = c.PublishCloudEvents(newCTX, endpoint, topicName, []*azeventgrid.CloudEvent{
		//{Data: "Hello World", Datacontenttype: to.Ptr("text/plain")},
		// TODO: ID should be auto-assigned?
		{
			Data: "Hello World",

			// So these fields need to also be filled in but it looks like other
			// implementations do this by default.
			ID:          to.Ptr("hello"),
			Source:      to.Ptr("hello-source"),
			Type:        to.Ptr("world"),
			Specversion: to.Ptr("1.0")},
	}, nil)
	require.NoError(t, err)

	resp, err := c.ReceiveCloudEvents(context.Background(), endpoint, topicName, subscriptionName, nil)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Value)

	ackArgs := azeventgrid.AcknowledgeOptions{}

	for _, e := range resp.Value {
		if e.BrokerProperties.LockToken == nil {
			panic("You can't add a nil lockToken!")
		}

		ackArgs.LockTokens = append(ackArgs.LockTokens, e.BrokerProperties.LockToken)
	}

	// TODO: it's weird that these are marked optional.
	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), endpoint, topicName, subscriptionName, ackArgs, nil)
	require.NoError(t, err)

	// TODO: this is already a pain - now I have to go through the original array and find
	// all of these matches, which means I have to now make a dictionary out of the strings
	// so I can efficiently determine which ones need to be tried again.
	//
	// If the result was just a simple parallel array with null's then it'd just be index <-> index
	// and no lookup/search would be required.

	// why do we have this one? Is there some middle ground here?
	// TODO: the formatting on the `FailedLockTokens` is messed up - I think it's doing something weird
	// with the newlines or similar.

	for _, flt := range ackResp.FailedLockTokens {
		// error code here is not an enumerated set of values, but it's kind of useless without some
		// idea of what error codes are valid, and which aren't.

		// TODO: also, will ErrorCode always be there? Or do I need to _also_ do a nil check for that?
		switch *flt.ErrorCode {
		case "BadToken":
		case "TokenLost":
		case "InternalServerError":
		default:
			// unknown
		}
	}

	require.Empty(t, ackResp.FailedLockTokens)
	require.NotEmpty(t, ackResp.SucceededLockTokens)
}
