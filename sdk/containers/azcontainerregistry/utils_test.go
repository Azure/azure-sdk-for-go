//go:build go1.18
// +build go1.18

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
	fakeLoginServer     = recording.SanitizedValue + ".azurecr.io"
	fakeRepository      = fakeLoginServer + "/fake"
	recordingDirectory  = "sdk/containers/azcontainerregistry/testdata"
)

var testConfig = struct {
	cloud       azcloud.Configuration
	credential  azcore.TokenCredential
	loginServer string
}{
	cloud:       azcloud.AzurePublic,
	credential:  &credential.Fake{},
	loginServer: fakeLoginServer,
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
		testConfig.credential, err = credential.New(nil)
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

func pushImage(t *testing.T) (string, string) {
	repository := strings.ReplaceAll(strings.ToLower(t.Name()), "/", "_")
	if recording.GetRecordMode() == recording.PlaybackMode {
		return repository, fakeDigest
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	image := testConfig.loginServer + "/" + repository

	cmd := exec.CommandContext(ctx, "docker", "build", "-t", image, "--build-arg", "ID="+repository, ".")
	cmd.Dir = "testdata"
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	cmd = exec.CommandContext(ctx, "docker", "push", image)
	out, err = cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	digest := string(regexp.MustCompile("(sha256:[0-9a-f]{64})").Find(out))
	require.NotEmpty(t, digest, "failed to find digest in "+string(out))

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()
		require.NoError(t, exec.CommandContext(ctx, "docker", "rmi", image).Run())
	})

	_, hash, found := strings.Cut(digest, ":")
	require.True(t, found)
	if recording.GetRecordMode() == recording.RecordingMode {
		// require.NoError(t, recording.AddGeneralRegexSanitizer("fake", repository, nil))
		require.NoError(t, recording.AddGeneralRegexSanitizer("00", hash, nil))
	}
	return repository, digest
}
