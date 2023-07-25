// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	initEnvVars()
	os.Exit(m.Run())
}
