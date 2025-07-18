// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package enum

type EnumAdd string

const (
	EnumAddA EnumAdd = "A"
	EnumAddB EnumAdd = "B"
)

func PossibleEnumAddValues() []EnumAdd {
	return []EnumAdd{
		EnumAddA,
		EnumAddB,
	}
}

type EnumExist string

const (
	EnumExistA EnumExist = "A"
	EnumExistB EnumExist = "B"
)

func PossibleEnumExistValues() []EnumExist {
	return []EnumExist{
		EnumExistA,
		EnumExistB,
	}
}
