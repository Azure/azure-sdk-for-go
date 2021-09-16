// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamespaceWithUserAgentOption(t *testing.T) {
	userAgent := "custom-user-agent"
	nsUserAgentOption := NamespaceWithUserAgent(userAgent)
	ns, err := NewNamespace(nsUserAgentOption)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%s/%s", rootUserAgent, userAgent), ns.getUserAgent())
}

func TestNamespaceWithoutUserAgentOption(t *testing.T) {
	userAgent := ""
	nsUserAgentOption := NamespaceWithUserAgent(userAgent)
	ns, err := NewNamespace(nsUserAgentOption)
	assert.Nil(t, err)
	assert.Equal(t, rootUserAgent, ns.getUserAgent())
}
