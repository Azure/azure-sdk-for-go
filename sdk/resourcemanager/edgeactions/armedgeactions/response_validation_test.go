// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armedgeactions_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/edgeactions/armedgeactions"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type validatorKey struct{}
type skipValidationKey struct{}

type pendingOperation struct {
	op                   *specOperation
	statusURL, resultURL string
	lastStatus           string
	polls                int
	resultSeen           bool
	initialSucceeded     bool
}

// Scenarios execute one operation at a time. This policy validates the wire
// response before deserialization, including fields the SDK would otherwise ignore.
type responseValidator struct {
	spec     *responseSpec
	pipeline runtime.Pipeline
	pending  *pendingOperation
}

func (v *responseValidator) Do(req *policy.Request) (*http.Response, error) {
	method, address := req.Raw().Method, req.Raw().URL.String()
	resp, err := req.Next()
	if err != nil || req.Raw().Context().Value(skipValidationKey{}) != nil {
		return resp, err
	}
	body, err := runtime.Payload(resp)
	if err != nil {
		return resp, fmt.Errorf("read response for contract validation: %w", err)
	}
	if p := v.pending; p != nil && method == http.MethodGet {
		if sameURL(address, p.statusURL) {
			if terminalStatus(p.lastStatus) {
				return resp, fmt.Errorf("%s: status polled after termination", p.op.ID)
			}
			state, err := validateStatus(resp.StatusCode, body)
			if err != nil {
				return resp, err
			}
			p.lastStatus = state
			p.polls++
			return resp, nil
		}
		if sameURL(address, p.resultURL) {
			if p.statusURL != "" && !strings.EqualFold(p.lastStatus, "Succeeded") && !p.initialSucceeded {
				return resp, fmt.Errorf("%s: final result requested before status succeeded", p.op.ID)
			}
			if resp.StatusCode == http.StatusAccepted {
				if p.statusURL != "" {
					return resp, fmt.Errorf("%s: final Location is still pending", p.op.ID)
				}
				return resp, v.spec.validate(p.op, resp.StatusCode, body, resp.Header)
			}
			// Only the known swapDefault final Location response is excluded.
			// Its initial response, URLs, polling, and subsequent resource GETs are checked.
			if !swapLocationException(p.op, resp.StatusCode, body) {
				if err := v.spec.validateFinal(p.op, resp.StatusCode, body, resp.Header); err != nil {
					return resp, err
				}
			}
			p.resultSeen = true
			return resp, nil
		}
	}
	// Expected negative service tests retain the SDK's normal ResponseError handling.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, nil
	}
	op, err := v.spec.operation(method, address)
	if err != nil {
		return resp, err
	}
	if op == nil {
		if strings.Contains(strings.ToLower(address), "/providers/microsoft.cdn/edgeactions") {
			return resp, fmt.Errorf("undocumented Edge Actions operation")
		}
		return resp, nil
	}
	if err := v.spec.validate(op, resp.StatusCode, body, resp.Header); err != nil {
		return resp, err
	}
	if !op.LRO {
		return resp, nil
	}
	statusURL, resultURL := resp.Header.Get("Azure-AsyncOperation"), resp.Header.Get("Location")
	if statusURL == "" && resultURL == "" {
		if resp.StatusCode == http.StatusAccepted {
			return resp, fmt.Errorf("%s: HTTP 202 has no polling URL", op.ID)
		}
		return resp, nil
	}
	for _, address := range []string{statusURL, resultURL} {
		if address != "" {
			if err := validatePollingURL(address); err != nil {
				return resp, fmt.Errorf("%s: %w", op.ID, err)
			}
		}
	}
	if statusURL != "" && sameURL(statusURL, resultURL) {
		return resp, fmt.Errorf("%s: status and final Location URLs must differ", op.ID)
	}
	if v.pending != nil {
		return resp, fmt.Errorf("%s: previous LRO was not completed and checked", op.ID)
	}
	var resource struct {
		Properties struct{ ProvisioningState string }
	}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &resource); err != nil {
			return resp, fmt.Errorf("%s: invalid initial resource body", op.ID)
		}
	}
	v.pending = &pendingOperation{
		op: op, statusURL: statusURL, resultURL: resultURL,
		initialSucceeded: (method == http.MethodPut || method == http.MethodPatch) &&
			strings.EqualFold(resource.Properties.ProvisioningState, "Succeeded"),
	}
	return resp, nil
}

