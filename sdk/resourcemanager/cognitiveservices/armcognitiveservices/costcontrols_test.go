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
	testProjectName     = "test-project"
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
	created := costControlResponse(25, `"etag-1"`)
	transport := &queueTransport{handlers: []responseHandler{staticResponse(http.StatusCreated, created, `"etag-1"`)}}
	client := newWireClient(t, transport)

	response, err := client.CreateOrUpdate(context.Background(), testResourceGroup, testAccountName, testCostControlName,
		newCostControl(25),
		&armcognitiveservices.AccountCostControlsClientCreateOrUpdateOptions{IfNoneMatch: ptr("*")})
	require.NoError(t, err)
	require.Equal(t, azcore.ETag(`"etag-1"`), *response.Etag)
	require.Equal(t, azcore.ETag(`"etag-1"`), *response.AccountCostControl.Etag)

	request := transport.requests[0]
	require.Equal(t, http.MethodPut, request.method)
	require.Contains(t, request.url, "/costControls/"+testCostControlName)
	require.Equal(t, "2026-11-15-preview", queryValue(t, request.url, "api-version"))
	require.Equal(t, "*", request.headers.Get("If-None-Match"))

	var body map[string]any
	require.NoError(t, json.Unmarshal(request.body, &body))
	properties := body["properties"].(map[string]any)
	require.Equal(t, "Cost", properties["category"])
	require.Equal(t, "Monthly", properties["timeGrain"])
	require.Equal(t, float64(25), properties["amount"])
	require.Equal(t, "2026-11-01T00:00:00Z", properties["timePeriod"].(map[string]any)["startDate"])
	notification := properties["notifications"].(map[string]any)["usageAlert"].(map[string]any)
	require.Equal(t, "Percentage", notification["thresholdType"])
	require.Equal(t, float64(80), notification["threshold"])
	require.Equal(t, true, notification["enabled"])
	require.Equal(t, "GreaterThanOrEqualTo", notification["operator"])
	require.Contains(t, notification["eventGridConnectionId"], "/connections/budget-alerts")
	action := properties["providerProperties"].(map[string]any)["actions"].(map[string]any)["stopRequests"].(map[string]any)
	require.Equal(t, "Block", action["type"])
	require.Equal(t, "Percentage", action["thresholdType"])
	require.Equal(t, float64(100), action["threshold"])
}

func TestCostControlsWireGet(t *testing.T) {
	transport := &queueTransport{handlers: []responseHandler{
		staticResponse(http.StatusOK, costControlResponse(25, `"etag-1"`), `"etag-1"`),
	}}
	client := newWireClient(t, transport)

	response, err := client.Get(context.Background(), testResourceGroup, testAccountName, testCostControlName, nil)
	require.NoError(t, err)
	require.Equal(t, armcognitiveservices.CostControlCategoryCost, *response.Properties.Category)
	require.Equal(t, float64(25), *response.Properties.Amount)
	require.Equal(t, armcognitiveservices.CostControlTimeGrainMonthly, *response.Properties.TimeGrain)
	require.Equal(t, http.MethodGet, transport.requests[0].method)
	require.Equal(t, "2026-11-15-preview", queryValue(t, transport.requests[0].url, "api-version"))
}

