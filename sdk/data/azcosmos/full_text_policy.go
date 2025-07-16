// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// FullTextPolicy represents a full-text policy for a container.
// This policy defines how text properties are indexed for full-text search operations.
// For more information see https://docs.microsoft.com/azure/cosmos-db/gen-ai/full-text-search
type FullTextPolicy struct {
	// DefaultLanguage specifies the default language for full-text indexing and search.
	// Supported languages include: en-US (English), de-DE (German), es-ES (Spanish), fr-FR (French).
	DefaultLanguage string `json:"defaultLanguage"`
	// FullTextPaths defines the text properties and their languages for full-text indexing.
	FullTextPaths []FullTextPath `json:"fullTextPaths"`
}

// FullTextPath represents a path to a text property with its associated language for full-text indexing.
type FullTextPath struct {
	// Path to the text property in the document.
	Path string `json:"path"`
	// Language specifies the language for this specific text property.
	// This can override the default language specified in the FullTextPolicy.
	Language string `json:"language"`
}
