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
	//  modify Restore to use implementation with custom poller handler
	regexReplace("client.go", `\[PreFullRestoreResponse\], error\) \{\s(?:.+\s)+\}`, "[PreFullRestoreResponse], error) {return client.beginPreFullRestore(ctx, preRestoreOperationParameters, options)}")
	regexReplace("client.go", `\[FullRestoreResponse\], error\) \{\s(?:.+\s)+\}`, "[FullRestoreResponse], error) {return client.beginFullRestore(ctx, restoreBlobDetails, options)}")
	regexReplace("client.go", `\[SelectiveKeyRestoreResponse\], error\) \{\s(?:.+\s)+\}`, "[SelectiveKeyRestoreResponse], error) {return client.beginSelectiveKeyRestore(ctx, keyName, restoreBlobDetails, options)}")

	// replace Error with ErrorInfo
	regexReplace("models.go", `Error \*string`, "Error *ErrorInfo")
}
