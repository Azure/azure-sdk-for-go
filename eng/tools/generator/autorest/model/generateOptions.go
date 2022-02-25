// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// GenerateOptions deserialize from eng/generate_options.json file in root directory of azure-sdk-for-go
type GenerateOptions struct {
	AutorestArguments []string `json:"autorestArguments"`
	AdditionalOptions []string `json:"additionalOptions,omitempty"`
}

// NewGenerateOptionsFrom reads GenerateOptions from the given io.Reader
func NewGenerateOptionsFrom(reader io.Reader) (*GenerateOptions, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var raw GenerateOptions
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}

	return &raw, nil
}