func (v *responseValidator) complete(ctx context.Context) error {
	p := v.pending
	if p == nil {
		return nil
	}
	defer func() { v.pending = nil }()
	if p.statusURL != "" && p.polls == 0 && p.initialSucceeded {
		// The SDK correctly skips polling an already-completed PUT. Still check
		// the advertised monitor; do not require the SDK to make redundant requests.
		if err := v.get(ctx, p.statusURL); err != nil {
			return err
		}
	}
	if p.statusURL != "" && (p.polls == 0 || !strings.EqualFold(p.lastStatus, "Succeeded")) {
		return fmt.Errorf("%s: missing successful Azure-AsyncOperation polling", p.op.ID)
	}
	if p.resultURL != "" && !p.resultSeen && p.op.ID != "EdgeActionVersions_SwapDefault" {
		// PUT/PATCH pollers can GET the original resource instead of Location.
		// Check the advertised endpoint through the same authenticated, recorded pipeline.
		if err := v.get(ctx, p.resultURL); err != nil {
			return err
		}

		if !p.resultSeen {
			return fmt.Errorf("%s: final Location was not validated", p.op.ID)
		}
	}
	return nil
}

func (v *responseValidator) get(ctx context.Context, address string) (err error) {
	request, err := runtime.NewRequest(ctx, http.MethodGet, address)
	if err != nil {
		return fmt.Errorf("invalid LRO validation request URL")
	}
	response, err := v.pipeline.Do(request)
	if response != nil && response.Body != nil {
		defer func() { err = errors.Join(err, response.Body.Close()) }()
	}
	return err
}

func sameURL(left, right string) bool {
	if left == "" || right == "" {
		return false
	}
	a, err := url.Parse(left)
	if err != nil {
		return false
	}
	b, err := url.Parse(right)
	return err == nil && strings.EqualFold(a.Scheme, b.Scheme) && strings.EqualFold(a.Host, b.Host) &&
		strings.EqualFold(a.Path, b.Path) && a.Query().Encode() == b.Query().Encode()
}

func terminalStatus(state string) bool {
	return strings.EqualFold(state, "Succeeded") || strings.EqualFold(state, "Failed") || strings.EqualFold(state, "Canceled")
}

func pollForScenario[T any](ctx context.Context, poller *runtime.Poller[T]) (*T, error) {
	// Each LRO gets a fresh budget; earlier deployments cannot consume a later
	// execution filter's timeout. Parent cancellation still bounds the entire test.
	ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	result, err := testutil.PollForTest(ctx, poller)
	if validator, ok := ctx.Value(validatorKey{}).(*responseValidator); ok {
		if err != nil {
			validator.pending = nil
		} else {
			err = validator.complete(ctx)
		}
	}
	return result, err
}

func TestContractAlreadyCompletedPUT(t *testing.T) {
	for _, status := range []int{200, 201} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			spec, err := loadResponseSpec()
			require.NoError(t, err)
			validator := &responseValidator{spec: spec}
			factory, server := newMockFactory(t, validator)
			statusURL := "https://management.azure.com/operations/one?api-version=" + apiVersion
			resultURL := "https://management.azure.com" + versionPath + "?api-version=" + apiVersion
			expectRequest(t, server, "PUT", versionPath, "", status, versionJSON,
				mock.WithHeader("Azure-AsyncOperation", statusURL), mock.WithHeader("Location", resultURL))
			expectRequest(t, server, "GET", versionPath, "", 200, versionJSON)
			expectRequest(t, server, "GET", "/operations/one", "", 200, `{"status":"Succeeded"}`)
			ctx := context.WithValue(t.Context(), validatorKey{}, validator)
			poller, err := factory.NewEdgeActionVersionsClient().BeginCreate(ctx, "test-rg", "test-action", "v1",
				armedgeactions.EdgeActionVersion{Location: to.Ptr("global")}, nil)
			require.NoError(t, err)
			result, err := pollForScenario(ctx, poller)
			require.NoError(t, err)
			assert.True(t, poller.Done())
			assert.Equal(t, to.Ptr(versionPath), result.ID)
			assert.Equal(t, 3, server.Requests(), "validate advertised monitor even when the SDK skips polling")
			assert.Nil(t, validator.pending)
		})
	}
}

