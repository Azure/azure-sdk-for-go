// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armapimanagement_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const (
	pathToPackage = "sdk/resourcemanager/apimanagement/armapimanagement/testdata"
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	if recording.GetRecordMode() == recording.LiveMode {
		return m.Run()
	}
	// Live tests are disabled; set AZURE_RECORD_MODE=live to run them.
	fmt.Println("Skipping: live tests are disabled, set AZURE_RECORD_MODE=live to run live tests")
	return 0
}
