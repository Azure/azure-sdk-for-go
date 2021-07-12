// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model_ext

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"io"
	"io/ioutil"
	"strings"
)

type RawOptions struct {
	AutorestArguments []string `json:"autorestArguments,omitempty"`
}

func NewRawOptionsFrom(reader io.Reader) (*RawOptions, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var result RawOptions
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r RawOptions) Parse(absSDK string) (model.Options, error) {
	// replace go-sdk-folder value by the absolute path
	var argument []model.Option
	for _, v := range r.AutorestArguments {
		if strings.HasPrefix(v, "--go-sdk-folder") {
			continue
		}
		if v == "--multiapi" {
			continue
		}
		o, err := model.NewOption(v)
		if err != nil {
			return nil, err
		}
		argument = append(argument, o)
	}
	argument = append(argument, model.NewKeyValueOption("go-sdk-folder", absSDK))
	o := model.NewOptions(argument...)
	return o, nil
}

type Options struct {
	autorestArguments []model.Option
}
