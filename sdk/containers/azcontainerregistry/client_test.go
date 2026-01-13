// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestClient_DeleteManifest(t *testing.T) {
	repository, _ := buildImage(t)
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteManifest(ctx, "", "digest", nil)
	require.Error(t, err)
	_, err = client.DeleteManifest(ctx, "name", "", nil)
	require.Error(t, err)
	resp, err := client.GetTagProperties(ctx, repository, "latest", nil)
	require.NoError(t, err)
	_, err = client.DeleteManifest(ctx, repository, *resp.Tag.Digest, nil)
	require.NoError(t, err)
}

func TestClient_DeleteRepository(t *testing.T) {
	repository, _ := buildImage(t)
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteRepository(ctx, "", nil)
	require.Error(t, err)
	_, err = client.DeleteRepository(ctx, repository, nil)
	require.NoError(t, err)
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
	_, err = client.DeleteRepository(ctx, "test", nil)
	require.Error(t, err)
}

func TestClient_DeleteTag(t *testing.T) {
	repository, _ := buildImage(t)
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteTag(ctx, "", "tag", nil)
	require.Error(t, err)
	_, err = client.DeleteTag(ctx, "name", "", nil)
	require.Error(t, err)
	_, err = client.DeleteTag(ctx, repository, "latest", nil)
	require.NoError(t, err)
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
	_, err = client.DeleteTag(ctx, "name", "tag", nil)
	require.Error(t, err)
}

