// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorest_ext

import (
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest_ext/model_ext"
	"strings"
)

func GetAdditionalOptions(metadata autorest.GenerationMetadata) model.Options {
	// additional options
	additionalOptions, _ := model_ext.ParseOptions(strings.Split(metadata.AdditionalProperties.AdditionalOptions, " "))

	return *additionalOptions
}
