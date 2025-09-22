// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"log"
	"os"
	"regexp"
)

// removing client prefix from types
func regexReplace(fileName string, regex string, replace string) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	r := regexp.MustCompile(regex)
	file = r.ReplaceAll(file, []byte(replace))

	err = os.WriteFile(fileName, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// change type of Rows from []byte to []Row
	regexReplace("models.go", `Rows \[\]\[\]\[\]byte`, "Rows []Row")

	// remove generated error types to replace with custom error types
	regexReplace("models.go", `(?:\/\/.*\s)+type Error(Info|Detail) struct.+\{(?:\s.+\s)+\}`, "")
	regexReplace("models_serde.go", `(?:\/\/.*\s)+func \(\w \*?Error(Info|Detail)\).*\{\s(?:.+\s)+\}\s`, "")

	// change type of Timespan from *string to TimeInterval
	regexReplace("models.go", `Timespan \*string`, "Timespan *TimeInterval")

	// fix up options type
	regexReplace("options.go", `Options \*string`, "Options *QueryOptions")
	regexReplace("client.go", `\*options\.Options`, "options.Options.preferHeader()")
	regexReplace("fake/server.go", `Options\: optionsParam`, "Options: preferHeaderToQueryOptions(*optionsParam)")

	// Adjust URL path handling in fake/server.go to remove the "/v1" prefix from req.URL.EscapedPath().
	regexReplace("fake/server.go", `req\.URL\.EscapedPath\(\)`, `strings.TrimPrefix(req.URL.EscapedPath(), "/v1")`)
	regexReplace("fake/server.go", `import \(`, `import ( "strings"`)
}
