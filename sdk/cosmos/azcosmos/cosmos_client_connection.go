// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// cosmosClientConnection maintains a Pipeline for the client.
// The Pipeline is build based on the CosmosClientOptions.
type cosmosClientConnection struct {
	Pipeline azcore.Pipeline
}

func (cc *cosmosClientConnection) getPath(parentPath string, pathSegment string, id string) string {
	var completePath strings.Builder
	parentPathLength := len(parentPath)
	completePath.Grow(parentPathLength + 2 + len(pathSegment) + len(id))
	if parentPathLength > 0 {
		completePath.WriteString(parentPath)
		completePath.WriteString("/")
	}
	completePath.WriteString(pathSegment)
	completePath.WriteString("/")
	completePath.WriteString(url.QueryEscape(id))
	return completePath.String()
}
