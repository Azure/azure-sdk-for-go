package model

import "fmt"

type Option interface {
	Format() string
}

type argument struct {
	value string
}

func (f argument) Format() string {
	return f.value
}

type flagOption struct {
	value string
}

func (f flagOption) Format() string {
	return fmt.Sprintf("--%s", f.value)
}

type keyValueOption struct {
	key   string
	value string
}

func (f keyValueOption) Format() string {
	return fmt.Sprintf("--%s=%s", f.key, f.value)
}

func NewArgument(value string) Option {
	return argument{
		value: value,
	}
}

func NewFlagOption(flag string) Option {
	return flagOption{
		value: flag,
	}
}

func NewKeyValueOption(key, value string) Option {
	return keyValueOption{
		key:   key,
		value: value,
	}
}
