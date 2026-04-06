// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

// StartProxy starts the test proxy with the path to store test recording file.
// It should be used in the module test preparation stage only once.
// It will return a delegate function to stop test proxy.
func StartProxy(pathToPackage string) func() {
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		proxy, err := recording.StartTestProxy(pathToPackage, nil)
		if err != nil {
			panic(fmt.Sprintf("Failed to start recording proxy: %v", err))
		}

		// sanitizer for any uuid string, e.g., subscriptionID
		err = recording.AddGeneralRegexSanitizer("00000000-0000-0000-0000-000000000000", `[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`, proxy.Options)
		if err != nil {
			panic(fmt.Sprintf("Failed to add uuid sanitizer: %v", err))
		}
		// consolidate resource group name for recording and playback
		err = recording.AddGeneralRegexSanitizer(recording.SanitizedValue, `go-sdk-test-\d+`, proxy.Options)
		if err != nil {
			panic(fmt.Sprintf("Failed to add resource group name sanitizer: %v", err))
		}
		// disable location header sanitizer
		err = recording.RemoveRegisteredSanitizers([]string{"AZSDK2003", "AZSDK2030"}, proxy.Options)
		if err != nil {
			panic(fmt.Sprintf("Failed to remove location header sanitizer: %v", err))
		}

		return func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(fmt.Sprintf("Failed to stop recording proxy: %v", err))
			}
		}
	}
	return func() {}
}

// StartRecording starts the recording with the path to store recording file.
// It will return a delegate function to stop recording.
func StartRecording(t *testing.T, pathToPackage string) func() {
	err := recording.Start(t, pathToPackage, nil)
	if err != nil {
		t.Fatalf("Failed to start recording: %v", err)
	}
	return func() { StopRecording(t) }
}

// StopRecording stops the recording.
func StopRecording(t *testing.T) {
	err := recording.Stop(t, nil)
	if err != nil {
		t.Fatalf("Failed to stop recording: %v", err)
	}
}
