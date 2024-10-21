//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const typeSpecDir = "../testdata/TempTypeSpecFiles/OpenAI.Inference"

const modelsGoFile = "../models.go"

var goModelsFiles = []string{modelsGoFile, "../custom_models.go", "../custom_client_embeddings.go"}

var typeSpecModelRE = regexp.MustCompile(`(?m)^model\s+([^\s]+)`)
var goModelRE = regexp.MustCompile(`(?m)^type\s+([^\s]+)\s+struct`)
var byteFieldRE = regexp.MustCompile(`\s+([^\s]+)\s+\[\]byte`)

// Tests that all the models have been generated but doesn't check fields.
func TestAllModelsHaveBeenGenerated(t *testing.T) {
	// OpenAI uses a few constructs that aren't used in Azure services like
	// unions that involve primitive types. The current emitters will silently
	// drop those types since they can't be represented.
	//
	// Until we resolve all of that we just check to make sure that hasn't happened,
	// and manually correct it. See the [TestNoUntypedFields] below for how to fix it.
	//
	if _, err := os.Stat(typeSpecDir); err != nil {
		t.Skipf("Skipping model/typespec tests: %s doesn't exist - run `go generate` to create it.", typeSpecDir)
	}

	typeSpecModels, err := getAllTypeSpecModelNames()
	require.NoError(t, err)
	require.NotEmpty(t, typeSpecModels)

	// grab all the models that we have in our project as well
	goModels := map[string]bool{}

	for _, path := range goModelsFiles {
		models, err := getFirstCaptureForRE(path, goModelRE)
		require.NoError(t, err)

		for _, model := range models {
			goModels[model] = true
		}
	}

	renames := map[string]string{
		// types that aren't used (either were subsumed into other types, or expanded into required arguments)
		"azuremachinelearningindexchatextensionconfiguration": "", // removed
		"azuremachinelearningindexchatextensionparameters":    "", // removed
		"batchimagegenerationoperationresponse":               "",
		"batchcreateresponse":                                 "",
		"getchatcompletionsbody":                              "",
		"getcompletionsbody":                                  "",
		"generatespeechfromtextbody":                          "",
		"deployment":                                          "",
		"getembeddingsbody":                                   "",
		"getimagegenerationsbody":                             "",
		"openaifile":                                          "",
		"openaipageablelistof<t>":                             "",
	}

	for _, typespecModel := range typeSpecModels {
		v, ok := renames[typespecModel]

		if v == "" {
			if ok {
				continue
			}
		} else {
			typespecModel = v
		}

		assert.Truef(t, goModels[typespecModel], "TypeSpec name: %q", typespecModel)
	}
}

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

		data, err := os.ReadFile("../custom_client.go")
		require.NoError(t, err)

		matches := re.FindStringSubmatch(string(data))
		require.NotEmpty(t, matches)

		require.Equal(t, openAPI.Info.Version, matches[1], "update the client_shared_test.go to use the API version we just generated from")
	})
}

func getGoModelsWithByteSliceFields(goFile string, allowed map[string]bool) ([]string, error) {
	file, err := os.Open(goFile)

	if err != nil {
		return nil, err
	}

	defer file.Close()

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

func getAllTypeSpecModelNames() ([]string, error) {
	var modelNames []string

	err := filepath.WalkDir(typeSpecDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".tsp" {
			return nil
		}

		models, err := getFirstCaptureForRE(path, typeSpecModelRE)

		if err != nil {
			return err
		}

		modelNames = append(modelNames, models...)
		return nil
	})

	return modelNames, err
}
