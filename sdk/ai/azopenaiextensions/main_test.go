// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiextensions_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const RecordingDirectory = "sdk/ai/azopenaiextensions/testdata"

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		defaultOptions := getRecordingOptions(nil)
		proxy, err := recording.StartTestProxy(RecordingDirectory, defaultOptions)
		proxy.Cmd.Stderr = os.Stderr
		proxy.Cmd.Stdout = os.Stdout
		if err != nil {
			panic(err)
		}

		defer func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(err)
			}
		}()
	}
	os.Setenv("AOAI_OYD_ENDPOINT", os.Getenv("AOAI_ENDPOINT_USEAST"))
	os.Setenv("AOAI_OYD_MODEL", "gpt-4-0613")

	return m.Run()
}
