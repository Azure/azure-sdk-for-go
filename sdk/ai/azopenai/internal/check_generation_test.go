// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const modelsGoFile = "../models.go"

var goModelRE = regexp.MustCompile(`(?m)^type\s+([^\s]+)\s+struct`)
var byteFieldRE = regexp.MustCompile(`\s+([^\s]+)\s+\[\]byte`)

// Tests to see if any of our fields looks like one where the generator defaulted to
// just accepting JSON, which is typical when TypeSpec uses a union type that is
// not polymorphic (ie, string | someObject).
func TestNoUntypedFields(t *testing.T) {
	// these types are allowed as they're intended to be []byte fields.
	allowed := map[string]bool{
		"AddUploadPartRequest.Data":                                true,
		"AudioTranscriptionOptions.File":                           true,
		"AudioTranslationOptions.File":                             true,
		"ChatCompletionsFunctionToolDefinitionFunction.Parameters": true, // user intentionally passes their own serialized JSON bytes
		"ChatCompletionsJSONSchemaResponseFormatJSONSchema.Schema": true, // user intentionally passes their own serialized JSON bytes
		"FunctionDefinition.Parameters":                            true, // user intentionally passes their own serialized JSON bytes
		"SpeechGenerationResponse.Audio":                           true,
	}

	withByteFields, err := getGoModelsWithByteSliceFields(modelsGoFile, allowed)
	require.NoError(t, err)

	// To fix this, you'll need manually create a union input type:
	//
	// 1. Create the union type and it's associated functions. Look at custom_models.go and [MongoDBChatExtensionParametersEmbeddingDependency]
	//    to see what you'll need:
	//    - MongoDBChatExtensionParametersEmbeddingDependency (the union type - naming is "object that has field" + "field name")
	//    - NewMongoDBChatExtensionParametersEmbeddingDependency (the function the user calls to construct the MongoDBChatExtensionParametersEmbeddingDependency)
	//    - MongoDBChatExtensionParametersEmbeddingDependency.MarshalJSON
	//
	// 2. Add in the an autorest.md snippet in "## Unions" section. This will make it so the Go emitter will reference
	//    your custom type. See 'MongoDBChatExtensionParametersEmbeddingDependency's block within there for a sample.
	require.Empty(t, withByteFields, "no new []byte fields. If this test fails see the test for details on how to fix it.")
}

func TestAllOYDModelsAreGenerated(t *testing.T) {
	if _, err := os.Stat("../testdata/generated/openapi.json"); err != nil {
		t.Skip("openapi.json isn't there, not doing codegen tests")
	}

	// we do a little autorest hackery to trim out models that aren't used, just check that we didn't
	// miss something new. If we did, just add it to the "Keep only "Azure OpenAI On Your Data"
	// models, or enhancements."
	// yaml block.

	// grab all the models that we have in our project as well
	goModels := map[string]bool{}

	models, err := getFirstCaptureForRE("../models.go", goModelRE)
	require.NoError(t, err)

	for _, model := range models {
		goModels[model] = true
	}

	/*
		Example:

		definitions.AzureCosmosDBChatExtensionConfiguration: {
			"allOf": [{
				"$ref": "#/definitions/AzureChatExtensionConfiguration"
			}],
		}
	*/

	var openAPI *struct {
		Definitions map[string]struct {
			AllOf []*struct {
				Ref string `json:"$ref"`
			}
		}
	}

	data, err := os.ReadFile("../testdata/generated/openapi.json")
	require.NoError(t, err)

	err = json.Unmarshal(data, &openAPI)
	require.NoError(t, err)

	for name, defn := range openAPI.Definitions {
		if len(defn.AllOf) == 0 || len(defn.AllOf) > 1 || defn.AllOf[0].Ref != "#/definitions/AzureChatExtensionConfiguration" {
			continue
		}

		assert.True(t, goModels[strings.ToLower(name)], "%s exists in the swagger, but didn't get generated", name)
	}
}

func TestAPIVersionIsBumped(t *testing.T) {
	if _, err := os.Stat("../testdata/generated/openapi.json"); err != nil {
		t.Skip("openapi.json isn't there, not doing codegen tests")
	}

	var openAPI *struct {
		Info struct {
			Version string
		}
	}

	data, err := os.ReadFile("../testdata/generated/openapi.json")
	require.NoError(t, err)

	err = json.Unmarshal(data, &openAPI)
	require.NoError(t, err)

	t.Run("TestsUseNewAPIVersion", func(t *testing.T) {
		// ex: const apiVersion = "2024-07-01-preview"
		re := regexp.MustCompile(`const apiVersion = "(.+?)"`)

		data, err := os.ReadFile("../client_shared_test.go")
		require.NoError(t, err)

		matches := re.FindStringSubmatch(string(data))
		require.NotEmpty(t, matches)

		require.Equal(t, openAPI.Info.Version, matches[1], "update the client_shared_test.go to use the API version we just generated from")
	})

	// check examples
	t.Run("ExamplesUseNewAPIVersion", func(t *testing.T) {
		// ex: azure.WithEndpoint(endpoint, "2024-07-01-preview"),
		re := regexp.MustCompile(`azure\.WithEndpoint\(.+?, "(.+?)"\),`)

		paths, err := filepath.Glob("../example*.go")
		require.NoError(t, err)
		require.NotEmpty(t, paths)

		for _, path := range paths {
			t.Logf("Checking example %s", path)

			file, err := os.ReadFile(path)
			require.NoError(t, err)

			matches := re.FindAllStringSubmatch(string(file), -1)
			require.NotEmpty(t, matches)

			for _, m := range matches {
				assert.Equalf(t, openAPI.Info.Version, m[1], "api-version out of date in %s", path)
			}
		}
	})
}

func getGoModelsWithByteSliceFields(goFile string, allowed map[string]bool) ([]string, error) {
	file, err := os.Open(goFile)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)

	var byteFields []string
	currentStruct := ""

	for scanner.Scan() {
		line := scanner.Text()

		if matches := goModelRE.FindStringSubmatch(line); len(matches) > 0 {
			currentStruct = matches[1]
			continue
		}

		if matches := byteFieldRE.FindStringSubmatch(line); len(matches) > 0 {
			key := fmt.Sprintf("%s.%s", currentStruct, matches[1])
			if allowed[key] {
				continue
			}

			byteFields = append(byteFields, key)
		}
	}

	sort.Strings(byteFields)
	return byteFields, nil
}

func getFirstCaptureForRE(file string, re *regexp.Regexp) ([]string, error) {
	var modelNames []string

	data, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	for _, match := range re.FindAllStringSubmatch(string(data), -1) {
		modelName := strings.ToLower(match[1])
		modelNames = append(modelNames, modelName)
	}

	return modelNames, nil
}