func TestCostControlsWireListPaging(t *testing.T) {
	created := costControlResponse(25, `"etag-1"`)
	updated := costControlResponse(50, `"etag-2"`)
	transport := &queueTransport{handlers: []responseHandler{
		func(request *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[%s],"nextLink":%q}`, created, request.URL.String()+"&$skiptoken=next")
			return response(request, http.StatusOK, body, ""), nil
		},
		staticResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, updated), ""),
	}}
	client := newWireClient(t, transport)

	pager := client.NewListPager(testResourceGroup, testAccountName, nil)
	var amounts []float64
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		for _, item := range page.Value {
			amounts = append(amounts, *item.Properties.Amount)
		}
	}
	require.Equal(t, []float64{25, 50}, amounts)
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
		staticResponse(http.StatusOK, costControlResponse(50, `"etag-2"`), `"etag-2"`),
	}}
	client := newWireClient(t, transport)

	response, err := client.Update(context.Background(), testResourceGroup, testAccountName, testCostControlName,
		armcognitiveservices.CostControlPatch{Properties: &armcognitiveservices.CostControlPatchProperties{
			Amount: ptr(50.0),
		}}, &armcognitiveservices.AccountCostControlsClientUpdateOptions{IfMatch: ptr(`"etag-1"`)})
	require.NoError(t, err)
	require.Equal(t, float64(50), *response.Properties.Amount)

	request := transport.requests[0]
	require.Equal(t, http.MethodPatch, request.method)
	require.Equal(t, `"etag-1"`, request.headers.Get("If-Match"))
	var body map[string]any
	require.NoError(t, json.Unmarshal(request.body, &body))
	require.Equal(t, map[string]any{"properties": map[string]any{"amount": float64(50)}}, body)
}

func TestCostControlsWireDeleteAndNotFound(t *testing.T) {
	transport := &queueTransport{handlers: []responseHandler{
		staticResponse(http.StatusNoContent, "", ""),
		staticResponse(http.StatusNotFound,
			`{"error":{"code":"ResourceNotFound","message":"Cost Control not found."}}`, ""),
	}}
	client := newWireClient(t, transport)

	_, err := client.Delete(context.Background(), testResourceGroup, testAccountName, testCostControlName,
		&armcognitiveservices.AccountCostControlsClientDeleteOptions{IfMatch: ptr(`"etag-2"`)})
	require.NoError(t, err)
	require.Equal(t, http.MethodDelete, transport.requests[0].method)
	require.Equal(t, `"etag-2"`, transport.requests[0].headers.Get("If-Match"))

	_, err = client.Get(context.Background(), testResourceGroup, testAccountName, testCostControlName, nil)
	var responseError *azcore.ResponseError
	require.ErrorAs(t, err, &responseError)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)
	require.Equal(t, http.MethodGet, transport.requests[1].method)
}

func TestCostControlsWireCheckExistence(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusNotFound} {
		for _, project := range []bool{false, true} {
			t.Run(fmt.Sprintf("status=%d/project=%t", status, project), func(t *testing.T) {
				transport := &queueTransport{handlers: []responseHandler{staticResponse(status, "", `"etag-1"`)}}
				var success bool
				var etag *azcore.ETag
				var err error
				if project {
					client := newProjectWireClient(t, transport)
					result, callErr := client.CheckExistence(context.Background(), testResourceGroup, testAccountName, testProjectName, testCostControlName, nil)
					success, etag, err = result.Success, result.Etag, callErr
				} else {
					client := newWireClient(t, transport)
					result, callErr := client.CheckExistence(context.Background(), testResourceGroup, testAccountName, testCostControlName, nil)
					success, etag, err = result.Success, result.Etag, callErr
				}
				if status == http.StatusNotFound {
					var responseError *azcore.ResponseError
					require.ErrorAs(t, err, &responseError)
					require.Equal(t, status, responseError.StatusCode)
					require.False(t, success)
				} else {
					require.NoError(t, err)
					require.True(t, success)
					require.NotNil(t, etag)
					require.Equal(t, azcore.ETag(`"etag-1"`), *etag)
				}
				require.Len(t, transport.requests, 1)
				request := transport.requests[0]
				require.Equal(t, http.MethodHead, request.method)
				ownerPath := "/accounts/" + testAccountName
				if project {
					ownerPath += "/projects/" + testProjectName
				}
				require.Contains(t, request.url, ownerPath+"/costControls/"+testCostControlName)
				require.Equal(t, "2026-11-15-preview", queryValue(t, request.url, "api-version"))
				require.Empty(t, request.body)
			})
		}
	}
}

func TestCostControlsWirePatchCollections(t *testing.T) {
	transport := &queueTransport{handlers: []responseHandler{
		staticResponse(http.StatusOK, costControlResponse(25, `"etag-2"`), `"etag-2"`),
	}}
	client := newWireClient(t, transport)
	_, err := client.Update(context.Background(), testResourceGroup, testAccountName, testCostControlName,
		armcognitiveservices.CostControlPatch{Properties: &armcognitiveservices.CostControlPatchProperties{
			Grouping: []*armcognitiveservices.CostControlGrouping{},
			Notifications: map[string]*armcognitiveservices.CostControlNotificationPatch{
				"usageAlert": {Enabled: ptr(false)},
			},
			ProviderProperties: &armcognitiveservices.CostControlProviderPropertiesPatch{
				Actions: map[string]*armcognitiveservices.CostControlActionPatch{"stopRequests": nil},
			},
		}}, nil)
	require.NoError(t, err)
	require.JSONEq(t, `{"properties":{"grouping":[],"notifications":{"usageAlert":{"enabled":false}},"providerProperties":{"actions":{"stopRequests":null}}}}`, string(transport.requests[0].body))
}

func TestCostControlsWireProjectLifecycle(t *testing.T) {
	created := scopedCostControlResponse(testProjectName, 1000000, `"etag-1"`)
	updated := scopedCostControlResponse(testProjectName, 2000000, `"etag-2"`)
	transport := &queueTransport{handlers: []responseHandler{
		staticResponse(http.StatusCreated, created, `"etag-1"`),
		staticResponse(http.StatusOK, created, `"etag-1"`),
		func(request *http.Request) (*http.Response, error) {
			return response(request, http.StatusOK, fmt.Sprintf(`{"value":[%s],"nextLink":%q}`, created, request.URL.String()+"&$skiptoken=next"), ""), nil
		},
		staticResponse(http.StatusOK, `{"value":[]}`, ""),
		staticResponse(http.StatusOK, updated, `"etag-2"`),
		staticResponse(http.StatusNoContent, "", ""),
		staticResponse(http.StatusNotFound, `{"error":{"code":"ResourceNotFound","message":"Cost control not found."}}`, ""),
	}}
	client := newProjectWireClient(t, transport)
	ctx := context.Background()
	properties := newCostControl(1000000).Properties
	properties.Category = ptr(armcognitiveservices.CostControlCategoryTokens)
	properties.TimeGrain = ptr(armcognitiveservices.CostControlTimeGrainDaily)
	properties.Filter = &armcognitiveservices.CostControlFilter{
		AdditionalInfo: &armcognitiveservices.CostControlAdditionalInfoFilter{
			Name:   ptr(armcognitiveservices.CostControlFilterNameTokenType),
			Values: []*string{ptr("input"), ptr("output")},
		},
	}
	properties.Grouping = []*armcognitiveservices.CostControlGrouping{{Name: ptr(armcognitiveservices.CostControlGroupingNameSessionID)}}
	createResult, err := client.CreateOrUpdate(ctx, testResourceGroup, testAccountName, testProjectName, testCostControlName,
		armcognitiveservices.ProjectCostControl{Properties: properties},
		&armcognitiveservices.ProjectCostControlsClientCreateOrUpdateOptions{IfNoneMatch: ptr("*")})
	require.NoError(t, err)
	require.Equal(t, azcore.ETag(`"etag-1"`), *createResult.Etag)
	require.Equal(t, azcore.ETag(`"etag-1"`), *createResult.ProjectCostControl.Etag)
	require.Equal(t, "*", transport.requests[0].headers.Get("If-None-Match"))
	var body map[string]any
	require.NoError(t, json.Unmarshal(transport.requests[0].body, &body))
	wireProperties := body["properties"].(map[string]any)
	require.Equal(t, "Tokens", wireProperties["category"])
	require.Equal(t, "Daily", wireProperties["timeGrain"])
	require.Equal(t, float64(1000000), wireProperties["amount"])
	require.Equal(t, map[string]any{"additionalInfo": map[string]any{
		"name": "tokenType", "operator": "In", "values": []any{"input", "output"},
	}}, wireProperties["filter"])
	require.Equal(t, []any{map[string]any{"name": "sessionId", "type": "AdditionalInfoKey"}}, wireProperties["grouping"])

	fetched, err := client.Get(ctx, testResourceGroup, testAccountName, testProjectName, testCostControlName, nil)
	require.NoError(t, err)
	require.Equal(t, armcognitiveservices.CostControlCategoryTokens, *fetched.Properties.Category)
	require.Equal(t, float64(1000000), *fetched.Properties.Amount)
	pager := client.NewListPager(testResourceGroup, testAccountName, testProjectName, nil)
	var amounts []float64
	for pager.More() {
		page, pageErr := pager.NextPage(ctx)
		require.NoError(t, pageErr)
		for _, item := range page.Value {
			amounts = append(amounts, *item.Properties.Amount)
		}
	}
	require.Equal(t, []float64{1000000}, amounts)
	require.Equal(t, "next", queryValue(t, transport.requests[3].url, "$skiptoken"))
	updatedResult, err := client.Update(ctx, testResourceGroup, testAccountName, testProjectName, testCostControlName,
		armcognitiveservices.CostControlPatch{Properties: &armcognitiveservices.CostControlPatchProperties{Amount: ptr(2000000.0)}},
		&armcognitiveservices.ProjectCostControlsClientUpdateOptions{IfMatch: ptr(string(*fetched.Etag))})
	require.NoError(t, err)
	require.Equal(t, float64(2000000), *updatedResult.Properties.Amount)
	require.Equal(t, `"etag-1"`, transport.requests[4].headers.Get("If-Match"))
	require.JSONEq(t, `{"properties":{"amount":2000000}}`, string(transport.requests[4].body))
	_, err = client.Delete(ctx, testResourceGroup, testAccountName, testProjectName, testCostControlName,
		&armcognitiveservices.ProjectCostControlsClientDeleteOptions{IfMatch: ptr(string(*updatedResult.Etag))})
	require.NoError(t, err)
	require.Equal(t, `"etag-2"`, transport.requests[5].headers.Get("If-Match"))
	_, err = client.Get(ctx, testResourceGroup, testAccountName, testProjectName, testCostControlName, nil)
	var responseError *azcore.ResponseError
	require.ErrorAs(t, err, &responseError)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)
	require.Empty(t, transport.handlers)
	require.Len(t, transport.requests, 7)
	for index, method := range []string{http.MethodPut, http.MethodGet, http.MethodGet, http.MethodGet, http.MethodPatch, http.MethodDelete, http.MethodGet} {
		request := transport.requests[index]
		require.Equal(t, method, request.method)
		require.Contains(t, request.url, "/accounts/"+testAccountName+"/projects/"+testProjectName+"/costControls")
		require.Equal(t, "2026-11-15-preview", queryValue(t, request.url, "api-version"))
	}
}

func newProjectWireClient(t *testing.T, transport policy.Transporter) *armcognitiveservices.ProjectCostControlsClient {
	t.Helper()
	client, err := armcognitiveservices.NewProjectCostControlsClient(testSubscriptionID, fakeCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)
	return client
}

func newWireClient(t *testing.T, transport policy.Transporter) *armcognitiveservices.AccountCostControlsClient {
	t.Helper()
	client, err := armcognitiveservices.NewAccountCostControlsClient(testSubscriptionID, fakeCredential{}, &arm.ClientOptions{
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

func costControlResponse(amount float64, etag string) string {
	return scopedCostControlResponse("", amount, etag)
}

func scopedCostControlResponse(projectName string, amount float64, etag string) string {
	ownerID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CognitiveServices/accounts/%s", testSubscriptionID, testResourceGroup, testAccountName)
	resourceType, category, timeGrain := "Microsoft.CognitiveServices/accounts/costControls", "Cost", "Monthly"
	if projectName != "" {
		ownerID += "/projects/" + projectName
		resourceType, category, timeGrain = "Microsoft.CognitiveServices/accounts/projects/costControls", "Tokens", "Daily"
	}
	return fmt.Sprintf(`{"id":%q,"name":%q,"type":%q,"etag":%q,"properties":{"provider":"Microsoft.CognitiveServices","category":%q,"amount":%g,"timeGrain":%q,"timePeriod":{"startDate":"2026-11-01T00:00:00Z","endDate":"2027-11-01T00:00:00Z"}}}`,
		ownerID+"/costControls/"+testCostControlName, testCostControlName, resourceType, etag, category, amount, timeGrain)
}

func newCostControl(amount float64) armcognitiveservices.AccountCostControl {
	return armcognitiveservices.AccountCostControl{Properties: &armcognitiveservices.CostControlProperties{
		Amount:    ptr(amount),
		Category:  ptr(armcognitiveservices.CostControlCategoryCost),
		TimeGrain: ptr(armcognitiveservices.CostControlTimeGrainMonthly),
		TimePeriod: &armcognitiveservices.CostControlTimePeriod{
			StartDate: ptr(time.Date(2026, time.November, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:   ptr(time.Date(2027, time.November, 1, 0, 0, 0, 0, time.UTC)),
		},
		Notifications: map[string]*armcognitiveservices.CostControlNotification{
			"usageAlert": {
				ThresholdType:         ptr(armcognitiveservices.CostControlThresholdTypePercentage),
				Threshold:             ptr(80.0),
				Enabled:               ptr(true),
				EventGridConnectionID: ptr(fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CognitiveServices/accounts/%s/connections/budget-alerts", testSubscriptionID, testResourceGroup, testAccountName)),
			},
		},
		ProviderProperties: &armcognitiveservices.CostControlProviderProperties{
			Actions: map[string]*armcognitiveservices.CostControlAction{
				"stopRequests": {
					ThresholdType: ptr(armcognitiveservices.CostControlThresholdTypePercentage),
					Threshold:     ptr(100.0),
				},
			},
		},
	}}
}

func queryValue(t *testing.T, rawURL, key string) string {
	t.Helper()
	request, err := http.NewRequest(http.MethodGet, rawURL, nil)
	require.NoError(t, err)
	return request.URL.Query().Get(key)
}

func ptr[T any](value T) *T {
	return &value
}
