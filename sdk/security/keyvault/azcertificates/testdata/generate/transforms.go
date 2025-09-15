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
	// allow `version` to be optional (TypeSpec doesn't allow optional path parameters)
	regexReplace("client.go", `\sif version == "" \{\s+.+version cannot be empty"\)\s+\}\s`, "")
	regexReplace("fake/server.go", `(\(\?P<certificate_version\>(.*?)\))`, `?$1?`)

	regexReplace("models.go", `(type (?:Deleted)?Certificate(?:Properties|Policy|Operation)? struct \{(?:\s.+\s)+\sID \*)string`, "$1 ID")
	regexReplace("models.go", `(\/\/ READ-ONLY; The certificate id\.\n\tID \*)string`, `$1 ID`)
	regexReplace("models.go", `\sKID \*string(\s+.*)`, "KID *ID$1")
	regexReplace("models.go", `\sSID \*string(\s+.*)`, "SID *ID$1")

	// replace Error with ErrorInfo
	regexReplace("models.go", `Error \*string`, `Error *ErrorInfo`)

	// clean up doc comments
	regexReplace("models.go", `For valid values\, see JsonWebKeyCurveName\.`, "")
	regexReplace("constants.go", `For valid values\, see JsonWebKeyCurveName\.`, "")
}
