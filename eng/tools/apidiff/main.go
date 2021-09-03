// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"github.com/Azure/azure-sdk-for-go/eng/tools/apidiff/cmd"
)

// breaking changes definitions
//
// const
//
// 1. removal/renaming of a const
// 2. changing the type of a const
//
// func/method
//
// 1. removal/renaming of a func
// 2. addition/removal of params and/or return values
// 3. changing param/retval types
// 4. changing receiver/param/retval byval/byref
//
// struct
// 1. removal/renaming of a field
// 2. chaning a field's type
// 3. changing a field from anonymous to explicit or explicit to anonymous
// 4. changing byval/byref
//

func main() {
	cmd.Execute()
}
