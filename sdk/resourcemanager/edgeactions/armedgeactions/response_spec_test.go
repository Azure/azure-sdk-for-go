// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armedgeactions_test

//go:generate go run testdata/refresh_spec.go

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const apiVersion = "2026-10-01"

type responseSpec struct {
	Source      map[string]string
	Paths       map[string]map[string]*specOperation
	Definitions map[string]*responseSchema
}

type specOperation struct {
	ID        string `json:"operationId"`
	LRO       bool   `json:"x-ms-long-running-operation"`
	Responses map[string]struct {
		Schema  *responseSchema
		Headers map[string]*responseSchema
	}
}

// This intentionally supports only the schema vocabulary in the pinned response
// snapshot. Unknown validation keywords fail loading rather than passing unchecked.
type responseSchema struct {
	Ref                  string `json:"$ref"`
	Type                 string
	Format               string
	Enum                 []any
	Required             []string
	Properties           map[string]*responseSchema
	Items                *responseSchema
	AllOf                []*responseSchema
	AdditionalProperties *responseSchema
}

func (s *responseSchema) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	for key := range fields {
		switch key {
		case "$ref", "type", "format", "enum", "required", "properties", "items", "allOf", "additionalProperties",
			"description", "title", "readOnly":
		default:
			if !strings.HasPrefix(key, "x-ms-") {
				return fmt.Errorf("unsupported response schema keyword %q", key)
			}
		}
	}
	type plain responseSchema
	return json.Unmarshal(data, (*plain)(s))
}

func loadResponseSpec() (*responseSpec, error) {
	data, err := os.ReadFile("testdata/response_spec.json")
	if err != nil {
		return nil, err
	}
	var spec responseSpec
	if err = json.Unmarshal(data, &spec); err != nil {
		return nil, err
	}
	location, err := os.ReadFile("tsp-location.yaml")
	if err != nil {
		return nil, err
	}
	for _, key := range []string{"repo", "commit", "directory"} {
		match := regexp.MustCompile(`(?m)^` + key + `:\s*(\S+)\s*$`).FindSubmatch(location)
		if len(match) != 2 || strings.Trim(string(match[1]), `"'`) != spec.Source[key] {
			return nil, fmt.Errorf("response snapshot %s is stale; run go generate", key)
		}
	}
	return &spec, nil
}

func (s *responseSpec) operation(method, address string) (*specOperation, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("invalid operation URL")
	}
	segments := strings.Split(strings.Trim(u.Path, "/"), "/")
	for path, methods := range s.Paths {
		template := strings.Split(strings.Trim(path, "/"), "/")
		if len(template) != len(segments) {
			continue
		}
		matches := true
		for i, part := range template {
			if strings.HasPrefix(part, "{") && segments[i] != "" {
				continue
			}
			if !strings.EqualFold(part, segments[i]) {
				matches = false
				break
			}
		}
		if !matches {
			continue
		}
		op := methods[strings.ToLower(method)]
		if op == nil {
			return nil, fmt.Errorf("undocumented Edge Actions HTTP method")
		}
		if u.Query().Get("api-version") != apiVersion {
			return nil, fmt.Errorf("%s: API version mismatch", op.ID)
		}
		return op, nil
	}
	return nil, nil // Resource group setup uses a separate ARM API.
}

func (s *responseSpec) validate(op *specOperation, status int, body []byte, headers http.Header) error {
	response, ok := op.Responses[strconv.Itoa(status)]
	if !ok && status >= 400 && status <= 599 {
		response, ok = op.Responses["default"]
	}
	if !ok {
		return fmt.Errorf("%s: HTTP %d is not a documented response", op.ID, status)
	}
	var failures []error
	headerNames := make([]string, 0, len(response.Headers))
	for key := range response.Headers {
		headerNames = append(headerNames, key)
	}
	sort.Strings(headerNames)
	for _, key := range headerNames {
		schema := response.Headers[key]
		if value := headers.Get(key); value != "" {
			var typed any = value
			if schema.Type == "integer" {
				number, err := strconv.ParseInt(value, 10, 32)
				if err != nil {
					failures = append(failures, fmt.Errorf("%s must be an integer", key))
					continue
				}
				typed = float64(number)
			}
			if err := s.validateSchema(schema, typed, key); err != nil {
				failures = append(failures, err)
			}
		}
	}
	if err := s.validateBody(response.Schema, body, status); err != nil {
		failures = append(failures, err)
	}
	if err := errors.Join(failures...); err != nil {
		return fmt.Errorf("%s: %w", op.ID, err)
	}
	return nil
}

func (s *responseSpec) validateFinal(op *specOperation, status int, body []byte, headers http.Header) error {
	var statusErr error
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		statusErr = fmt.Errorf("%s: final Location returned non-success HTTP %d", op.ID, status)
	}
	return errors.Join(statusErr, s.validate(op, status, body, headers))
}

