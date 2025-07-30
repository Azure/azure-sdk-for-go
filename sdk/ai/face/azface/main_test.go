// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const RecordingDirectory = "sdk/ai/face/azface/testdata"

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	recordMode := recording.GetRecordMode()
	
	// Only start test proxy for recording mode or if we have recordings for playback
	if recordMode == recording.RecordingMode {
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
	} else if recordMode == recording.PlaybackMode {
		// For playback mode, only start proxy if we have recordings
		// This allows the tests to run without recordings for initial development
		if _, err := os.Stat("testdata"); err == nil {
			proxy, err := recording.StartTestProxy(RecordingDirectory, nil)
			if err != nil {
				// If we can't start the proxy (e.g., no recordings), just skip proxy setup
				// The individual tests will handle this gracefully
				return m.Run()
			}

			defer func() {
				err := recording.StopTestProxy(proxy)
				if err != nil {
					panic(err)
				}
			}()
		}
	}

	return m.Run()
}