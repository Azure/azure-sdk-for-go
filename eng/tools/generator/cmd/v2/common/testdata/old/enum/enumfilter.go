// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package enum

type EnumRemove string

const (
	EnumRemoveA EnumRemove = "A"
	EnumRemoveB EnumRemove = "B"
	EnumRemoveC EnumRemove = "C"
)

func PossibleEnumRemoveValues() []EnumRemove {
	return []EnumRemove{
		EnumRemoveA,
		EnumRemoveB,
		EnumRemoveC,
	}
}

type EnumExist string

const (
	EnumExistA EnumExist = "A"
)

func PossibleEnumExistValues() []EnumExist {
	return []EnumExist{
		EnumExistA,
	}
}
