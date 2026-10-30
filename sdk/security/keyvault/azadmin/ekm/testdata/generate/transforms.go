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
	// the generated doc comment describes the whole administration service rather than this package
	regexReplace("keyvault_client.go", `(?s)// KeyVaultClient - .*?\n// Don't use this type directly`, "// Client manages the External Key Manager (EKM) connection of a Managed HSM.\n// Don't use this type directly")

	for _, f := range []string{"keyvault_client.go", "models.go", "models_serde.go", "options.go", "responses.go", "fake/keyvault_server.go"} {
		// name the client and fake server like the other azadmin packages: ekm.Client and fake.Server
		regexReplace(f, `keyVaultServerTransportInterceptor`, "serverTransportInterceptor")
		regexReplace(f, `KeyVaultServer`, "Server")
		regexReplace(f, `KeyVaultClient`, "Client")
		// remove the redundant package name from method, option and response names e.g. CreateEkmConnection -> CreateConnection
		regexReplace(f, `EkmConnection`, "Connection")
		regexReplace(f, `EkmCertificate`, "Certificate")
		// uppercase the EKM acronym in field names
		regexReplace(f, `\bEkm(Product|Vendor)\b`, "EKM$1")
	}
}
