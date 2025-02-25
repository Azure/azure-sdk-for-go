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

	// remove the DeletionRecoveryLevel type
	regexReplace("models.go", "DeletionRecoveryLevel", "string")
	regexReplace("constants.go", `(?:\/\/.*\s)+type DeletionRecoveryLevel string`, "")
	regexReplace("constants.go", `(?:\/\/.*\s)+func PossibleDeletionRecovery(?:.+\s)+\}`, "")
	regexReplace("constants.go", `const \(\n\/\/ DeletionRecoveryLevel(?:.+\s)+\)`, "")

	// remove Max Results parameter
	regexReplace("options.go", `(?:\/\/.*\s)+\sMaxresults \*int32`, ``)
	regexReplace("client.go", `\sif options != nil && options.Maxresults != nil \{\s+.+\)\s+\}\s`, "")
	regexReplace("client.go", `options \*ListIssuerPropertiesOptions\) \(\*policy`, "_ *ListIssuerPropertiesOptions) (*policy")
	regexReplace("client.go", `options \*ListCertificatePropertiesVersionsOptions\) \(\*policy`, "_ *ListCertificatePropertiesVersionsOptions) (*policy")
	regexReplace("options.go", `{\n\n}`, `{
		// placeholder for future optional parameters
	}`)

	// replace Error with ErrorInfo
	regexReplace("models.go", `Error \*KeyVaultErrorError`, `Error *ErrorInfo`)
	regexReplace("models.go", `type KeyVaultErrorError struct.+\{(?:\s.+\s)+\}`, "")
	regexReplace("models_serde.go", `(?:\/\/.*\s)+func \(\w \*?KeyVaultErrorError\).*\{\s(?:.+\s)+\}\s`, "")

	// clean up doc comments
	regexReplace("models.go", `For valid values\, see JsonWebKeyCurveName\.`, "")
	regexReplace("constants.go", `For valid values\, see JsonWebKeyCurveName\.`, "")
}
