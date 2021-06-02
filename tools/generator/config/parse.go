// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type parser struct {
	reader io.Reader
}

func FromReader(reader io.Reader) *parser {
	return &parser{
		reader: reader,
	}
}

func (p *parser) Parse() (*Config, error) {
	b, err := ioutil.ReadAll(p.reader)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, fmt.Errorf("content of configs must be a valid JSON: %+v", err)
	}

	return &config, nil
}
