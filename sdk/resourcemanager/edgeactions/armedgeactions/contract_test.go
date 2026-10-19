// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armedgeactions_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/edgeactions/armedgeactions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	actionPath  = "/subscriptions/" + subscriptionID + "/resourceGroups/test-rg/providers/Microsoft.Cdn/edgeActions/test-action"
	versionPath = actionPath + "/versions/v1"
	filterPath  = actionPath + "/executionFilters/sdkfilter"
	actionJSON  = `{"id":"` + actionPath + `","name":"test-action","location":"global","sku":{"name":"Standard","tier":"Standard"},"properties":{"provisioningState":"Succeeded","attachments":[]}}`
	versionJSON = `{"id":"` + versionPath + `","name":"v1","location":"global","properties":{"deploymentType":"file","isDefaultVersion":"False","provisioningState":"Succeeded"}}`
	filterJSON  = `{"id":"` + filterPath + `","name":"sdkfilter","location":"global","properties":{"versionId":"` + versionPath + `","executionFilterIdentifierHeaderName":"x-sdk-test","executionFilterIdentifierHeaderValue":"initial","provisioningState":"Succeeded"}}`
)

func newMockFactory(t *testing.T, validation ...*responseValidator) (*armedgeactions.ClientFactory, *mock.Server) {
	t.Helper()
	server, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	t.Cleanup(close)
	options := &arm.ClientOptions{
		DisableRPRegistration: true,
		ClientOptions:         policy.ClientOptions{Transport: server, Retry: policy.RetryOptions{MaxRetries: -1}},
	}
	require.LessOrEqual(t, len(validation), 1)
	if len(validation) == 1 {
		options.PerCallPolicies = append(options.PerCallPolicies, validation[0])
		client, err := arm.NewClient("edgeactions-tests", "v0.0.0", &credential.Fake{}, options)
		require.NoError(t, err)
		validation[0].pipeline = client.Pipeline()
	}
	factory, err := armedgeactions.NewClientFactory(subscriptionID, &credential.Fake{}, options)
	require.NoError(t, err)
	return factory, server
}

func finish[T any](poller *runtime.Poller[T], err error) (any, error) {
	if err != nil {
		return nil, err
	}
	return poller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{Frequency: time.Millisecond})
}

func expectRequest(t *testing.T, server *mock.Server, method, path, body string, status int, response string, headers ...mock.ResponseOption) {
	t.Helper()
	server.AppendResponse(mock.WithPredicate(func(request *http.Request) bool {
		assert.Equal(t, method, request.Method)
		assert.Equal(t, path, request.URL.Path)
		assert.Equal(t, "2026-10-01", request.URL.Query().Get("api-version"))
		if body != "" {
			assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
			data, err := io.ReadAll(request.Body)
			if assert.NoError(t, err) {
				assert.JSONEq(t, body, string(data))
			}
		}
		return false
	}))
	options := []mock.ResponseOption{
		mock.WithStatusCode(status), mock.WithBody([]byte(response)),
		mock.WithHeader("Content-Type", "application/json"),
	}
	server.AppendResponse(append(options, headers...)...)
}

