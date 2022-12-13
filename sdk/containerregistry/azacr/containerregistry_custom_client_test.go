//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContainerRegistryClient_DeleteManifest(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteManifest(ctx, "hello-world", "latest", nil)
	require.NoError(t, err)
	_, err = client.DeleteManifest(ctx, "hello-world-test", "sha256:f54a58bc1aac5ea1a25d796ae155dc228b3f0e11d046ae276b39c4bf2f13d8c4", nil)
	require.NoError(t, err)
}

func TestContainerRegistryClient_GetDigest(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:f54a58bc1aac5ea1a25d796ae155dc228b3f0e11d046ae276b39c4bf2f13d8c4"
	tag := "latest"
	digestReturn, err := client.GetDigest(ctx, "hello-world", "sha256:f54a58bc1aac5ea1a25d796ae155dc228b3f0e11d046ae276b39c4bf2f13d8c4", nil)
	require.NoError(t, err)
	require.Equal(t, digestReturn, digest)
	tagReturn, err := client.GetDigest(ctx, "hello-world", tag, nil)
	require.NoError(t, err)
	require.Equal(t, tagReturn, digest)
}

func TestContainerRegistryClient_GetManifestProperties(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:f54a58bc1aac5ea1a25d796ae155dc228b3f0e11d046ae276b39c4bf2f13d8c4"
	tag := "latest"
	digestRes, err := client.GetManifestProperties(ctx, "hello-world", digest, nil)
	require.NoError(t, err)
	require.Equal(t, *digestRes.Manifest.Digest, digest)
	tagRes, err := client.GetManifestProperties(ctx, "hello-world", tag, nil)
	require.NoError(t, err)
	require.Equal(t, *tagRes.Manifest.Digest, digest)
}

func TestContainerRegistryClient_UpdateManifestProperties(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:f54a58bc1aac5ea1a25d796ae155dc228b3f0e11d046ae276b39c4bf2f13d8c4"
	tag := "latest"
	res, err := client.UpdateManifestProperties(ctx, "hello-world", tag, &ContainerRegistryClientUpdateManifestPropertiesOptions{Value: &ManifestWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.NoError(t, err)
	require.Equal(t, *res.Manifest.ChangeableAttributes.CanWrite, false)
	res, err = client.UpdateManifestProperties(ctx, "hello-world", digest, &ContainerRegistryClientUpdateManifestPropertiesOptions{Value: &ManifestWriteableProperties{
		CanWrite: to.Ptr(true),
	},
	})
	require.NoError(t, err)
	require.Equal(t, *res.Manifest.ChangeableAttributes.CanWrite, true)
}
