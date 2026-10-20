// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armedgeactions_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type contractRecording struct {
	Entries []contractEntry
}

type contractEntry struct {
	RequestURI      string `json:"RequestUri"`
	RequestMethod   string
	StatusCode      int
	ResponseHeaders map[string]recordedHeader
	ResponseBody    json.RawMessage
}

func (e contractEntry) headers() http.Header {
	headers := make(http.Header)
	for key, values := range e.ResponseHeaders {
		headers[http.CanonicalHeaderKey(key)] = []string(values)
	}
	return headers
}

func (e contractEntry) body() ([]byte, error) {
	data := bytes.TrimSpace(e.ResponseBody)
	// The proxy stores an absent body as null and text payloads as JSON strings.
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil, nil
	}
	if data[0] == '"' {
		var text string
		if json.Unmarshal(data, &text) != nil {
			return nil, fmt.Errorf("invalid recording body encoding")
		}
		return []byte(text), nil
	}
	return data, nil
}

type recordingIssue struct {
	entry int
	err   error
}

type recordingAudit struct {
	issues     []recordingIssue
	checked    int
	fixtures   int
	exceptions int
	operations map[string]int
}

func (a *recordingAudit) check(entry int, err error) {
	if err != nil {
		a.issues = append(a.issues, recordingIssue{entry: entry, err: err})
	}
}

type recordedOperation struct {
	pendingOperation
	entry int
}

// Audit every stored response independently of lifecycle execution, so assertions
// or cleanup bypasses cannot hide later failures in the same recording.
func auditRecording(spec *responseSpec, entries []contractEntry) recordingAudit {
	audit := recordingAudit{operations: make(map[string]int)}
	var operations []*recordedOperation
	var fixtureURLs []string
	for i, entry := range entries {
		index := i + 1
		headers := entry.headers()
		body, err := entry.body()
		audit.check(index, err)
		if err != nil {
			continue
		}
		address, err := url.Parse(entry.RequestURI)
		if err != nil {
			audit.check(index, fmt.Errorf("invalid recorded request URL"))
			continue
		}
		// Resource-group setup/teardown is a different API, not part of the
		// pinned Edge Actions spec. Only its explicit routes/advertised URLs are excluded.
		fixture := isResourceGroupRoute(address.Path)
		if strings.EqualFold(entry.RequestMethod, http.MethodGet) {
			for _, target := range fixtureURLs {
				fixture = fixture || sameRecordedURL(entry.RequestURI, target)
			}
		}
		if fixture {
			audit.fixtures++
			for _, key := range []string{"Azure-AsyncOperation", "Location"} {
				if target := headers.Get(key); target != "" {
					fixtureURLs = append(fixtureURLs, target)
				}
			}
			continue
		}
		audit.checked++
		matched := false
		if strings.EqualFold(entry.RequestMethod, http.MethodGet) {
			for j := len(operations) - 1; j >= 0; j-- {
				op := operations[j]
				if sameRecordedURL(entry.RequestURI, op.statusURL) {
					matched = true
					state, err := validateStatus(entry.StatusCode, body)
					audit.check(index, err)
					if terminalStatus(op.lastStatus) {
						audit.check(index, fmt.Errorf("%s: status polled after termination", op.op.ID))
					}
					op.polls++
					op.lastStatus = state
					break
				}
				if !op.resultSeen && sameRecordedURL(entry.RequestURI, op.resultURL) {
					matched = true
					if op.statusURL != "" && !op.initialSucceeded && !strings.EqualFold(op.lastStatus, "Succeeded") {
						audit.check(index, fmt.Errorf("%s: final result requested before status succeeded", op.op.ID))
					}
					if entry.StatusCode == http.StatusAccepted {
						if op.statusURL != "" {
							audit.check(index, fmt.Errorf("%s: final Location is still pending", op.op.ID))
						}
						audit.check(index, spec.validate(op.op, entry.StatusCode, body, headers))
					} else {
						// The exception applies to this LRO result, never a later resource GET.
						if swapLocationException(op.op, entry.StatusCode, body) {
							audit.exceptions++
						} else {
							audit.check(index, spec.validateFinal(op.op, entry.StatusCode, body, headers))
						}
						op.resultSeen = entry.StatusCode >= http.StatusOK && entry.StatusCode < http.StatusMultipleChoices
					}
					break
				}
			}
		}
		if matched {
			continue
		}
		op, err := spec.operation(entry.RequestMethod, entry.RequestURI)
		audit.check(index, err)
		if err != nil {
			continue
		}
		if op == nil {
			audit.check(index, fmt.Errorf("response has no matching Edge Actions operation or advertised polling endpoint"))
			continue
		}
		audit.operations[op.ID]++
		audit.check(index, spec.validate(op, entry.StatusCode, body, headers))
		// Record LRO associations even when its initial schema is invalid.
		if !op.LRO || entry.StatusCode < 200 || entry.StatusCode >= 300 {
			continue
		}
		statusURL, resultURL := headers.Get("Azure-AsyncOperation"), headers.Get("Location")
		if statusURL == "" && resultURL == "" {
			if entry.StatusCode == http.StatusAccepted {
				audit.check(index, fmt.Errorf("%s: HTTP 202 has no polling URL", op.ID))
			}
			continue
		}
		for _, target := range []string{statusURL, resultURL} {
			if target != "" {
				audit.check(index, validatePollingURL(target))
			}
		}
		if sameURL(statusURL, resultURL) {
			audit.check(index, fmt.Errorf("%s: status and final Location URLs must differ", op.ID))
		}
		var resource struct {
			Properties struct{ ProvisioningState string }
		}
		// Schema validation above reports malformed initial bodies.
		decoded := json.Unmarshal(body, &resource) == nil
		operations = append(operations, &recordedOperation{
			entry: index,
			pendingOperation: pendingOperation{
				op: op, statusURL: statusURL, resultURL: resultURL,
				initialSucceeded: decoded && (strings.EqualFold(entry.RequestMethod, "PUT") || strings.EqualFold(entry.RequestMethod, "PATCH")) &&
					strings.EqualFold(resource.Properties.ProvisioningState, "Succeeded"),
			},
		})
	}
	for _, op := range operations {
		if op.statusURL != "" && (op.polls == 0 || !strings.EqualFold(op.lastStatus, "Succeeded")) {
			audit.check(op.entry, fmt.Errorf("%s: recording lacks successful Azure-AsyncOperation completion", op.op.ID))
		}
		if op.resultURL != "" && !op.resultSeen && (op.op.ID != "EdgeActionVersions_SwapDefault" || op.statusURL == "") {
			audit.check(op.entry, fmt.Errorf("%s: recording lacks a completed final Location response", op.op.ID))
		}
	}
	return audit
}

