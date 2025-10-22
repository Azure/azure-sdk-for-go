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
	// fix up the scope parameter for fakes
	regexReplace("fake/server.go", `\, scopeParam\,`, ", rbac.RoleScope(`/`+scopeParam),")
	regexReplace("fake/server.go", `\(scopeParam\, `, "(rbac.RoleScope(`/`+scopeParam), ")
	regexReplace("fake/server.go", `(\(\?P<scope\>(.*?)\))`, `?$1?`)
}
