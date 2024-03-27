//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents/internal/gopls"
)

// DeleteType removes a type from models.go and it's associated marshalling functions from models_serde.go.
func DeleteType(model string) error {
	deleteModelTypeRE := regexp.MustCompile(fmt.Sprintf(`(?s)// %s -.+?\n}\n`, model))
	deleteModelFuncsRE := regexp.MustCompile(fmt.Sprintf(`(?s)// (UnmarshalJSON) implements the (json\.Unmarshaller) interface for type %s.+?\n}\n`, model))

	modelsGoBytes, err := os.ReadFile(modelsGoFile)

	if err != nil {
		return err
	}

	newModelsGoBytes := deleteModelTypeRE.ReplaceAll(modelsGoBytes, nil)

	if bytes.Equal(modelsGoBytes, newModelsGoBytes) {
		return fmt.Errorf("no match or changes in %s for model %s. Is the type name correct?", modelsGoFile, model)
	}

	if err := os.WriteFile(modelsGoFile, newModelsGoBytes, 0700); err != nil {
		return err
	}

	modelsSerdeGoBytes, err := os.ReadFile(modelsSerdeGoFile)

	if err != nil {
		return err
	}

	newModelsSerdeGoBytes := deleteModelFuncsRE.ReplaceAll(modelsSerdeGoBytes, nil)

	if bytes.Equal(modelsSerdeGoBytes, newModelsSerdeGoBytes) {
		return fmt.Errorf("no match or changes in %s for model %s. Is the type name correct?", modelsGoFile, model)
	}

	return os.WriteFile(modelsSerdeGoFile, newModelsSerdeGoBytes, 0700)
}

// SwapType changes the declared type for the symbol, which is expected to be a field
func SwapType(sym *gopls.Symbol, newType string) error {
	if sym.Type != gopls.SymbolTypeField {
		return fmt.Errorf("can only swap types for a field, not a %s", sym.Type)
	}

	lineNum, err := sym.StartLine()

	if err != nil {
		return err
	}

	file, err := os.Open(sym.File)

	if err != nil {
		return err
	}

	defer file.Close()

	b := strings.Builder{}
	scanner := bufio.NewScanner(file)

	i := int64(0)
	for scanner.Scan() {
		i++

		if lineNum == i {
			// splitting something like "FieldName string"
			parts := strings.SplitN(scanner.Text(), " ", 2)
			b.WriteString(parts[0] + " " + newType + "\n")
		} else {
			b.WriteString(scanner.Text() + "\n")
		}
	}

	if err := file.Close(); err != nil {
		return err
	}

	return os.WriteFile(sym.File, []byte(b.String()), 0600)
}

func UseCustomUnpopulate(filename string, symbolName string, newFuncCall string) error {
	parts := strings.Split(symbolName, ".")
	buff, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	re := regexp.MustCompile(
		// ex: 'func (a *AcsAdvancedMessageReceivedEventData) UnmarshalJSON(data []byte) error {'
		`(?s)(func \([a-zA-Z]+? \*` + parts[0] + `\) UnmarshalJSON\(data \[\]byte\) error \{.+?` +
			// ex: 'err = unpopulate(val, "Content", &a.Content)'
			`err = )unpopulate(\(val, "` + parts[1] + `")`)

	newBuff := re.ReplaceAll(buff, []byte("$1 "+newFuncCall+"$2"))

	return os.WriteFile(filename, newBuff, 0600)
}
