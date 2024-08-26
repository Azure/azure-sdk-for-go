// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
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

func ParseTrack2(config *Config, specRoot string) (armServices map[string][]common.PackageInfo, errResult error) {
	armServices = make(map[string][]common.PackageInfo)
	for readme, track2Request := range config.Track2Requests {
		for _, request := range track2Request {
			service, err := common.ReadV2ModuleNameToGetNamespace(filepath.Join(specRoot, strings.ReplaceAll(readme, "readme.md", "readme.go.md")))
			if err != nil {
				errResult = errors.Join(errResult, fmt.Errorf("cannot get readme.go.md content: %+v", err))
				continue
			}

			var subService string
			_, after, _ := strings.Cut(request.PackageFlag, "package-")
			s := strings.Split(after, "-")
			if _, err = strconv.Atoi(s[0]); err != nil && s[0] != "preview" {
				subService = s[0]
			}

			for arm, packageInfos := range service {
				if subService == "" && len(packageInfos) > 0 {
					packageInfos[0].RequestLink = request.RequestLink
					packageInfos[0].ReleaseDate = request.TargetDate
					packageInfos[0].Tag = fmt.Sprintf("tag: %s", request.PackageFlag)
					armServices[arm] = append(armServices[arm], packageInfos[0])
				}

				if subService != "" && len(packageInfos) == 1 {
					packageInfos[0].RequestLink = request.RequestLink
					packageInfos[0].ReleaseDate = request.TargetDate
					packageInfos[0].Tag = fmt.Sprintf("tag: %s", request.PackageFlag)
					armServices[arm] = append(armServices[arm], packageInfos[0])
				} else if subService != "" && len(packageInfos) > 1 {
					for _, packageInfo := range packageInfos {
						if strings.Contains(packageInfo.Config, subService) {
							packageInfo.RequestLink = request.RequestLink
							packageInfo.ReleaseDate = request.TargetDate
							packageInfo.Tag = fmt.Sprintf("tag: %s", request.PackageFlag)
							armServices[arm] = append(armServices[arm], packageInfo)
							break
						}
					}
				}
			}
		}
	}

	return armServices, errResult
}