func sameRecordedURL(requestURL, headerURL string) bool {
	// The proxy redacts management.azure.com in RequestUri, but the ARM
	// sanitizers preserve LRO header authorities. Normalize only that known alias.
	request, err := url.Parse(requestURL)
	if err != nil {
		return false
	}
	if strings.EqualFold(request.Host, "Sanitized.azure.com") {
		request.Host = "management.azure.com"
	}
	return sameURL(request.String(), headerURL)
}

func isResourceGroupRoute(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return len(parts) == 4 && strings.EqualFold(parts[0], "subscriptions") && parts[1] != "" &&
		strings.EqualFold(parts[2], "resourceGroups") && parts[3] != ""
}

func TestContractRecordingAuditContinues(t *testing.T) {
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	versionURL := "https://management.azure.com" + versionPath + "?api-version=" + apiVersion
	entries := []contractEntry{
		{RequestMethod: "GET", RequestURI: versionURL, StatusCode: 200,
			ResponseBody: json.RawMessage(strings.Replace(versionJSON, `"file"`, `42`, 1))},
		{RequestMethod: "GET", RequestURI: versionURL, StatusCode: 200,
			ResponseBody: json.RawMessage(strings.Replace(versionJSON, `"global"`, `false`, 1))},
		{RequestMethod: "DELETE", RequestURI: versionURL, StatusCode: 200, ResponseBody: json.RawMessage(`{"unexpected":true}`)},
		{RequestMethod: "GET", RequestURI: "https://management.azure.com/unknown", StatusCode: 200, ResponseBody: json.RawMessage(`{}`)},
	}
	audit := auditRecording(spec, entries)
	require.Len(t, audit.issues, 4)
	assert.Equal(t, 4, audit.checked)
	for i, issue := range audit.issues {
		assert.Equal(t, i+1, issue.entry)
	}
	assert.ErrorContains(t, audit.issues[0].err, "deploymentType")
	assert.ErrorContains(t, audit.issues[1].err, "location")
	assert.ErrorContains(t, audit.issues[2].err, "does not declare")
	assert.ErrorContains(t, audit.issues[3].err, "no matching")
}

