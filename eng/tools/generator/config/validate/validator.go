// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package validate

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config"
	"github.com/hashicorp/go-multierror"
)

type Validator interface {
	Validate(cfg config.Config) error
}

func NewLocalValidator(specRoot string) Validator {
	return &localValidator{
		specRoot: specRoot,
	}
}

func NewRemoteValidator(ctx context.Context, client *query.Client) Validator {
	return &remoteValidator{
		ctx:    ctx,
		client: client,
	}
}
func ParseTrack2(config *config.Config, specRoot string) (armServices map[string][]common.PackageInfo, errResult error) {
	armServices = make(map[string][]common.PackageInfo)
	for readme, track2Request := range config.Track2Requests {
		for _, request := range track2Request {
			service, err := common.ReadV2ModuleNameToGetNamespace(filepath.Join(specRoot, getReadmeGoFromReadme(readme)))
			if err != nil {
				errResult = multierror.Append(errResult, fmt.Errorf("cannot get readme.go.md content: %+v", err))
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
