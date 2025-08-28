package typespec_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
	"github.com/stretchr/testify/assert"
)

func TestDebugParsing(t *testing.T) {
	yaml := `parameters:
  service-dir:
    default: "sdk/messaging/eventgrid"

options:
  "@azure-tools/typespec-go":
    module: "github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents"
    emitter-output-dir: "{output-dir}/"`

	// Create temporary file with YAML content
	tmpFile, err := os.CreateTemp("", "tspconfig_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yaml)
	assert.NoError(t, err)
	tmpFile.Close()

	// Parse config from file
	config, err := typespec.ParseTypeSpecConfig(tmpFile.Name())
	if err != nil {
		t.Logf("Parse error: %v", err)
		return
	}

	t.Logf("Module: %s", config.Options.GoConfig.Module)
	t.Logf("ServiceDir: %s", config.Options.GoConfig.ServiceDir)
	t.Logf("PackageDir: %s", config.Options.GoConfig.PackageDir)
	t.Logf("EmitterOutputDir: %s", config.Options.GoConfig.EmitterOutputDir)
	t.Logf("GetPackageRelativePath: %s", config.GetPackageRelativePath())
	t.Logf("GetModuleRelativePath: %s", config.GetModuleRelativePath())
}
