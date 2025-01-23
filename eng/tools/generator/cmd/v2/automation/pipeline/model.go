// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package pipeline

import (
	"encoding/json"
	"io"
)

// GenerateInput ...
type GenerateInput struct {
	DryRun                       bool                          `json:"dryRun,omitempty"`
	SpecFolder                   string                        `json:"specFolder,omitempty"`
	HeadSha                      string                        `json:"headSha,omitempty"`
	HeadRef                      string                        `json:"headRef,omitempty"`
	RepoHTTPSURL                 string                        `json:"repoHttpsUrl,omitempty"`
	Trigger                      string                        `json:"trigger,omitempty"`
	ChangedFiles                 []string                      `json:"changedFiles,omitempty"`
	RelatedReadmeMdFile          string                        `json:"relatedReadmeMdFile,omitempty"`
	RelatedReadmeMdFiles         []string                      `json:"relatedReadmeMdFiles,omitempty"`
	InstallInstructionInput      InstallInstructionScriptInput `json:"installInstructionInput,omitempty"`
	RelatedTypeSpecProjectFolder []string                      `json:"relatedTypeSpecProjectFolder"`
}

// NewGenerateInputFrom ...
func NewGenerateInputFrom(reader io.Reader) (*GenerateInput, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var result GenerateInput
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// String ...
func (i GenerateInput) String() string {
	b, _ := json.MarshalIndent(i, "", "  ")
	return string(b)
}

// InstallInstructionScriptInput ...
type InstallInstructionScriptInput struct {
	PackageName             string   `json:"packageName,omitempty"`
	Artifacts               []string `json:"artifacts,omitempty"`
	IsPublic                bool     `json:"isPublic,omitempty"`
	DownloadURLPrefix       string   `json:"downloadUrlPrefix,omitempty"`
	DownloadCommandTemplate string   `json:"downloadCommandTemplate,omitempty"`
	Trigger                 string   `json:"trigger,omitempty"`
}

// GenerateOutput ...
type GenerateOutput struct {
	Packages []PackageResult `json:"packages"`
}

// String ...
func (o GenerateOutput) String() string {
	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}

// WriteTo ...
func (o GenerateOutput) WriteTo(writer io.Writer) (int64, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return 0, err
	}
	i, err := writer.Write(b)
	return int64(i), err
}

// PackageResult ...
type PackageResult struct {
	Version             string                          `json:"version,omitempty"`
	PackageName         string                          `json:"packageName,omitempty"`
	Path                []string                        `json:"path"`
	PackageFolder       string                          `json:"packageFolder"`
	ReadmeMd            []string                        `json:"readmeMd,omitempty"`
	Changelog           *Changelog                      `json:"changelog,omitempty"`
	Artifacts           []string                        `json:"artifacts,omitempty"`
	InstallInstructions *InstallInstructionScriptOutput `json:"installInstructions,omitempty"`
	APIViewArtifact     string                          `json:"apiViewArtifact,omitempty"`
	Language            string                          `json:"language,omitempty"`
	TypespecProject     []string                        `json:"typespecProject,omitempty"`
	HasExceptions       bool                            `json:"hasExceptions,omitempty"`
}

// Changelog ...
type Changelog struct {
	Content             *string   `json:"content,omitempty"`
	HasBreakingChange   *bool     `json:"hasBreakingChange,omitempty"`
	BreakingChangeItems *[]string `json:"breakingChangeItems,omitempty"`
}

// InstallInstructionScriptOutput ...
type InstallInstructionScriptOutput struct {
	Full string `json:"full,omitempty"`
	Lite string `json:"lite,omitempty"`
}
