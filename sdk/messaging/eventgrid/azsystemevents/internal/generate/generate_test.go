//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents/internal/gopls"
	"github.com/stretchr/testify/require"
)

func checkGopls(t *testing.T) {
	cmd := exec.Command("gopls", "--help")

	if err := cmd.Run(); err != nil {
		t.Skipf("Skipping gopls based test since gopls is not on the path. Install with `go install golang.org/x/tools/gopls@latest`")
	}
}

func TestUseCustomUnpopulate(t *testing.T) {
	goCode := `// UnmarshalJSON implements the json.Unmarshaller interface for type MyTypeName.
func (a *MyTypeName) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", a, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "context":
			err = unpopulate(val, "Context", &a.Context)
			delete(rawMsg, key)
		case "error":
			err = unpopulate(val, "Error", &a.Error)
			delete(rawMsg, key)
		case "from":
			err = unpopulate(val, "From", &a.From)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", a, err)
		}
	}
	return nil
}
`

	expectedGoCode := `// UnmarshalJSON implements the json.Unmarshaller interface for type MyTypeName.
func (a *MyTypeName) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", a, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "context":
			err = unpopulate(val, "Context", &a.Context)
			delete(rawMsg, key)
		case "error":
			err =  customUnmarshaller(val, "Error", &a.Error)
			delete(rawMsg, key)
		case "from":
			err = unpopulate(val, "From", &a.From)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", a, err)
		}
	}
	return nil
}
`

	tmpFileName := writeTempFile(t, goCode)

	err := UseCustomUnpopulate(tmpFileName, "MyTypeName.Error", "customUnmarshaller")
	require.NoError(t, err)

	buff, err := os.ReadFile(tmpFileName)
	require.NoError(t, err)

	require.Equal(t, expectedGoCode, string(buff))
}

func TestSwapType(t *testing.T) {
	checkGopls(t)

	goCode := `
package main

type MyStruct struct {
	TestField string
	SomeOtherField string
}
`

	expectedGoCode := `
package main

type MyStruct struct {
	TestField *ANewType
	SomeOtherField string
}
`

	tmpFileName := writeTempFile(t, goCode)

	symbols, err := gopls.Symbols(tmpFileName)
	require.NoError(t, err)
	require.NotEmpty(t, symbols)

	sym := symbols.Get("MyStruct.TestField")
	err = SwapType(sym, "*ANewType")
	require.NoError(t, err)

	newBuff, err := os.ReadFile(tmpFileName)
	require.NoError(t, err)

	require.Equal(t, expectedGoCode, string(newBuff))
}

func writeTempFile(t *testing.T, contents string) string {
	tmpFile, err := os.CreateTemp("", "tempfile*.go")
	require.NoError(t, err)

	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	_, err = tmpFile.Write([]byte(contents))
	require.NoError(t, err)

	err = tmpFile.Close()
	require.NoError(t, err)

	return tmpFile.Name()
}
