// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding

import (
	"encoding/json"
	"errors"
)

// AnalysisInput identifies a URL or base64-encoded binary input.
// Set exactly one of URL or Data. Data is encoded as base64 by encoding/json.
type AnalysisInput struct {
	URL      string `json:"url,omitempty"`
	Data     []byte `json:"data,omitempty"`
	Name     string `json:"name,omitempty"`
	MIMEType string `json:"mimeType,omitempty"`
}

// AnalyzeRequest contains inputs to the JSON analysis endpoints.
type AnalyzeRequest struct {
	Inputs []AnalysisInput `json:"inputs"`
}

// AnalyzeResponse retains the completed service envelope without dropping unknown properties.
// Result and Usage are JSON subtrees. Use encoding/json to decode only the properties
// required by an application, without losing integer precision or polymorphic fields.
type AnalyzeResponse struct {
	ID      string
	Status  string
	Result  json.RawMessage
	Usage   json.RawMessage
	RawJSON json.RawMessage
}

// UnmarshalJSON decodes a completed successful analysis envelope.
func (r *AnalyzeResponse) UnmarshalJSON(data []byte) error {
	var wire struct {
		ID     string          `json:"id"`
		Status string          `json:"status"`
		Result json.RawMessage `json:"result"`
		Usage  json.RawMessage `json:"usage"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if wire.Status != "Succeeded" || len(wire.Result) == 0 || string(wire.Result) == "null" || wire.Result[0] != '{' {
		return errors.New("analysis response must contain a Succeeded status and an object result")
	}
	*r = AnalyzeResponse{
		ID: wire.ID, Status: wire.Status, Result: wire.Result, Usage: wire.Usage,
		RawJSON: append(json.RawMessage(nil), data...),
	}
	return nil
}

// Analyzer retains an analyzer definition and its common identifying properties.
// RawJSON includes arbitrary schemas, configuration, categories and model metadata.
type Analyzer struct {
	AnalyzerID string
	Status     string
	RawJSON    json.RawMessage
}

// UnmarshalJSON decodes an analyzer without discarding unknown properties.
func (a *Analyzer) UnmarshalJSON(data []byte) error {
	var wire struct {
		AnalyzerID string `json:"analyzerId"`
		Status     string `json:"status"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if wire.AnalyzerID == "" {
		return errors.New("analyzer response is missing analyzerId")
	}
	*a = Analyzer{AnalyzerID: wire.AnalyzerID, Status: wire.Status, RawJSON: append(json.RawMessage(nil), data...)}
	return nil
}

// ListAnalyzersResponse is one page of analyzer definitions.
type ListAnalyzersResponse struct {
	Value    []Analyzer `json:"value"`
	NextLink string     `json:"nextLink"`
}

// Defaults contains model-to-deployment mappings.
type Defaults struct {
	ModelDeployments map[string]string `json:"modelDeployments"`
}

// CopyAnalyzerRequest identifies the analyzer to copy.
type CopyAnalyzerRequest struct {
	SourceAnalyzerID      string `json:"sourceAnalyzerId"`
	SourceAzureResourceID string `json:"sourceAzureResourceId,omitempty"`
	SourceRegion          string `json:"sourceRegion,omitempty"`
}

// GrantCopyAuthorizationRequest identifies the destination allowed to copy an analyzer.
type GrantCopyAuthorizationRequest struct {
	TargetAzureResourceID string `json:"targetAzureResourceId"`
	TargetRegion          string `json:"targetRegion,omitempty"`
}

// CopyAuthorization is a service-side, time-limited authorization record.
type CopyAuthorization struct {
	Source                string `json:"source"`
	TargetAzureResourceID string `json:"targetAzureResourceId"`
	ExpiresAt             string `json:"expiresAt"`
}