func (s *responseSpec) validateBody(schema *responseSchema, body []byte, status int) error {
	if schema == nil {
		if len(bytes.TrimSpace(body)) != 0 {
			return fmt.Errorf("HTTP %d does not declare a response body", status)
		}
		return nil
	}
	var value any
	if len(bytes.TrimSpace(body)) == 0 {
		return fmt.Errorf("HTTP %d requires a response body", status)
	}
	if json.Unmarshal(body, &value) != nil {
		return fmt.Errorf("response body is not valid JSON")
	}
	return s.validateSchema(schema, value, "$")
}

func (s *responseSpec) validateSchema(schema *responseSchema, value any, path string) error {
	fail := func(rule string) error { return fmt.Errorf("%s violates OpenAPI %s", path, rule) }
	var failures []error
	if schema.Ref != "" {
		ref, ok := s.Definitions[strings.TrimPrefix(schema.Ref, "#/definitions/")]
		if !ok || ref == nil {
			return fail("reference")
		}
		if err := s.validateSchema(ref, value, path); err != nil {
			failures = append(failures, err)
		}
	}
	for _, part := range schema.AllOf {
		if err := s.validateSchema(part, value, path); err != nil {
			failures = append(failures, err)
		}
	}
	switch schema.Type {
	case "":
	case "object":
		if _, ok := value.(map[string]any); !ok {
			return fail("type")
		}
	case "array":
		if _, ok := value.([]any); !ok {
			return fail("type")
		}
	case "string":
		if _, ok := value.(string); !ok {
			return fail("type")
		}
	case "integer", "number":
		number, ok := value.(float64)
		if !ok || (schema.Type == "integer" && math.Trunc(number) != number) {
			return fail("type")
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fail("type")
		}
	default:
		return fail("unsupported type")
	}
	if len(schema.Enum) > 0 {
		found := false
		for _, allowed := range schema.Enum {
			found = found || reflect.DeepEqual(value, allowed)
		}
		if !found {
			failures = append(failures, fail("enum"))
		}
	}
	if text, ok := value.(string); ok {
		switch schema.Format {
		case "":
		case "date-time":
			if _, err := time.Parse(time.RFC3339Nano, text); err != nil {
				failures = append(failures, fail("format"))
			}
		case "uri":
			u, err := url.Parse(text)
			if err != nil || !u.IsAbs() {
				failures = append(failures, fail("format"))
			}
		case "arm-id":
			if !strings.HasPrefix(text, "/") {
				failures = append(failures, fail("format"))
			}
		default:
			return fail("unsupported format")
		}
	}
	if object, ok := value.(map[string]any); ok {
		for _, key := range schema.Required {
			if deferredValidationStatus(schema, key) {
				continue
			}
			if _, present := object[key]; !present {
				failures = append(failures, fail("required ("+key+")"))
			}
		}
		keys := make([]string, 0, len(object))
		for key := range object {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if deferredValidationStatus(schema, key) {
				continue
			}
			field := object[key]
			if property, ok := schema.Properties[key]; ok {
				if err := s.validateSchema(property, field, path+"."+key); err != nil {
					failures = append(failures, err)
				}
			} else if schema.AdditionalProperties != nil {
				// Do not include arbitrary user-provided property names in diagnostics.
				if err := s.validateSchema(schema.AdditionalProperties, field, path+".*"); err != nil {
					failures = append(failures, err)
				}
			}
		}
	}
	if items, ok := value.([]any); ok && schema.Items != nil {
		for i, item := range items {
			if err := s.validateSchema(schema.Items, item, fmt.Sprintf("%s[%d]", path, i)); err != nil {
				failures = append(failures, err)
			}
		}
	}
	return errors.Join(failures...)
}

func deferredValidationStatus(schema *responseSchema, key string) bool {
	// User-approved temporary exception pending the optional, non-nullable RP fix.
	property := schema.Properties[key]
	return key == "validationStatus" && property != nil &&
		strings.HasSuffix(property.Ref, "_EdgeActionVersionValidationStatus")
}

func validatePollingURL(address string) error {
	u, err := url.Parse(address)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.Fragment != "" {
		return fmt.Errorf("LRO header must contain an absolute HTTPS URL without credentials or fragment")
	}
	// Do not erase an advertised API version when sanitizing or replaying URLs.
	if versions, ok := u.Query()["api-version"]; ok && (len(versions) != 1 || versions[0] != apiVersion) {
		return fmt.Errorf("LRO header API version mismatch")
	}
	return nil
}

func validateStatus(status int, body []byte) (string, error) {
	if status != http.StatusOK {
		return "", fmt.Errorf("Azure-AsyncOperation: expected HTTP 200")
	}
	var object map[string]any
	if json.Unmarshal(body, &object) != nil {
		return "", fmt.Errorf("Azure-AsyncOperation: expected a status object")
	}
	state, ok := object["status"].(string)
	if !ok || strings.TrimSpace(state) == "" {
		return "", fmt.Errorf("Azure-AsyncOperation: expected a nonempty string status")
	}
	return state, nil
}

func swapLocationException(op *specOperation, status int, body []byte) bool {
	return op.ID == "EdgeActionVersions_SwapDefault" && status == http.StatusOK && len(bytes.TrimSpace(body)) > 0
}

