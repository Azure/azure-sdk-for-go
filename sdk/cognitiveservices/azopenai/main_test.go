// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"os"
	"testing"
)

const RecordingDirectory = "sdk/cognitiveservices/azopenai/testdata"

func TestMain(m *testing.M) {
	initEnvVars()
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		proxy, err := recording.StartTestProxy(RecordingDirectory, nil)
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

	return m.Run()
}
