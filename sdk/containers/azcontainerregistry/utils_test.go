//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azcloud "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"github.com/stretchr/testify/require"
)

const (
	fakeACRRefreshToken = ".eyJqdGkiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJzdWIiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJuYmYiOjQ2NzA0MTEyMTIsImV4cCI6NDY3MDQyMjkxMiwiaWF0Ijo0NjcwNDExMjEyLCJpc3MiOiJBenVyZSBDb250YWluZXIgUmVnaXN0cnkiLCJhdWQiOiJhemFjcmxpdmV0ZXN0LmF6dXJlY3IuaW8iLCJ2ZXJzaW9uIjoiMS4wIiwicmlkIjoiMDAwMCIsImdyYW50X3R5cGUiOiJyZWZyZXNoX3Rva2VuIiwiYXBwaWQiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJwZXJtaXNzaW9ucyI6eyJBY3Rpb25zIjpbInJlYWQiLCJ3cml0ZSIsImRlbGV0ZSIsImRlbGV0ZWQvcmVhZCIsImRlbGV0ZWQvcmVzdG9yZS9hY3Rpb24iXSwiTm90QWN0aW9ucyI6bnVsbH0sInJvbGVzIjpbXX0=."
	fakeLoginServer     = recording.SanitizedValue + ".azurecr.io"
	recordingDirectory  = "sdk/containers/azcontainerregistry/testdata"
)

var testConfig = struct {
	cloud       azcloud.Configuration
	credential  azcore.TokenCredential
	loginServer string
}{
	cloud:       azcloud.AzurePublic,
	credential:  &FakeCredential{},
	loginServer: fakeLoginServer,
}

// FakeCredential is an empty credential for testing.
type FakeCredential struct {
}

// GetToken provide a fake access token.
func (c *FakeCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: recording.SanitizedValue, ExpiresOn: time.Now().Add(time.Hour * 24).UTC()}, nil
}

// getEndpointCredAndClientOptions will create a credential and a client options for test application.
// The client options will initialize the transport for recording client add recording policy to the pipeline.
// In the record mode, the credential will be a DefaultAzureCredential which combines several common credentials.
// In the playback mode, the credential will be a fake credential which will bypass truly authorization.
func getEndpointCredAndClientOptions(t *testing.T) (string, azcore.TokenCredential, azcore.ClientOptions) {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	options := azcore.ClientOptions{
		Transport: transport,
	}
	return "https://" + testConfig.loginServer, testConfig.credential, options
}

// startRecording starts the recording.
func startRecording(t *testing.T) {
	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
}

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	if recording.GetRecordMode() != recording.PlaybackMode {
		var err error
		testConfig.credential, err = azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			panic(err)
		}
		if testConfig.loginServer = os.Getenv("LOGIN_SERVER"); testConfig.loginServer == "" {
			panic("no value for LOGIN_SERVER")
		}
		env := os.Getenv("AZCONTAINERREGISTRY_ENVIRONMENT")
		switch {
		case strings.EqualFold(env, "AzureUSGovernment"):
			testConfig.cloud = azcloud.AzureGovernment
		case strings.EqualFold(env, "AzureChinaCloud"):
			testConfig.cloud = azcloud.AzureChina
		case len(env) > 0:
			panic("unexpected value for AZCONTAINERREGISTRY_ENVIRONMENT: " + env)
		}
		err = importTestImages()
		if err != nil {
			panic(err)
		}
	}
	if recording.GetRecordMode() != recording.LiveMode {
		proxy, err := recording.StartTestProxy(recordingDirectory, nil)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(err)
			}
		}()
		err = recording.RemoveRegisteredSanitizers([]string{
			"AZSDK2003", // Location header
			"AZSDK3401", // $..refresh_token (client needs a JWT; the sanitizer added below substitutes a static fake)
		}, nil)
		if err != nil {
			panic(err)
		}
		err = recording.AddBodyKeySanitizer("$..refresh_token", fakeACRRefreshToken, "", nil)
		if err != nil {
			panic(err)
		}
		err = recording.AddGeneralRegexSanitizer(fakeLoginServer, testConfig.loginServer, nil)
		if err != nil {
			panic(err)
		}
	}
	return m.Run()
}

func importTestImages() error {
	subID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if subID == "" {
		return errors.New("no value for AZURE_SUBSCRIPTION_ID")
	}
	rg := os.Getenv("AZCONTAINERREGISTRY_RESOURCE_GROUP")
	if rg == "" {
		return errors.New("no value for AZCONTAINERREGISTRY_RESOURCE_GROUP")
	}
	registryName := os.Getenv("REGISTRY_NAME")
	if registryName == "" {
		return errors.New("no value for REGISTRY_NAME")
	}

	ctx := context.Background()
	client, err := armcontainerregistry.NewRegistriesClient(subID, testConfig.credential, &arm.ClientOptions{ClientOptions: azcore.ClientOptions{Cloud: testConfig.cloud}})
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	images := []string{"hello-world:latest", "alpine:3.17.1", "alpine:3.16.3", "alpine:3.15.6", "alpine:3.14.8", "ubuntu:20.04", "nginx:latest"}
	ec := make(chan error, len(images))
	for _, img := range images {
		wg.Add(1)
		go func(image string) {
			poller, err := client.BeginImportImage(ctx, rg, registryName, armcontainerregistry.ImportImageParameters{
				Source: &armcontainerregistry.ImportSource{
					SourceImage: to.Ptr("library/" + image),
					RegistryURI: to.Ptr("docker.io"),
				},
				TargetTags: []*string{to.Ptr(image)},
				Mode:       to.Ptr(armcontainerregistry.ImportModeForce),
			}, nil)
			if err == nil {
				_, err = poller.PollUntilDone(ctx, nil)
			}
			if err != nil {
				ec <- err
			}
			wg.Done()
		}(img)
	}
	wg.Wait()
	select {
	case err = <-ec:
	default:
	}
	return err
}
