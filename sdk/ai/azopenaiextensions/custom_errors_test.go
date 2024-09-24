//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiextensions_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiextensions"
	"github.com/stretchr/testify/require"
)

func TestExtractContentFilterError(t *testing.T) {
	t.Run("NilError", func(t *testing.T) {
		require.False(t, azopenaiextensions.ExtractContentFilterError(nil, nil))

		var contentFilterErr *azopenaiextensions.ContentFilterError
		require.False(t, azopenaiextensions.ExtractContentFilterError(nil, &contentFilterErr))
		require.Nil(t, contentFilterErr)
	})
}
