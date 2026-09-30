// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcognitiveservices_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices/v4"
	"github.com/stretchr/testify/require"
)

const (
	testSubscriptionID  = "00000000-0000-0000-0000-000000000000"
	testResourceGroup   = "test-resource-group"
	testAccountName     = "test-account"
	testCostControlName = "test-cost-control"
)

type fakeCredential struct{}

func (fakeCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "unit-test-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type capturedRequest struct {
	method  string
	url     string
	headers http.Header
	body    []byte
}

type responseHandler func(*http.Request) (*http.Response, error)

type queueTransport struct {
	handlers []responseHandler
	requests []capturedRequest
}

func (transport *queueTransport) Do(request *http.Request) (*http.Response, error) {
	var body []byte
	if request.Body != nil {
		var err error
		body, err = io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
	}
	transport.requests = append(transport.requests, capturedRequest{
		method:  request.Method,
		url:     request.URL.String(),
		headers: request.Header.Clone(),
		body:    body,
	})
	if len(transport.handlers) == 0 {
		return nil, fmt.Errorf("unexpected request: %s %s", request.Method, request.URL)
	}
	handler := transport.handlers[0]
	transport.handlers = transport.handlers[1:]
	return handler(request)
}

func TestCostControlsWireCreate(t *testing.T) {
	created := costControlResponse("Cost Control", 25, `"etag-1"`)
	transport := &queueTransport{handlers: []responseHandler{staticResponse(http.StatusCreated, created, `"etag-1"`)}}
	client := newWireClient(t, transport)

	response, err := client.CreateOrUpdate(context.Background(), testResourceGroup, testAccountName, testCostControlName,
		newCostControl("Cost Control", 25),
		&armcognitiveservices.CostControlsClientCreateOrUpdateOptions{IfNoneMatch: ptr("*")})
	require.NoError(t, err)
	require.Equal(t, `"etag-1"`, *response.Etag)
	require.Equal(t, "etag-1", *response.CostControl.Etag)

	request := transport.requests[0]
	require.Equal(t, http.MethodPut, request.method)
	require.Contains(t, request.url, "/costControls/"+testCostControlName)
	require.NotEmpty(t, queryValue(t, request.url, "api-version"))
	require.Equal(t, "*", request.headers.Get("If-None-Match"))

	var body map[string]any
	require.NoError(t, json.Unmarshal(request.body, &body))
	properties := body["properties"].(map[string]any)
	rule := properties["rules"].([]any)[0].(map[string]any)
	threshold := rule["thresholds"].([]any)[0].(map[string]any)
	require.Equal(t, "account", rule["counterKey"].([]any)[0].(map[string]any)["type"])
	require.Equal(t, "usd", rule["unit"])
	require.Equal(t, float64(25), rule["amount"])
	require.Equal(t, "audit", threshold["action"])
	require.Equal(t, float64(0), threshold["value"])
}

func TestCostControlsWireGet(t *testing.T) {
	transport := &queueTransport{handlers: []responseHandler{
		staticResponse(http.StatusOK, costControlResponse("Cost Control", 25, `"etag-1"`), `"etag-1"`),
	}}
	client := newWireClient(t, transport)

	response, err := client.Get(context.Background(), testResourceGroup, testAccountName, testCostControlName, nil)
	require.NoError(t, err)
	require.Equal(t, armcognitiveservices.CostControlThresholdActionAudit,
		*response.Properties.Rules[0].Thresholds[0].Action)
	require.Equal(t, float64(25), *response.Properties.Rules[0].Amount)
	require.Equal(t, http.MethodGet, transport.requests[0].method)
	require.NotEmpty(t, queryValue(t, transport.requests[0].url, "api-version"))
}

