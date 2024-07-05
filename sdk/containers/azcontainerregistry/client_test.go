//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

const alpineManifestDigest = "sha256:f271e74b17ced29b915d351685fd4644785c6d1559dd1f2d4189a5e851ef753a"

func TestClient_DeleteManifest(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	resp, err := client.GetTagProperties(ctx, "busybox", "1.36.1-uclibc", nil)
	require.NoError(t, err)
	_, err = client.DeleteManifest(ctx, "busybox", *resp.Tag.Digest, nil)
	require.NoError(t, err)
	_, err = client.DeleteManifest(ctx, "hello-world", "sha256:sha256:aa0cc8055b82dc2509bed2e19b275c8f463506616377219d9642221ab53cf9fe", nil)
	require.NoError(t, err)
}

func TestClient_DeleteManifest_wrongDigest(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteManifest(ctx, "alpine", "error-digest", nil)
	require.Error(t, err)
}

func TestClient_DeleteManifest_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.DeleteManifest(ctx, "", "digest", nil)
	require.Error(t, err)
	_, err = client.DeleteManifest(ctx, "name", "", nil)
	require.Error(t, err)
}

func TestClient_DeleteRepository(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteRepository(ctx, "eclipse-mosquitto", nil)
	require.NoError(t, err)
}

func TestClient_DeleteRepository_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.DeleteRepository(ctx, "", nil)
	require.Error(t, err)
}

func TestClient_DeleteRepository_error(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))

	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	require.NoError(t, err)
	client := &Client{
		azcoreClient,
		srv.URL(),
	}
	ctx := context.Background()
	_, err = client.DeleteRepository(ctx, "test", nil)
	require.Error(t, err)
}

func TestClient_DeleteTag(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteTag(ctx, "alpine", "3.14.8", nil)
	require.NoError(t, err)
}

func TestClient_DeleteTag_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.DeleteTag(ctx, "", "tag", nil)
	require.Error(t, err)
	_, err = client.DeleteTag(ctx, "name", "", nil)
	require.Error(t, err)
}

func TestClient_DeleteTag_error(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))

	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	require.NoError(t, err)
	client := &Client{
		azcoreClient,
		srv.URL(),
	}
	ctx := context.Background()
	_, err = client.DeleteTag(ctx, "name", "tag", nil)
	require.Error(t, err)
}

func TestClient_GetManifest(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.GetManifest(ctx, "alpine", "3.17.1", &ClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	require.NoError(t, err)
	reader, err := NewDigestValidationReader(*res.DockerContentDigest, res.ManifestData)
	require.NoError(t, err)
	manifest, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.NotEmpty(t, manifest)
	fmt.Printf("manifest content: %s\n", manifest)
}

func TestClient_GetManifest_wrongServerDigest(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("test")), mock.WithHeader("Docker-Content-Digest", "sha256:wrong"))

	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	require.NoError(t, err)
	client := &Client{
		azcoreClient,
		srv.URL(),
	}
	ctx := context.Background()
	resp, err := client.GetManifest(ctx, "name", "sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", nil)
	require.NoError(t, err)
	reader, err := NewDigestValidationReader(*resp.DockerContentDigest, resp.ManifestData)
	require.NoError(t, err)
	_, err = io.ReadAll(reader)
	require.Error(t, err, ErrMismatchedHash)
}

func TestClient_GetManifest_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.GetManifest(ctx, "", "tag", nil)
	require.Error(t, err)
	_, err = client.GetManifest(ctx, "name", "", nil)
	require.Error(t, err)
}

func TestClient_GetManifest_wrongTag(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetManifest(ctx, "alpine", "wrong-tag", &ClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	require.Error(t, err)
}

func TestClient_GetManifestProperties(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digestRes, err := client.GetManifestProperties(ctx, "alpine", alpineManifestDigest, nil)
	require.NoError(t, err)
	require.Equal(t, *digestRes.Manifest.Digest, alpineManifestDigest)
	resp, err := client.GetTagProperties(ctx, "alpine", "3.17.1", nil)
	require.NoError(t, err)
	tagRes, err := client.GetManifestProperties(ctx, "alpine", *resp.Tag.Digest, nil)
	require.NoError(t, err)
	require.Equal(t, alpineManifestDigest, *tagRes.Manifest.Digest)
}

func TestClient_GetManifestProperties_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.GetManifestProperties(ctx, "", "digest", nil)
	require.Error(t, err)
	_, err = client.GetManifestProperties(ctx, "name", "", nil)
	require.Error(t, err)
}

func TestClient_GetManifestProperties_wrongDigest(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetManifestProperties(ctx, "alpine", "wrong-digest", nil)
	require.Error(t, err)
}

func TestClient_GetRepositoryProperties(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.GetRepositoryProperties(ctx, "alpine", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.Name)
	fmt.Printf("repository name: %s\n", *res.Name)
	require.NotEmpty(t, *res.RegistryLoginServer)
	fmt.Printf("registry login server of the repository: %s\n", *res.RegistryLoginServer)
	require.NotEmpty(t, *res.ManifestCount)
	fmt.Printf("repository manifest count: %d\n", *res.ManifestCount)
}

