//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func loadFile(t *testing.T, path string) string {
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	str := string(data)

	return strings.ReplaceAll(str, "\r\n", "\n")
}

func TestFunction(t *testing.T) {
	text := loadFile(t, "testdata/update_func.txt")

	newText, err := updateFunction(text, "Client", "uploadFileCreateRequest", func(_ string) (string, error) {
		return "MIDDLE", nil
	}, &updateFunctionOptions{
		IgnoreComment: true,
	})
	require.NoError(t, err)
	require.Equal(t, "BEGIN\n// uploadFileCreateRequest creates the UploadFile request.\n// another line of documentation.\nMIDDLE\nEND\n", newText)

	newText, err = updateFunction(text, "Client", "uploadFileCreateRequest", func(_ string) (string, error) {
		return "MIDDLE", nil
	}, &updateFunctionOptions{
		IgnoreComment: false,
	})
	require.NoError(t, err)
	require.Equal(t, "BEGIN\nMIDDLE\nEND\n", newText)
}

func TestFunctionRemove(t *testing.T) {
	text := loadFile(t, "testdata/remove_func.txt")

	newText, err := removeFunctions(text, "Client", "uploadFileCreateRequest")
	require.NoError(t, err)

	require.Equal(t, "SOME TEXT BEFORE\n\n"+
		"// uploadFileHandleResponse handles the UploadFile response.\n"+
		"func (client *Client) uploadFileHandleResponse() error {\n"+
		"    // another little function\n"+
		"}\n", newText)
}

func TestSnipModel(t *testing.T) {
	modelsText := loadFile(t, "testdata/remove_type_models.txt")
	modelsSerdeText := loadFile(t, "testdata/remove_type_models_serde.txt")

	fileCache := &FileCache{
		files: map[string]string{
			"models.go":       modelsText,
			"models_serde.go": modelsSerdeText,
			"responses.go":    "ignored",
			"options.go":      "ignored",
		},
	}

	err := removeTypes(fileCache, []string{"Paths1Filz8PFilesPostRequestbodyContentMultipartFormDataSchema"}, &removeTypesOptions{
		IgnoreComment: true,
	})
	require.NoError(t, err)

	require.Equal(t, map[string]string{
		"models.go":       "//Before that function\n\n\n//After that function\n",
		"models_serde.go": "import (\n\t\"encoding/json\"\n\t\"fmt\"\n)\n\n//Before that model\n\n\n\n\n//After that model\n",
		"responses.go":    "ignored",
		"options.go":      "ignored",
	}, fileCache.files)
}

func TestSnipResponseType(t *testing.T) {
	fileCache := &FileCache{
		files: map[string]string{
			"models.go":       "ignored",
			"models_serde.go": "ignored",
			"responses.go": "hello\n// GetFileContentResponse contains the response from method Client.GetFileContent.\n" +
				"type GetFileContentResponse struct {\n" +
				"  Value []byte\n" +
				"}\n" +
				" world",
			"options.go": "ignored",
		},
	}

	err := removeTypes(fileCache, []string{"GetFileContentResponse"}, nil)
	require.NoError(t, err)

	require.Equal(t, map[string]string{
		"models.go":       "ignored",
		"models_serde.go": "ignored",
		"responses.go":    "hello\n\n world",
		"options.go":      "ignored",
	}, fileCache.files)
}
