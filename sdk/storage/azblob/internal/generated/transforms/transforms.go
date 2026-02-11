// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"log"
	"os"
	"regexp"
)

type replacement struct {
	regex   string
	replace string
}

// applyReplacements applies multiple regex replacements to the given files
func applyReplacements(fileNames []string, replacements []replacement) {
	for _, fileName := range fileNames {
		file, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}

		for _, r := range replacements {
			re := regexp.MustCompile(r.regex)
			file = re.ReplaceAll(file, []byte(r.replace))
		}

		err = os.WriteFile(fileName, file, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	// Replace ETag fields in models and responses
	applyReplacements([]string{"models.go"}, []replacement{
		{`Etag\s+\*string`, `ETag *azcore.ETag`},
	})

	applyReplacements([]string{"responses.go"}, []replacement{
		{`"time"`, "\"time\"\n\t\"github.com/Azure/azure-sdk-for-go/sdk/azcore\""},
		{`ETag\s+\*string`, `ETag *azcore.ETag`},
	})

	// Replace ETag fields in options
	applyReplacements([]string{"options.go"}, []replacement{
		{`import "time"`, "import (\n\t\"time\"\n\t\"github.com/Azure/azure-sdk-for-go/sdk/azcore\"\n)"},
		{`IfMatch\s+\*string`, `IfMatch *azcore.ETag`},
		{`IfNoneMatch\s+\*string`, `IfNoneMatch *azcore.ETag`},
		{`SourceIfMatch\s+\*string`, `SourceIfMatch *azcore.ETag`},
		{`SourceIfNoneMatch\s+\*string`, `SourceIfNoneMatch *azcore.ETag`},
	})

	// Apply common client transformations for ETag handling
	clientFiles := []string{
		"appendblob_client.go",
		"blob_client.go",
		"blockblob_client.go",
		"container_client.go",
		"pageblob_client.go",
	}

	applyReplacements(clientFiles, []replacement{
		{`result\.ETag\s+=\s+&val`, `result.ETag = (*azcore.ETag)(&val)`},
		{`\*options.IfMatch`, `string(*options.IfMatch)`},
		{`\*options.IfNoneMatch`, `string(*options.IfNoneMatch)`},
		{`\*options.SourceIfMatch`, `string(*options.SourceIfMatch)`},
		{`\*options.SourceIfNoneMatch`, `string(*options.SourceIfNoneMatch)`},
	})
}
