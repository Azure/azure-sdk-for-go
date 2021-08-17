// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	pathSegmentDatabase   string = "dbs"
	pathSegmentCollection string = "colls"
)

func getResourcePath(resourceType resourceType) (string, error) {
	switch resourceType {
	case resourceTypeDatabase:
		return pathSegmentDatabase, nil
	case resourceTypeCollection:
		return pathSegmentCollection, nil
	default:
		return "", fmt.Errorf("%v is not a valid resource type", resourceType)
	}
}

func getPath(parentPath string, pathSegment string, id string) string {
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
