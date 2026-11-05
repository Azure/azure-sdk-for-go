// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcognitiveservices_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
)

const pathToPackage = "sdk/resourcemanager/cognitiveservices/armcognitiveservices/testdata"

func TestMain(m *testing.M) {
	stopProxy := testutil.StartProxy(pathToPackage)
	code := m.Run()
	stopProxy()
	os.Exit(code)
}
