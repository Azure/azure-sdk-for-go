// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model

import "fmt"

// Option describes an option of a autorest command line
type Option interface {
	// Format returns the actual option in string
	Format() string
}

type argument struct {
	value string
}

// Format ...
func (f argument) Format() string {
	return f.value
}

type flagOption struct {
	value string
}

// Format ...
func (f flagOption) Format() string {
	return fmt.Sprintf("--%s", f.value)
}

type keyValueOption struct {
	key   string
	value string
}

// Format ...
func (f keyValueOption) Format() string {
	return fmt.Sprintf("--%s=%s", f.key, f.value)
}

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
