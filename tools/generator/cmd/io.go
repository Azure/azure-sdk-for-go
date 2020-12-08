package cmd

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/tools/generator/pipeline"
)

func readInputFrom(inputPath string) (*pipeline.GenerateInput, error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	return pipeline.NewGenerateInputFrom(inputFile)
}

func writeOutputTo(outputPath string, output *pipeline.GenerateOutput) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := output.WriteTo(file); err != nil {
		return err
	}
	return nil
}
