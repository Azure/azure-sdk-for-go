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
	// settings
	// fix up span names
	regexReplace("settings/client.go", `StartSpan\(ctx, "Client`, `StartSpan(ctx, "settings.Client`)

	// rbac
	// change type of scope parameter from string to RoleScope
	regexReplace("rbac/client.go", `scope string`, "scope RoleScope")
	regexReplace("rbac/client.go", `scope\)`, "string(scope))")
	regexReplace("rbac/fake/server.go", `scope string`, "scope rbac.RoleScope")
	regexReplace("rbac/fake/server.go", `\, scopeParam\,`, ", rbac.RoleScope(`/`+scopeParam),")
	regexReplace("rbac/fake/server.go", `\(scopeParam\, `, "(rbac.RoleScope(`/`+scopeParam), ")

	// fix up span names
	regexReplace("rbac/client.go", `StartSpan\(ctx, "Client`, `StartSpan(ctx, "rbac.Client`)

	// backup restore
	// change type of Error from Error to ErrorInfo
	// delete error struct
	regexReplace("backup/models.go", `Error \*Error`, "Error *ErrorInfo")
	regexReplace("backup/models.go", `(?:\/\/.*\s)+type Error struct.+\{(?:\s.+\s)+\}`, "")
	regexReplace("backup/models_serde.go", `(?:\/\/.*\s)+func \(\w \*?Error\).*\{\s(?:.+\s)+\}\s`, "")

	//  modify Restore to use implementation with custom poller handler
	regexReplace("backup/client.go", `\[FullRestoreResponse\], error\) \{\s(?:.+\s)+\}`, "[FullRestoreResponse], error) {return client.beginFullRestore(ctx, restoreBlobDetails, options)}")
	regexReplace("backup/client.go", `\[SelectiveKeyRestoreResponse\], error\) \{\s(?:.+\s)+\}`, "[SelectiveKeyRestoreResponse], error) {return client.beginSelectiveKeyRestore(ctx, keyName, restoreBlobDetails, options)}")

	// remove OperationStatus
	regexReplace("backup/models.go", `OperationStatus`, "string")

	// replace FullBackupOperationError with ErrorInfo
	regexReplace("backup/models.go", `type FullBackupOperationError.+\{(?:\s.+\s)+\}\s`, "")
	regexReplace("backup/models_serde.go", `(?:\/\/.*\s)+func \(\w \*?FullBackupOperationError\).*\{\s(?:.+\s)+\}\s`, "")
	regexReplace("backup/models.go", `FullBackupOperationError`, "ErrorInfo")

	// fix up span names
	regexReplace("backup/client.go", `StartSpan\(ctx, "Client`, `StartSpan(ctx, "backup.Client`)

	// fix fakes regex to allow scope to be optional
	regexReplace("rbac/fake/server.go", `(\(\?P<scope\>(.*?)\))`, `?$1?`)
}
