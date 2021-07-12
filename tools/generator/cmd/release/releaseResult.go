// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package release

import (
	"fmt"
	"strings"

	"github.com/ahmetb/go-linq/v3"
)

func GetReleaseResult(packages []GenerateResult) []string {
	selectFunc := getReleaseResultPerPackage
	if linq.From(packages).All(func(r interface{}) bool {
		return len(r.(GenerateResult).Info) > 0
	}) {
		selectFunc = getReleaseResultPerRequest
	}
	var result []string
	linq.From(packages).Select(selectFunc).ToSlice(&result)
	return result
}

func getReleaseResultPerPackage(r interface{}) interface{} {
	return fmt.Sprintf(r.(GenerateResult).Package.PackageName)
}

func getReleaseResultPerRequest(r interface{}) interface{} {
	var lines []string
	for _, info := range r.(GenerateResult).Info {
		lines = append(lines, fmt.Sprintf("Fixes %s", info.RequestLink))
	}
	return strings.Join(lines, "\n")
}
