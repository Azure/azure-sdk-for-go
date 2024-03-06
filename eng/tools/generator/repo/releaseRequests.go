// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package repo

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/common"
)

var (
	cache map[string]autorest.GenerationMetadata
)

func ContainsPackage(root, readme, tag string) (autorest.GenerationMetadata, bool) {
	if cache == nil {
		if err := initCache(root); err != nil {
			panic(err)
		}
	}

	v, ok := cache[fmt.Sprintf("%s|%s", readme, tag)]
	return v, ok
}

func initCache(root string) error {
	cache = make(map[string]autorest.GenerationMetadata)
	m, err := autorest.CollectGenerationMetadata(common.ServicesPath(root))
	if err != nil {
		return err
	}

	for _, metadata := range m {
		cache[fmt.Sprintf("%s|%s", metadata.RelativeReadme(), metadata.Tag)] = metadata
	}

	return nil
}
