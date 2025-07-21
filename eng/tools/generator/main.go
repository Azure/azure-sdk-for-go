// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd"
)

func main() {
	if err := cmd.Command().Execute(); err != nil {
		os.Exit(1)
	}
}
