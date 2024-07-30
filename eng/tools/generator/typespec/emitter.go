// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package typespec

type TypeSpecEmitters string

const (
	TypeSpec_GO       TypeSpecEmitters = "@azure-tools/typespec-go"
	TypeSpec_AUTOREST TypeSpecEmitters = "@azure-tools/typespec-autorest"
	TypeSpec_CSHARP   TypeSpecEmitters = "@azure-tools/typespec-csharp"
	TypeSpec_PYTHON   TypeSpecEmitters = "@azure-tools/typespec-python"
	TypeSpec_TS       TypeSpecEmitters = "@azure-tools/typespec-ts"
	TypeSpec_JAVA     TypeSpecEmitters = "@azure-tools/typespec-java"
)
