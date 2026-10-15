// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package ekm_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/ekm"
	"github.com/stretchr/testify/require"
)

// These tests cover the EKM proxy private endpoint operations and the
// connectivity_mode connection property against a mock transport. Payloads are
// taken from the 2026-07-01-preview REST API examples.

const (
	mockPEHSMURL       = "https://fake-hsm.managedhsm.azure.net"
	mockPEChallenge    = `Bearer authorization="https://login.microsoftonline.com/tenant", resource="https://managedhsm.azure.net"`
	mockPEName         = "my-ekm-pe"
	mockPEOperationID  = "45aacd568ab049a2803861e8dd3ae21f"
	mockPLSAlias       = "myekmproxy-pls.4845dda6-8f58-4643-b3af-1161abcdb077.eastus.azure.privatelinkservice"
	mockPLSResourceID  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.Network/privateLinkServices/myEkmProxyPLS"
	mockPERequestMsg   = "Please approve this private endpoint connection for the EKM proxy."
	mockPEOperationURL = mockPEHSMURL + "/ekm/privateendpoint/operations/" + mockPEOperationID
)

const mockPrivateEndpointBody = `{
	"name":"my-ekm-pe",
	"location":"eastus",
	"provisioningState":"Succeeded",
	"properties":{
		"privateLinkServiceId":"` + mockPLSResourceID + `",
		"privateLinkServiceConnectionName":"my-ekm-pe-connection"
	},
	"privateLinkServiceConnectionState":{
		"status":"Approved",
		"description":"Approved by the Private Link Service owner.",
		"actionsRequired":"None"
	}
}`

const mockPECreateRunningBody = `{
	"jobId":"45aacd568ab049a2803861e8dd3ae21f",
	"privateEndpointName":"my-ekm-pe",
	"operationType":"Create",
	"status":"Running",
	"statusDetails":"Private endpoint creation is in progress.",
	"startTime":1490790000
}`

const mockPECreateSucceededBody = `{
	"jobId":"45aacd568ab049a2803861e8dd3ae21f",
	"privateEndpointName":"my-ekm-pe",
	"operationType":"Create",
	"status":"Succeeded",
	"statusDetails":"Private endpoint creation completed successfully.",
	"startTime":1490790000,
	"endTime":1490790120
}`

const mockPECreateFailedBody = `{
	"jobId":"45aacd568ab049a2803861e8dd3ae21f",
	"privateEndpointName":"my-ekm-pe",
	"operationType":"Create",
	"status":"Failed",
	"statusDetails":"Private endpoint creation failed.",
	"startTime":1490790000,
	"endTime":1490790120,
	"error":{
		"code":"PrivateEndpointCreationFailed",
		"message":"The Private Link Service could not be found.",
		"innererror":{"code":"PrivateLinkServiceNotFound","message":"alias not found"}
	}
}`

const mockPEDeleteRunningBody = `{
	"jobId":"7b2c1e9d4f8a4c3b9e6d5a1f2c3b4d5e",
	"privateEndpointName":"my-ekm-pe",
	"operationType":"Delete",
	"status":"Running",
	"statusDetails":"Private endpoint deletion is in progress.",
	"startTime":1490790000
}`

const mockPEDeleteSucceededBody = `{
	"jobId":"7b2c1e9d4f8a4c3b9e6d5a1f2c3b4d5e",
	"privateEndpointName":"my-ekm-pe",
	"operationType":"Delete",
	"status":"Succeeded",
	"statusDetails":"Private endpoint deletion completed successfully.",
	"startTime":1490790000,
	"endTime":1490790060
}`

type recordedRequest struct {
	method     string
	path       string
	apiVersion string
	authorized bool
	body       []byte
}

// recordingTransport records every request sent through the pipeline before
// forwarding it to the mock server. Unlike mock.WithPredicate, it doesn't
// affect which queued response the mock server returns.
type recordingTransport struct {
	inner policy.Transporter
	mu    sync.Mutex
	reqs  []recordedRequest
}

