//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPossibleArtifactArchitectureValues(t *testing.T) {
	require.Equal(t, 13, len(PossibleArtifactArchitectureValues()))
}

func TestPossibleArtifactManifestOrderByValues(t *testing.T) {
	require.Equal(t, 3, len(PossibleArtifactManifestOrderByValues()))
}

func TestPossibleArtifactOperatingSystemValues(t *testing.T) {
	require.Equal(t, 14, len(PossibleArtifactOperatingSystemValues()))
}

func TestPossibleArtifactTagOrderByValues(t *testing.T) {
	require.Equal(t, 3, len(PossibleArtifactTagOrderByValues()))
}

func TestPossibleContentTypeValues(t *testing.T) {
	require.Equal(t, 2, len(PossibleContentTypeValues()))
}

func Test_possiblePostContentSchemaGrantTypeValues(t *testing.T) {
	require.Equal(t, 3, len(possiblePostContentSchemaGrantTypeValues()))
}

func Test_possibleTokenGrantTypeValues(t *testing.T) {
	require.Equal(t, 2, len(possibleTokenGrantTypeValues()))
}
