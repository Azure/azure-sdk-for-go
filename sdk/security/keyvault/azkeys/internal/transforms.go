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
	// delete the version path param check (TypeSpec doesn't allow optional path parameters)
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

	// delete DeletionRecoveryLevel
	regexReplace("models.go", `RecoveryLevel \*DeletionRecoveryLevel`, "RecoveryLevel *string")
	regexReplace("constants.go", `(?:\/\/.*\s)+type DeletionRecoveryLevel string`, "")
	regexReplace("constants.go", `(?:\/\/.*\s)+func PossibleDeletionRecovery(?:.+\s)+\}`, "")
	regexReplace("constants.go", `const \(\n\/\/ DeletionRecoveryLevel(?:.+\s)+\)`, "")

	// fix up doc comments
	regexReplace("models.go", `DeletedKeyBundle`, `DeletedKey`)
	regexReplace("responses.go", `DeletedKeyBundle`, `DeletedKey`)
	regexReplace("models.go", `For.*?, see((.|\n\/\/)*)\.`, "")
	regexReplace("constants.go", `For.*?, see((.|\n\/\/)*)\.`, "")
}
