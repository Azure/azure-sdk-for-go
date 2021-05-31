// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Options ...
type Options interface {
	// Arguments returns the argument defined in this options
	Arguments() []Option
	// CodeGeneratorVersion returns the code generator version defined in this options
	CodeGeneratorVersion() string
	// MergeOptions merges the current options with the given options
	MergeOptions(other ...Option) Options
}

// MergeOptions will merge the given options and new option slice, and return a new Options instance
func MergeOptions(options Options, other ...Option) Options {
	arguments := options.Arguments()
	for _, no := range other {
		i := indexOfOptions(arguments, no)
		if i >= 0 {
			arguments[i] = no
		} else {
			arguments = append(arguments, no)
		}
	}

	return localOptions(arguments)
}

type localOptions []Option

// ParseOptions returns an Options instance by parsing the given raw option strings
func ParseOptions(raw []string) (Options, error) {
	var options []Option
	for _, r := range raw {
		o, err := NewOption(r)
		if err != nil {
			return nil, err
		}
		options = append(options, o)
	}

	return NewOptions(options...), nil
}

// NewOptions returns a new instance of Options with the give slice of Option
func NewOptions(options ...Option) Options {
	return localOptions(options)
}

// Arguments ...
func (o localOptions) Arguments() []Option {
	return o
}

// CodeGeneratorVersion ...
func (o localOptions) CodeGeneratorVersion() string {
	for _, argument := range o.Arguments() {
		if v, ok := argument.(KeyValueOption); ok {
			if v.Key() == "use" {
				return v.Value()
			}
		}
	}
	return ""
}

// MergeOptions ...
func (o localOptions) MergeOptions(other ...Option) Options {
	return MergeOptions(o, other...)
}

// String ...
func (o localOptions) String() string {
	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}

func indexOfOptions(options []Option, option Option) int {
	for i, o := range options {
		if matchOption(o, option) {
			return i
		}
	}
	return -1
}

func matchOption(left, right Option) bool {
	if left.Type() == Argument && right.Type() == Argument {
		// we always identify arguments as different entities, even they have the same content.
		return false
	}

	// we identity flags and key value options as map keys therefore they are unique
	return getKey(left) == getKey(right)
}

func getKey(o Option) string {
	switch n := o.(type) {
	case FlagOption:
		return n.Flag()
	case KeyValueOption:
		return n.Key()
	case ArgumentOption:
		return ""
	default:
		panic(fmt.Sprintf("unknown type of option %v", reflect.TypeOf(o)))
	}
}
