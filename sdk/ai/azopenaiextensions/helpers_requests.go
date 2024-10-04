//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiextensions

import (
	"github.com/openai/openai-go/option"
)

// WithDataSources adds in Azure data sources to be used with the "Azure OpenAI On Your Data" feature.
func WithDataSources(dataSources ...AzureChatExtensionConfigurationClassification) option.RequestOption {
	return option.WithJSONSet("data_sources", dataSources)
}

// WithEnhancements configures Azure OpenAI enhancements, optical character recognition (OCR).
func WithEnhancements(config AzureChatEnhancementConfiguration) option.RequestOption {
	return option.WithJSONSet("enhancements", config)
}
