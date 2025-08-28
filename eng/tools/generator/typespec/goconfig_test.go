// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package typespec_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
	"github.com/stretchr/testify/assert"
)

func TestGoEmitterOptionsAreValid(t *testing.T) {
	modules := []string{
		"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces",
		"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus",
	}

	for _, module := range modules {
		t.Run("Check "+module, func(t *testing.T) {
			// module format is correct
			goEmitOptions := typespec.GoEmitterOptions{
				Module: module,
			}
			err := goEmitOptions.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestGoEmitterOptionsValidate(t *testing.T) {
	goEmitOptions := typespec.GoEmitterOptions{
		Module: "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/xxx/armxxx",
	}
	err := goEmitOptions.Validate()
	assert.NoError(t, err)

	// module format is wrong
	goEmitOptions = typespec.GoEmitterOptions{
		Module: "github.com/Azure/azure-sdk-for-go/sdk/xxx/armxxx",
	}
	err = goEmitOptions.Validate()
	assert.EqualError(t, err, typespec.ErrModuleFormat.Error())

	goEmitOptions = typespec.GoEmitterOptions{
		Module: "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/xxx/xxx",
	}
	err = goEmitOptions.Validate()
	assert.EqualError(t, err, typespec.ErrModuleFormat.Error())

	// module is empty
	goEmitOptions = typespec.GoEmitterOptions{
		Module: "",
	}
	err = goEmitOptions.Validate()
	assert.EqualError(t, err, typespec.ErrModuleEmpty.Error())

	// containing-module is empty
	goEmitOptions = typespec.GoEmitterOptions{
		ContainingModule: "",
	}
	err = goEmitOptions.Validate()
	assert.EqualError(t, err, typespec.ErrModuleEmpty.Error())

	// empty config
	goEmitOptions = typespec.GoEmitterOptions{}
	err = goEmitOptions.Validate()
	assert.EqualError(t, err, typespec.ErrModuleEmpty.Error())
}
