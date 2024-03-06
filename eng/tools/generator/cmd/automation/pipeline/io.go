// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package pipeline

import "os"

// ReadInput returns a *GenerateInput using the filepath.
func ReadInput(inputPath string) (*GenerateInput, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewGenerateInputFrom(file)
}

// WriteOutput writes the GenerateOutput to the specific file.
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
