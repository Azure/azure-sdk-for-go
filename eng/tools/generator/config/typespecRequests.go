// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package config

import (
	"encoding/json"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
)

type TypeSpecReleaseRequests map[string][]Track2Request

func (c TypeSpecReleaseRequests) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c TypeSpecReleaseRequests) Add(tspConfig string, info Track2Request) {
	if !c.Contains(tspConfig) {
		c[tspConfig] = make([]Track2Request, 0)
	}
	c[tspConfig] = append(c[tspConfig], info)
}

func (c TypeSpecReleaseRequests) Contains(tspConfig string) bool {
	_, ok := c[tspConfig]
	return ok
}

type TypeSpecPakcageInfo struct {
	common.PackageInfo
	TspConfigPath string
}

func GetTypeSpecProjectsFromConfig(config *Config, specRoot string) (tspProjects map[string][]TypeSpecPakcageInfo, errResult error) {
	tspProjects = make(map[string][]TypeSpecPakcageInfo)
	specRootAbs, err := filepath.Abs(specRoot)
	if err != nil {
		return nil, err
	}
	for tspConfigPath, typespecRequests := range config.TypeSpecRequests {
		for _, releaseRequestInfo := range typespecRequests {
			localTspConfigPath := filepath.Join(specRootAbs, tspConfigPath)
			tspConfig, err := typespec.ParseTypeSpecConfig(localTspConfigPath)
			if err != nil {
				return nil, err
			}
			module, err := tspConfig.GetRpAndPackageName()
			if err != nil {
				return nil, err
			}

			tspProjects[module[0]] = append(tspProjects[localTspConfigPath], TypeSpecPakcageInfo{
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
