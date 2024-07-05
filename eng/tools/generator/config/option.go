// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package config

import (
	"fmt"
	"regexp"
	"strings"
)

// OptionType describes the type of the option, possible values are 'Argument', 'Flag' or 'KeyValue'
type OptionType string

const (
	// Argument ...
	Argument OptionType = "Argument"
	// Flag ...
	Flag OptionType = "Flag"
	// KeyValue ...
	KeyValue OptionType = "KeyValue"
)

// Option describes an option of a autorest command line
type Option interface {
	// Type returns the type of this option
	Type() OptionType
	// Format returns the actual option in string
	Format() string
}

// ArgumentOption is an option not starting with '--' or '-'
type ArgumentOption interface {
	Option
	Argument() string
}

// FlagOption is an option that starts with '--' but do not have a value
type FlagOption interface {
	Option
	Flag() string
}

// KeyValueOption is an option that starts with '--' and have a value
type KeyValueOption interface {
	Option
	Key() string
	Value() string
}

// NewOption returns an Option from a formatted string. We should hold this result using this function: input == result.Format()
func NewOption(input string) (Option, error) {
	if strings.HasPrefix(input, "--") {
		// this option is either a flag or a key value pair option
		return parseFlagOrKeyValuePair(strings.TrimPrefix(input, "--"))
	}
	// this should be an argument
	return argument{value: input}, nil
}

var flagRegex = regexp.MustCompile(`^[a-zA-Z]`)

func parseFlagOrKeyValuePair(input string) (Option, error) {
	if !flagRegex.MatchString(input) {
		return nil, fmt.Errorf("cannot parse flag '%s', a flag or option should only start with letters", input)
	}
	segments := strings.Split(input, "=")
	if len(segments) <= 1 {
		// this should be a flag
		return flagOption{value: input}, nil
	}
	return keyValueOption{
		key:   segments[0],
		value: strings.Join(segments[1:], "="),
	}, nil
}

type argument struct {
	value string
}

func (f argument) Type() OptionType {
	return Argument
}

// Format ...
func (f argument) Format() string {
	return f.value
}

func (f argument) Argument() string {
	return f.value
}

func (f argument) String() string {
	return f.Format()
}

var _ ArgumentOption = (*argument)(nil)

type flagOption struct {
	value string
}

func (f flagOption) Type() OptionType {
	return Flag
}

// Format ...
func (f flagOption) Format() string {
	return fmt.Sprintf("--%s", f.value)
}

func (f flagOption) String() string {
	return f.Format()
}

func (f flagOption) Flag() string {
	return f.value
}

var _ FlagOption = (*flagOption)(nil)

type keyValueOption struct {
	key   string
	value string
}

func (f keyValueOption) Type() OptionType {
	return KeyValue
}

// Format ...
func (f keyValueOption) Format() string {
	return fmt.Sprintf("--%s=%s", f.key, f.value)
}

func (f keyValueOption) String() string {
	return f.Format()
}

func (f keyValueOption) Key() string {
	return f.key
}

func (f keyValueOption) Value() string {
	return f.value
}

var _ KeyValueOption = (*keyValueOption)(nil)

// NewArgument returns a new argument option (without "--")
func NewArgument(value string) Option {
	return argument{
		value: value,
	}
}

// NewFlagOption returns a flag option (with "--", without value)
func NewFlagOption(flag string) Option {
	return flagOption{
		value: flag,
	}
}

// NewKeyValueOption returns a key-value option like "--tag=something"
func NewKeyValueOption(key, value string) Option {
	return keyValueOption{
		key:   key,
		value: value,
	}
}
