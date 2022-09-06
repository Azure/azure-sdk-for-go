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

			for _, packageInfos := range service {
				for i := range packageInfos {
					packageInfos[i].RequestLink = request.RequestLink
				}
			}

			for arm, packageInfos := range service {
				armServices[arm] = make([]common.PackageInfo, 0)
				if subService == "" || len(packageInfos) == 1 {
					armServices[arm] = packageInfos
				} else {
					for _, info := range packageInfos {
						if strings.Contains(info.Config, subService) {
							armServices[arm] = append(armServices[arm], info)
							break
						}
					}
				}
			}
		}
	}

	return armServices, errResult
}
