// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestDefaultSanitizers(t *testing.T) {
	before := recordMode
	defer func() { recordMode = before }()
	recordMode = RecordingMode

	t.Setenv(proxyManualStartEnv, "false")
	o := *defaultOptions()
	o.insecure = true
	proxy, err := StartTestProxy("", &o)
	require.NoError(t, err)
	defer func() {
		err := StopTestProxy(proxy)
		require.NoError(t, err)
		_ = os.Remove(filepath.Join("testdata", "recordings", t.Name()+".json"))
	}()

	client, err := NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	srv, close := mock.NewTLSServer()
	defer close()

	// build a request and response containing all the values that should be sanitized by default
	fail := "FAIL"
	failSAS := strings.ReplaceAll("sv=*&sig=*&se=*&srt=*&ss=*&sp=*", "*", fail)
	q := "?sig=" + fail
	req, err := http.NewRequest(http.MethodGet, srv.URL()+q, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resOpts := []mock.ResponseOption{mock.WithStatusCode(http.StatusOK)}
	body := map[string]any{}
	for _, s := range defaultSanitizers {
		switch s.Name {
		case "BodyKeySanitizer":
			k := strings.TrimLeft(s.Body.JSONPath, "$.")
			var v any = fail
			if before, after, found := strings.Cut(k, "."); found {
				// path is e.g. $..foo.bar, so this value would be in a nested object
				k = before
				if strings.HasSuffix(k, "[*]") {
					// path is e.g. $..foo[*].bar, so this value would be in an object array
					k = strings.TrimSuffix(k, "[*]")
					v = []map[string]string{{after: fail}}
				} else {
					v = map[string]string{after: fail}
				}
			}
			body[k] = v
		case "HeaderRegexSanitizer":
			// if there's no group specified, we can generate a matching value because this sanitizer
			// performs a simple replacement (this works provided the default regex sanitizers continue
			// to follow the convention of always naming a group)
			if s.Body.GroupForReplace == "" {
				req.Header.Set(s.Body.Key, fail)
				resOpts = append(resOpts, mock.WithHeader(s.Body.Key, fail))
			}
		default:
			// handle regex sanitizers below because generating matching values is tricky
		}
	}
	// add values matching body regex sanitizers
	for i, v := range []string{
		"client_secret=" + fail + "&client_assertion=" + fail,
		strings.ReplaceAll("-----BEGIN PRIVATE KEY-----\n*\n*\n*\n-----END PRIVATE KEY-----\n", "*", fail),
		failSAS,
		strings.Join([]string{"AccessKey", "accesskey", "Accesskey", "AccountKey", "SharedAccessKey"}, "="+fail+";") + "=" + fail,
		strings.ReplaceAll("<UserDelegationKey><SignedOid>*</SignedOid><SignedTid>*</SignedTid><Value>*</Value></UserDelegationKey>", "*", fail),
		fmt.Sprintf("<PrimaryKey>%s</PrimaryKey>", failSAS),
	} {
		k := fmt.Sprint(i)
		require.NotContains(t, body, k, "test bug: body already has key %q", k)
		body[k] = v
	}
	// add values matching header regex sanitizers
	for _, h := range []string{"ServiceBusDlqSupplementaryAuthorization", "ServiceBusSupplementaryAuthorization", "SupplementaryAuthorization"} {
		req.Header.Set(h, failSAS)
	}

	// set request and response bodies
	j, err := json.Marshal(body)
	require.NoError(t, err)
	req.Body = io.NopCloser(bytes.NewReader(j))
	srv.SetResponse(append(resOpts, mock.WithBody(j))...)

	err = Start(t, packagePath, nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	err = Stop(t, nil)
	require.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		t.Fatal(string(b))
	}

	b, err := os.ReadFile(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	if bytes.Contains(b, []byte(fail)) {
		var buf bytes.Buffer
		require.NoError(t, json.Indent(&buf, b, "", "  "))
		t.Fatalf("%q shouldn't appear in this recording:\n%s%q shouldn't appear in the above recording", fail, buf.String(), fail)
	}
}
