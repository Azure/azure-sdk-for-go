// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armplatformvalidation_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/platformvalidation/armplatformvalidation"
	"github.com/stretchr/testify/require"
)

const (
	pathToPackage = "sdk/resourcemanager/platformvalidation/armplatformvalidation/testdata"
	bootTestName  = "compute-vmimagevalidation-vm-boot-test"
	bootTestID    = "/providers/Microsoft.PlatformValidation/validationTests/" + bootTestName
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	stop := testutil.StartProxy(pathToPackage)
	defer stop()
	if err := configureSanitizers(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return m.Run()
}

func configureSanitizers() error {
	if recording.GetRecordMode() == recording.LiveMode {
		return nil
	}
	// Resource identities and LRO routes are reused by subsequent SDK requests.
	if err := recording.RemoveRegisteredSanitizers([]string{"AZSDK3430", "AZSDK3493", "AZSDK4001"}, nil); err != nil {
		return err
	}
	endpoint := recording.GetEnvVariable("PLATFORMVALIDATION_ARM_ENDPOINT", "https://management.azure.com")
	if err := recording.AddGeneralRegexSanitizer("https://management.azure.com", regexp.QuoteMeta(endpoint), nil); err != nil {
		return err
	}
	if err := recording.AddGeneralRegexSanitizer("sanitized@example.invalid", `(?i)[a-z0-9._%+-]+@microsoft\.com`, nil); err != nil {
		return err
	}
	for _, key := range []string{"planConfigurationJson", "planConfigurationSnapshot", "inputsJson"} {
		if err := recording.AddBodyKeySanitizer("$.."+key, "https://example.invalid/image.vhd", `https://[^"\s]+?\.vhd(?:\?[^"\s]*)?`, nil); err != nil {
			return err
		}
	}
	for _, key := range []string{"createdBy", "lastModifiedBy"} {
		if err := recording.AddBodyKeySanitizer("$..systemData."+key, "sanitized@example.invalid", "", nil); err != nil {
			return err
		}
	}
	for _, key := range []string{"tags.*", "properties.description", "name"} {
		if err := recording.AddBodyKeySanitizer(
			"$.value[?(@.type =~ /cloudvalidations$/i && @.name =~ /^(?!gopv).*/)]."+key,
			"other-cloud-validation", "", nil); err != nil {
			return err
		}
	}
	if err := recording.AddBodyRegexSanitizer("/resourceGroups/other-group", `(?i)/resourceGroups/(?!rg-gopv)[^/"\\]+`, nil); err != nil {
		return err
	}
	return recording.AddBodyRegexSanitizer("other-cloud-validation",
		`(?i)/subscriptions/[^/"\\]+/resourceGroups/[^/"\\]+/providers/Microsoft\.PlatformValidation/cloudValidations/(?<name>(?!gopv)[^/"\\?]+)`,
		&recording.RecordingOptions{UseHTTPS: true, GroupForReplace: "name"})
}

type testClients struct {
	factory        *armplatformvalidation.ClientFactory
	credential     azcore.TokenCredential
	options        *arm.ClientOptions
	subscriptionID string
	location       string
}

func newTestClients(t *testing.T) (*testClients, context.Context) {
	t.Helper()
	t.Cleanup(testutil.StartRecording(t, pathToPackage))
	cred, options := testutil.GetCredAndClientOptions(t)
	endpoint := recording.GetEnvVariable("PLATFORMVALIDATION_ARM_ENDPOINT", "https://management.azure.com")
	options.DisableRPRegistration = true
	options.Cloud = cloud.Configuration{
		ActiveDirectoryAuthorityHost: cloud.AzurePublic.ActiveDirectoryAuthorityHost,
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {Endpoint: endpoint, Audience: "https://management.azure.com/"},
		},
	}
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	if recording.GetRecordMode() != recording.PlaybackMode {
		require.NotEmpty(t, os.Getenv("AZURE_SUBSCRIPTION_ID"), "Set the live-test subscription explicitly")
	}
	factory, err := armplatformvalidation.NewClientFactory(subscriptionID, cred, options)
	requireNoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Minute)
	t.Cleanup(cancel)
	return &testClients{
		factory: factory, credential: cred, options: options, subscriptionID: subscriptionID,
		location: recording.GetEnvVariable("LOCATION", "eastus2euap"),
	}, ctx
}

// Service error bodies can echo a plan's signed image URL.
func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	var response *azcore.ResponseError
	if errors.As(err, &response) {
		require.NoError(t, fmt.Errorf("HTTP %d: %s (response body omitted)", response.StatusCode, response.ErrorCode))
	} else {
		require.NoError(t, fmt.Errorf("%T (details omitted to protect signed input)", err))
	}
}

func requireNotFound(t *testing.T, err error) {
	t.Helper()
	var response *azcore.ResponseError
	require.True(t, errors.As(err, &response), "Expected an ARM response error")
	require.Equal(t, http.StatusNotFound, response.StatusCode)
}

func stringValue[T ~string](value *T) string {
	if value == nil {
		return ""
	}
	return string(*value)
}
