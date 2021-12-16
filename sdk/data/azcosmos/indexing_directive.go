// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// Specifies whether or not the resource in the Azure Cosmos DB database is to be indexed.
type IndexingDirective string

const (
	// Use any pre-defined/pre-configured defaults.
	IndexingDirectiveDefault IndexingDirective = "Default"
	// Index the resource.
	IndexingDirectiveInclude IndexingDirective = "Include"
	// Do not index the resource.
	IndexingDirectiveExclude IndexingDirective = "Exclude"
)

// Returns a list of available indexing directives
func IndexingDirectives() []IndexingDirective {
	return []IndexingDirective{IndexingDirectiveDefault, IndexingDirectiveInclude, IndexingDirectiveExclude}
}

func (c IndexingDirective) ToPtr() *IndexingDirective {
	return &c
}
