// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armedgeactions_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizerRegistration(t *testing.T) {
	if recording.GetRecordMode() == recording.LiveMode {
		t.Skip("proxy sanitizer registration applies only in record/playback modes")
	}
	type registration struct {
		kind                string
		value, regex, group string
	}
	var registrations []registration
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Value string
			Regex string
			Group string `json:"groupForReplace"`
		}
		if !assert.NoError(t, json.NewDecoder(r.Body).Decode(&body)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		registrations = append(registrations, registration{r.Header.Get("x-abstraction-identifier"), body.Value, body.Regex, body.Group})
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	_, port, err := net.SplitHostPort(server.Listener.Addr().String())
	require.NoError(t, err)
	portNumber, err := strconv.Atoi(port)
	require.NoError(t, err)
	secret := "fake.secret+[regex]"
	require.NoError(t, registerSanitizers(func(key string) string {
		if key == "AZURE_CLIENT_SECRET" {
			return secret
		}
		return ""
	}, &recording.RecordingOptions{ProxyPort: portNumber}))
	var general []registration
	for _, registration := range registrations {
		if registration.kind == "GeneralRegexSanitizer" {
			general = append(general, registration)
		}
	}
	require.Len(t, general, 2, "empty environment values must not register a match-all sanitizer")
	assert.Equal(t, regexp.QuoteMeta(secret), general[0].regex)
	assert.Equal(t, subscriptionID, general[0].value)
	assert.Equal(t, pollingSecretRegex, general[1].regex)
	assert.Equal(t, "secret", general[1].group)
}

type recordedHeader []string

func (h *recordedHeader) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '[' {
		var values []string
		if err := json.Unmarshal(data, &values); err != nil {
			return err
		}
		*h = values
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*h = []string{value}
	return nil
}

type proxyRecording struct {
	Entries []struct {
		RequestURI      string `json:"RequestUri"`
		RequestHeaders  map[string]recordedHeader
		ResponseHeaders map[string]recordedHeader
		ResponseBody    struct {
			ID           string
			Name         string
			SKU          struct{ Name string }
			PollingURL   string
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			IDToken      string `json:"id_token"`
			Secret       string
			Identity     struct {
				TenantID     string `json:"tenantId"`
				ClientID     string `json:"clientId"`
				PrincipalID  string `json:"principalId"`
				ClientSecret string `json:"client_secret"`
			}
			SystemData struct{ CreatedBy, LastModifiedBy string }
		}
	}
}

func TestSanitizerProxyRoundTrip(t *testing.T) {
	if recording.GetRecordMode() == recording.LiveMode {
		t.Skip("local proxy record/playback regression does not apply to live mode")
	}
	require.NoError(t, registerSanitizers(func(key string) string {
		if key == "AZURE_CLIENT_SECRET" {
			return "fake.secret+[regex]"
		}
		return ""
	}, &recording.RecordingOptions{UseHTTPS: true}))
	for _, query := range []string{
		"?api-version=2026-10-01&t=fake-time&c=fake-context%26value&s=fake-signature%3D&h=fake-hash&count=2&st=unchanged",
		"?api-version=2026-10-01&C=fake-context&c=&count=2&st=unchanged",
	} {
		for _, recordedURL := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/recorded-url=%t", query, recordedURL), func(t *testing.T) {
				testSanitizerRoundTrip(t, query, recordedURL)
			})
		}
	}
}

