// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package interfacetoany

type Interface2Any struct {
	StringType interface{}
	ArrayType  []interface{}
	MapType    map[string]interface{}

	NewType interface{}
}
