//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestContainerRegistryClient_CheckDockerV2Support(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.CheckDockerV2Support(ctx, nil)
	require.NoError(t, err)
}

func TestContainerRegistryClient_DeleteRepository(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteRepository(ctx, "hello-world", nil)
	require.NoError(t, err)
}

func TestContainerRegistryClient_DeleteTag(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteTag(ctx, "hello-world", "latest", nil)
	require.NoError(t, err)
}

func TestContainerRegistryClient_GetManifest(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.GetManifest(ctx, "hello-world", "latest", &ContainerRegistryClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	require.NoError(t, err)
	manifest, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NotEmpty(t, manifest)
	fmt.Printf("manifest content: %s\n", manifest)
}

func TestContainerRegistryClient_GetRepositoryProperties(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.GetRepositoryProperties(ctx, "hello-world", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.Name)
	fmt.Printf("repository name: %s\n", *res.Name)
	require.NotEmpty(t, *res.RegistryLoginServer)
	fmt.Printf("registry login server of the repository: %s\n", *res.RegistryLoginServer)
	require.NotEmpty(t, *res.ManifestCount)
	fmt.Printf("repository manifest count: %d\n", *res.ManifestCount)
}

func TestContainerRegistryClient_GetTagProperties(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.GetTagProperties(ctx, "hello-world", "latest", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.Tag.Name)
	fmt.Printf("tag name: %s\n", *res.Tag.Name)
	require.NotEmpty(t, *res.Tag.Digest)
	fmt.Printf("tag digest: %s\n", *res.Tag.Digest)
}

func TestContainerRegistryClient_NewListManifestsPager(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	pager := client.NewListManifestsPager("hello-world", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.Manifests.Manifests)
		for i, v := range page.Manifests.Manifests {
			fmt.Printf("manifest %d: %s\n", i+1, *v.Digest)
		}
	}
}

func TestContainerRegistryClient_NewListRepositoriesPager(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	pager := client.NewListRepositoriesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.Repositories.Repositories)
		for i, v := range page.Repositories.Repositories {
			fmt.Printf("repository %d: %s\n", i+1, *v)
		}
	}
}

func TestContainerRegistryClient_NewListTagsPager(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	pager := client.NewListTagsPager("hello-world", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.TagAttributeBases)
		for i, v := range page.TagAttributeBases {
			fmt.Printf("tag %d: %s\n", i+1, *v.Name)
		}
	}
}

func TestContainerRegistryClient_UpdateRepositoryProperties(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.UpdateRepositoryProperties(ctx, "hello-world", &ContainerRegistryClientUpdateRepositoryPropertiesOptions{Value: &RepositoryWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.NoError(t, err)
	require.Equal(t, *res.ContainerRepositoryProperties.ChangeableAttributes.CanWrite, false)
	res, err = client.UpdateRepositoryProperties(ctx, "hello-world", &ContainerRegistryClientUpdateRepositoryPropertiesOptions{Value: &RepositoryWriteableProperties{
		CanWrite: to.Ptr(true),
	},
	})
	require.NoError(t, err)
	require.Equal(t, *res.ContainerRepositoryProperties.ChangeableAttributes.CanWrite, true)
}

func TestContainerRegistryClient_UpdateTagProperties(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.UpdateTagProperties(ctx, "hello-world", "latest", &ContainerRegistryClientUpdateTagPropertiesOptions{Value: &TagWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.NoError(t, err)
	require.Equal(t, *res.Tag.ChangeableAttributes.CanWrite, false)
	res, err = client.UpdateTagProperties(ctx, "hello-world", "latest", &ContainerRegistryClientUpdateTagPropertiesOptions{Value: &TagWriteableProperties{
		CanWrite: to.Ptr(true),
	},
	})
	require.NoError(t, err)
	require.Equal(t, *res.Tag.ChangeableAttributes.CanWrite, true)
}

func TestContainerRegistryClient_UploadManifest(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryClientOptions{ClientOptions: options})
	require.NoError(t, err)
	getRes, err := client.GetManifest(ctx, "hello-world", "latest", &ContainerRegistryClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	require.NoError(t, err)
	manifest, err := io.ReadAll(getRes.Body)
	require.NoError(t, err)
	uploadRes, err := client.UploadManifest(ctx, "hello-world", "test", "application/vnd.docker.distribution.manifest.v2+json", streaming.NopCloser(bytes.NewReader(manifest)), nil)
	require.NoError(t, err)
	require.NotEmpty(t, *uploadRes.DockerContentDigest)
	fmt.Printf("uploaded manifest digest: %s\n", *uploadRes.DockerContentDigest)
}
