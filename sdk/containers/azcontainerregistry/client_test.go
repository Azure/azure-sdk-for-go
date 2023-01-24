//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

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

func TestClient_DeleteManifest(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	resp, err := client.GetTagProperties(ctx, "hello-world", "latest", nil)
	require.NoError(t, err)
	_, err = client.DeleteManifest(ctx, "hello-world", *resp.Tag.Digest, nil)
	require.NoError(t, err)
	_, err = client.DeleteManifest(ctx, "hello-world-test", "sha256:f54a58bc1aac5ea1a25d796ae155dc228b3f0e11d046ae276b39c4bf2f13d8c4", nil)
	require.NoError(t, err)
}

func TestClient_DeleteManifest_wrongDigest(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteManifest(ctx, "hello-world", "error-digest", nil)
	require.Error(t, err)
}

func TestClient_DeleteRepository(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteRepository(ctx, "hello-world", nil)
	require.NoError(t, err)
}

func TestClient_DeleteTag(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteTag(ctx, "hello-world", "latest", nil)
	require.NoError(t, err)
}

func TestClient_GetManifest(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.GetManifest(ctx, "hello-world", "latest", &ClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	require.NoError(t, err)
	manifest, err := io.ReadAll(res.ManifestData)
	require.NoError(t, err)
	require.NotEmpty(t, manifest)
	fmt.Printf("manifest content: %s\n", manifest)
}

func TestClient_GetManifest_wrongTag(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetManifest(ctx, "hello-world", "wrong-tag", &ClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	require.Error(t, err)
}

func TestClient_GetManifestProperties(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:f54a58bc1aac5ea1a25d796ae155dc228b3f0e11d046ae276b39c4bf2f13d8c4"
	tag := "latest"
	digestRes, err := client.GetManifestProperties(ctx, "hello-world", digest, nil)
	require.NoError(t, err)
	require.Equal(t, *digestRes.Manifest.Digest, digest)
	resp, err := client.GetTagProperties(ctx, "hello-world", tag, nil)
	require.NoError(t, err)
	tagRes, err := client.GetManifestProperties(ctx, "hello-world", *resp.Tag.Digest, nil)
	require.NoError(t, err)
	require.Equal(t, *tagRes.Manifest.Digest, digest)
}

func TestClient_GetManifestProperties_wrongDigest(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetManifestProperties(ctx, "hello-world", "wrong-digest", nil)
	require.Error(t, err)
}

func TestClient_GetRepositoryProperties(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
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

func TestClient_GetRepositoryProperties_wrongName(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetRepositoryProperties(ctx, "wrong-name", nil)
	require.Error(t, err)
}

func TestClient_GetTagProperties(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.GetTagProperties(ctx, "hello-world", "latest", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.Tag.Name)
	fmt.Printf("tag name: %s\n", *res.Tag.Name)
	require.NotEmpty(t, *res.Tag.Digest)
	fmt.Printf("tag digest: %s\n", *res.Tag.Digest)
}

func TestClient_GetTagProperties_wrongTag(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetTagProperties(ctx, "hello-world", "wrong-tag", nil)
	require.Error(t, err)
}

func TestClient_NewListManifestsPager(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	pager := client.NewListManifestsPager("hello-world", &ClientListManifestsOptions{
		MaxNum: to.Ptr[int32](1),
	})
	pages := 0
	items := 0
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.Manifests.Attributes)
		pages++
		for i, v := range page.Manifests.Attributes {
			fmt.Printf("page %d manifest %d: %s\n", pages, i+1, *v.Digest)
			items++
		}
	}
	require.Equal(t, pages, 2)
	require.Equal(t, items, 2)

	pager = client.NewListManifestsPager("hello-world", &ClientListManifestsOptions{
		OrderBy: to.Ptr(ArtifactManifestOrderByLastUpdatedOnDescending),
	})
	var descendingItems []*ManifestAttributes
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.Manifests.Attributes)
		for i, v := range page.Manifests.Attributes {
			fmt.Printf("manifest order by last updated on descending %d: %s\n", i+1, *v.Digest)
			descendingItems = append(descendingItems, v)
		}
	}
	pager = client.NewListManifestsPager("hello-world", &ClientListManifestsOptions{
		OrderBy: to.Ptr(ArtifactManifestOrderByLastUpdatedOnAscending),
	})
	var ascendingItems []*ManifestAttributes
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.Manifests.Attributes)
		for i, v := range page.Manifests.Attributes {
			fmt.Printf("manifest order by last updated on descending %d: %s\n", i+1, *v.Digest)
			ascendingItems = append(ascendingItems, v)
		}
	}
	for i := range descendingItems {
		require.Equal(t, descendingItems[i].Digest, ascendingItems[len(ascendingItems)-1-i].Digest)
	}
}

func TestClient_NewListManifestsPager_wrongRepositoryName(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	pager := client.NewListManifestsPager("wrong-name", nil)
	for pager.More() {
		_, err := pager.NextPage(ctx)
		require.Error(t, err)
		break
	}
}

func TestClient_NewListRepositoriesPager(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	pager := client.NewListRepositoriesPager(&ClientListRepositoriesOptions{
		MaxNum: to.Ptr[int32](1),
	})
	pages := 0
	items := 0
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.Repositories.Names)
		pages++
		for i, v := range page.Repositories.Names {
			fmt.Printf("page %d repository %d: %s\n", pages, i+1, *v)
			items++
		}
	}
	require.Equal(t, pages, 3)
	require.Equal(t, items, 3)
}

