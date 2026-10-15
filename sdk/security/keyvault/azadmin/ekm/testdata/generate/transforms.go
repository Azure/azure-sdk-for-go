// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"log"
	"os"
	"regexp"
)

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
	// the TypeSpec error union generates an undocumented KeyVaultErrorError type; rename it to ErrorInfo (matching the backup package)
	regexReplace("models.go", `\ntype KeyVaultErrorError struct`, "\n// ErrorInfo - Internal error from Azure Key Vault server.\ntype ErrorInfo struct")
	regexReplace("models.go", `\bKeyVaultErrorError\b`, "ErrorInfo")
	regexReplace("models_serde.go", `\bKeyVaultErrorError\b`, "ErrorInfo")
}
