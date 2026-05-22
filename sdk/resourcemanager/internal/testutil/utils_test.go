// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"os"
	"testing"
)

const (
	pathToPackage = "sdk/resourcemanager/internal/testdata"
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	f := StartProxy(pathToPackage)
	defer f()
	return m.Run()
}
