package azopenai_test

import (
	"bufio"
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

var goModelsFiles = []string{"../models.go", "../custom_models.go", "../custom_client_embeddings.go"}

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
	// and manually correct it.
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

// Tests to see if any of our fields looks like one where the generator defaulted to just returning/accepting
// JSON, which is typical when TypeSpec uses a union type.
func TestNoUntypedFields(t *testing.T) {
	withByteFields, err := getGoModelsWithByteSliceFields("../models.go")
	require.NoError(t, err)

	require.Empty(t, withByteFields, "all []byte fields should be accounted for")
}

func getGoModelsWithByteSliceFields(goFile string) ([]string, error) {
	allowed := map[string]bool{
		"AddUploadPartRequest.Data":                                true,
		"AudioTranscriptionOptions.File":                           true,
		"AudioTranslationOptions.File":                             true,
		"ChatCompletionsFunctionToolDefinitionFunction.Parameters": true,
		"ChatCompletionsJSONSchemaResponseFormatJSONSchema.Schema": true,
		"FunctionDefinition.Parameters":                            true,
		"SpeechGenerationResponse.Audio":                           true,
	}

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
