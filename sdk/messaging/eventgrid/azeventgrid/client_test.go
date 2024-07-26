//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azeventgrid"
	"github.com/stretchr/testify/require"
)

// TestPublishEvent publishes an event using the EventGrid format.
func TestPublishEvent(t *testing.T) {
	testPublish := func(t *testing.T, client *azeventgrid.Client) {
		_, err := client.PublishEvents(context.Background(), []azeventgrid.Event{
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
		vars := newTestVars(t)
		client, err := azeventgrid.NewClientWithSAS(vars.EG.Endpoint, azcore.NewSASCredential(vars.EG.SAS), newClientOptionsForTest(t).EG)
		require.NoError(t, err)
		testPublish(t, client)
	})

	t.Run("sharedkey", func(t *testing.T) {
		vars := newTestVars(t)
		client, err := azeventgrid.NewClientWithSharedKeyCredential(vars.EG.Endpoint, azcore.NewKeyCredential(vars.EG.Key), newClientOptionsForTest(t).EG)
		require.NoError(t, err)
		testPublish(t, client)
	})

	t.Run("tokencredential", func(t *testing.T) {
		vars := newTestVars(t)

		// note you need the "Event Grid sender" role.
		cred, err := credential.New(nil)
		require.NoError(t, err)

		client, err := azeventgrid.NewClient(vars.EG.Endpoint, cred, newClientOptionsForTest(t).EG)
		require.NoError(t, err)
		testPublish(t, client)
	})
}

// TestPublishCloudEvent publishes an event using the CloudEvent format.
func TestPublishCloudEvent(t *testing.T) {
	testPublish := func(t *testing.T, client *azeventgrid.Client) {
		ce, err := messaging.NewCloudEvent("source", "eventType", map[string]string{
			"hello": "world",
		}, nil)
		require.NoError(t, err)

		_, err = client.PublishCloudEvents(context.Background(), []messaging.CloudEvent{ce}, nil)
		require.NoError(t, err)
	}

	t.Run("sas", func(t *testing.T) {
		vars := newTestVars(t)

		client, err := azeventgrid.NewClientWithSAS(vars.CE.Endpoint, azcore.NewSASCredential(vars.CE.SAS), newClientOptionsForTest(t).EG)
		require.NoError(t, err)
		testPublish(t, client)
	})

	t.Run("sharedkey", func(t *testing.T) {
		vars := newTestVars(t)

		client, err := azeventgrid.NewClientWithSharedKeyCredential(vars.CE.Endpoint, azcore.NewKeyCredential(vars.CE.Key), newClientOptionsForTest(t).EG)
		require.NoError(t, err)
		testPublish(t, client)
	})

	t.Run("tokencredential", func(t *testing.T) {
		vars := newTestVars(t)

		tokenCred, err := credential.New(nil)
		require.NoError(t, err)

		client, err := azeventgrid.NewClient(vars.CE.Endpoint, tokenCred, newClientOptionsForTest(t).EG)
		require.NoError(t, err)
		testPublish(t, client)
	})
}

func generateSAS(endpoint string, sharedKey string, baseTime time.Time) string {
	ttl := baseTime.UTC().Add(time.Hour).Format(time.RFC3339)
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
