// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package delta

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
)

func TestParamsEqual(t *testing.T) {
	tests := []struct {
		name     string
		lhs      []exports.Param
		rhs      []exports.Param
		expected bool
	}{
		{
			name: "Simple order change",
			lhs: []exports.Param{
				{Name: "resourceGroupName", Type: "string"},
				{Name: "serviceName", Type: "string"},
				{Name: "options", Type: "*Options"},
			},
			rhs: []exports.Param{
				{Name: "serviceName", Type: "string"},
				{Name: "resourceGroupName", Type: "string"},
				{Name: "options", Type: "*Options"},
			},
			expected: false,
		},
		{
			name: "Three params order change",
			lhs: []exports.Param{
				{Name: "resourceGroupName", Type: "string"},
				{Name: "serviceName", Type: "string"},
				{Name: "subscriptionID", Type: "string"},
			},
			rhs: []exports.Param{
				{Name: "serviceName", Type: "string"},
				{Name: "subscriptionID", Type: "string"},
				{Name: "resourceGroupName", Type: "string"},
			},
			expected: false,
		},
		{
			name: "No change",
			lhs: []exports.Param{
				{Name: "ctx", Type: "context.Context"},
				{Name: "name", Type: "string"},
				{Name: "value", Type: "int"},
			},
			rhs: []exports.Param{
				{Name: "ctx", Type: "context.Context"},
				{Name: "name", Type: "string"},
				{Name: "value", Type: "int"},
			},
			expected: true,
		},
		{
			name: "Name change only",
			lhs: []exports.Param{
				{Name: "oldName", Type: "string"},
				{Name: "newName", Type: "string"},
			},
			rhs: []exports.Param{
				{Name: "firstName", Type: "string"},
				{Name: "lastName", Type: "string"},
			},
			expected: false,
		},
		{
			name: "Type change",
			lhs: []exports.Param{
				{Name: "value", Type: "int"},
			},
			rhs: []exports.Param{
				{Name: "value", Type: "string"},
			},
			expected: false,
		},
		{
			name: "Different number of params",
			lhs: []exports.Param{
				{Name: "a", Type: "string"},
				{Name: "b", Type: "string"},
			},
			rhs: []exports.Param{
				{Name: "a", Type: "string"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := paramsEqual(tt.lhs, tt.rhs)
			if result != tt.expected {
				t.Errorf("paramsEqual() = %v, expected %v", result, tt.expected)
				t.Logf("lhs: %+v", tt.lhs)
				t.Logf("rhs: %+v", tt.rhs)
			}
		})
	}
}