func TestContractResourceRequests(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name, method, path, body, response string
		run                                func(*armedgeactions.ClientFactory) (any, error)
	}{
		{"action-get", "GET", actionPath, "", actionJSON, func(f *armedgeactions.ClientFactory) (any, error) {
			return f.NewClient().Get(ctx, "test-rg", "test-action", nil)
		}},
		{"version-get", "GET", versionPath, "", versionJSON, func(f *armedgeactions.ClientFactory) (any, error) {
			return f.NewEdgeActionVersionsClient().Get(ctx, "test-rg", "test-action", "v1", nil)
		}},
		{"filter-get", "GET", filterPath, "", filterJSON, func(f *armedgeactions.ClientFactory) (any, error) {
			return f.NewEdgeActionExecutionFiltersClient().Get(ctx, "test-rg", "test-action", "sdkfilter", nil)
		}},
		{"action-create", "PUT", actionPath, `{"location":"global","sku":{"name":"Standard","tier":"Standard"}}`, actionJSON, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewClient().BeginCreate(ctx, "test-rg", "test-action", armedgeactions.EdgeAction{
				Location: to.Ptr("global"), SKU: &armedgeactions.SKUType{Name: to.Ptr("Standard"), Tier: to.Ptr("Standard")},
			}, nil))
		}},
		{"version-create", "PUT", versionPath, `{"location":"global","properties":{"deploymentType":"file","isDefaultVersion":"False"}}`, versionJSON, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionVersionsClient().BeginCreate(ctx, "test-rg", "test-action", "v1", armedgeactions.EdgeActionVersion{
				Location: to.Ptr("global"), Properties: &armedgeactions.EdgeActionVersionProperties{
					DeploymentType:   to.Ptr(armedgeactions.EdgeActionVersionDeploymentTypeFile),
					IsDefaultVersion: to.Ptr(armedgeactions.EdgeActionIsDefaultVersionFalse),
				},
			}, nil))
		}},
		{"filter-create", "PUT", filterPath, `{"location":"global","properties":{"versionId":"` + versionPath + `","executionFilterIdentifierHeaderName":"x-sdk-test","executionFilterIdentifierHeaderValue":"initial"}}`, filterJSON, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionExecutionFiltersClient().BeginCreate(ctx, "test-rg", "test-action", "sdkfilter", armedgeactions.EdgeActionExecutionFilter{
				Location: to.Ptr("global"), Properties: &armedgeactions.EdgeActionExecutionFilterProperties{
					VersionID: to.Ptr(versionPath), ExecutionFilterIdentifierHeaderName: to.Ptr("x-sdk-test"), ExecutionFilterIdentifierHeaderValue: to.Ptr("initial"),
				},
			}, nil))
		}},
		{"version-update", "PATCH", versionPath, `{"properties":{"isDefaultVersion":"False"}}`, versionJSON, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionVersionsClient().BeginUpdate(ctx, "test-rg", "test-action", "v1", armedgeactions.EdgeActionVersionUpdate{
				Properties: &armedgeactions.EdgeActionVersionUpdateProperties{IsDefaultVersion: to.Ptr(armedgeactions.EdgeActionIsDefaultVersionFalse)},
			}, nil))
		}},
		{"filter-update", "PATCH", filterPath, `{"tags":{"purpose":"updated"},"properties":{"versionId":"` + versionPath + `","executionFilterIdentifierHeaderName":"x-sdk-test","executionFilterIdentifierHeaderValue":"updated"}}`, filterJSON, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionExecutionFiltersClient().BeginUpdate(ctx, "test-rg", "test-action", "sdkfilter", armedgeactions.EdgeActionExecutionFilterUpdate{
				Tags: map[string]*string{"purpose": to.Ptr("updated")},
				Properties: &armedgeactions.EdgeActionExecutionFilterUpdateProperties{
					VersionID: to.Ptr(versionPath), ExecutionFilterIdentifierHeaderName: to.Ptr("x-sdk-test"), ExecutionFilterIdentifierHeaderValue: to.Ptr("updated"),
				},
			}, nil))
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			factory, server := newMockFactory(t)
			expectRequest(t, server, tc.method, tc.path, tc.body, http.StatusOK, tc.response)
			result, err := tc.run(factory)
			require.NoError(t, err)
			data, err := json.Marshal(result)
			require.NoError(t, err)
			assert.JSONEq(t, tc.response, string(data))
			assert.Equal(t, 1, server.Requests())
		})
	}
}

func TestContractPatchSemantics(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		update     armedgeactions.EdgeActionUpdate
		status     int
	}{
		{"replace", `{"tags":{"purpose":"updated"}}`, armedgeactions.EdgeActionUpdate{Tags: map[string]*string{"purpose": to.Ptr("updated")}}, 200},
		{"clear", `{"tags":{}}`, armedgeactions.EdgeActionUpdate{Tags: map[string]*string{}}, 200},
		{"empty", `{}`, armedgeactions.EdgeActionUpdate{}, 400},
		{"null-tags", `{"tags":null}`, armedgeactions.EdgeActionUpdate{Tags: azcore.NullValue[map[string]*string]()}, 400},
		{"null-sku", `{"sku":null}`, armedgeactions.EdgeActionUpdate{SKU: azcore.NullValue[*armedgeactions.SKUTypeUpdate]()}, 400},
		{"existing-sku", `{"sku":{"name":"Standard","tier":"Standard"}}`, armedgeactions.EdgeActionUpdate{SKU: &armedgeactions.SKUTypeUpdate{Name: to.Ptr("Standard"), Tier: to.Ptr("Standard")}}, 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			factory, server := newMockFactory(t)
			body := actionJSON
			if tc.status == 400 {
				body = `{"error":{"code":"BadRequest","message":"PATCH rejected"}}`
			}
			expectRequest(t, server, "PATCH", actionPath, tc.body, tc.status, body)
			_, err := finish(factory.NewClient().BeginUpdate(context.Background(), "test-rg", "test-action", tc.update, nil))
			if tc.status == 400 {
				requireStatus(t, err, 400)
				assert.Contains(t, err.Error(), "PATCH rejected")
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, 1, server.Requests())
		})
	}
}

