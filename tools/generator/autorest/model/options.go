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

type localOptions struct {
	AutorestArguments []string `json:"autorestArguments"`
}

func (o localOptions) Arguments() []string {
	return o.AutorestArguments
}

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

// String ...
func (o localOptions) String() string {
	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}

// CodeGeneratorVersion ...
func (o localOptions) CodeGeneratorVersion() string {
	return ""
}
