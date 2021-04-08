// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Options ...
type Options interface {
	// Arguments returns the argument defined in this options
	Arguments() []string
	// CodeGeneratorVersion returns the code generator version defined in this options
	CodeGeneratorVersion() string
}

type localOptions struct {
	AutorestArguments []string `json:"autorestArguments"`
}

// NewOptionsFrom returns a new instance of Options from the given io.Reader
func NewOptionsFrom(reader io.Reader) (Options, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var result localOptions
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Argument ...
func (o localOptions) Arguments() []string {
	return o.AutorestArguments
}

// String ...
func (o localOptions) String() string {
	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}

// CodeGeneratorVersion ...
func (o localOptions) CodeGeneratorVersion() string {
	return ""
}
