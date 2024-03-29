//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmanagementpartner

const (
	moduleName    = "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementpartner/armmanagementpartner"
	moduleVersion = "v0.7.0"
)

// ManagementPartnerState - this is the management partner state: Active or Deleted
type ManagementPartnerState string

const (
	ManagementPartnerStateActive  ManagementPartnerState = "Active"
	ManagementPartnerStateDeleted ManagementPartnerState = "Deleted"
)

// PossibleManagementPartnerStateValues returns the possible values for the ManagementPartnerState const type.
func PossibleManagementPartnerStateValues() []ManagementPartnerState {
	return []ManagementPartnerState{
		ManagementPartnerStateActive,
		ManagementPartnerStateDeleted,
	}
}
