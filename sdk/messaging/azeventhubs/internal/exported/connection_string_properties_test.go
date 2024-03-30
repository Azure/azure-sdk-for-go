// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
)

var (
	namespace = "mynamespace"
	keyName   = "keyName"
	secret    = "superSecret="
	hubName   = "myhub"
)

func TestNewConnectionStringProperties(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		var happyConnStr = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret + ";EntityPath=" + hubName

		props, err := exported.ParseConnectionString(happyConnStr)
		require.NoError(t, err)

		require.Equal(t, exported.ConnectionStringProperties{
			EntityPath:              &hubName,
			Endpoint:                "sb://" + namespace + ".servicebus.windows.net/",
			FullyQualifiedNamespace: namespace + ".servicebus.windows.net",
			SharedAccessKeyName:     &keyName,
			SharedAccessKey:         &secret,
			SharedAccessSignature:   nil,
			Emulator:                false,
		}, props)
	})

	t.Run("CaseIndifference", func(t *testing.T) {
		var lowerCase = "endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccesskeyName=" + keyName + ";sharedAccessKey=" + secret + ";Entitypath=" + hubName

		props, err := exported.ParseConnectionString(lowerCase)
		require.NoError(t, err)

		require.Equal(t, exported.ConnectionStringProperties{
			EntityPath:              &hubName,
			Endpoint:                "sb://" + namespace + ".servicebus.windows.net/",
			FullyQualifiedNamespace: namespace + ".servicebus.windows.net",
			SharedAccessKeyName:     &keyName,
			SharedAccessKey:         &secret,
			SharedAccessSignature:   nil,
		}, props)
	})

	t.Run("NoEntityPath", func(t *testing.T) {
		var noEntityPath = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret

		props, err := exported.ParseConnectionString(noEntityPath)
		require.NoError(t, err)

		require.Equal(t, exported.ConnectionStringProperties{
			EntityPath:              nil,
			Endpoint:                "sb://" + namespace + ".servicebus.windows.net/",
			FullyQualifiedNamespace: namespace + ".servicebus.windows.net",
			SharedAccessKeyName:     &keyName,
			SharedAccessKey:         &secret,
			SharedAccessSignature:   nil,
		}, props)
	})

	t.Run("EmbeddedSAS", func(t *testing.T) {
		var withEmbeddedSAS = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessSignature=SharedAccessSignature sr=" + namespace + ".servicebus.windows.net&sig=<base64-sig>&se=<expiry>&skn=<keyname>"

		props, err := exported.ParseConnectionString(withEmbeddedSAS)
		require.NoError(t, err)

		require.Equal(t, exported.ConnectionStringProperties{
			EntityPath:              nil,
			Endpoint:                "sb://" + namespace + ".servicebus.windows.net/",
			FullyQualifiedNamespace: namespace + ".servicebus.windows.net",
			SharedAccessKeyName:     nil,
			SharedAccessKey:         nil,
			SharedAccessSignature:   to.Ptr("SharedAccessSignature sr=" + namespace + ".servicebus.windows.net&sig=<base64-sig>&se=<expiry>&skn=<keyname>"),
		}, props)
	})

	t.Run("WithoutEndpoint", func(t *testing.T) {
		_, err := exported.ParseConnectionString("NoEndpoint=Blah")
		require.EqualError(t, err, "key \"Endpoint\" must not be empty")
	})

	t.Run("NoSASOrKeyName", func(t *testing.T) {
		_, err := exported.ParseConnectionString("Endpoint=sb://" + namespace + ".servicebus.windows.net/")
		require.EqualError(t, err, "key \"SharedAccessKeyName\" must not be empty")
	})

	t.Run("NoSASOrKeyValue", func(t *testing.T) {
		var s = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessKeyName=" + keyName + ";EntityPath=" + hubName

		_, err := exported.ParseConnectionString(s)
		require.EqualError(t, err, "key \"SharedAccessKey\" or \"SharedAccessSignature\" cannot both be empty")
	})

	t.Run("UseDevelopmentEmulator", func(t *testing.T) {
		cs := "Endpoint=sb://localhost:6765;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret + ";UseDevelopmentEmulator=true"
		parsed, err := exported.ParseConnectionString(cs)
		require.NoError(t, err)
		require.True(t, parsed.Emulator)
		require.Equal(t, "sb://localhost:6765", parsed.Endpoint)

		// emulator can give connection strings that have a trailing ';'
		cs = "Endpoint=sb://localhost:6765;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret + ";UseDevelopmentEmulator=true;"
		parsed, err = exported.ParseConnectionString(cs)
		require.NoError(t, err)
		require.True(t, parsed.Emulator)
		require.Equal(t, "sb://localhost:6765", parsed.Endpoint)

		// UseDevelopmentEmulator only works for localhost
		cs = "Endpoint=sb://myserver.com:6765;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret + ";UseDevelopmentEmulator=true"
		parsed, err = exported.ParseConnectionString(cs)
		require.EqualError(t, err, "UseEmulator=true can only be used with sb://localhost:<port>, not sb://myserver.com:6765")

		// there's no reason for a person to pass False, but it's allowed.
		// If they're not using the dev emulator then there's no special behavior, it's like a normal connection string
		cs = "Endpoint=sb://localhost:6765;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret + ";UseDevelopmentEmulator=false"
		parsed, err = exported.ParseConnectionString(cs)
		require.NoError(t, err)
		require.False(t, parsed.Emulator)
		require.Equal(t, "sb://localhost:6765", parsed.Endpoint)
	})
}
