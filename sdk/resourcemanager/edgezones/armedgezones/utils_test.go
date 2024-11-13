// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armedgezones_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
)

const (
	pathToPackage = "sdk/resourcemanager/edgezones/armedgezones/testdata"
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	f := testutil.StartProxy(pathToPackage)
	defer f()
	return m.Run()
}
