// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package refresh

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest"
	"github.com/ahmetb/go-linq/v3"
)

type GenerationInfo struct {
	PackageFullPath string
	autorest.GenerationMetadata
}

func (g GenerationInfo) String() string {
	return fmt.Sprintf("[path: %s, readme: %s, tag: %s, commit: %s]", g.PackageFullPath, g.Readme, g.Tag, g.CommitHash)
}

type GenerationMap map[string][]GenerationInfo

func (m GenerationMap) Add(info GenerationInfo) {
	// check if the commit hash was in the map
	if info.CommitHash == "" {
		log.Printf("[WARNING] Empty commit hash found in package '%s', ignoring", info.PackageFullPath)
		return
	}
	if l, ok := m[info.CommitHash]; ok {
		m[info.CommitHash] = append(l, info)
	} else {
		m[info.CommitHash] = []GenerationInfo{info}
	}
}

func (m GenerationMap) Sort() {
	for _, v := range m {
		sort.SliceStable(v, func(i, j int) bool {
			return v[i].PackageFullPath < v[j].PackageFullPath
		})
	}
}

func (m GenerationMap) String() string {
	builder := strings.Builder{}
	for commit, info := range m {
		var infoMessages []string
		linq.From(info).SelectT(func(item GenerationInfo) string {
			return item.String()
		}).ToSlice(&infoMessages)
		builder.WriteString(fmt.Sprintf("%s: \n%s\n", commit, strings.Join(infoMessages, "\n")))
	}
	return builder.String()
}

func (m GenerationMap) Count() int {
	count := 0
	for _, l := range m {
		count += len(l)
	}
	return count
}

func NewGenerationMap(m map[string]autorest.GenerationMetadata) GenerationMap {
	result := GenerationMap{}

	for path, metadata := range m {
		result.Add(GenerationInfo{
			PackageFullPath:    path,
			GenerationMetadata: metadata,
		})
	}

	result.Sort()
	return result
}
