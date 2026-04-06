// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"testing"
)

func TestStartStopRecording(t *testing.T) {
	stop := StartRecording(t, pathToPackage)
	defer stop()
}