func TestContractRecordingLROAssociation(t *testing.T) {
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	const statusURL = "https://management.azure.com/operations/one?api-version=" + apiVersion
	const resultURL = "https://management.azure.com/results/one?api-version=" + apiVersion
	for _, action := range []string{"deployVersionCode", "swapDefault"} {
		t.Run(action, func(t *testing.T) {
			body := versionJSON
			if action == "deployVersionCode" {
				body = `{"status":"Succeeded"}`
			}
			entries := []contractEntry{
				{RequestMethod: "POST", RequestURI: "https://management.azure.com" + versionPath + "/" + action + "?api-version=" + apiVersion,
					StatusCode: 202, ResponseHeaders: map[string]recordedHeader{"Azure-AsyncOperation": {statusURL}, "Location": {resultURL}}},
				{RequestMethod: "GET", RequestURI: statusURL, StatusCode: 200, ResponseBody: json.RawMessage(`{"status":"Succeeded","error":42,"unmodeled":false}`)},
				{RequestMethod: "GET", RequestURI: resultURL, StatusCode: 200, ResponseBody: json.RawMessage(body)},
				{RequestMethod: "GET", RequestURI: "https://management.azure.com" + versionPath + "?api-version=" + apiVersion, StatusCode: 200, ResponseBody: json.RawMessage(versionJSON)},
			}
			audit := auditRecording(spec, entries)
			assert.Equal(t, 4, audit.checked)
			if action == "swapDefault" {
				assert.Empty(t, audit.issues)
				assert.Equal(t, 1, audit.exceptions)
				entries[3].ResponseBody = json.RawMessage(`{}`)
				audit = auditRecording(spec, entries)
				require.Len(t, audit.issues, 1)
				assert.Equal(t, 4, audit.issues[0].entry, "standalone GET is never covered by swap exception")
			} else {
				require.Len(t, audit.issues, 1)
				assert.Equal(t, 3, audit.issues[0].entry)
				assert.ErrorContains(t, audit.issues[0].err, "required")
			}
			audit = auditRecording(spec, entries[:1])
			assert.NotEmpty(t, audit.issues, "missing monitor/final requests must fail the recording audit")
		})
	}
}

func TestContractRecordingSwapLocationCompletion(t *testing.T) {
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	const resultURL = "https://management.azure.com/results/one?api-version=" + apiVersion
	for _, tc := range []struct {
		name     string
		statuses []int
		complete bool
	}{
		{"initial-only", nil, false},
		{"pending-only", []int{http.StatusAccepted, http.StatusAccepted}, false},
		{"completed", []int{http.StatusOK}, true},
		{"pending-then-completed", []int{http.StatusAccepted, http.StatusOK}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			entries := []contractEntry{
				{RequestMethod: "POST", RequestURI: "https://management.azure.com" + versionPath + "/swapDefault?api-version=" + apiVersion,
					StatusCode: http.StatusAccepted, ResponseHeaders: map[string]recordedHeader{"Location": {resultURL}}},
			}
			for _, status := range tc.statuses {
				entry := contractEntry{RequestMethod: "GET", RequestURI: resultURL, StatusCode: status}
				if status == http.StatusOK {
					entry.ResponseBody = json.RawMessage(versionJSON)
				}
				entries = append(entries, entry)
			}
			audit := auditRecording(spec, entries)
			assert.Equal(t, len(entries), audit.checked)
			if tc.complete {
				assert.Empty(t, audit.issues)
				assert.Equal(t, 1, audit.exceptions, "retain the HTTP-200-with-body schema exception")
			} else {
				require.Len(t, audit.issues, 1)
				assert.Equal(t, 1, audit.issues[0].entry)
				assert.ErrorContains(t, audit.issues[0].err, "lacks a completed final Location response")
				assert.Zero(t, audit.exceptions)
			}
		})
	}
}

