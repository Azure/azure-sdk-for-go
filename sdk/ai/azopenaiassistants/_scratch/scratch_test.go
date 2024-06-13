package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func findTSPUnions() (map[string]bool, error) {
	unions := map[string]bool{}

	unionRE := regexp.MustCompile(`(?m)^union ([A-Za-z0-9]+).+?$`)

	err := filepath.WalkDir("../testdata/TempTypeSpecFiles/OpenAI.Assistants", func(path string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(path, ".tsp") {
			return nil
		}

		tspBytes, err := os.ReadFile(path)

		if err != nil {
			return err
		}

		allMatches := unionRE.FindAllSubmatch(tspBytes, -1)

		if allMatches == nil {
			return nil
		}

		for _, match := range allMatches {
			unionName := match[1]
			unions[string(unionName)] = true
		}

		return nil
	})

	return unions, err
}

func getSwaggerType(typeName string) (map[string]any, error) {
	swaggerBytes, err := os.ReadFile("../testdata/generated/openapi.json")

	if err != nil {
		return nil, err
	}

	var swaggerObj *struct {
		Definitions map[string]map[string]any
	}

	if err := json.Unmarshal(swaggerBytes, &swaggerObj); err != nil {
		return nil, err
	}

	for defnTypeName, defn := range swaggerObj.Definitions {

		if strings.EqualFold(defnTypeName, typeName) {
			return defn, nil
		}
	}

	return nil, fmt.Errorf("no type named %s", typeName)
}

func TestAnalysis(t *testing.T) {
	// check how our type generation is doing
	// 1. check for any fields that are typed as 'any'
	// 2. check if the swagger and the typespec look good/similar. Are we
	//     missing types, or generating types incorrectly?
	unions, err := findTSPUnions()
	require.NoError(t, err)

	jsonBytes, err := json.MarshalIndent(unions, "  ", "  ")
	require.NoError(t, err)

	t.Logf("Unions: %s\n", string(jsonBytes))

	ok := 0

	// they all exist in the swagger - are the properties all there?
	// okay, all the unions are empty. So I'll need to translate all 33 to something :|
	for union := range unions {
		swagger, err := getSwaggerType(union)
		require.NoError(t, err)

		if assert.NotEmptyf(t, swagger["properties"], "%s should have fields and not just be an empty definition", union) {
			ok++
		}
	}

	t.Logf("# of unions that were okay: %d/%d", ok, len(unions)) // the answer is 0/33
}

func TestAnyAnys(t *testing.T) {
	filepath.WalkDir("../", func(path string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(".go") {
			return nil
		}

	})
}
