package model

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Options interface {
	Arguments() []string
	String() string
	CodeGeneratorVersion() string
}

type LocalOptions struct {
	arguments []string
}

func (o LocalOptions) Arguments() []string {
	return o.arguments
}

func NewOptionsFrom(reader io.Reader) (Options, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var result LocalOptions
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// String ...
func (o LocalOptions) String() string {
	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}

// CodeGeneratorVersion ...
func (o LocalOptions) CodeGeneratorVersion() string {
	return ""
}
