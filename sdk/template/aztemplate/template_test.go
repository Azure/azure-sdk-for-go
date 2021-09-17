//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztemplate

import (
	"fmt"
	"testing"
)

func TestOutput(t *testing.T) {

	client := TemplateClient()

	fmt.Println(client.ClientVersion())
}
