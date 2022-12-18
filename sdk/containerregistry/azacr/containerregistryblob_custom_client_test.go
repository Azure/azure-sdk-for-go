//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	"bytes"
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestContainerRegistryBlobClient_UploadBlob(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryBlobClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryBlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	getRes, err := client.GetBlob(ctx, "hello-world", digest, nil)
	require.NoError(t, err)
	blob, err := io.ReadAll(getRes.Body)
	require.NoError(t, err)
	res, err := client.UploadBlob(ctx, "hello-world", streaming.NopCloser(bytes.NewReader(blob)), nil)
	require.NoError(t, err)
	require.NotEmpty(t, res.Digest)
}
