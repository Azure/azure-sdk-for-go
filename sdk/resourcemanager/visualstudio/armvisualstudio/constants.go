//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armvisualstudio

const (
	moduleName    = "armvisualstudio"
	moduleVersion = "v0.5.0"
)

// AccountResourceRequestOperationType - The type of the operation.
type AccountResourceRequestOperationType string

const (
	AccountResourceRequestOperationTypeUnknown AccountResourceRequestOperationType = "unknown"
	AccountResourceRequestOperationTypeCreate  AccountResourceRequestOperationType = "create"
	AccountResourceRequestOperationTypeUpdate  AccountResourceRequestOperationType = "update"
	AccountResourceRequestOperationTypeLink    AccountResourceRequestOperationType = "link"
)

// PossibleAccountResourceRequestOperationTypeValues returns the possible values for the AccountResourceRequestOperationType const type.
func PossibleAccountResourceRequestOperationTypeValues() []AccountResourceRequestOperationType {
	return []AccountResourceRequestOperationType{
		AccountResourceRequestOperationTypeUnknown,
		AccountResourceRequestOperationTypeCreate,
		AccountResourceRequestOperationTypeUpdate,
		AccountResourceRequestOperationTypeLink,
	}
}
