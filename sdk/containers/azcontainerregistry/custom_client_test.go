//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("test", nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	wrongCloudConfig := cloud.Configuration{
		ActiveDirectoryAuthorityHost: "test", Services: map[cloud.ServiceName]cloud.ServiceConfiguration{},
	}
	_, err = NewClient("test", nil, &ClientOptions{ClientOptions: azcore.ClientOptions{Cloud: wrongCloudConfig}})
	require.Errorf(t, err, "provided Cloud field is missing Azure Container Registry configuration")
}

func TestClient_GetManifest(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.GetManifest(ctx, "alpine", "3.17.1", &ClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	require.NoError(t, err)
	manifest, err := io.ReadAll(res.ManifestData)
	require.NoError(t, err)
	require.NotEmpty(t, manifest)
	fmt.Printf("manifest content: %s\n", manifest)
}

func TestClient_GetManifest_wrongReference(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("test")))

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	client := &Client{
		srv.URL(),
		pl,
	}
	ctx := context.Background()
	_, err := client.GetManifest(ctx, "name", "sha256:wrong", nil)
	require.Error(t, err)
}

func TestClient_GetManifest_wrongServerDigest(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("test")), mock.WithHeader("Docker-Content-Digest", "sha256:wrong"))

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	client := &Client{
		srv.URL(),
		pl,
	}
	ctx := context.Background()
	_, err := client.GetManifest(ctx, "name", "sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", nil)
	require.Error(t, err)
}

func TestClient_UploadManifest(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewClient(endpoint, cred, &ClientOptions{ClientOptions: options})
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

func TestClient_UploadManifest_wrongClientDigest(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	client := &Client{
		srv.URL(),
		pl,
	}
	ctx := context.Background()
	_, err := client.UploadManifest(ctx, "name", "sha256:wrong", "contentType", streaming.NopCloser(bytes.NewReader([]byte("test"))), nil)
	require.Error(t, err)
}

func TestClient_UploadManifest_wrongServerDigest(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusCreated), mock.WithHeader("Docker-Content-Digest", "sha256:wrong"))

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	client := &Client{
		srv.URL(),
		pl,
	}
	ctx := context.Background()
	_, err := client.UploadManifest(ctx, "name", "sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", "contentType", streaming.NopCloser(bytes.NewReader([]byte("test"))), nil)
	require.Error(t, err)
}
