//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFixAnonymousModels(t *testing.T) {

}

func TestFunction(t *testing.T) {
	fileBytes, err := os.ReadFile("testdata/update_func.txt")
	require.NoError(t, err)

	newText, err := updateFunction(string(fileBytes), "Client", "uploadFileCreateRequest", func(text string) (string, error) {
		return "MIDDLE", nil
	}, &updateFunctionOptions{
		IgnoreComment: true,
	})
	require.NoError(t, err)
	require.Equal(t, "BEGIN\n// uploadFileCreateRequest creates the UploadFile request.\n// another line of documentation.\nMIDDLE\nEND\n", newText)

	newText, err = updateFunction(string(fileBytes), "Client", "uploadFileCreateRequest", func(text string) (string, error) {
		return "MIDDLE", nil
	}, &updateFunctionOptions{
		IgnoreComment: false,
	})
	require.NoError(t, err)
	require.Equal(t, "BEGIN\nMIDDLE\nEND\n", newText)
}

func TestFunctionRemove(t *testing.T) {
	fileBytes, err := os.ReadFile("testdata/remove_func.txt")
	require.NoError(t, err)

	newText, err := removeFunction(string(fileBytes), "Client", "uploadFileCreateRequest")
	require.NoError(t, err)

	require.Equal(t, "SOME TEXT BEFORE\n\n"+
		"// uploadFileHandleResponse handles the UploadFile response.\n"+
		"func (client *Client) uploadFileHandleResponse() error {\n"+
		"    // another little function\n"+
		"}\n", newText)
}

func TestSnipModel(t *testing.T) {
	modelsBytes, err := os.ReadFile("testdata/remove_type_models.txt")
	require.NoError(t, err)

	modelsSerdeBytes, err := os.ReadFile("testdata/remove_type_models_serde.txt")
	require.NoError(t, err)

	fileCache := &FileCache{
		files: map[string]string{
			"models.go":       string(modelsBytes),
			"models_serde.go": string(modelsSerdeBytes),
		},
	}
	err = removeType(fileCache, "Paths1Filz8PFilesPostRequestbodyContentMultipartFormDataSchema")
	require.NoError(t, err)

	require.Equal(t, map[string]string{
		"models.go":       "//Before that function\n\n\n//After that function\n",
		"models_serde.go": "import (\n\t\"encoding/json\"\n\t\"fmt\"\n)\n\n//Before that model\n\n\n\n\n//After that model\n",
	}, fileCache.files)
}
