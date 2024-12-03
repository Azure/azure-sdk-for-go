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
			goOption := map[string]any{
				"module": module,
			}
			goEmitOptions, err := typespec.NewGoEmitterOptions(goOption)
			assert.NoError(t, err)
			err = goEmitOptions.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestGoEmitterOptionsValidate(t *testing.T) {
	goOption := map[string]any{
		"module": "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/xxx/armxxx",
	}
	goEmitOptions, err := typespec.NewGoEmitterOptions(goOption)
	assert.NoError(t, err)
	err = goEmitOptions.Validate()
	assert.NoError(t, err)

	// module format is wrong
	goOption = map[string]any{
		"module": "github.com/Azure/azure-sdk-for-go/sdk/xxx/armxxx",
	}
	goEmitOptions, err = typespec.NewGoEmitterOptions(goOption)
	assert.NoError(t, err)
	err = goEmitOptions.Validate()
	assert.EqualError(t, err, typespec.ErrModuleFormat.Error())

	goOption = map[string]any{
		"module": "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/xxx/xxx",
	}
	goEmitOptions, err = typespec.NewGoEmitterOptions(goOption)
	assert.NoError(t, err)
	err = goEmitOptions.Validate()
	assert.EqualError(t, err, typespec.ErrModuleFormat.Error())

	// module is empty
	goOption = map[string]any{
		"module": "",
	}
	goEmitOptions, err = typespec.NewGoEmitterOptions(goOption)
	assert.NoError(t, err)
	err = goEmitOptions.Validate()
	assert.EqualError(t, err, typespec.ErrModuleEmpty.Error())
}
