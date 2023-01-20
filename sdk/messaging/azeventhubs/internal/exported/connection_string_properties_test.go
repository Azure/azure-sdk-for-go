// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
)

const (
	namespace = "mynamespace"
	keyName   = "keyName"
	secret    = "superSecret="
	hubName   = "myhub"
)

func TestNewConnectionStringProperties(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		const happyConnStr = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret + ";EntityPath=" + hubName

		props, err := exported.NewConnectionStringProperties(happyConnStr)
		require.NoError(t, err)

		require.Equal(t, exported.ConnectionStringProperties{
			EventHubName:            hubName,
			Endpoint:                "sb://" + namespace + ".servicebus.windows.net/",
			FullyQualifiedNamespace: namespace + ".servicebus.windows.net",
			SharedAccessKeyName:     keyName,
			SharedAccessKey:         secret,
			SharedAccessSignature:   "",
		}, props)
	})

	t.Run("CaseIndifference", func(t *testing.T) {
		const lowerCase = "endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccesskeyName=" + keyName + ";sharedAccessKey=" + secret + ";Entitypath=" + hubName

		props, err := exported.NewConnectionStringProperties(lowerCase)
		require.NoError(t, err)

		require.Equal(t, exported.ConnectionStringProperties{
			EventHubName:            hubName,
			Endpoint:                "sb://" + namespace + ".servicebus.windows.net/",
			FullyQualifiedNamespace: namespace + ".servicebus.windows.net",
			SharedAccessKeyName:     keyName,
			SharedAccessKey:         secret,
			SharedAccessSignature:   "",
		}, props)
	})

	t.Run("NoEntityPath", func(t *testing.T) {
		const noEntityPath = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret

		props, err := exported.NewConnectionStringProperties(noEntityPath)
		require.NoError(t, err)

		require.Equal(t, exported.ConnectionStringProperties{
			EventHubName:            "",
			Endpoint:                "sb://" + namespace + ".servicebus.windows.net/",
			FullyQualifiedNamespace: namespace + ".servicebus.windows.net",
			SharedAccessKeyName:     keyName,
			SharedAccessKey:         secret,
			SharedAccessSignature:   "",
		}, props)
	})

	t.Run("EmbeddedSAS", func(t *testing.T) {
		const withEmbeddedSAS = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessSignature=SharedAccessSignature sr=" + namespace + ".servicebus.windows.net&sig=<base64-sig>&se=<expiry>&skn=<keyname>"

		props, err := exported.NewConnectionStringProperties(withEmbeddedSAS)
		require.NoError(t, err)

		require.Equal(t, exported.ConnectionStringProperties{
			EventHubName:            "",
			Endpoint:                "sb://" + namespace + ".servicebus.windows.net/",
			FullyQualifiedNamespace: namespace + ".servicebus.windows.net",
			SharedAccessKeyName:     "",
			SharedAccessKey:         "",
			SharedAccessSignature:   "SharedAccessSignature sr=" + namespace + ".servicebus.windows.net&sig=<base64-sig>&se=<expiry>&skn=<keyname>",
		}, props)
	})

	t.Run("WithoutEndpoint", func(t *testing.T) {
		_, err := exported.NewConnectionStringProperties("NoEndpoint=Blah")
		require.EqualError(t, err, "key \"Endpoint\" must not be empty")
	})

	t.Run("NoSASOrKeyName", func(t *testing.T) {
		_, err := exported.NewConnectionStringProperties("Endpoint=sb://" + namespace + ".servicebus.windows.net/")
		require.EqualError(t, err, "key \"SharedAccessKeyName\" must not be empty")
	})

	t.Run("NoSASOrKeyValue", func(t *testing.T) {
		const s = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessKeyName=" + keyName + ";EntityPath=" + hubName

		_, err := exported.NewConnectionStringProperties(s)
		require.EqualError(t, err, "key \"SharedAccessKey\" or \"SharedAccessSignature\" cannot both be empty")
	})
}
