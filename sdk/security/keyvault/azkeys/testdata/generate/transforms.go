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
	regexReplace("fake/server.go", `(\(\?P<key_version\>(.*?)\))`, `?$1?`)

	// make secret IDs a convenience type so we can add parsing methods
	regexReplace("models.go", `\sKID \*string(\s+.*)`, "KID *ID$1")

	// delete SignatureAlgorithmRSNULL
	regexReplace("constants.go", `.*(\bSignatureAlgorithmRSNULL\b).*`, "")

	// delete KeyOperationExport
	regexReplace("constants.go", `.*(\bKeyOperationExport\b).*`, "")

	// fix up doc comments
	regexReplace("models.go", `DeletedKeyBundle`, `DeletedKey`)
	regexReplace("responses.go", `DeletedKeyBundle`, `DeletedKey`)
	regexReplace("models.go", `For.*?, see((.|\n\/\/)*)\.`, "")
	regexReplace("constants.go", `For.*?, see((.|\n\/\/)*)\.`, "")
}