func (r *recordingTransport) Do(req *http.Request) (*http.Response, error) {
	rec := recordedRequest{
		method:     req.Method,
		path:       req.URL.Path,
		apiVersion: req.URL.Query().Get("api-version"),
		authorized: req.Header.Get("Authorization") != "",
	}
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		_ = req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(b))
		rec.body = b
	}
	r.mu.Lock()
	r.reqs = append(r.reqs, rec)
	r.mu.Unlock()
	return r.inner.Do(req)
}

func (r *recordingTransport) requests() []recordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]recordedRequest(nil), r.reqs...)
}

// newChallengeMockClient returns a client whose first request receives a Key Vault
// authentication challenge. That makes the challenge policy send the complete
// request (including its body) to the mock server, so tests can inspect it.
func newChallengeMockClient(t *testing.T) (*ekm.KeyVaultClient, *mock.Server, *recordingTransport) {
	t.Helper()
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	t.Cleanup(closeSrv)
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusUnauthorized),
		mock.WithHeader("WWW-Authenticate", mockPEChallenge),
	)
	rt := &recordingTransport{inner: srv}
	client, err := ekm.NewClient(mockPEHSMURL, &azcred.Fake{}, &ekm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: rt,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		},
	})
	require.NoError(t, err)
	return client, srv, rt
}

var fastPoll = &runtime.PollUntilDoneOptions{Frequency: time.Millisecond}

func samplePECreateParameters() ekm.PrivateEndpointCreateParameters {
	return ekm.PrivateEndpointCreateParameters{
		PrivateLinkServiceID: to.Ptr(mockPLSAlias),
		RequestMessage:       to.Ptr(mockPERequestMsg),
	}
}

func requireSamplePrivateEndpoint(t *testing.T, pe ekm.PrivateEndpoint) {
	t.Helper()
	require.Equal(t, mockPEName, *pe.Name)
	require.Equal(t, "eastus", *pe.Location)
	require.Equal(t, ekm.EkmPrivateEndpointProvisioningStateSucceeded, *pe.ProvisioningState)
	require.NotNil(t, pe.Properties)
	require.Equal(t, mockPLSResourceID, *pe.Properties.PrivateLinkServiceID)
	require.Equal(t, "my-ekm-pe-connection", *pe.Properties.PrivateLinkServiceConnectionName)
	require.NotNil(t, pe.PrivateLinkServiceConnectionState)
	require.Equal(t, ekm.EkmPrivateEndpointConnectionStatusApproved, *pe.PrivateLinkServiceConnectionState.Status)
	require.Equal(t, "Approved by the Private Link Service owner.", *pe.PrivateLinkServiceConnectionState.Description)
	require.Equal(t, "None", *pe.PrivateLinkServiceConnectionState.ActionsRequired)
}

func TestGetEkmPrivateEndpoint_Mock(t *testing.T) {
	client, srv, rt := newChallengeMockClient(t)
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(mockPrivateEndpointBody)))

	res, err := client.GetEkmPrivateEndpoint(context.Background(), mockPEName, nil)
	require.NoError(t, err)
	requireSamplePrivateEndpoint(t, res.PrivateEndpoint)

	reqs := rt.requests()
	require.Len(t, reqs, 2)
	require.Equal(t, http.MethodGet, reqs[1].method)
	require.Equal(t, "/ekm/privateendpoint/"+mockPEName, reqs[1].path)
	require.Equal(t, "2026-07-01-preview", reqs[1].apiVersion)
	require.True(t, reqs[1].authorized)
}

func TestGetEkmPrivateEndpoint_MockError(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusNotFound),
		mock.WithBody([]byte(`{"error":{"code":"PrivateEndpointNotFound","message":"not found"}}`)),
	)

	client := newMockClient(t, srv)
	_, err := client.GetEkmPrivateEndpoint(context.Background(), mockPEName, nil)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusNotFound, respErr.StatusCode)
	require.Equal(t, "PrivateEndpointNotFound", respErr.ErrorCode)
}

