// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/stretchr/testify/require"
)

type captureSubscriptionPolicy struct {
	body           string
	method         string
	path           string
	responseStatus EntityStatus
}

func (p *captureSubscriptionPolicy) Do(req *policy.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Raw().Body)
	if err != nil {
		return nil, err
	}
	p.body = string(body)
	p.method = req.Raw().Method
	p.path = req.Raw().URL.Path
	response := fmt.Sprintf(`<entry xmlns="http://www.w3.org/2005/Atom"><title>my-subscription</title><content type="application/xml"><SubscriptionDescription xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect"><Status>%s</Status></SubscriptionDescription></content></entry>`, p.responseStatus)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(response)),
		Request:    req.Raw(),
	}, nil
}

func TestClientSubscriptionStatus(t *testing.T) {
	for _, operation := range []string{"create", "update"} {
		for _, tc := range []struct {
			name   string
			status *EntityStatus
		}{
			{name: "disabled", status: to.Ptr(EntityStatusDisabled)},
			{name: "active", status: to.Ptr(EntityStatusActive)},
			{name: "omitted"},
		} {
			t.Run(operation+"/"+tc.name, func(t *testing.T) {
				capture := &captureSubscriptionPolicy{responseStatus: EntityStatusActive}
				if tc.status != nil {
					capture.responseStatus = *tc.status
				}
				client, err := NewClient("fake.servicebus.windows.net", fakeAPIVersionCredential{}, &ClientOptions{
					ClientOptions: azcore.ClientOptions{PerCallPolicies: []policy.Policy{capture}},
				})
				require.NoError(t, err)
				props := SubscriptionProperties{
					Status:                  tc.status,
					EnableBatchedOperations: to.Ptr(true),
					AutoDeleteOnIdle:        to.Ptr("PT10M"),
				}
				var returnedStatus *EntityStatus
				if operation == "create" {
					response, err := client.CreateSubscription(context.Background(), "my-topic", "my-subscription", &CreateSubscriptionOptions{Properties: &props})
					require.NoError(t, err)
					returnedStatus = response.Status
				} else {
					response, err := client.UpdateSubscription(context.Background(), "my-topic", "my-subscription", props, nil)
					require.NoError(t, err)
					returnedStatus = response.Status
				}
				require.Equal(t, http.MethodPut, capture.method)
				require.Equal(t, "/my-topic/Subscriptions/my-subscription", capture.path)
				var envelope atom.SubscriptionEnvelope
				require.NoError(t, xml.Unmarshal([]byte(capture.body), &envelope))
				require.Equal(t, (*atom.EntityStatus)(tc.status), envelope.Content.SubscriptionDescription.Status)
				if tc.status == nil {
					require.NotContains(t, capture.body, "<Status")
				} else {
					require.Contains(t, capture.body, fmt.Sprintf("<EnableBatchedOperations>true</EnableBatchedOperations><Status>%s</Status>", *tc.status))
				}
				require.Equal(t, to.Ptr(capture.responseStatus), returnedStatus)
			})
		}
	}
}
