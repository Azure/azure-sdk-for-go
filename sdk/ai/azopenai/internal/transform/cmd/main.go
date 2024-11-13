// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai/internal/transform"
)

func main() {
	var (
		operation  = flag.String("op", "", "Operation to perform (rename-struct, rename-method, copy-struct)")
		filename   = flag.String("file", "", "File containing to modify")
		name       = flag.String("name", "", "Symbol name")
		newName    = flag.String("new-name", "", "New symbol name (for rename)")
		structName = flag.String("struct", "", "Struct name (for remove-field)")
		field      = flag.String("field", "", "Field name (for remove-field)")
	)

	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	transformer, err := transform.New(wd)
	if err != nil {
		log.Fatalf("Failed to initialize transformer: %v", err)
	}

	switch *operation {
	case "rename-struct":
		if err := transformer.RenameStruct(*filename, *name, *newName); err != nil {
			log.Fatalf("Failed to rename struct: %v", err)
		}
	case "rename-method":
		if err := transformer.RenameMethod(*filename, *name, *newName); err != nil {
			log.Fatalf("Failed to rename method: %v", err)
		}
	case "copy-struct":
		if err := transformer.CopyStruct(*filename, *name, *newName); err != nil {
			log.Fatalf("Failed to copy struct: %v", err)
		}
	case "remove-field":
		if err := transformer.RemoveField(*filename, *structName, *field); err != nil {
			log.Fatalf("Failed to remove field: %v", err)
		}
	default:
		log.Fatalf("Unknown operation: %s", *operation)
	}
}
