// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstructAtomPath(t *testing.T) {
	basePath := constructAtomPath("/something", 1, 2)

	// I'm assuming the ordering is non-deterministic since the underlying values are just a map
	assert.Truef(t, basePath == "/something?%24skip=1&%24top=2" || basePath == "/something?%24top=2&%24skip=1", "%s wasn't one of our two variations", basePath)

	basePath = constructAtomPath("/something", 0, -1)
	assert.EqualValues(t, "/something", basePath, "Values <= 0 are ignored")

	basePath = constructAtomPath("/something", -1, 0)
	assert.EqualValues(t, "/something", basePath, "Values <= 0 are ignored")
}

// sanity check to make sure my error conforms to azcore's interface
func TestResponseError(t *testing.T) {
	var err azcore.HTTPResponse = ResponseError{}
	require.NotNil(t, err)
}
