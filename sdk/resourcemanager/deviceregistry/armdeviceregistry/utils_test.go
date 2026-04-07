// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdeviceregistry_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const (
	pathToPackage = "sdk/resourcemanager/deviceregistry/armdeviceregistry/testdata"
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	if recording.GetRecordMode() == recording.LiveMode {
		return m.Run()
	}
	// No recordings (assets.json) exist yet; skip in playback/record mode.
	fmt.Println("Skipping: no recordings available, set AZURE_RECORD_MODE=live to run live tests")
	return 0
}
