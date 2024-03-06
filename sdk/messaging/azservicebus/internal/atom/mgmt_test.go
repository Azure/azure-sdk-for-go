// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueueUnmarshal(t *testing.T) {
	bytes, err := os.ReadFile("testdata/queue.xml")
	require.NoError(t, err)

	var env *QueueEnvelope
	err = xml.Unmarshal(bytes, &env)
	require.NoError(t, err)

	require.Equal(t, "test-queue-name", env.Title)
	desc := env.Content.QueueDescription

	// authorization rules is a bit special in that it's a collection with an xmlns that has to be written out
	// with its type name.
	require.NotEmpty(t, desc.AuthorizationRules)

	authRule := desc.AuthorizationRules[0]
	require.Equal(t, "SharedAccessAuthorizationRule", authRule.Type)
	require.Equal(t, "redacted-primary-key", *authRule.PrimaryKey)
	require.Equal(t, "redacted-secondary-key", *authRule.SecondaryKey)
	require.Equal(t, "TestSharedKeyName", *authRule.KeyName)
	require.Equal(t, "2022-01-12T01:28:26.1670445Z", authRule.CreatedTime.Format(time.RFC3339Nano))
	require.Equal(t, "2022-02-12T01:28:26.1670445Z", authRule.ModifiedTime.Format(time.RFC3339Nano))
	require.Equal(t, []string{"Manage", "Listen", "Send"}, authRule.Rights)

	indentedXML, err := xml.MarshalIndent(env, "  ", "  ")
	require.NoError(t, err)

	fmt.Printf("%s\n", indentedXML)
}