func testSanitizerRoundTrip(t *testing.T, query string, recordedURL bool) {
	t.Helper()
	var calls atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		pollURL := "http://" + r.Host + "/operations/one" + query + "#fragment"
		w.Header().Set("Content-Type", "application/json")
		for _, header := range []string{"Location", "Azure-AsyncOperation", "operation-location"} {
			w.Header().Set(header, pollURL)
		}
		for _, header := range traceHeaders {
			w.Header().Set(header, "fake-trace")
		}
		w.Header().Set("Set-Cookie", "fake-cookie")
		body := map[string]any{
			"id": actionPath, "name": "test-action", "sku": map[string]string{"name": "Standard"},
			"pollingUrl": pollURL, "access_token": "fake-access", "refresh_token": "fake-refresh",
			"id_token": "fake-id", "secret": "fake.secret+[regex]",
			"identity":   map[string]string{"tenantId": "fake-tenant", "clientId": "fake-client", "principalId": "fake-principal", "client_secret": "fake-unconfigured-secret"},
			"systemData": map[string]string{"createdBy": "fake-creator", "lastModifiedBy": "fake-modifier"},
		}
		assert.NoError(t, json.NewEncoder(w).Encode(body))
	}))
	defer upstream.Close()
	client, err := recording.GetHTTPClient(t)
	require.NoError(t, err)
	client.Timeout = 30 * time.Second
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	defer client.CloseIdleConnections()
	dummy, err := http.NewRequest("GET", "https://example.test/", nil)
	require.NoError(t, err)
	proxyURL := (&recording.RecordingOptions{UseHTTPS: true}).ReplaceAuthority(t, dummy).URL
	base := proxyURL.Scheme + "://" + proxyURL.Host
	path := filepath.Join(t.TempDir(), "sanitizers.json")
	start := func(mode string) string {
		body, err := json.Marshal(map[string]string{"x-recording-file": path})
		require.NoError(t, err)
		response, err := client.Post(base+"/"+mode+"/start", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer response.Body.Close()
		data, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode, "%s", data)
		id := response.Header.Get(recording.IDHeader)
		require.NotEmpty(t, id)
		return id
	}
	stop := func(mode, id string) {
		req, err := http.NewRequest("POST", base+"/"+mode+"/stop", strings.NewReader("{}"))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(recording.IDHeader, id)
		req.Header.Set("x-recording-save", "true")
		response, err := client.Do(req)
		require.NoError(t, err)
		defer response.Body.Close()
		data, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode, "%s", data)
	}
	send := func(mode, id, query string) {
		req, err := http.NewRequest("GET", base+"/operations/one"+query, nil)
		require.NoError(t, err)
		req.Header.Set(recording.IDHeader, id)
		req.Header.Set(recording.ModeHeader, mode)
		req.Header.Set(recording.UpstreamURIHeader, upstream.URL)
		req.Header.Set("Authorization", "Bearer fake-authorization")
		req.Header.Set("Cookie", "fake-request-cookie")
		response, err := client.Do(req)
		require.NoError(t, err)
		defer response.Body.Close()
		data, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode, "%s", data)
		assert.Contains(t, string(data), "test-action")
	}
	recordID := start("record")
	recordStopped := false
	t.Cleanup(func() {
		if !recordStopped {
			stop("record", recordID)
		}
	})
	send("record", recordID, query)
	stop("record", recordID)
	recordStopped = true
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.NotContains(t, string(data), "fake-")
	assert.NotContains(t, string(data), "fake.secret")
	var saved proxyRecording
	require.NoError(t, json.Unmarshal(bytes.TrimPrefix(data, []byte{0xef, 0xbb, 0xbf}), &saved))
	require.Len(t, saved.Entries, 1)
	entry := saved.Entries[0]
	assert.Equal(t, actionPath, entry.ResponseBody.ID)
	assert.Equal(t, "test-action", entry.ResponseBody.Name)
	assert.Equal(t, "Standard", entry.ResponseBody.SKU.Name)
	assert.Equal(t, "access_token", entry.ResponseBody.AccessToken)
	assert.Equal(t, "refresh_token", entry.ResponseBody.RefreshToken)
	assert.Equal(t, "id_token", entry.ResponseBody.IDToken)
	assert.Equal(t, subscriptionID, entry.ResponseBody.Secret)
	assert.Equal(t, recording.SanitizedValue, entry.ResponseBody.Identity.TenantID)
	assert.Equal(t, recording.SanitizedValue, entry.ResponseBody.Identity.ClientID)
	assert.Equal(t, recording.SanitizedValue, entry.ResponseBody.Identity.PrincipalID)
	assert.Equal(t, recording.SanitizedValue, entry.ResponseBody.Identity.ClientSecret)
	assert.Equal(t, recording.SanitizedValue, entry.ResponseBody.SystemData.CreatedBy)
	assert.Equal(t, recording.SanitizedValue, entry.ResponseBody.SystemData.LastModifiedBy)
	requestHeaders := http.Header{}
	for key, value := range entry.RequestHeaders {
		requestHeaders[http.CanonicalHeaderKey(key)] = value
	}
	assert.Equal(t, "cookie;", requestHeaders.Get("Cookie"))
	headers := http.Header{}
	for key, value := range entry.ResponseHeaders {
		headers[http.CanonicalHeaderKey(key)] = value
	}
	for _, header := range traceHeaders {
		assert.Equal(t, recording.SanitizedValue, headers.Get(header))
	}
	assert.Equal(t, "[set-cookie;]", headers.Get("Set-Cookie"))
	original, err := url.ParseQuery(strings.TrimPrefix(query, "?"))
	require.NoError(t, err)
	for key, values := range original {
		if strings.Contains("tcsh", strings.ToLower(key)) && len(key) == 1 && values[0] != "" {
			original.Set(key, recording.SanitizedValue)
		}
	}
	urls := []string{entry.RequestURI, entry.ResponseBody.PollingURL}
	for _, header := range []string{"Location", "Azure-AsyncOperation", "operation-location"} {
		urls = append(urls, headers.Get(header))
	}
	for _, value := range urls {
		parsed, err := url.Parse(value)
		require.NoError(t, err)
		assert.Equal(t, "/operations/one", parsed.Path)
		assert.Equal(t, original, parsed.Query(), "preserve query keys, API version, and unrelated values")
	}
	assert.True(t, strings.HasSuffix(entry.ResponseBody.PollingURL, "#fragment"))
	upstream.Close()
	playbackQuery := strings.ReplaceAll(query, "fake-", "different-")
	if recordedURL {
		pollURL, err := url.Parse(headers.Get("Azure-AsyncOperation"))
		require.NoError(t, err)
		playbackQuery = "?" + pollURL.RawQuery
	}
	playbackID := start("playback")
	defer stop("playback", playbackID)
	send("playback", playbackID, playbackQuery)
	assert.Equal(t, int32(1), calls.Load(), "playback must not contact the upstream")
}
