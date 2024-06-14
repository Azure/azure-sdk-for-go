package typespec

import (
	"os"

	"github.com/jinzhu/copier"
)

func TempTspConfigFile(tsc *TypeSpecConfig, emit string, emitOption map[string]any) (*TypeSpecConfig, error) {
	tempFile, err := os.CreateTemp("", "temp-*-tspconfig.yaml")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())

	var tempTspConfig TypeSpecConfig
	err = copier.Copy(&tempTspConfig, tsc)
	if err != nil {
		return nil, err
	}

	tempTspConfig.Path = tempFile.Name()
	// tempTspConfig.OnlyEmit(emit)
	tempTspConfig.EditOptions(emit, emitOption, true)

	err = tempTspConfig.Write()
	if err != nil {
		return nil, err
	}

	return &tempTspConfig, nil
}
