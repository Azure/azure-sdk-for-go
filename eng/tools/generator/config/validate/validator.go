// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package validate

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config"
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
func ParseTrack2(config *config.Config, specRoot string) (armServices map[string][]string, errResult error) {
	var i int
	armServices = make(map[string][]string)
	for readme, _ := range config.Track2Requests {
		contentOfReadmeGo, err := getReadmeContent(specRoot, getReadmeGoFromReadme(readme))
		if err != nil {
			errResult = multierror.Append(errResult, fmt.Errorf("cannot get readme.go.md content: %+v", err))
			continue
		}
		// 获取specRoot service name --spec-rp-name="sagger service name"
		splits := strings.Split(readme, "/")
		//specRpName := fmt.Sprintf("--spec-rp-name=%s", splits[1])
		service, armService := GetModuleName(contentOfReadmeGo)
		// 得到需要执行release的arm信息
		//armServices = append(armServices, []string{})
		armServices[readme] = make([]string, 0)
		armServices[readme] = append(armServices[readme], service, armService, splits[1])
		i++
	}
	return armServices, errResult
}
