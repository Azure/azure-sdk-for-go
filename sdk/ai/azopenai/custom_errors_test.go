// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/stretchr/testify/require"
)

func TestExtractContentFilterError(t *testing.T) {
	t.Run("NilError", func(t *testing.T) {
		require.False(t, azopenai.ExtractContentFilterError(nil, nil))

		var contentFilterErr *azopenai.ContentFilterError
		require.False(t, azopenai.ExtractContentFilterError(nil, &contentFilterErr))
		require.Nil(t, contentFilterErr)
	})
}
