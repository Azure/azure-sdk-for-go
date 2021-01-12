package pipeline

import "os"

func ReadInput(inputPath string) (*GenerateInput, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewGenerateInputFrom(file)
}

func WriteOutput(outputPath string, output *GenerateOutput) error {
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
