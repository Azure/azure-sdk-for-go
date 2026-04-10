// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"context"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azcloud "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/stretchr/testify/require"
)

const (
	fakeACRRefreshToken = ".eyJqdGkiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJzdWIiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJuYmYiOjQ2NzA0MTEyMTIsImV4cCI6NDY3MDQyMjkxMiwiaWF0Ijo0NjcwNDExMjEyLCJpc3MiOiJBenVyZSBDb250YWluZXIgUmVnaXN0cnkiLCJhdWQiOiJhemFjcmxpdmV0ZXN0LmF6dXJlY3IuaW8iLCJ2ZXJzaW9uIjoiMS4wIiwicmlkIjoiMDAwMCIsImdyYW50X3R5cGUiOiJyZWZyZXNoX3Rva2VuIiwiYXBwaWQiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJwZXJtaXNzaW9ucyI6eyJBY3Rpb25zIjpbInJlYWQiLCJ3cml0ZSIsImRlbGV0ZSIsImRlbGV0ZWQvcmVhZCIsImRlbGV0ZWQvcmVzdG9yZS9hY3Rpb24iXSwiTm90QWN0aW9ucyI6bnVsbH0sInJvbGVzIjpbXX0=."
	fakeDigest          = "sha256:00"
	fakeLoginServer     = fakeRegistry + ".azurecr.io"
	fakeRegistry        = recording.SanitizedValue
	recordingDirectory  = "sdk/containers/azcontainerregistry/testdata"
)

var (
	ctx = context.Background()

	testConfig = struct {
		cloud                     azcloud.Configuration
		credential                azcore.TokenCredential
		loginServer, registryName string
	}{
		cloud:        azcloud.AzurePublic,
		credential:   &credential.Fake{},
		loginServer:  fakeLoginServer,
		registryName: fakeRegistry,
	}
)

// getEndpointCredAndClientOptions will create a credential and a client options for test application.
// The client options will initialize the transport for recording client add recording policy to the pipeline.
// In the record mode, the credential will be a DefaultAzureCredential which combines several common credentials.
// In the playback mode, the credential will be a fake credential which will bypass truly authorization.
func getEndpointCredAndClientOptions(t *testing.T) (string, azcore.TokenCredential, azcore.ClientOptions) {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	options := azcore.ClientOptions{
		Cloud:     testConfig.cloud,
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
		testConfig.credential, err = credential.New(nil)
		if err != nil {
			panic(err)
		}
		if testConfig.loginServer = os.Getenv("LOGIN_SERVER"); testConfig.loginServer == "" {
			panic("no value for LOGIN_SERVER")
		}
		if testConfig.registryName = os.Getenv("REGISTRY_NAME"); testConfig.registryName == "" {
			panic("no value for REGISTRY_NAME")
		}
		env := os.Getenv("AZCONTAINERREGISTRY_ENVIRONMENT")
		switch {
		case strings.EqualFold(env, "AzureUSGovernment"):
			testConfig.cloud = azcloud.AzureGovernment
		case strings.EqualFold(env, "AzureCloud"):
			testConfig.cloud = azcloud.AzurePublic
		case strings.EqualFold(env, "AzureChinaCloud"):
			testConfig.cloud = azcloud.AzureChina
		case len(env) > 0:
			panic("unexpected value for AZCONTAINERREGISTRY_ENVIRONMENT: " + env)
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

// buildImage invokes the Azure CLI to build a new image in ACR for the given test. It returns the image's repository and digest.
func buildImage(t *testing.T) (string, string) {
	repository := strings.ReplaceAll(strings.ToLower(t.Name()), "/", "_")
	if recording.GetRecordMode() == recording.PlaybackMode {
		return repository, fakeDigest
	}

	// build images in parallel, in separate goroutines, because building can be slow and may require retries in CI
	t.Parallel()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	ch := make(chan struct{})
	var (
		err error
		out []byte
	)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			default:
				cmd := exec.CommandContext(ctx, "az", "acr", "build", "-r", testConfig.registryName, "--image", repository, "--build-arg", "ID="+repository, ".")
				cmd.Dir = "testdata"
				out, err = cmd.CombinedOutput()
				if err == nil || strings.Contains(string(out), "az login") {
					return
				}
			}
		}
	}()
	<-ch
	require.NoError(t, err, string(out))

	// this assumes the image has one layer digest i.e., it's FROM scratch and the Dockerfile touches the filesystem once
	digest := string(regexp.MustCompile("(sha256:[0-9a-f]{64})").Find(out))
	require.NotEmpty(t, digest, "failed to find digest in "+string(out))
	if recording.GetRecordMode() == recording.RecordingMode {
		_, sum, found := strings.Cut(digest, ":")
		require.True(t, found)
		require.NoError(t, recording.AddGeneralRegexSanitizer("00", sum, nil))
	}
	return repository, digest
}
