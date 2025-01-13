// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package config

import (
	"encoding/json"
	"errors"
	"strings"
)

const Root = "github.com/Azure/azure-sdk-for-go"

type RefreshInfo struct {
	// AdditionalFlags are the additional options that will be used in the general refresh
	AdditionalFlags []string `json:"additionalOptions,omitempty"`
	// Packages are the full package identifier of the packages to refresh, eg 'github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight'
	Packages []string `json:"packages,omitempty"`
}

func (r RefreshInfo) AdditionalOptions() ([]Option, error) {
	return parseAdditionalOptions(r.AdditionalFlags)
}

func (r RefreshInfo) RelativePackages() []string {
	var packages []string
	for _, p := range r.Packages {
		l := strings.TrimPrefix(strings.TrimPrefix(p, Root), "/")
		packages = append(packages, l)
	}

	return packages
}

func (r RefreshInfo) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func parseAdditionalOptions(input []string) ([]Option, error) {
	var errResult error
	var options []Option
	for _, f := range input {
		o, err := NewOption(f)
		if err != nil {
			errResult = errors.Join(errResult, err)
			continue
		}
		options = append(options, o)
	}

	return options, errResult
}
