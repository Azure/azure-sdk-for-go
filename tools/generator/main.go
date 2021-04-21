// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd"
)

func main() {
	if err := cmd.Command().Execute(); err != nil {
		for _, line := range strings.Split(err.Error(), "\n") {
			log.Printf("[ERROR] %s", line)
		}
		os.Exit(1)
	}
}