func TestClient_GetManifest(t *testing.T) {
	repository, _ := buildImage(t)
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetManifest(ctx, repository, "wrong-tag", &ClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	require.Error(t, err)
	res, err := client.GetManifest(ctx, repository, "latest", &ClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	require.NoError(t, err)
	reader, err := NewDigestValidationReader(*res.DockerContentDigest, res.ManifestData)
	require.NoError(t, err)
	if recording.GetRecordMode() == recording.PlaybackMode {
		reader.digestValidator = &sha256Validator{&fakeHash{}}
	}
	manifest, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.NotEmpty(t, manifest)
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
	resp, err := client.GetManifest(ctx, "name", "sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", nil)
	require.NoError(t, err)
	reader, err := NewDigestValidationReader(*resp.DockerContentDigest, resp.ManifestData)
	require.NoError(t, err)
	_, err = io.ReadAll(reader)
	require.Error(t, err, ErrMismatchedHash)
}

func TestClient_GetManifest_empty(t *testing.T) {
	client, err := NewClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.GetManifest(ctx, "", "tag", nil)
	require.Error(t, err)
	_, err = client.GetManifest(ctx, "name", "", nil)
	require.Error(t, err)
}

func TestClient(t *testing.T) {
	repository, digest := buildImage(t)

	t.Run("GetManifestProperties", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
		require.NoError(t, err)
		_, err = client.GetManifestProperties(ctx, "", "digest", nil)
		require.Error(t, err)
		_, err = client.GetManifestProperties(ctx, "name", "", nil)
		require.Error(t, err)
		_, err = client.GetManifestProperties(ctx, repository, "wrong-digest", nil)
		require.Error(t, err)
		digestRes, err := client.GetManifestProperties(ctx, repository, digest, nil)
		require.NoError(t, err)
		require.Equal(t, *digestRes.Manifest.Digest, digest)
		resp, err := client.GetTagProperties(ctx, repository, "latest", nil)
		require.NoError(t, err)
		tagRes, err := client.GetManifestProperties(ctx, repository, *resp.Tag.Digest, nil)
		require.NoError(t, err)
		require.Equal(t, digest, *tagRes.Manifest.Digest)
	})

	t.Run("GetRepositoryProperties", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
		require.NoError(t, err)
		_, err = client.GetRepositoryProperties(ctx, "", nil)
		require.Error(t, err)
		_, err = client.GetRepositoryProperties(ctx, "wrong-name", nil)
		require.Error(t, err)
		res, err := client.GetRepositoryProperties(ctx, repository, nil)
		require.NoError(t, err)
		require.NotEmpty(t, *res.Name)
		require.NotEmpty(t, *res.RegistryLoginServer)
		require.NotEmpty(t, *res.ManifestCount)
	})

	t.Run("GetTagProperties", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
		require.NoError(t, err)
		_, err = client.GetTagProperties(ctx, "", "", nil)
		require.Error(t, err)
		_, err = client.GetTagProperties(ctx, "name", "", nil)
		require.Error(t, err)
		_, err = client.GetTagProperties(ctx, repository, "wrong-tag", nil)
		require.Error(t, err)
		res, err := client.GetTagProperties(ctx, repository, "latest", nil)
		require.NoError(t, err)
		require.NotEmpty(t, *res.Tag.Name)
		require.NotEmpty(t, *res.Tag.Digest)
	})

	t.Run("NewListManifestsPager", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
		require.NoError(t, err)
		pager := client.NewListManifestsPager(repository, &ClientListManifestsOptions{
			MaxNum: to.Ptr[int32](1),
		})
		items := 0
		for pager.More() {
			page, err := pager.NextPage(ctx)
			require.NoError(t, err)
			require.NotEmpty(t, page.Attributes)
			items += len(page.Attributes)
		}
		require.NotZero(t, items)

		pager = client.NewListManifestsPager(repository, &ClientListManifestsOptions{
			OrderBy: to.Ptr(ArtifactManifestOrderByLastUpdatedOnDescending),
		})
		var descendingItems []*ManifestAttributes
		for pager.More() {
			page, err := pager.NextPage(ctx)
			require.NoError(t, err)
			require.NotEmpty(t, page.Attributes)
			descendingItems = append(descendingItems, page.Attributes...)
		}
		pager = client.NewListManifestsPager(repository, &ClientListManifestsOptions{
			OrderBy: to.Ptr(ArtifactManifestOrderByLastUpdatedOnAscending),
		})
		var ascendingItems []*ManifestAttributes
		for pager.More() {
			page, err := pager.NextPage(ctx)
			require.NoError(t, err)
			require.NotEmpty(t, page.Attributes)
			ascendingItems = append(ascendingItems, page.Attributes...)
		}
		for i := range descendingItems {
			require.Equal(t, descendingItems[i].Digest, ascendingItems[len(ascendingItems)-1-i].Digest)
		}
	})

	t.Run("NewListTagsPager", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
		require.NoError(t, err)
		pager := client.NewListTagsPager(repository, &ClientListTagsOptions{
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

		pager = client.NewListTagsPager(repository, &ClientListTagsOptions{
			OrderBy: to.Ptr(ArtifactTagOrderByLastUpdatedOnDescending),
		})
		var descendingItems []*TagAttributes
		for pager.More() {
			page, err := pager.NextPage(ctx)
			require.NoError(t, err)
			require.NotEmpty(t, page.Tags)
			descendingItems = append(descendingItems, page.Tags...)
		}
		pager = client.NewListTagsPager(repository, &ClientListTagsOptions{
			OrderBy: to.Ptr(ArtifactTagOrderByLastUpdatedOnAscending),
		})
		var ascendingItems []*TagAttributes
		for pager.More() {
			page, err := pager.NextPage(ctx)
			require.NoError(t, err)
			require.NotEmpty(t, page.Tags)
			ascendingItems = append(ascendingItems, page.Tags...)
		}
		for i := range descendingItems {
			require.Equal(t, descendingItems[i].Name, ascendingItems[len(ascendingItems)-1-i].Name)
		}
	})
}

func TestClient_NewListManifestsPager_empty(t *testing.T) {
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
	// ensure the registry contains at least one repository
	_, _ = buildImage(t)
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
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
		pages++
		items += len(page.Names)
	}
	require.NotZero(t, pages)
	require.NotZero(t, items)
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
	pager := client.NewListRepositoriesPager(nil)
	for pager.More() {
		_, err := pager.NextPage(ctx)
		require.Error(t, err)
		break
	}
}

func TestClient_NewListTagsPager_empty(t *testing.T) {
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
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	pager := client.NewListTagsPager("wrong-name", nil)
	for pager.More() {
		_, err := pager.NextPage(ctx)
		require.Error(t, err)
		break
	}
}

