// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"testing"
)

func TestContainerCreateAccessContainer(t *testing.T) {
	// TODO
	cred, err := NewSharedKeyCredential("foo", "Kg==")
	if err != nil {
		t.Fatal(err)
	}

	NewTableClient("https://foo", cred, &TableClientOptions{})
}