func TestListEkmPrivateEndpoints_Mock(t *testing.T) {
	client, srv, rt := newChallengeMockClient(t)
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"value":[`+mockPrivateEndpointBody+`]}`)))

	res, err := client.ListEkmPrivateEndpoints(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, res.Value, 1)
	requireSamplePrivateEndpoint(t, *res.Value[0])

	reqs := rt.requests()
	require.Len(t, reqs, 2)
	require.Equal(t, http.MethodGet, reqs[1].method)
	require.Equal(t, "/ekm/privateendpoint", reqs[1].path)
}

func TestListEkmPrivateEndpoints_MockEmpty(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"value":[]}`)))

	client := newMockClient(t, srv)
	res, err := client.ListEkmPrivateEndpoints(context.Background(), nil)
	require.NoError(t, err)
	require.Empty(t, res.Value)
}

func TestListEkmPrivateEndpoints_MockError(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusForbidden),
		mock.WithBody([]byte(`{"error":{"code":"Forbidden","message":"missing ekm/read"}}`)),
	)

	client := newMockClient(t, srv)
	_, err := client.ListEkmPrivateEndpoints(context.Background(), nil)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusForbidden, respErr.StatusCode)
}

func TestGetEkmPrivateEndpointOperationStatus_Mock(t *testing.T) {
	client, srv, rt := newChallengeMockClient(t)
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(mockPECreateSucceededBody)))

	res, err := client.GetEkmPrivateEndpointOperationStatus(context.Background(), mockPEOperationID, nil)
	require.NoError(t, err)
	require.Equal(t, mockPEOperationID, *res.JobID)
	require.Equal(t, mockPEName, *res.PrivateEndpointName)
	require.Equal(t, ekm.EkmPrivateEndpointOperationTypeCreate, *res.OperationType)
	require.Equal(t, ekm.OperationStateSucceeded, *res.Status)
	require.Equal(t, "Private endpoint creation completed successfully.", *res.StatusDetails)
	require.True(t, time.Unix(1490790000, 0).Equal(*res.StartTime))
	require.True(t, time.Unix(1490790120, 0).Equal(*res.EndTime))
	require.Nil(t, res.Error)

	reqs := rt.requests()
	require.Len(t, reqs, 2)
	require.Equal(t, http.MethodGet, reqs[1].method)
	require.Equal(t, "/ekm/privateendpoint/operations/"+mockPEOperationID, reqs[1].path)
}

func TestGetEkmPrivateEndpointOperationStatus_MockFailedOperation(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(mockPECreateFailedBody)))

	client := newMockClient(t, srv)
	res, err := client.GetEkmPrivateEndpointOperationStatus(context.Background(), mockPEOperationID, nil)
	require.NoError(t, err)
	require.Equal(t, ekm.OperationStateFailed, *res.Status)
	require.NotNil(t, res.Error)
	require.Equal(t, "PrivateEndpointCreationFailed", *res.Error.Code)
	require.Equal(t, "The Private Link Service could not be found.", *res.Error.Message)
	require.NotNil(t, res.Error.InnerError)
	require.Equal(t, "PrivateLinkServiceNotFound", *res.Error.InnerError.Code)
	require.Equal(t, "alias not found", *res.Error.InnerError.Message)
	require.Nil(t, res.Error.InnerError.InnerError)
}

func TestGetEkmPrivateEndpointOperationStatus_MockError(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusNotFound),
		mock.WithBody([]byte(`{"error":{"code":"OperationNotFound","message":"not found"}}`)),
	)

	client := newMockClient(t, srv)
	_, err := client.GetEkmPrivateEndpointOperationStatus(context.Background(), mockPEOperationID, nil)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusNotFound, respErr.StatusCode)
}

