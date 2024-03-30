package conn

//	MIT License
//
//	Copyright (c) Microsoft Corporation. All rights reserved.
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	namespace       = "mynamespace"
	keyName         = "keyName"
	secret          = "superSecret="
	hubName         = "myhub"
	happyConnStr    = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret + ";EntityPath=" + hubName
	noEntityPath    = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret
	lowerCase       = "endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccesskeyName=" + keyName + ";sharedAccessKey=" + secret + ";Entitypath=" + hubName
	withEmbeddedSAS = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessSignature=SharedAccessSignature sr=" + namespace + ".servicebus.windows.net&sig=<base64-sig>&se=<expiry>&skn=<keyname>"
	noEndpoint      = "NoEndpoint=Blah"
)

func TestParsedConnectionFromStr(t *testing.T) {
	parsed, err := ParseConnectionString(happyConnStr)
	if assert.NoError(t, err) {
		assert.Equal(t, namespace+".servicebus.windows.net", parsed.FullyQualifiedNamespace)
		assert.Equal(t, keyName, *parsed.SharedAccessKeyName)
		assert.Equal(t, secret, *parsed.SharedAccessKey)
		assert.Equal(t, hubName, *parsed.EntityPath)
		require.False(t, parsed.Emulator)
	}
}

func TestParsedConnectionFromStrCaseIndifference(t *testing.T) {
	parsed, err := ParseConnectionString(lowerCase)
	if assert.NoError(t, err) {
		assert.Equal(t, namespace+".servicebus.windows.net", parsed.FullyQualifiedNamespace)
		assert.Equal(t, keyName, *parsed.SharedAccessKeyName)
		assert.Equal(t, secret, *parsed.SharedAccessKey)
		assert.Equal(t, hubName, *parsed.EntityPath)
	}
}

func TestParsedConnectionFromStrWithoutEntityPath(t *testing.T) {
	parsed, err := ParseConnectionString(noEntityPath)
	if assert.NoError(t, err) {
		assert.Equal(t, namespace+".servicebus.windows.net", parsed.FullyQualifiedNamespace)
		assert.Equal(t, keyName, *parsed.SharedAccessKeyName)
		assert.Equal(t, secret, *parsed.SharedAccessKey)
		require.Nil(t, parsed.EntityPath)
	}
}

func TestParsedConnectionFromStrWithEmbeddedSAS(t *testing.T) {
	parsed, err := ParseConnectionString(withEmbeddedSAS)
	require.NoError(t, err)

	require.Equal(t, ConnectionStringProperties{
		Endpoint:                "sb://" + namespace + ".servicebus.windows.net/",
		FullyQualifiedNamespace: namespace + ".servicebus.windows.net",
		SharedAccessSignature:   to.Ptr("SharedAccessSignature sr=" + namespace + ".servicebus.windows.net&sig=<base64-sig>&se=<expiry>&skn=<keyname>"),
	}, parsed)

}

func TestFailedParsedConnectionFromStrWithoutEndpoint(t *testing.T) {
	_, err := ParseConnectionString(noEndpoint)
	assert.Error(t, err)
}

func TestUseDevelopmentEmulatorProperty(t *testing.T) {
	cs := "Endpoint=sb://localhost:6765;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret + ";UseDevelopmentEmulator=true"
	parsed, err := ParseConnectionString(cs)
	require.NoError(t, err)
	require.True(t, parsed.Emulator)
	require.Equal(t, "sb://localhost:6765", parsed.Endpoint)

	// emulator can give connection strings that have a trailing ';'
	cs = "Endpoint=sb://localhost:6765;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret + ";UseDevelopmentEmulator=true;"
	parsed, err = ParseConnectionString(cs)
	require.NoError(t, err)
	require.True(t, parsed.Emulator)
	require.Equal(t, "sb://localhost:6765", parsed.Endpoint)

	// UseDevelopmentEmulator only works for localhost
	cs = "Endpoint=sb://myserver.com:6765;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret + ";UseDevelopmentEmulator=true"
	parsed, err = ParseConnectionString(cs)
	require.EqualError(t, err, "UseEmulator=true can only be used with sb://localhost:<port>, not sb://myserver.com:6765")

	// there's no reason for a person to pass False, but it's allowed.
	// If they're not using the dev emulator then there's no special behavior, it's like a normal connection string
	cs = "Endpoint=sb://localhost:6765;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret + ";UseDevelopmentEmulator=false"
	parsed, err = ParseConnectionString(cs)
	require.NoError(t, err)
	require.False(t, parsed.Emulator)
	require.Equal(t, "sb://localhost:6765", parsed.Endpoint)
}