func TestContractVersionActions(t *testing.T) {
	for _, deploymentType := range []armedgeactions.EdgeActionVersionDeploymentType{
		armedgeactions.EdgeActionVersionDeploymentTypeFile, armedgeactions.EdgeActionVersionDeploymentTypeZip,
	} {
		t.Run(string(deploymentType), func(t *testing.T) {
			factory, server := newMockFactory(t)
			client := factory.NewEdgeActionVersionsClient()
			content, err := versionCode(deploymentType)
			require.NoError(t, err)
			body, err := json.Marshal(armedgeactions.VersionCode{Name: to.Ptr("v1"), Content: to.Ptr(content)})
			require.NoError(t, err)
			expectRequest(t, server, "POST", versionPath+"/deployVersionCode", string(body), 200, versionJSON)
			deploy, err := client.BeginDeployVersionCode(t.Context(), "test-rg", "test-action", "v1", armedgeactions.VersionCode{Name: to.Ptr("v1"), Content: to.Ptr(content)}, nil)
			require.NoError(t, err)
			deployed := await(t, t.Context(), deploy)
			require.NotNil(t, deployed.Properties)
			assert.Equal(t, to.Ptr(versionPath), deployed.ID)
			assert.Equal(t, to.Ptr("v1"), deployed.Name)
			assert.Equal(t, to.Ptr("global"), deployed.Location)
			assert.Equal(t, to.Ptr(armedgeactions.ProvisioningStateSucceeded), deployed.Properties.ProvisioningState)
			expectRequest(t, server, "POST", versionPath+"/getVersionCode", "", 200, string(body))
			poller, err := client.BeginGetVersionCode(context.Background(), "test-rg", "test-action", "v1", nil)
			require.NoError(t, err)
			result := await(t, context.Background(), poller)
			assert.Equal(t, to.Ptr(content), result.Content)
			expectRequest(t, server, "POST", versionPath+"/swapDefault", "", 204, "")
			_, err = finish(client.BeginSwapDefault(context.Background(), "test-rg", "test-action", "v1", nil))
			require.NoError(t, err)
			assert.Equal(t, 3, server.Requests())
		})
	}
}