func TestBeginCreateEkmPrivateEndpoint_Mock(t *testing.T) {
	client, srv, rt := newChallengeMockClient(t)
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusAccepted),
		mock.WithHeader("Azure-AsyncOperation", mockPEOperationURL),
		mock.WithBody([]byte(mockPECreateRunningBody)),
	)
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(mockPECreateRunningBody)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(mockPECreateSucceededBody)))

	poller, err := client.BeginCreateEkmPrivateEndpoint(context.Background(), mockPEName, samplePECreateParameters(), nil)
	require.NoError(t, err)
	require.False(t, poller.Done())

	res, err := poller.PollUntilDone(context.Background(), fastPoll)
	require.NoError(t, err)
	require.Equal(t, mockPEOperationID, *res.JobID)
	require.Equal(t, ekm.EkmPrivateEndpointOperationTypeCreate, *res.OperationType)
	require.Equal(t, ekm.OperationStateSucceeded, *res.Status)
	require.True(t, time.Unix(1490790120, 0).Equal(*res.EndTime))

	reqs := rt.requests()
	require.Len(t, reqs, 4)
	// reqs[0] is the challenge policy's body-less probe
	create := reqs[1]
	require.Equal(t, http.MethodPost, create.method)
	require.Equal(t, "/ekm/privateendpoint/"+mockPEName, create.path)
	require.Equal(t, "2026-07-01-preview", create.apiVersion)
	require.True(t, create.authorized)
	require.JSONEq(t, `{"privateLinkServiceId":"`+mockPLSAlias+`","requestMessage":"`+mockPERequestMsg+`"}`, string(create.body))
	for _, poll := range reqs[2:] {
		require.Equal(t, http.MethodGet, poll.method)
		require.Equal(t, "/ekm/privateendpoint/operations/"+mockPEOperationID, poll.path)
	}
}

func TestBeginCreateEkmPrivateEndpoint_MockResumeToken(t *testing.T) {
	client, srv, rt := newChallengeMockClient(t)
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusAccepted),
		mock.WithHeader("Azure-AsyncOperation", mockPEOperationURL),
		mock.WithBody([]byte(mockPECreateRunningBody)),
	)
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(mockPECreateSucceededBody)))

	poller, err := client.BeginCreateEkmPrivateEndpoint(context.Background(), mockPEName, samplePECreateParameters(), nil)
	require.NoError(t, err)
	token, err := poller.ResumeToken()
	require.NoError(t, err)

	resumed, err := client.BeginCreateEkmPrivateEndpoint(context.Background(), mockPEName, samplePECreateParameters(), &ekm.BeginCreateEkmPrivateEndpointOptions{ResumeToken: token})
	require.NoError(t, err)
	res, err := resumed.PollUntilDone(context.Background(), fastPoll)
	require.NoError(t, err)
	require.Equal(t, ekm.OperationStateSucceeded, *res.Status)

	// resuming must poll the operation rather than send another create request
	reqs := rt.requests()
	require.Len(t, reqs, 3)
	require.Equal(t, http.MethodGet, reqs[2].method)
	require.Equal(t, "/ekm/privateendpoint/operations/"+mockPEOperationID, reqs[2].path)
}

func TestBeginCreateEkmPrivateEndpoint_MockOperationFailed(t *testing.T) {
	client, srv, _ := newChallengeMockClient(t)
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusAccepted),
		mock.WithHeader("Azure-AsyncOperation", mockPEOperationURL),
		mock.WithBody([]byte(mockPECreateRunningBody)),
	)
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(mockPECreateFailedBody)))

	poller, err := client.BeginCreateEkmPrivateEndpoint(context.Background(), mockPEName, samplePECreateParameters(), nil)
	require.NoError(t, err)
	_, err = poller.PollUntilDone(context.Background(), fastPoll)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, "PrivateEndpointCreationFailed", respErr.ErrorCode)
}