func TestContractResponseSpec(t *testing.T) {
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	op, err := spec.operation("POST", "https://management.azure.com"+versionPath+"/deployVersionCode?api-version="+apiVersion)
	require.NoError(t, err)
	require.NotNil(t, op)
	require.NoError(t, spec.validate(op, 200, []byte(versionJSON), nil))
	for _, tc := range []struct{ name, body, failure string }{
		{"empty", "", "requires a response body"},
		{"null", "null", "type"},
		{"status", `{"status":"Succeeded"}`, "required"},
		{"properties-only", `{"deploymentType":"file","isDefaultVersion":"False"}`, "required"},
		{"missing-deployment-type", `{"location":"global","properties":{"isDefaultVersion":"False"}}`, "required"},
		{"wrong-type", strings.Replace(versionJSON, `"file"`, `42`, 1), "type"},
		{"wrong-enum", strings.Replace(versionJSON, `"file"`, `"invalid"`, 1), "enum"},
		{"wrong-date", strings.Replace(versionJSON, `"deploymentType"`, `"lastPackageUpdateTime":"invalid","deploymentType"`, 1), "format"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.ErrorContains(t, spec.validate(op, 200, []byte(tc.body), nil), tc.failure)
		})
	}
	require.ErrorContains(t, spec.validate(op, 204, nil, nil), "not a documented response")
	require.ErrorContains(t, spec.validate(op, 202, []byte("null"), nil), "does not declare")
	require.ErrorContains(t, spec.validate(op, 202, nil, http.Header{"Retry-After": {"invalid"}}), "must be an integer")
	require.ErrorContains(t, spec.validate(op, 202, nil, http.Header{"Azure-Asyncoperation": {"invalid"}}), "format")
	var schema responseSchema
	require.ErrorContains(t, json.Unmarshal([]byte(`{"minLength":1}`), &schema), "unsupported")
	_, err = spec.operation("GET", "https://management.azure.com"+versionPath+"?api-version=old")
	require.ErrorContains(t, err, "API version mismatch")
}

func TestContractValidationStatusException(t *testing.T) {
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	before, err := json.Marshal(spec)
	require.NoError(t, err)
	op, err := spec.operation("GET", "https://management.azure.com"+versionPath+"?api-version="+apiVersion)
	require.NoError(t, err)
	for _, value := range []string{`null`, `"Succeeded"`, `42`, `{}`} {
		body := strings.Replace(versionJSON, `"deploymentType"`, `"validationStatus":`+value+`,"deploymentType"`, 1)
		require.NoError(t, spec.validate(op, 200, []byte(body), nil), "field validation is explicitly deferred")
	}
	require.NoError(t, spec.validate(op, 200, []byte(versionJSON), nil), "optional field may be absent")
	body := strings.Replace(versionJSON, `"Succeeded"`, `null`, 1)
	require.ErrorContains(t, spec.validate(op, 200, []byte(body), nil), "provisioningState")
	after, err := json.Marshal(spec)
	require.NoError(t, err)
	assert.Equal(t, before, after, "exceptions must not alter the pinned contract")
}

func TestContractReportsAllResponseViolations(t *testing.T) {
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	op, err := spec.operation("PUT", "https://management.azure.com"+versionPath+"?api-version="+apiVersion)
	require.NoError(t, err)
	err = spec.validate(op, 201, []byte(`{"location":false,"properties":{"deploymentType":42,"isDefaultVersion":null}}`), http.Header{"Retry-After": {"invalid"}})
	require.Error(t, err)
	for _, field := range []string{"Retry-After", "location", "deploymentType", "isDefaultVersion"} {
		assert.Contains(t, err.Error(), field)
	}
	require.NoError(t, spec.validate(op, 400, []byte(`{"error":{"code":"BadRequest","message":"invalid input"}}`), nil))
	require.ErrorContains(t, spec.validate(op, 400, []byte(`{"error":{"code":42,"message":false}}`), nil), "type")
	require.ErrorContains(t, spec.validate(op, 204, nil, nil), "not a documented response")
}

func TestContractStatusEnvelope(t *testing.T) {
	for _, state := range []string{"Accepted", "Running", "Succeeded", "Failed", "Canceled"} {
		got, err := validateStatus(200, []byte(`{"status":"`+state+`","error":null,"operationType":"Deploy"}`))
		require.NoError(t, err)
		require.Equal(t, state, got)
	}

	for _, body := range []string{"", "null", "[]", `{}`, `{"status":0}`, `{"status":""}`, `{"status":" "}`} {
		_, err := validateStatus(200, []byte(body))
		require.Error(t, err)
	}
	for _, status := range []int{201, 202, 204, 400, 500} {
		_, err := validateStatus(status, []byte(`{"status":"Succeeded"}`))
		require.Error(t, err)
	}
	for _, address := range []string{"http://management.azure.com/result", "https:/result", "https://user:password@example.test/result", "https://example.test/result?api-version=old"} {
		require.Error(t, validatePollingURL(address))
	}
	require.NoError(t, validatePollingURL("https://management.azure.com/result?api-version="+apiVersion))
}