func TestContractPagination(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		name, path, resource string
		list                 func(*armedgeactions.ClientFactory) ([]string, error)
	}{
		{"actions-in-group", actionPath[:len(actionPath)-len("/test-action")], actionJSON, func(f *armedgeactions.ClientFactory) ([]string, error) {
			p := f.NewClient().NewListByResourceGroupPager("test-rg", nil)
			var names []string
			for p.More() {
				page, err := p.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, item := range page.Value {
					names = append(names, *item.Name)
				}
			}
			return names, nil
		}},
		{"actions-in-subscription", "/subscriptions/" + subscriptionID + "/providers/Microsoft.Cdn/edgeActions", actionJSON, func(f *armedgeactions.ClientFactory) ([]string, error) {
			p := f.NewClient().NewListBySubscriptionPager(nil)
			var names []string
			for p.More() {
				page, err := p.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, item := range page.Value {
					names = append(names, *item.Name)
				}
			}
			return names, nil
		}},
		{"versions", actionPath + "/versions", versionJSON, func(f *armedgeactions.ClientFactory) ([]string, error) {
			p := f.NewEdgeActionVersionsClient().NewListByEdgeActionPager("test-rg", "test-action", nil)
			var names []string
			for p.More() {
				page, err := p.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, item := range page.Value {
					names = append(names, *item.Name)
				}
			}
			return names, nil
		}},
		{"filters", actionPath + "/executionFilters", filterJSON, func(f *armedgeactions.ClientFactory) ([]string, error) {
			p := f.NewEdgeActionExecutionFiltersClient().NewListByEdgeActionPager("test-rg", "test-action", nil)
			var names []string
			for p.More() {
				page, err := p.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, item := range page.Value {
					names = append(names, *item.Name)
				}
			}
			return names, nil
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			factory, server := newMockFactory(t)
			next := "https://management.azure.com" + tc.path + "?api-version=2026-10-01&$skiptoken=next"
			expectRequest(t, server, "GET", tc.path, "", 200, `{"value":[`+tc.resource+`],"nextLink":"`+next+`"}`)
			server.AppendResponse(mock.WithPredicate(func(req *http.Request) bool {
				assert.Equal(t, tc.path, req.URL.Path)
				assert.Equal(t, "next", req.URL.Query().Get("$skiptoken"))
				assert.Equal(t, "2026-10-01", req.URL.Query().Get("api-version"))
				return false
			}))
			server.AppendResponse(mock.WithStatusCode(200), mock.WithBody([]byte(`{"value":[]}`)), mock.WithHeader("Content-Type", "application/json"))
			names, err := tc.list(factory)
			require.NoError(t, err)
			assert.Len(t, names, 1)
			assert.Equal(t, 2, server.Requests())
		})
	}
}

func TestContractLongRunningOperations(t *testing.T) {
	ctx := context.Background()
	for _, operation := range []struct {
		name, method, path, result string
		status                     int
		run                        func(*armedgeactions.ClientFactory) (any, error)
	}{
		{"action-create", "PUT", actionPath, actionJSON, 201, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewClient().BeginCreate(ctx, "test-rg", "test-action", armedgeactions.EdgeAction{Location: to.Ptr("global"), SKU: &armedgeactions.SKUType{Name: to.Ptr("Standard"), Tier: to.Ptr("Standard")}}, nil))
		}},
		{"version-create", "PUT", versionPath, versionJSON, 201, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionVersionsClient().BeginCreate(ctx, "test-rg", "test-action", "v1", armedgeactions.EdgeActionVersion{Location: to.Ptr("global")}, nil))
		}},
		{"filter-create", "PUT", filterPath, filterJSON, 201, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionExecutionFiltersClient().BeginCreate(ctx, "test-rg", "test-action", "sdkfilter", armedgeactions.EdgeActionExecutionFilter{Location: to.Ptr("global")}, nil))
		}},
		{"action-update", "PATCH", actionPath, actionJSON, 202, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewClient().BeginUpdate(ctx, "test-rg", "test-action", armedgeactions.EdgeActionUpdate{Tags: map[string]*string{}}, nil))
		}},
		{"version-update", "PATCH", versionPath, versionJSON, 202, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionVersionsClient().BeginUpdate(ctx, "test-rg", "test-action", "v1", armedgeactions.EdgeActionVersionUpdate{
				Tags: map[string]*string{"purpose": to.Ptr("updated")},
			}, nil))
		}},
		{"filter-update", "PATCH", filterPath, filterJSON, 202, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionExecutionFiltersClient().BeginUpdate(ctx, "test-rg", "test-action", "sdkfilter", armedgeactions.EdgeActionExecutionFilterUpdate{
				Tags: map[string]*string{"purpose": to.Ptr("updated")},
			}, nil))
		}},
		{"deploy-code", "POST", versionPath + "/deployVersionCode", versionJSON, 202, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionVersionsClient().BeginDeployVersionCode(ctx, "test-rg", "test-action", "v1", armedgeactions.VersionCode{Name: to.Ptr("v1"), Content: to.Ptr("aGVsbG8=")}, nil))
		}},
		{"get-code", "POST", versionPath + "/getVersionCode", `{"name":"v1","content":"aGVsbG8="}`, 202, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionVersionsClient().BeginGetVersionCode(ctx, "test-rg", "test-action", "v1", nil))
		}},
		{"swap-default", "POST", versionPath + "/swapDefault", "{}", 202, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionVersionsClient().BeginSwapDefault(ctx, "test-rg", "test-action", "v1", nil))
		}},
	} {
		statuses := []int{operation.status}
		if operation.method == "PUT" {
			statuses = []int{200, 201}
		}
		for _, status := range statuses {
			for _, terminal := range []string{"Succeeded", "Failed", "Canceled"} {
				t.Run(operation.name+"/"+strconv.Itoa(status)+"/"+terminal, func(t *testing.T) {
					spec, err := loadResponseSpec()
					require.NoError(t, err)
					validator := &responseValidator{spec: spec}
					factory, server := newMockFactory(t, validator)
					pollURL := "https://management.azure.com/operations/one?api-version=2026-10-01"
					resultURL := "https://management.azure.com/results/one?api-version=2026-10-01"
					initialBody := ""
					if operation.method == "PUT" {
						initialBody = strings.Replace(operation.result, `"Succeeded"`, `"Provisioning"`, 1)
					}
					expectRequest(t, server, operation.method, operation.path, "", status, initialBody,
						mock.WithHeader("Azure-AsyncOperation", pollURL), mock.WithHeader("Location", resultURL))
					expectRequest(t, server, "GET", "/operations/one", "", 200, `{"status":"InProgress"}`)
					if terminal != "Succeeded" {
						expectRequest(t, server, "GET", "/operations/one", "", 200, `{"status":"`+terminal+`","error":{"code":"DeploymentFailed","message":"deployment rejected"}}`)
					} else {
						expectRequest(t, server, "GET", "/operations/one", "", 200, `{"status":"Succeeded"}`)
						if operation.method != "POST" {
							expectRequest(t, server, "GET", operation.path, "", 200, operation.result)
						}
						finalBody := operation.result
						if operation.name == "swap-default" {
							finalBody = versionJSON
						}
						expectRequest(t, server, "GET", "/results/one", "", 200, finalBody)
					}
					result, err := operation.run(factory)
					if terminal != "Succeeded" {
						var responseError *azcore.ResponseError
						require.ErrorAs(t, err, &responseError)
						assert.Equal(t, "DeploymentFailed", responseError.ErrorCode)
						assert.Contains(t, err.Error(), "deployment rejected")
						assert.Equal(t, 3, server.Requests())
					} else {
						require.NoError(t, err)
						require.NoError(t, validator.complete(t.Context()))
						data, err := json.Marshal(result)
						require.NoError(t, err)
						assert.JSONEq(t, operation.result, string(data))
						expected := 4
						if operation.method != "POST" {
							expected++
						}
						assert.Equal(t, expected, server.Requests())
					}
				})
			}
		}
	}
}