func TestBeginCreateEkmPrivateEndpoint_MockError(t *testing.T) {
	client, srv, _ := newChallengeMockClient(t)
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusConflict),
		mock.WithBody([]byte(`{"error":{"code":"Conflict","message":"a pool may have up to two private endpoints"}}`)),
	)

	_, err := client.BeginCreateEkmPrivateEndpoint(context.Background(), mockPEName, samplePECreateParameters(), nil)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusConflict, respErr.StatusCode)
	require.Equal(t, "Conflict", respErr.ErrorCode)
}

func TestBeginDeleteEkmPrivateEndpoint_Mock(t *testing.T) {
	client, srv, rt := newChallengeMockClient(t)
	deleteOperationURL := mockPEHSMURL + "/ekm/privateendpoint/operations/7b2c1e9d4f8a4c3b9e6d5a1f2c3b4d5e"
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusAccepted),
		mock.WithHeader("Azure-AsyncOperation", deleteOperationURL),
		mock.WithBody([]byte(mockPEDeleteRunningBody)),
	)
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(mockPEDeleteSucceededBody)))

	poller, err := client.BeginDeleteEkmPrivateEndpoint(context.Background(), mockPEName, nil)
	require.NoError(t, err)
	res, err := poller.PollUntilDone(context.Background(), fastPoll)
	require.NoError(t, err)
	require.Equal(t, "7b2c1e9d4f8a4c3b9e6d5a1f2c3b4d5e", *res.JobID)
	require.Equal(t, ekm.EkmPrivateEndpointOperationTypeDelete, *res.OperationType)
	require.Equal(t, ekm.OperationStateSucceeded, *res.Status)

	reqs := rt.requests()
	require.Len(t, reqs, 3)
	require.Equal(t, http.MethodDelete, reqs[1].method)
	require.Equal(t, "/ekm/privateendpoint/"+mockPEName, reqs[1].path)
	require.Equal(t, "2026-07-01-preview", reqs[1].apiVersion)
	require.Equal(t, http.MethodGet, reqs[2].method)
	require.Equal(t, "/ekm/privateendpoint/operations/7b2c1e9d4f8a4c3b9e6d5a1f2c3b4d5e", reqs[2].path)
}

func TestBeginDeleteEkmPrivateEndpoint_MockError(t *testing.T) {
	client, srv, _ := newChallengeMockClient(t)
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusConflict),
		mock.WithBody([]byte(`{"error":{"code":"Conflict","message":"private endpoint is referenced by the EKM connection"}}`)),
	)

	_, err := client.BeginDeleteEkmPrivateEndpoint(context.Background(), mockPEName, nil)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusConflict, respErr.StatusCode)
}

func TestPrivateEndpointRequiredParameters(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	client := newMockClient(t, srv)
	ctx := context.Background()

	_, err := client.BeginCreateEkmPrivateEndpoint(ctx, "", samplePECreateParameters(), nil)
	require.ErrorContains(t, err, "peName cannot be empty")
	_, err = client.BeginDeleteEkmPrivateEndpoint(ctx, "", nil)
	require.ErrorContains(t, err, "peName cannot be empty")
	_, err = client.GetEkmPrivateEndpoint(ctx, "", nil)
	require.ErrorContains(t, err, "peName cannot be empty")
	_, err = client.GetEkmPrivateEndpointOperationStatus(ctx, "", nil)
	require.ErrorContains(t, err, "operationID cannot be empty")
	require.Zero(t, srv.Requests(), "no request should be sent when a required parameter is empty")
}

