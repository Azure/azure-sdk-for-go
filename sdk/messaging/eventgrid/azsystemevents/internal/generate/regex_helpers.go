//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// DeleteType removes a type from models.go and it's associated marshalling functions from models_serde.go.
func DeleteType(model string) error {
	deleteModelTypeRE := regexp.MustCompile(fmt.Sprintf(`(?s)// %s -.+?\n}\n`, model))

	/* Example:
	// MarshalJSON implements the json.Marshaller interface for type AvsClusterEventData.
	*/
	deleteModelFuncsRE := regexp.MustCompile(fmt.Sprintf(`(?s)// (UnmarshalJSON|MarshalJSON) implements the (json\.Unmarshaller|json\.Marshaller) interface for type %s.+?\n}\n`, model))

	err := Replace(modelsGoFile, fmt.Sprintf("Remove model %s", model), func(s string) (string, error) {
		return deleteModelTypeRE.ReplaceAllString(s, ""), nil
	})

	if err != nil {
		return err
	}

	return Replace(modelsSerdeGoFile, fmt.Sprintf("Remove %s serde functions", model), func(s string) (string, error) {
		return deleteModelFuncsRE.ReplaceAllString(s, ""), nil
	})
}

// UseCustomUnpopulate replaces the serde code for a type with your own custom unpopulate function
// Used to handle the cases where we want to convert references to a custom error type to our package level `Error`
// type.
// symbolName - the name of the type and the field (ACSMessageReceivedEventData.Error)
// newFuncCall - the name of the Go function you're replacing the unmarshalling code with. This'll be a real Go
// function you have defined in custom_events.go, like `unmarshalInternalACSMessageChannelEventError`
func UseCustomUnpopulate(modelsSerdeGo string, symbolName string, newFuncCall string) string {
	parts := strings.Split(symbolName, ".")

	re := regexp.MustCompile(
		// ex: 'func (a *AcsAdvancedMessageReceivedEventData) UnmarshalJSON(data []byte) error {'
		`(?is)(func \([a-zA-Z]+? \*` + parts[0] + `\) UnmarshalJSON\(data \[\]byte\) error \{.+?` +
			// ex: 'err = unpopulate(val, "Content", &a.Content)'
			`err = )unpopulate(\(val, "` + parts[1] + `")`)

	return re.ReplaceAllString(modelsSerdeGo, "$1 "+newFuncCall+"$2")
}

func Replace(filename string, purpose string, repl func(string) (string, error)) error {
	data, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	newData, err := repl(string(data))

	if err != nil {
		return err
	}

	if newData == string(data) {
		return fmt.Errorf("no replacements made for %s", purpose)
	}

	return os.WriteFile(filename, []byte(newData), 0600)
}