func TestContractDelete(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		name, path string
		run        func(*armedgeactions.ClientFactory) (any, error)
	}{
		{"action", actionPath, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewClient().BeginDelete(ctx, "test-rg", "test-action", nil))
		}},
		{"version", versionPath, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionVersionsClient().BeginDelete(ctx, "test-rg", "test-action", "v1", nil))
		}},
		{"filter", filterPath, func(f *armedgeactions.ClientFactory) (any, error) {
			return finish(f.NewEdgeActionExecutionFiltersClient().BeginDelete(ctx, "test-rg", "test-action", "sdkfilter", nil))
		}},
	} {
		for _, status := range []int{200, 202, 204} {
			t.Run(tc.name+"/"+http.StatusText(status), func(t *testing.T) {
				spec, err := loadResponseSpec()
				require.NoError(t, err)
				validator := &responseValidator{spec: spec}
				factory, server := newMockFactory(t, validator)
				var headers []mock.ResponseOption
				if status == 202 {
					headers = append(headers, mock.WithHeader("Location", "https://management.azure.com/operations/delete?api-version=2026-10-01"))
				}
				expectRequest(t, server, "DELETE", tc.path, "", status, "", headers...)
				if status == 202 {
					expectRequest(t, server, "GET", "/operations/delete", "", 204, "")
				}
				_, err = tc.run(factory)
				require.NoError(t, err)
				require.NoError(t, validator.complete(t.Context()))
				expected := 1
				if status == 202 {
					expected++
				}
				assert.Equal(t, expected, server.Requests())
			})
		}
	}
}