func TestContractInvalidLROResponses(t *testing.T) {
	const monitor = "https://management.azure.com/operations/one?api-version=" + apiVersion
	const final = "https://management.azure.com/results/one?api-version=" + apiVersion
	const initiate = "https://management.azure.com" + versionPath + "/deployVersionCode?api-version=" + apiVersion
	for _, tc := range []struct {
		name, statusURL, resultURL, pollBody, resultBody, failure string
		pollStatus, resultStatus                                  int
	}{
		{"missing-headers", "", "", "", "", "no polling URL", 0, 0},
		{"insecure-monitor", "http://example.test/monitor", final, "", "", "HTTPS", 0, 0},
		{"insecure-result", monitor, "http://example.test/result", "", "", "HTTPS", 0, 0},
		{"same-urls", monitor, monitor, "", "", "must differ", 0, 0},
		{"monitor-api-version", strings.Replace(monitor, apiVersion, "old", 1), final, "", "", "API version", 0, 0},
		{"missing-status", monitor, final, `{}`, "", "nonempty string status", 200, 0},
		{"null-status", monitor, final, `{"status":null}`, "", "nonempty string status", 200, 0},
		{"monitor-http202", monitor, final, `{"status":"Succeeded"}`, "", "expected HTTP 200", 202, 0},
		{"early-result", monitor, final, `{"status":"Running"}`, versionJSON, "before status succeeded", 200, 200},
		{"final-status-document", monitor, final, `{"status":"Succeeded"}`, `{"status":"Succeeded"}`, "required", 200, 200},
		{"final-null", monitor, final, `{"status":"Succeeded"}`, "null", "type", 200, 200},
		{"final-empty", monitor, final, `{"status":"Succeeded"}`, "", "requires a response body", 200, 200},
		{"final-http202", monitor, final, `{"status":"Succeeded"}`, "", "still pending", 200, 202},
		{"final-http204", monitor, final, `{"status":"Succeeded"}`, "", "not a documented response", 200, 204},
		{"final-http404", monitor, final, `{"status":"Succeeded"}`, `{"error":{"code":"NotFound","message":"Missing result"}}`, "non-success HTTP 404", 200, 404},
		{"final-http500", monitor, final, `{"status":"Succeeded"}`, `{"error":{"code":"InternalError","message":"Failed result"}}`, "non-success HTTP 500", 200, 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			spec, err := loadResponseSpec()
			require.NoError(t, err)
			validator := &responseValidator{spec: spec}
			_, server := newMockFactory(t, validator)
			send := func(method, address string, status int, body string, headers ...mock.ResponseOption) error {
				options := []mock.ResponseOption{mock.WithStatusCode(status), mock.WithBody([]byte(body)), mock.WithHeader("Content-Type", "application/json")}
				server.AppendResponse(append(options, headers...)...)
				request, err := runtime.NewRequest(t.Context(), method, address)
				require.NoError(t, err)
				response, err := validator.pipeline.Do(request)
				if response != nil && response.Body != nil {
					defer func() { assert.NoError(t, response.Body.Close()) }()
				}
				return err
			}
			err = send("POST", initiate, 202, "", mock.WithHeader("Azure-AsyncOperation", tc.statusURL), mock.WithHeader("Location", tc.resultURL))
			if tc.pollStatus != 0 {
				require.NoError(t, err)
				err = send("GET", monitor, tc.pollStatus, tc.pollBody)
			}
			if tc.resultStatus != 0 {
				require.NoError(t, err)
				err = send("GET", final, tc.resultStatus, tc.resultBody)
			}
			require.ErrorContains(t, err, tc.failure)
			if tc.resultStatus != 0 {
				require.NotNil(t, validator.pending)
				assert.False(t, validator.pending.resultSeen)
			}
		})
	}
}

func TestContractCompleteRejectsFailedLocation(t *testing.T) {
	for _, method := range []string{http.MethodPut, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			spec, err := loadResponseSpec()
			require.NoError(t, err)
			op, err := spec.operation(method, "https://management.azure.com"+actionPath+"?api-version="+apiVersion)
			require.NoError(t, err)
			require.NotNil(t, op)
			pending := &pendingOperation{
				op: op, statusURL: "https://management.azure.com/operations/one?api-version=" + apiVersion,
				resultURL:  "https://management.azure.com/results/one?api-version=" + apiVersion,
				lastStatus: "Succeeded", polls: 1,
			}
			validator := &responseValidator{spec: spec, pending: pending}
			_, server := newMockFactory(t, validator)
			expectRequest(t, server, "GET", "/results/one", "", http.StatusNotFound,
				`{"error":{"code":"NotFound","message":"Missing result"}}`)
			require.ErrorContains(t, validator.complete(t.Context()), "non-success HTTP 404")
			assert.False(t, pending.resultSeen)
			assert.Nil(t, validator.pending)
			assert.Equal(t, 1, server.Requests())
		})
	}
}

func TestContractSwapExceptionIsScoped(t *testing.T) {
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	swap, err := spec.operation("POST", "https://management.azure.com"+versionPath+"/swapDefault?api-version="+apiVersion)
	require.NoError(t, err)
	require.NotNil(t, swap)
	require.ErrorContains(t, spec.validate(swap, 200, []byte(versionJSON), nil), "not a documented response")
	require.NoError(t, spec.validate(swap, 204, nil, nil))
	require.ErrorContains(t, spec.validate(swap, 204, []byte(versionJSON), nil), "does not declare")
	for _, operation := range []string{"deployVersionCode", "getVersionCode", "swapDefault"} {
		op, err := spec.operation("POST", "https://management.azure.com"+versionPath+"/"+operation+"?api-version="+apiVersion)
		require.NoError(t, err)
		validator := &responseValidator{
			spec:    spec,
			pending: &pendingOperation{op: op, statusURL: "https://management.azure.com/status", resultURL: "https://management.azure.com/result"},
		}
		require.ErrorContains(t, validator.complete(t.Context()), "missing successful")
	}
}
