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
	Arguments() []Option
	// CodeGeneratorVersion returns the code generator version defined in this options
	CodeGeneratorVersion() string
}

type generateOptions struct {
	AutorestArguments []string `json:"autorestArguments"`
}

type localOptions struct {
	arguments      []Option
	codeGenVersion string
}

// NewOptionsFrom returns a new instance of Options from the given io.Reader
func NewOptionsFrom(reader io.Reader) (Options, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var raw generateOptions
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	var options []Option
	var codeGenVersion string
	for _, r := range raw.AutorestArguments {
		o, err := NewOption(r)
		if err != nil {
			return nil, err
		}
		options = append(options, o)
		if v, ok := o.(KeyValueOption); ok && v.Key() == "use" {
			codeGenVersion = v.Value()
		}
	}
	return &localOptions{
		arguments:      options,
		codeGenVersion: codeGenVersion,
	}, nil
}

// Arguments ...
func (o localOptions) Arguments() []Option {
	return o.arguments
}

// CodeGeneratorVersion ...
func (o localOptions) CodeGeneratorVersion() string {
	return o.codeGenVersion
}

// String ...
func (o localOptions) String() string {
	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}
