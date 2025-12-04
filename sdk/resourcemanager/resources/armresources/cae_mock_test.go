// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/require"
)

const caeChallenge = `Bearer realm="", error_description="Continuous access evaluation resulted in challenge", error="insufficient_claims", claims="eyJhY2Nlc3NfdG9rZW4iOnsibmJmIjp7ImVzc2VudGlhbCI6dHJ1ZSwgInZhbHVlIjoiMTcyNjI1ODEyMiJ9fX0=" ` // this is not a real token, does not contain any sensitive info, just for test.
const invalidCaeChallenge = `Bearer realm="", error_description="", error="insufficient_claims", claims=""`

func TestCaeSupportForManagementLibrary(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusUnauthorized),
		mock.WithHeader("WWW-Authenticate", caeChallenge),
	)
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusNoContent),
	)
	c := cloud.Configuration{
		ActiveDirectoryAuthorityHost: srv.URL(),
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {
				Audience: "testAudience",
				Endpoint: srv.URL(),
			},
		},
	}
	client, err := armresources.NewClient("subscriptionID", credential.Fake{}, &arm.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: srv, Cloud: c}})
	require.NoError(t, err)
	resp, err := client.CheckExistence(context.Background(), "resourceGroupName", "resourceProviderNamespace", "parentResourcePath", "resourceType", "resourceName", "apiVersion", nil)
	require.NoError(t, err)
	require.Equal(t, true, resp.Success)
}

func TestCaeSupportForManagementLibrary_fail(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusUnauthorized),
		mock.WithHeader("WWW-Authenticate", invalidCaeChallenge),
	)
	c := cloud.Configuration{
		ActiveDirectoryAuthorityHost: srv.URL(),
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {
				Audience: "testAudience",
				Endpoint: srv.URL(),
			},
		},
	}
	client, err := armresources.NewClient("subscriptionID", credential.Fake{}, &arm.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: srv, Cloud: c}})
	require.NoError(t, err)
	_, err = client.CheckExistence(context.Background(), "resourceGroupName", "resourceProviderNamespace", "parentResourcePath", "resourceType", "resourceName", "apiVersion", nil)
	require.Error(t, err)
}
