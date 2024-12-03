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
	r, err := regexp.Compile(regex)
	if err != nil {
		panic(err)
	}

	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	file = r.ReplaceAll(file, []byte(replace))

	err = os.WriteFile(fileName, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// delete the version path param check (version == "" is legal for Key Vault but indescribable by OpenAPI)
	regexReplace("client.go", `\sif version == "" \{\s+.+version cannot be empty"\)\s+\}\s`, "")

	// make secret IDs a convenience type so we can add parsing methods
	regexReplace("models.go", `\sKID \*string(\s+.*)`, "KID *ID$1")

	// remove Max Results parameter
	regexReplace("options.go", `(?:\/\/.*\s)+\sMaxresults \*int32`, `// placeholder for future optional parameters`)
	regexReplace("client.go", `\sif options != nil && options.Maxresults != nil \{\s+.+\)\s+\}\s`, "")

	// change type of KeyOps to KeyOperation
	regexReplace("models.go", `KeyOps \[\]\*string`, `KeyOps []*KeyOperation`)

	// delete SignatureAlgorithmRSNULL
	regexReplace("constants.go", `.*(\bSignatureAlgorithmRSNULL\b).*`, "")

	// delete KeyOperationExport
	regexReplace("constants.go", `.*(\bKeyOperationExport\b).*`, "")

	// delete strconv
	regexReplace("client.go", `\"strconv\"`, "")
}