func TestCostControlsWireListPaging(t *testing.T) {
	created := costControlResponse("Cost Control", 25, `"etag-1"`)
	updated := costControlResponse("Cost Control updated", 50, `"etag-2"`)
	transport := &queueTransport{handlers: []responseHandler{
		func(request *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[%s],"nextLink":%q}`, created, request.URL.String()+"&$skiptoken=next")
			return response(request, http.StatusOK, body, ""), nil
		},
		staticResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, updated), ""),
	}}
	client := newWireClient(t, transport)

	pager := client.NewListPager(testResourceGroup, testAccountName, nil)
	var displayNames []string
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		for _, item := range page.Value {
			displayNames = append(displayNames, *item.Properties.DisplayName)
		}
	}
	require.Equal(t, []string{"Cost Control", "Cost Control updated"}, displayNames)
	firstURL, err := url.Parse(transport.requests[0].url)
	require.NoError(t, err)
	nextURL, err := url.Parse(transport.requests[1].url)
	require.NoError(t, err)
	require.Equal(t, firstURL.Path, nextURL.Path)
	require.Equal(t, "next", nextURL.Query().Get("$skiptoken"))
	require.Equal(t, firstURL.Query().Get("api-version"), nextURL.Query().Get("api-version"))
}

func TestCostControlsWireUpdate(t *testing.T) {
	transport := &queueTransport{handlers: []responseHandler{
		staticResponse(http.StatusOK, costControlResponse("Cost Control updated", 50, `"etag-2"`), `"etag-2"`),
	}}
	client := newWireClient(t, transport)

	response, err := client.Update(context.Background(), testResourceGroup, testAccountName, testCostControlName,
		armcognitiveservices.CostControlPatch{Properties: &armcognitiveservices.CostControlPatchProperties{
			DisplayName: ptr("Cost Control updated"),
			Rules:       []*armcognitiveservices.CostControlRule{newRule(50)},
		}}, &armcognitiveservices.CostControlsClientUpdateOptions{IfMatch: ptr(`"etag-1"`)})
	require.NoError(t, err)
	require.Equal(t, float64(50), *response.Properties.Rules[0].Amount)

	request := transport.requests[0]
	require.Equal(t, http.MethodPatch, request.method)
	require.Equal(t, `"etag-1"`, request.headers.Get("If-Match"))
	var body map[string]any
	require.NoError(t, json.Unmarshal(request.body, &body))
	rule := body["properties"].(map[string]any)["rules"].([]any)[0].(map[string]any)
	require.Equal(t, float64(50), rule["amount"])
}

func TestCostControlsWireDeleteAndNotFound(t *testing.T) {
	transport := &queueTransport{handlers: []responseHandler{
		staticResponse(http.StatusNoContent, "", ""),
		staticResponse(http.StatusNotFound,
			`{"error":{"code":"ResourceNotFound","message":"Cost Control not found."}}`, ""),
	}}
	client := newWireClient(t, transport)

	_, err := client.Delete(context.Background(), testResourceGroup, testAccountName, testCostControlName,
		&armcognitiveservices.CostControlsClientDeleteOptions{IfMatch: ptr(`"etag-2"`)})
	require.NoError(t, err)
	require.Equal(t, http.MethodDelete, transport.requests[0].method)
	require.Equal(t, `"etag-2"`, transport.requests[0].headers.Get("If-Match"))

	_, err = client.Get(context.Background(), testResourceGroup, testAccountName, testCostControlName, nil)
	var responseError *azcore.ResponseError
	require.ErrorAs(t, err, &responseError)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)
	require.Equal(t, http.MethodGet, transport.requests[1].method)
}

func newWireClient(t *testing.T, transport policy.Transporter) *armcognitiveservices.CostControlsClient {
	client, err := armcognitiveservices.NewCostControlsClient(testSubscriptionID, fakeCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)
	return client
}

func staticResponse(status int, body, etag string) responseHandler {
	return func(request *http.Request) (*http.Response, error) {
		return response(request, status, body, etag), nil
	}
}

func response(request *http.Request, status int, body, etag string) *http.Response {
	headers := http.Header{"Content-Type": []string{"application/json"}}
	if etag != "" {
		headers.Set("Etag", etag)
	}
	return &http.Response{
		StatusCode: status,
		Header:     headers,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Request:    request,
	}
}

func costControlResponse(displayName string, amount float64, etag string) string {
	return fmt.Sprintf(`{"id":"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CognitiveServices/accounts/%s/costControls/%s","name":"%s","type":"Microsoft.CognitiveServices/accounts/costControls","etag":%s,"properties":{"displayName":%q,"rules":[{"name":"account-budget","counterKey":[{"type":"account"}],"unit":"usd","amount":%g,"period":"month","recurring":true,"thresholds":[{"type":"percentage","value":0,"action":"audit"}]}]}}`,
		testSubscriptionID, testResourceGroup, testAccountName, testCostControlName, testCostControlName, etag, displayName, amount)
}

func newCostControl(displayName string, amount float64) armcognitiveservices.CostControl {
	return armcognitiveservices.CostControl{Properties: &armcognitiveservices.CostControlProperties{
		DisplayName: ptr(displayName),
		Rules:       []*armcognitiveservices.CostControlRule{newRule(amount)},
	}}
}

func newRule(amount float64) *armcognitiveservices.CostControlRule {
	return &armcognitiveservices.CostControlRule{
		Amount: ptr(amount),
		CounterKey: []*armcognitiveservices.CostControlDimension{{
			Type: ptr(armcognitiveservices.CostControlDimensionTypeAccount),
		}},
		Name:      ptr("account-budget"),
		Unit:      ptr(armcognitiveservices.CostControlUnitUsd),
		Period:    ptr(armcognitiveservices.CostControlPeriodMonth),
		Recurring: ptr(true),
		Thresholds: []*armcognitiveservices.CostControlThreshold{{
			Type:   ptr(armcognitiveservices.CostControlThresholdTypePercentage),
			Value:  ptr(0.0),
			Action: ptr(armcognitiveservices.CostControlThresholdActionAudit),
		}},
	}
}

func queryValue(t *testing.T, rawURL, key string) string {
	request, err := http.NewRequest(http.MethodGet, rawURL, nil)
	require.NoError(t, err)
	return request.URL.Query().Get(key)
}

func ptr[T any](value T) *T {
	return &value
}
