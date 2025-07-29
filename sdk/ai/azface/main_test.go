// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const RecordingDirectory = "sdk/ai/azface/testdata"

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		defaultOptions := getRecordingOptions(nil)
		proxy, err := recording.StartTestProxy(RecordingDirectory, defaultOptions)
		if err != nil {
			panic(err)
		}

		if err = configureTestProxy(*defaultOptions); err != nil {
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

func getRecordingOptions(testName *string) *recording.RecordingOptions {
	options := &recording.RecordingOptions{}
	return options
}

func configureTestProxy(options recording.RecordingOptions) error {
	// Configure test proxy settings if needed
	return nil
}