func TestContractRecordingSwapMonitorWithoutLocationResult(t *testing.T) {
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	const statusURL = "https://management.azure.com/operations/one?api-version=" + apiVersion
	const resultURL = "https://management.azure.com/results/one?api-version=" + apiVersion
	for _, status := range []string{"Succeeded", "Running", ""} {
		t.Run("status="+status, func(t *testing.T) {
			entries := []contractEntry{
				{RequestMethod: "POST", RequestURI: "https://management.azure.com" + versionPath + "/swapDefault?api-version=" + apiVersion,
					StatusCode: http.StatusAccepted, ResponseHeaders: map[string]recordedHeader{"Azure-AsyncOperation": {statusURL}, "Location": {resultURL}}},
			}
			if status != "" {
				entries = append(entries, contractEntry{RequestMethod: "GET", RequestURI: statusURL,
					StatusCode: http.StatusOK, ResponseBody: json.RawMessage(fmt.Sprintf(`{"status":%q}`, status))})
			}
			audit := auditRecording(spec, entries)
			if status == "Succeeded" {
				assert.Empty(t, audit.issues, "a successful status monitor proves swap completion")
			} else {
				require.Len(t, audit.issues, 1)
				assert.ErrorContains(t, audit.issues[0].err, "lacks successful Azure-AsyncOperation completion")
			}
		})
	}
}

func TestContractRecordingFailedLocation(t *testing.T) {
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	const statusURL = "https://management.azure.com/operations/one?api-version=" + apiVersion
	const resultURL = "https://management.azure.com/results/one?api-version=" + apiVersion
	for _, tc := range []struct {
		name, body, schemaFailure string
		status                    int
	}{
		{"not-found", `{"error":{"code":"NotFound","message":"Missing result"}}`, "", http.StatusNotFound},
		{"server-error", `{"error":{"code":"InternalError","message":"Failed result"}}`, "", http.StatusInternalServerError},
		{"invalid-error-body", `{"error":{"code":42}}`, "error.code", http.StatusNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			audit := auditRecording(spec, []contractEntry{
				{RequestMethod: "POST", RequestURI: "https://management.azure.com" + versionPath + "/deployVersionCode?api-version=" + apiVersion,
					StatusCode: 202, ResponseHeaders: map[string]recordedHeader{"Azure-AsyncOperation": {statusURL}, "Location": {resultURL}}},
				{RequestMethod: "GET", RequestURI: statusURL, StatusCode: 200, ResponseBody: json.RawMessage(`{"status":"Succeeded"}`)},
				{RequestMethod: "GET", RequestURI: resultURL, StatusCode: tc.status, ResponseBody: json.RawMessage(tc.body)},
			})
			require.Len(t, audit.issues, 2)
			assert.Equal(t, 3, audit.issues[0].entry)
			assert.ErrorContains(t, audit.issues[0].err, fmt.Sprintf("non-success HTTP %d", tc.status))
			if tc.schemaFailure != "" {
				assert.ErrorContains(t, audit.issues[0].err, tc.schemaFailure)
			}
			assert.Equal(t, 1, audit.issues[1].entry)
			assert.ErrorContains(t, audit.issues[1].err, "lacks a completed final Location")
		})
	}
}

func TestContractRecordingAuditRetainsInvalidInitialLRO(t *testing.T) {
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	const statusURL = "https://management.azure.com/operations/one?api-version=" + apiVersion
	const resultURL = "https://management.azure.com/results/one?api-version=" + apiVersion
	audit := auditRecording(spec, []contractEntry{
		{RequestMethod: "PUT", RequestURI: "https://management.azure.com" + versionPath + "?api-version=" + apiVersion,
			StatusCode: 201, ResponseBody: json.RawMessage(`{}`), ResponseHeaders: map[string]recordedHeader{"Azure-AsyncOperation": {statusURL}, "Location": {resultURL}}},
		{RequestMethod: "GET", RequestURI: statusURL, StatusCode: 200, ResponseBody: json.RawMessage(`{"status":"Succeeded"}`)},
		{RequestMethod: "GET", RequestURI: resultURL, StatusCode: 200, ResponseBody: json.RawMessage(`{"location":false}`)},
	})
	require.Len(t, audit.issues, 2)
	assert.Equal(t, 1, audit.issues[0].entry)
	assert.Equal(t, 3, audit.issues[1].entry)
	assert.ErrorContains(t, audit.issues[1].err, "EdgeActionVersions_Create")
}

