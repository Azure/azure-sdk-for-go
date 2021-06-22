// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model_ext

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"io"
	"io/ioutil"
	"reflect"
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

//func NewOptions(options ...model.Option) model.Options {
//	return &Options{
//		autorestArguments: options,
//	}
//}

func ParseOptions(raw []string) (*model.Options, error) {
	var options []model.Option
	for _, r := range raw {
		o, err := model.NewOption(r)
		if err != nil {
			return nil, fmt.Errorf("cannot parse option '%s': %+v", r, err)
		}
		options = append(options, o)
	}
	o := model.NewOptions(options...)
	return &o, nil
}

func (o *Options) Arguments() []model.Option {
	return o.autorestArguments
}

func (o *Options) CodeGeneratorVersion() string {
	for _, argument := range o.autorestArguments {
		if v, ok := argument.(model.KeyValueOption); ok {
			if v.Key() == "use" {
				return v.Value()
			}
		}
	}
	return ""
}

func (o *Options) MergeOptions(newOptions ...model.Option) *Options {
	arguments := o.autorestArguments
	for _, no := range newOptions {
		i := indexOfOptions(o.Arguments(), no)
		if i >= 0 {
			arguments[i] = no
		} else {
			arguments = append(arguments, no)
		}
	}

	return &Options{
		autorestArguments: arguments,
	}
}

func indexOfOptions(options []model.Option, option model.Option) int {
	for i, o := range options {
		if matchOption(o, option) {
			return i
		}
	}
	return -1
}

func matchOption(left, right model.Option) bool {
	if left.Type() == model.Argument && right.Type() == model.Argument {
		// we always identify arguments as different entities, even they have the same content.
		return false
	}

	// we identity flags and key value options as map keys therefore they are unique
	return getKey(left) == getKey(right)
}

func getKey(o model.Option) string {
	switch n := o.(type) {
	case model.FlagOption:
		return n.Flag()
	case model.KeyValueOption:
		return n.Key()
	case model.ArgumentOption:
		return ""
	default:
		panic(fmt.Sprintf("unknown type of option %v", reflect.TypeOf(o)))
	}
}
