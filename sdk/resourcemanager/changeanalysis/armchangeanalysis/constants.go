//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armchangeanalysis

const (
	moduleName    = "armchangeanalysis"
	moduleVersion = "v1.1.0"
)

// ChangeCategory - The change category.
type ChangeCategory string

const (
	ChangeCategoryUser   ChangeCategory = "User"
	ChangeCategorySystem ChangeCategory = "System"
)

// PossibleChangeCategoryValues returns the possible values for the ChangeCategory const type.
func PossibleChangeCategoryValues() []ChangeCategory {
	return []ChangeCategory{
		ChangeCategoryUser,
		ChangeCategorySystem,
	}
}

// ChangeType - The type of the change.
type ChangeType string

const (
	ChangeTypeAdd    ChangeType = "Add"
	ChangeTypeRemove ChangeType = "Remove"
	ChangeTypeUpdate ChangeType = "Update"
)

// PossibleChangeTypeValues returns the possible values for the ChangeType const type.
func PossibleChangeTypeValues() []ChangeType {
	return []ChangeType{
		ChangeTypeAdd,
		ChangeTypeRemove,
		ChangeTypeUpdate,
	}
}

type Level string

const (
	LevelImportant Level = "Important"
	LevelNoisy     Level = "Noisy"
	LevelNormal    Level = "Normal"
)

// PossibleLevelValues returns the possible values for the Level const type.
func PossibleLevelValues() []Level {
	return []Level{
		LevelImportant,
		LevelNoisy,
		LevelNormal,
	}
}
