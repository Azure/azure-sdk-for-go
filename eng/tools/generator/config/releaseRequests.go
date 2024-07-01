// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
)

type Track2ReleaseRequests map[string][]Track2Request

type Track2Request struct {
	ReleaseRequestInfo
	PackageFlag string `json:"packageFlag,omitempty"`
}

func (c Track2ReleaseRequests) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c Track2ReleaseRequests) Add(readme string, info Track2Request) {
	if !c.Contains(readme) {
		c[readme] = make([]Track2Request, 0)
	}
	c[readme] = append(c[readme], info)
}

func (c Track2ReleaseRequests) Contains(readme string) bool {
	_, ok := c[readme]
	return ok
}

type Track1ReleaseRequests map[string]TagForRelease

type TagForRelease map[string][]ReleaseRequestInfo

type ReleaseRequestInfo struct {
	TargetDate  *time.Time `json:"targetDate,omitempty"`
	RequestLink string     `json:"requestLink,omitempty"`
}

func (info ReleaseRequestInfo) HasTargetDate() bool {
	return info.TargetDate != nil
}

func (info ReleaseRequestInfo) String() string {
	m := fmt.Sprintf("Release request '%s'", info.RequestLink)
	if info.HasTargetDate() {
		m = fmt.Sprintf("%s (Target date: %s)", m, *info.TargetDate)
	}
	return m
}

func (c Track1ReleaseRequests) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c Track1ReleaseRequests) Add(readme, tag string, info ReleaseRequestInfo) {
	if !c.Contains(readme) {
		c[readme] = TagForRelease{}
	}
	c[readme].Add(tag, info)
}

func (c Track1ReleaseRequests) Contains(readme string) bool {
	_, ok := c[readme]
	return ok
}

func (c TagForRelease) Add(tag string, info ReleaseRequestInfo) {
	if !c.Contains(tag) {
		c[tag] = []ReleaseRequestInfo{}
	}
	c[tag] = append(c[tag], info)
}

func (c TagForRelease) Contains(tag string) bool {
	_, ok := c[tag]
	return ok
}

type TypeSpecReleaseRequests map[string][]Track2Request

func (c TypeSpecReleaseRequests) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c TypeSpecReleaseRequests) Add(readme string, info Track2Request) {
	if !c.Contains(readme) {
		c[readme] = make([]Track2Request, 0)
	}
	c[readme] = append(c[readme], info)
}

func (c TypeSpecReleaseRequests) Contains(readme string) bool {
	_, ok := c[readme]
	return ok
}

type TypeSpecPakcageInfo struct {
	common.PackageInfo
	TspConfigPath string
}

func GetTypeSpecFromConfig(config *Config, specRoot string) (tspServices map[string][]TypeSpecPakcageInfo, errResult error) {
	tspServices = make(map[string][]TypeSpecPakcageInfo)
	for tspConfigPath, typespecRequests := range config.TypeSpecRequests {
		for _, releaseRequestInfo := range typespecRequests {
			localTspConfigPath := filepath.Join(specRoot, tspConfigPath)
			tspConfig, err := typespec.ParseTypeSpecConfig(localTspConfigPath)
			if err != nil {
				return nil, err
			}
			module, err := tspConfig.GetModuleName()
			if err != nil {
				return nil, err
			}

			tspServices[module[0]] = append(tspServices[localTspConfigPath], TypeSpecPakcageInfo{
				PackageInfo: common.PackageInfo{
					Name:        module[1],
					RequestLink: releaseRequestInfo.RequestLink,
					ReleaseDate: releaseRequestInfo.TargetDate,
					Tag:         releaseRequestInfo.PackageFlag,
				},
				TspConfigPath: localTspConfigPath,
			})
		}
	}

	return
}