func TestClient_NewListTagsPager(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	pager := client.NewListTagsPager("hello-world", &ClientListTagsOptions{
		MaxNum: to.Ptr[int32](1),
	})
	pages := 0
	items := 0
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.Tags)
		pages++
		require.Equal(t, len(page.Tags), 1)
		for i, v := range page.Tags {
			fmt.Printf("page %d tag %d: %s\n", pages, i+1, *v.Name)
			items++
		}
	}
	require.Equal(t, pages, 3)
	require.Equal(t, items, 3)

	pager = client.NewListTagsPager("hello-world", &ClientListTagsOptions{
		OrderBy: to.Ptr(ArtifactTagOrderByLastUpdatedOnDescending),
	})
	var descendingItems []*TagAttributes
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.Tags)
		for i, v := range page.Tags {
			fmt.Printf("tag order by last updated on descending %d: %s\n", i+1, *v.Name)
			descendingItems = append(descendingItems, v)
		}
	}
	pager = client.NewListTagsPager("hello-world", &ClientListTagsOptions{
		OrderBy: to.Ptr(ArtifactTagOrderByLastUpdatedOnAscending),
	})
	var ascendingItems []*TagAttributes
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.Tags)
		for i, v := range page.Tags {
			fmt.Printf("tag order by last updated on descending %d: %s\n", i+1, *v.Name)
			ascendingItems = append(ascendingItems, v)
		}
	}
	for i := range descendingItems {
		require.Equal(t, descendingItems[i].Name, ascendingItems[len(ascendingItems)-1-i].Name)
	}
}

func TestClient_NewListTagsPager_wrongRepositoryName(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	pager := client.NewListTagsPager("wrong-name", nil)
	for pager.More() {
		_, err := pager.NextPage(ctx)
		require.Error(t, err)
		break
	}
}

func TestClient_UpdateManifestProperties(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:f54a58bc1aac5ea1a25d796ae155dc228b3f0e11d046ae276b39c4bf2f13d8c4"
	tag := "latest"
	resp, err := client.GetTagProperties(ctx, "hello-world", tag, nil)
	require.NoError(t, err)
	res, err := client.UpdateManifestProperties(ctx, "hello-world", *resp.Tag.Digest, &ClientUpdateManifestPropertiesOptions{Value: &ManifestWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.NoError(t, err)
	require.Equal(t, *res.Manifest.ChangeableAttributes.CanWrite, false)
	res, err = client.UpdateManifestProperties(ctx, "hello-world", digest, &ClientUpdateManifestPropertiesOptions{Value: &ManifestWriteableProperties{
		CanWrite: to.Ptr(true),
	},
	})
	require.NoError(t, err)
	require.Equal(t, *res.Manifest.ChangeableAttributes.CanWrite, true)
}

func TestClient_UpdateManifestProperties_wrongDigest(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetTagProperties(ctx, "hello-world", "wrong-digest", nil)
	require.Error(t, err)
}

func TestClient_UpdateRepositoryProperties(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.UpdateRepositoryProperties(ctx, "hello-world", &ClientUpdateRepositoryPropertiesOptions{Value: &RepositoryWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.NoError(t, err)
	require.Equal(t, *res.ContainerRepositoryProperties.ChangeableAttributes.CanWrite, false)
	res, err = client.UpdateRepositoryProperties(ctx, "hello-world", &ClientUpdateRepositoryPropertiesOptions{Value: &RepositoryWriteableProperties{
		CanWrite: to.Ptr(true),
	},
	})
	require.NoError(t, err)
	require.Equal(t, *res.ContainerRepositoryProperties.ChangeableAttributes.CanWrite, true)
}

func TestClient_UpdateRepositoryProperties_wrongRepository(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.UpdateRepositoryProperties(ctx, "wrong-repository", &ClientUpdateRepositoryPropertiesOptions{Value: &RepositoryWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.Error(t, err)
}

func TestClient_UpdateTagProperties(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.UpdateTagProperties(ctx, "hello-world", "latest", &ClientUpdateTagPropertiesOptions{Value: &TagWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.NoError(t, err)
	require.Equal(t, *res.Tag.ChangeableAttributes.CanWrite, false)
	res, err = client.UpdateTagProperties(ctx, "hello-world", "latest", &ClientUpdateTagPropertiesOptions{Value: &TagWriteableProperties{
		CanWrite: to.Ptr(true),
	},
	})
	require.NoError(t, err)
	require.Equal(t, *res.Tag.ChangeableAttributes.CanWrite, true)
}

func TestClient_UpdateTagProperties_wrongTag(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.UpdateTagProperties(ctx, "hello-world", "wrong-tag", &ClientUpdateTagPropertiesOptions{Value: &TagWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.Error(t, err)
}

func TestClient_UploadManifest(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient("https://azacrlivetest.azurecr.io", cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	getRes, err := client.GetManifest(ctx, "hello-world", "latest", &ClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	require.NoError(t, err)
	manifest, err := io.ReadAll(getRes.ManifestData)
	require.NoError(t, err)
	uploadRes, err := client.UploadManifest(ctx, "hello-world", "test", "application/vnd.docker.distribution.manifest.v2+json", streaming.NopCloser(bytes.NewReader(manifest)), nil)
	require.NoError(t, err)
	require.NotEmpty(t, *uploadRes.DockerContentDigest)
	fmt.Printf("uploaded manifest digest: %s\n", *uploadRes.DockerContentDigest)
}
