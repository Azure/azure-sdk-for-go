// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// GetRawJSON returns an independent copy of a JSON response's complete payload.
// Use policy.WithCaptureResponse on the begin and polling context to capture it.
func GetRawJSON(response *http.Response) (json.RawMessage, error) {
	if response == nil {
		return nil, errors.New("response is required")
	}
	data, err := runtime.Payload(response)
	if err != nil {
		return nil, err
	}
	if !json.Valid(data) {
		return nil, errors.New("response does not contain valid JSON")
	}
	return append(json.RawMessage(nil), data...), nil
}

// GetUsageDetails returns usage from a completed LRO or inline response.
// Missing or null usage returns nil; malformed present usage returns an error.
func GetUsageDetails(response *http.Response) (*UsageDetails, error) {
	data, err := GetRawJSON(response)
	if err != nil {
		return nil, err
	}
	if !bytes.HasPrefix(bytes.TrimSpace(data), []byte("{")) {
		return nil, errors.New("usage response must be a JSON object")
	}
	var envelope struct {
		Status string          `json:"status"`
		Usage  json.RawMessage `json:"usage"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}
	if envelope.Status != "" && envelope.Status != "Succeeded" && envelope.Status != "Failed" && envelope.Status != "Canceled" {
		return nil, nil
	}
	if len(envelope.Usage) == 0 || string(envelope.Usage) == "null" {
		return nil, nil
	}
	var usage UsageDetails
	if err := json.Unmarshal(envelope.Usage, &usage); err != nil {
		return nil, fmt.Errorf("decoding analysis usage: %w", err)
	}
	return &usage, nil
}

// GetOperationID reads the operation identifier from the initial response's
// Operation-Location header, or a completed analysis envelope's id property.
func GetOperationID(response *http.Response) (string, error) {
	if response == nil {
		return "", errors.New("response is required")
	}
	if location := response.Header.Get("Operation-Location"); location != "" {
		u, err := url.Parse(location)
		if err != nil || !u.IsAbs() || u.Host == "" {
			return "", errors.New("invalid Operation-Location header")
		}
		path := strings.TrimRight(u.Path, "/")
		id := path[strings.LastIndex(path, "/")+1:]
		if id == "" {
			return "", errors.New("Operation-Location is missing an operation ID")
		}
		return id, nil
	}
	data, err := GetRawJSON(response)
	if err != nil {
		return "", err
	}
	var envelope struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return "", err
	}
	if envelope.ID == "" {
		return "", errors.New("operation ID is absent from the response")
	}
	return envelope.ID, nil
}

// GetAnalysisResult decodes a captured completed operation without losing the
// original raw response. It accepts the observed service KeyFrameTimesMs alias.
func GetAnalysisResult(response *http.Response) (*AnalysisResult, error) {
	data, err := GetRawJSON(response)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Status string          `json:"status"`
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}
	if envelope.Status != "Succeeded" || len(envelope.Result) == 0 || string(envelope.Result) == "null" {
		return nil, errors.New("analysis must have a successful status and a result")
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(envelope.Result, &object); err != nil {
		return nil, err
	}
	var contents []json.RawMessage
	if err := json.Unmarshal(object["contents"], &contents); err != nil {
		return nil, err
	}
	for i, rawContent := range contents {
		var content map[string]json.RawMessage
		if err := json.Unmarshal(rawContent, &content); err != nil {
			return nil, err
		}
		if string(content["kind"]) != `"audioVisual"` {
			continue
		}
		times, err := firstKeyFrameTimes(rawContent)
		if err != nil {
			return nil, err
		}
		if times != nil {
			content["keyFrameTimesMs"] = times
			contents[i], err = json.Marshal(content)
			if err != nil {
				return nil, err
			}
		}
	}
	object["contents"], err = json.Marshal(contents)
	if err != nil {
		return nil, err
	}
	normalized, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	var result AnalysisResult
	if err := json.Unmarshal(normalized, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func firstKeyFrameTimes(data []byte) (json.RawMessage, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}
	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			return nil, err
		}
		if (key == "keyFrameTimesMs" || key == "KeyFrameTimesMs") && string(value) != "null" {
			return value, nil
		}
	}
	return nil, nil
}