func TestClient_Update(t *testing.T) {
	repository, _ := buildImage(t)

	t.Run("UpdateManifestProperties", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
		require.NoError(t, err)
		_, err = client.UpdateManifestProperties(ctx, "", "digest", nil)
		require.Error(t, err)
		_, err = client.UpdateManifestProperties(ctx, "name", "", nil)
		require.Error(t, err)
		_, err = client.GetTagProperties(ctx, repository, "wrong-tag", nil)
		require.Error(t, err)
		resp, err := client.GetTagProperties(ctx, repository, "latest", nil)
		require.NoError(t, err)
		res, err := client.UpdateManifestProperties(ctx, repository, *resp.Tag.Digest, &ClientUpdateManifestPropertiesOptions{
			Value: &ManifestWriteableProperties{
				CanWrite: to.Ptr(false),
			},
		})
		require.NoError(t, err)
		require.False(t, *res.Manifest.ChangeableAttributes.CanWrite)
		res, err = client.UpdateManifestProperties(ctx, repository, *resp.Tag.Digest, &ClientUpdateManifestPropertiesOptions{
			Value: &ManifestWriteableProperties{
				CanWrite: to.Ptr(true),
			},
		})
		require.NoError(t, err)
		require.True(t, *res.Manifest.ChangeableAttributes.CanWrite)
	})

	t.Run("UpdateRepositoryProperties", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
		require.NoError(t, err)
		_, err = client.UpdateRepositoryProperties(ctx, "", nil)
		require.Error(t, err)
		_, err = client.UpdateRepositoryProperties(ctx, "wrong-repository", &ClientUpdateRepositoryPropertiesOptions{
			Value: &RepositoryWriteableProperties{
				CanWrite: to.Ptr(false),
			},
		})
		require.Error(t, err)
		res, err := client.UpdateRepositoryProperties(ctx, repository, &ClientUpdateRepositoryPropertiesOptions{
			Value: &RepositoryWriteableProperties{
				CanWrite: to.Ptr(false),
			},
		})
		require.NoError(t, err)
		require.False(t, *res.ChangeableAttributes.CanWrite)
		res, err = client.UpdateRepositoryProperties(ctx, repository, &ClientUpdateRepositoryPropertiesOptions{
			Value: &RepositoryWriteableProperties{
				CanWrite: to.Ptr(true),
			},
		})
		require.NoError(t, err)
		require.True(t, *res.ChangeableAttributes.CanWrite)
	})

	t.Run("UpdateTagProperties", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
		require.NoError(t, err)
		_, err = client.UpdateTagProperties(ctx, "name", "", nil)
		require.Error(t, err)
		_, err = client.UpdateTagProperties(ctx, "", "tag", nil)
		require.Error(t, err)
		_, err = client.UpdateTagProperties(ctx, repository, "wrong-tag", &ClientUpdateTagPropertiesOptions{
			Value: &TagWriteableProperties{
				CanWrite: to.Ptr(false),
			},
		})
		require.Error(t, err)
		res, err := client.UpdateTagProperties(ctx, repository, "latest", &ClientUpdateTagPropertiesOptions{
			Value: &TagWriteableProperties{
				CanWrite: to.Ptr(false),
			},
		})
		require.NoError(t, err)
		require.False(t, *res.Tag.ChangeableAttributes.CanWrite)
		res, err = client.UpdateTagProperties(ctx, repository, "latest", &ClientUpdateTagPropertiesOptions{
			Value: &TagWriteableProperties{
				CanWrite: to.Ptr(true),
			},
		})
		require.NoError(t, err)
		require.True(t, *res.Tag.ChangeableAttributes.CanWrite)
	})
}

func TestClient_UploadManifest(t *testing.T) {
	repository, _ := buildImage(t)
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	getRes, err := client.GetManifest(ctx, repository, "latest", &ClientGetManifestOptions{Accept: to.Ptr("application/vnd.oci.image.index.v1+json")})
	require.NoError(t, err)
	manifest, err := io.ReadAll(getRes.ManifestData)
	require.NoError(t, err)
	reader := bytes.NewReader(manifest)
	uploadRes, err := client.UploadManifest(ctx, repository, "test", "application/vnd.oci.image.index.v1+json", streaming.NopCloser(reader), nil)
	require.NoError(t, err)
	require.NotEmpty(t, *uploadRes.DockerContentDigest)
	_, err = reader.Seek(0, io.SeekStart)
	require.NoError(t, err)
	validateReader, err := NewDigestValidationReader(*uploadRes.DockerContentDigest, reader)
	require.NoError(t, err)
	_, err = io.ReadAll(validateReader)
	require.NoError(t, err)
}

func TestClient_UploadManifest_empty(t *testing.T) {
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
