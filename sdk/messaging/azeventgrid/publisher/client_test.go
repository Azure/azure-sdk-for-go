//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package publisher_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid/publisher"
	"github.com/stretchr/testify/require"
)

// TestPublishEvent publishes an event using the EventGrid format.
func TestPublishEvent(t *testing.T) {
	skipIntegration(t)

	testPublish := func(t *testing.T, client *publisher.Client) {
		_, err := client.PublishEvents(context.Background(), []publisher.Event{
			{
				Data: map[string]string{
					"hello": "world",
				},
				Subject:     to.Ptr("subjectA"),
				EventType:   to.Ptr("eventType"),
				ID:          to.Ptr("id"),
				EventTime:   to.Ptr(time.Now()),
				DataVersion: to.Ptr("1.0"),
			},
		}, nil)
		require.NoError(t, err)
	}

	t.Run("sas", func(t *testing.T) {
		sas := generateSAS(testVars.EG.Endpoint, testVars.EG.Key)
		client, err := publisher.NewClientWithSAS(testVars.EG.Endpoint, azcore.NewSASCredential(sas), newClientOptionsForTest(t, testVars.EG))
		require.NoError(t, err)
		testPublish(t, client)
	})

	t.Run("sharedkey", func(t *testing.T) {
		client, err := publisher.NewClientWithSharedKeyCredential(testVars.EG.Endpoint, azcore.NewKeyCredential(testVars.EG.Key), newClientOptionsForTest(t, testVars.EG))
		require.NoError(t, err)
		testPublish(t, client)
	})

	t.Run("tokencredential", func(t *testing.T) {
		// note you need the "Event Grid sender" role.
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		require.NoError(t, err)

		client, err := publisher.NewClient(testVars.EG.Endpoint, cred, newClientOptionsForTest(t, testVars.EG))
		require.NoError(t, err)
		testPublish(t, client)
	})
}

// TestPublishCloudEvent publishes an event using the CloudEvent format.
func TestPublishCloudEvent(t *testing.T) {
	skipIntegration(t)

	testPublish := func(t *testing.T, client *publisher.Client) {
		ce, err := messaging.NewCloudEvent("source", "eventType", map[string]string{
			"hello": "world",
		}, nil)
		require.NoError(t, err)

		_, err = client.PublishCloudEvents(context.Background(), []messaging.CloudEvent{ce}, nil)
		require.NoError(t, err)
	}

	t.Run("sas", func(t *testing.T) {
		sas := generateSAS(testVars.CE.Endpoint, testVars.CE.Key)
		client, err := publisher.NewClientWithSAS(testVars.CE.Endpoint, azcore.NewSASCredential(sas), newClientOptionsForTest(t, testVars.CE))
		require.NoError(t, err)
		testPublish(t, client)
	})

	t.Run("sharedkey", func(t *testing.T) {
		client, err := publisher.NewClientWithSharedKeyCredential(testVars.CE.Endpoint, azcore.NewKeyCredential(testVars.CE.Key), newClientOptionsForTest(t, testVars.CE))
		require.NoError(t, err)
		testPublish(t, client)
	})

	t.Run("tokencredential", func(t *testing.T) {
		tokenCred, err := azidentity.NewDefaultAzureCredential(nil)
		require.NoError(t, err)

		client, err := publisher.NewClient(testVars.CE.Endpoint, tokenCred, newClientOptionsForTest(t, testVars.CE))
		require.NoError(t, err)
		testPublish(t, client)
	})
}

func generateSAS(endpoint string, sharedKey string) string {
	ttl := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	text := fmt.Sprintf("r=%s&e=%s", url.QueryEscape(endpoint), url.QueryEscape(ttl))

	decodedKey, err := base64.StdEncoding.DecodeString(sharedKey)

	if err != nil {
		panic(err)
	}

	h := hmac.New(sha256.New, []byte(decodedKey))
	_, err = h.Write([]byte(text))

	if err != nil {
		panic(err)
	}

	b64Sig := base64.StdEncoding.EncodeToString(h.Sum(nil))
	sig := url.QueryEscape(b64Sig)

	sas := fmt.Sprintf("%s&s=%s", text, sig)
	return sas
}