func TestClient_GetRepositoryProperties_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.GetRepositoryProperties(ctx, "", nil)
	require.Error(t, err)
}

func TestClient_GetRepositoryProperties_wrongName(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetRepositoryProperties(ctx, "wrong-name", nil)
	require.Error(t, err)
}

func TestClient_GetTagProperties(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.GetTagProperties(ctx, "alpine", "3.17.1", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.Tag.Name)
	fmt.Printf("tag name: %s\n", *res.Tag.Name)
	require.NotEmpty(t, *res.Tag.Digest)
	fmt.Printf("tag digest: %s\n", *res.Tag.Digest)
}

func TestClient_GetTagProperties_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.GetTagProperties(ctx, "", "", nil)
	require.Error(t, err)
	_, err = client.GetTagProperties(ctx, "name", "", nil)
	require.Error(t, err)
}

func TestClient_GetTagProperties_wrongTag(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetTagProperties(ctx, "alpine", "wrong-tag", nil)
	require.Error(t, err)
}

func TestClient_NewListManifestsPager(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	pager := client.NewListManifestsPager("alpine", &ClientListManifestsOptions{
		MaxNum: to.Ptr[int32](1),
	})
	items := 0
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.Manifests.Attributes)
		items += len(page.Manifests.Attributes)
	}
	require.NotZero(t, items)

	pager = client.NewListManifestsPager("alpine", &ClientListManifestsOptions{
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
	pager = client.NewListManifestsPager("alpine", &ClientListManifestsOptions{
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

func TestClient_NewListManifestsPager_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	pager := client.NewListManifestsPager("", nil)
	for pager.More() {
		_, err := pager.NextPage(ctx)
		require.Error(t, err)
		break
	}
}

func TestClient_NewListManifestsPager_wrongRepositoryName(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
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
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
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
	require.Equal(t, 3, pages)
	require.Equal(t, 3, items)
}

func TestClient_NewListRepositoriesPager_error(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))

	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	require.NoError(t, err)
	client := &Client{
		azcoreClient,
		srv.URL(),
	}
	ctx := context.Background()
	pager := client.NewListRepositoriesPager(nil)
	for pager.More() {
		_, err := pager.NextPage(ctx)
		require.Error(t, err)
		break
	}
}

func TestClient_NewListTagsPager(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	pager := client.NewListTagsPager("alpine", &ClientListTagsOptions{
		MaxNum: to.Ptr[int32](1),
	})
	items := 0
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.Tags)
		require.Equal(t, 1, len(page.Tags))
		items += len(page.Tags)
	}
	require.NotZero(t, items)

	pager = client.NewListTagsPager("alpine", &ClientListTagsOptions{
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
	pager = client.NewListTagsPager("alpine", &ClientListTagsOptions{
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

func TestClient_NewListTagsPager_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	pager := client.NewListTagsPager("", nil)
	for pager.More() {
		_, err := pager.NextPage(ctx)
		require.Error(t, err)
		break
	}
}

func TestClient_NewListTagsPager_wrongRepositoryName(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
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
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	resp, err := client.GetTagProperties(ctx, "alpine", "3.17.1", nil)
	require.NoError(t, err)
	res, err := client.UpdateManifestProperties(ctx, "alpine", *resp.Tag.Digest, &ClientUpdateManifestPropertiesOptions{Value: &ManifestWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.NoError(t, err)
	require.False(t, *res.Manifest.ChangeableAttributes.CanWrite)
	res, err = client.UpdateManifestProperties(ctx, "alpine", alpineManifestDigest, &ClientUpdateManifestPropertiesOptions{Value: &ManifestWriteableProperties{
		CanWrite: to.Ptr(true),
	},
	})
	require.NoError(t, err)
	require.True(t, *res.Manifest.ChangeableAttributes.CanWrite)
}

func TestClient_UpdateManifestProperties_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.UpdateManifestProperties(ctx, "", "digest", nil)
	require.Error(t, err)
	_, err = client.UpdateManifestProperties(ctx, "name", "", nil)
	require.Error(t, err)
}

func TestClient_UpdateManifestProperties_wrongDigest(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetTagProperties(ctx, "alpine", "wrong-digest", nil)
	require.Error(t, err)
}

func TestClient_UpdateRepositoryProperties(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.UpdateRepositoryProperties(ctx, "busybox", &ClientUpdateRepositoryPropertiesOptions{Value: &RepositoryWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.NoError(t, err)
	require.False(t, *res.ContainerRepositoryProperties.ChangeableAttributes.CanWrite)
	res, err = client.UpdateRepositoryProperties(ctx, "busybox", &ClientUpdateRepositoryPropertiesOptions{Value: &RepositoryWriteableProperties{
		CanWrite: to.Ptr(true),
	},
	})
	require.NoError(t, err)
	require.True(t, *res.ContainerRepositoryProperties.ChangeableAttributes.CanWrite)
}

func TestClient_UpdateRepositoryProperties_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.UpdateRepositoryProperties(ctx, "", nil)
	require.Error(t, err)
}

func TestClient_UpdateRepositoryProperties_wrongRepository(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.UpdateRepositoryProperties(ctx, "wrong-repository", &ClientUpdateRepositoryPropertiesOptions{Value: &RepositoryWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.Error(t, err)
}

func TestClient_UpdateTagProperties(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.UpdateTagProperties(ctx, "alpine", "3.17.1", &ClientUpdateTagPropertiesOptions{Value: &TagWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.NoError(t, err)
	require.False(t, *res.Tag.ChangeableAttributes.CanWrite)
	res, err = client.UpdateTagProperties(ctx, "alpine", "3.17.1", &ClientUpdateTagPropertiesOptions{Value: &TagWriteableProperties{
		CanWrite: to.Ptr(true),
	},
	})
	require.NoError(t, err)
	require.True(t, *res.Tag.ChangeableAttributes.CanWrite)
}

func TestClient_UpdateTagProperties_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.UpdateTagProperties(ctx, "name", "", nil)
	require.Error(t, err)
	_, err = client.UpdateTagProperties(ctx, "", "tag", nil)
	require.Error(t, err)
}

func TestClient_UpdateTagProperties_wrongTag(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.UpdateTagProperties(ctx, "alpine", "wrong-tag", &ClientUpdateTagPropertiesOptions{Value: &TagWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	require.Error(t, err)
}

func TestClient_UploadManifest(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	getRes, err := client.GetManifest(ctx, "hello-world", "latest", &ClientGetManifestOptions{Accept: to.Ptr("application/vnd.oci.image.index.v1+json")})
	require.NoError(t, err)
	manifest, err := io.ReadAll(getRes.ManifestData)
	require.NoError(t, err)
	reader := bytes.NewReader(manifest)
	uploadRes, err := client.UploadManifest(ctx, "hello-world", "test", "application/vnd.oci.image.index.v1+json", streaming.NopCloser(reader), nil)
	require.NoError(t, err)
	require.NotEmpty(t, *uploadRes.DockerContentDigest)
	fmt.Printf("uploaded manifest digest: %s\n", *uploadRes.DockerContentDigest)
	_, err = reader.Seek(0, io.SeekStart)
	require.NoError(t, err)
	validateReader, err := NewDigestValidationReader(*uploadRes.DockerContentDigest, reader)
	require.NoError(t, err)
	_, err = io.ReadAll(validateReader)
	require.NoError(t, err)
}

func TestClient_UploadManifest_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.UploadManifest(ctx, "", "reference", "contentType", nil, nil)
	require.Error(t, err)
	_, err = client.UploadManifest(ctx, "name", "", "contentType", nil, nil)
	require.Error(t, err)
}

func TestClient_UploadManifest_error(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))

	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	require.NoError(t, err)
	client := &Client{
		azcoreClient,
		srv.URL(),
	}
	ctx := context.Background()
	_, err = client.UploadManifest(ctx, "name", "reference", "contentType", nil, nil)
	require.Error(t, err)
}

func TestClient_wrongEndpoint(t *testing.T) {
	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{}, nil)
	require.NoError(t, err)
	client := &Client{
		azcoreClient,
		"wrong-endpoint",
	}
	ctx := context.Background()
	_, err = client.DeleteManifest(ctx, "name", "digest", nil)
	require.Error(t, err)
	_, err = client.DeleteRepository(ctx, "name", nil)
	require.Error(t, err)
	_, err = client.DeleteTag(ctx, "name", "tag", nil)
	require.Error(t, err)
	_, err = client.GetManifest(ctx, "name", "reference", nil)
	require.Error(t, err)
	_, err = client.GetManifestProperties(ctx, "name", "digest", nil)
	require.Error(t, err)
	_, err = client.GetRepositoryProperties(ctx, "name", nil)
	require.Error(t, err)
	_, err = client.GetTagProperties(ctx, "name", "tag", nil)
	require.Error(t, err)
	_, err = client.NewListManifestsPager("name", nil).NextPage(ctx)
	require.Error(t, err)
	_, err = client.NewListRepositoriesPager(nil).NextPage(ctx)
	require.Error(t, err)
	_, err = client.NewListTagsPager("name", nil).NextPage(ctx)
	require.Error(t, err)
	_, err = client.UpdateManifestProperties(ctx, "name", "digest", nil)
	require.Error(t, err)
	_, err = client.UpdateRepositoryProperties(ctx, "name", nil)
	require.Error(t, err)
	_, err = client.UpdateTagProperties(ctx, "name", "tag", nil)
	require.Error(t, err)
	_, err = client.UploadManifest(ctx, "name", "reference", "contentType", nil, nil)
	require.Error(t, err)
}