func TestCreateEkmConnection_MockPrivateEndpointConnectivity(t *testing.T) {
	client, srv, rt := newChallengeMockClient(t)
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{
		"host":"my-ekm-pe",
		"connectivity_mode":"PrivateEndpoint",
		"server_subject_common_name":"ekm-proxy.contoso.com",
		"server_ca_certificates":["Y2EtY2VydC1ieXRlcw=="]
	}`)))

	conn := sampleConnection()
	conn.Host = to.Ptr(mockPEName)
	conn.ConnectivityMode = to.Ptr(ekm.EkmConnectivityModePrivateEndpoint)
	res, err := client.CreateEkmConnection(context.Background(), conn, nil)
	require.NoError(t, err)
	require.Equal(t, mockPEName, *res.Host)
	require.Equal(t, ekm.EkmConnectivityModePrivateEndpoint, *res.ConnectivityMode)

	reqs := rt.requests()
	require.Len(t, reqs, 2)
	var sent map[string]any
	require.NoError(t, json.Unmarshal(reqs[1].body, &sent))
	require.Equal(t, "PrivateEndpoint", sent["connectivity_mode"])
	require.Equal(t, mockPEName, sent["host"])
}

func TestGetEkmConnection_MockDefaultConnectivityMode(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(mockEKMConnectionBody)))

	client := newMockClient(t, srv)
	res, err := client.GetEkmConnection(context.Background(), nil)
	require.NoError(t, err)
	// the service omits connectivity_mode for connections created before it existed
	require.Nil(t, res.ConnectivityMode)
}

func TestPrivateEndpointSerde(t *testing.T) {
	pe := ekm.PrivateEndpoint{
		Location: to.Ptr("eastus"),
		Name:     to.Ptr(mockPEName),
		PrivateLinkServiceConnectionState: &ekm.PrivateEndpointConnectionState{
			ActionsRequired: to.Ptr("None"),
			Description:     to.Ptr("Approved"),
			Status:          to.Ptr(ekm.EkmPrivateEndpointConnectionStatusApproved),
		},
		Properties: &ekm.PrivateEndpointProperties{
			PrivateLinkServiceConnectionName: to.Ptr("my-ekm-pe-connection"),
			PrivateLinkServiceID:             to.Ptr(mockPLSResourceID),
		},
		ProvisioningState: to.Ptr(ekm.EkmPrivateEndpointProvisioningStateSucceeded),
	}
	testSerde(t, &pe)
	testSerde(t, pe.PrivateLinkServiceConnectionState)
	testSerde(t, pe.Properties)

	list := ekm.PrivateEndpointListResult{Value: []*ekm.PrivateEndpoint{&pe}}
	testSerde(t, &list)

	params := samplePECreateParameters()
	testSerde(t, &params)

	start := time.Unix(1490790000, 0).UTC()
	end := start.Add(2 * time.Minute)
	op := ekm.PrivateEndpointOperation{
		EndTime: &end,
		Error: &ekm.ErrorInfo{
			Code:       to.Ptr("PrivateEndpointCreationFailed"),
			Message:    to.Ptr("failed"),
			InnerError: &ekm.ErrorInfo{Code: to.Ptr("Inner"), Message: to.Ptr("inner failure")},
		},
		JobID:               to.Ptr(mockPEOperationID),
		OperationType:       to.Ptr(ekm.EkmPrivateEndpointOperationTypeCreate),
		PrivateEndpointName: to.Ptr(mockPEName),
		StartTime:           &start,
		Status:              to.Ptr(ekm.OperationStateFailed),
		StatusDetails:       to.Ptr("failed"),
	}
	data, err := json.Marshal(op)
	require.NoError(t, err)
	var roundTripped ekm.PrivateEndpointOperation
	require.NoError(t, json.Unmarshal(data, &roundTripped))
	require.True(t, start.Equal(*roundTripped.StartTime))
	require.True(t, end.Equal(*roundTripped.EndTime))
	roundTripped.StartTime, roundTripped.EndTime = op.StartTime, op.EndTime
	require.Equal(t, op, roundTripped)
	testSerde(t, &op)
	testSerde(t, op.Error)
}