func TestContractRecordingBodyAndFixtures(t *testing.T) {
	assert.True(t, sameRecordedURL("https://Sanitized.azure.com/result?a=1", "https://management.azure.com/result?a=1"))
	for _, target := range []string{
		"https://example.test/result?a=1", "http://management.azure.com/result?a=1",
		"https://management.azure.com/other?a=1", "https://management.azure.com/result?a=2",
	} {
		assert.False(t, sameRecordedURL("https://Sanitized.azure.com/result?a=1", target))
	}
	for _, tc := range []struct{ raw, expected string }{
		{"null", ""}, {`""`, ""}, {`"null"`, "null"}, {`"{\"status\":\"Succeeded\"}"`, `{"status":"Succeeded"}`},
		{`{"status":"Succeeded"}`, `{"status":"Succeeded"}`},
	} {
		body, err := (contractEntry{ResponseBody: json.RawMessage(tc.raw)}).body()
		require.NoError(t, err)
		assert.Equal(t, tc.expected, string(body))
	}
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	audit := auditRecording(spec, []contractEntry{
		{RequestMethod: "DELETE", RequestURI: "https://management.azure.com/subscriptions/" + subscriptionID + "/resourceGroups/rg?api-version=2021-04-01",
			StatusCode: 202, ResponseHeaders: map[string]recordedHeader{"Location": {"https://management.azure.com/fixture-monitor"}}},
		{RequestMethod: "GET", RequestURI: "https://management.azure.com/fixture-monitor", StatusCode: 200},
	})
	assert.Empty(t, audit.issues)
	assert.Equal(t, 2, audit.fixtures)
	assert.Zero(t, audit.checked)
}

func recordingDirectory(ctx context.Context) (string, error) {
	root, err := filepath.Abs(filepath.Join("..", "..", "..", ".."))
	if err != nil {
		return "", err
	}
	executable := "Azure.Sdk.Tools.TestProxy"
	if runtime.GOOS == "windows" {
		executable += ".exe"
	}
	config, err := filepath.Abs("assets.json")
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, filepath.Join(root, ".proxy", executable), "config", "locate", "-a", config).Output()
	if err != nil {
		return "", fmt.Errorf("locate recording assets with test proxy: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	assetRoot := filepath.FromSlash(strings.TrimSpace(lines[len(lines)-1]))
	if !filepath.IsAbs(assetRoot) {
		return "", fmt.Errorf("test proxy did not return an absolute asset directory")
	}
	return filepath.Join(assetRoot, filepath.FromSlash(pathToPackage), "recordings"), nil
}

func TestRecordingResponseContracts(t *testing.T) {
	if recording.GetRecordMode() == recording.LiveMode {
		t.Skip("recording response audit applies only in record/playback modes")
	}
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	directory, err := recordingDirectory(t.Context())
	require.NoError(t, err)
	var files []string
	err = filepath.WalkDir(directory, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.EqualFold(filepath.Ext(path), ".json") {
			files = append(files, path)
		}
		return nil
	})
	require.NoError(t, err, "service recordings must exist; no live fallback")
	require.NotEmpty(t, files, "no service recordings found")
	for _, path := range files {
		name, err := filepath.Rel(directory, path)
		require.NoError(t, err)
		t.Run(filepath.ToSlash(name), func(t *testing.T) {
			data, err := os.ReadFile(path)
			require.NoError(t, err)
			var saved contractRecording
			// Do not print a JSON parse error containing recording values.
			require.True(t, json.Unmarshal(bytes.TrimPrefix(data, []byte{0xef, 0xbb, 0xbf}), &saved) == nil, "invalid recording JSON")
			require.NotEmpty(t, saved.Entries)
			audit := auditRecording(spec, saved.Entries)
			assert.Positive(t, audit.checked, "recording contains no Edge Actions responses")
			for _, issue := range audit.issues {
				assert.NoError(t, issue.err, "recording entry %d", issue.entry)
			}
			t.Logf("checked=%d fixture responses excluded=%d swap Location exceptions=%d operations=%v",
				audit.checked, audit.fixtures, audit.exceptions, audit.operations)
		})
	}
}
