// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armedgeactions_test

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/joho/godotenv"
)

const (
	defaultTestParallelism = 7
	pathToPackage          = "sdk/resourcemanager/edgeactions/armedgeactions/testdata"
	subscriptionID         = "00000000-0000-0000-0000-000000000000"
	// The proxy uses .NET regular expressions. Replace only the named group, not the query key.
	pollingSecretRegex = `(?i)[?&][tcsh]=(?<secret>[^&\s"\\#]+)`
)

var traceHeaders = []string{
	"x-ms-correlation-request-id", "x-ms-routing-request-id",
	"x-ms-operation-identifier", "x-ms-original-request-ids", "X-MSEdge-Ref",
}

func TestMain(m *testing.M) {
	flag.Parse()
	if err := configureTestParallelism(flag.CommandLine); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(runTests(m))
}

func configureTestParallelism(flags *flag.FlagSet) error {
	explicit := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "test.parallel" {
			explicit = true
		}
	})
	if explicit {
		return nil
	}
	return flags.Set("test.parallel", fmt.Sprint(defaultTestParallelism))
}

func runTests(m *testing.M) int {
	if err := loadTestEnv(recording.GetRecordMode()); err != nil {
		// Dotenv parse errors can contain file contents, including credentials.
		fmt.Fprintln(os.Stderr, "Failed to load Edge Actions .env; check file access and dotenv syntax.")
		return 1
	}
	stop := testutil.StartProxy(pathToPackage)
	defer stop()
	if recording.GetRecordMode() != recording.LiveMode {
		if err := registerSanitizers(os.Getenv, &recording.RecordingOptions{UseHTTPS: true}); err != nil {
			panic(fmt.Sprintf("configure Edge Actions sanitizers: %v", err))
		}
	}
	return m.Run()
}

func loadTestEnv(mode string) error {
	if mode == recording.PlaybackMode {
		return nil
	}
	err := godotenv.Load()
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func registerSanitizers(getenv func(string) string, options *recording.RecordingOptions) error {
	// ARM testutil already preserves LRO headers and sanitizes UUIDs and resource group names.
	// Keep resource identity for version references, list assertions, and SKU comparisons.
	if err := recording.RemoveRegisteredSanitizers([]string{"AZSDK3430", "AZSDK3493"}, options); err != nil {
		return err
	}
	for _, key := range []string{"AZURE_SUBSCRIPTION_ID", "AZURE_TENANT_ID", "AZURE_CLIENT_ID", "AZURE_CLIENT_SECRET"} {
		if value := getenv(key); value != "" {
			if err := recording.AddGeneralRegexSanitizer(subscriptionID, regexp.QuoteMeta(value), options); err != nil {
				return err
			}
		}
	}
	queryOptions := *options
	queryOptions.GroupForReplace = "secret"
	if err := recording.AddGeneralRegexSanitizer(recording.SanitizedValue, pollingSecretRegex, &queryOptions); err != nil {
		return err
	}
	for _, header := range traceHeaders {
		if err := recording.AddHeaderRegexSanitizer(header, recording.SanitizedValue, "", options); err != nil {
			return err
		}
	}
	for header, value := range map[string]string{"Cookie": "cookie;", "Set-Cookie": "[set-cookie;]"} {
		if err := recording.AddHeaderRegexSanitizer(header, value, "", options); err != nil {
			return err
		}
	}
	for _, token := range []string{"access_token", "refresh_token", "id_token"} {
		if err := recording.AddBodyKeySanitizer("$.."+token, token, "", options); err != nil {
			return err
		}
		for _, field := range []string{"tenantId", "clientId", "principalId", "client_secret"} {
			if err := recording.AddBodyKeySanitizer("$.."+field, recording.SanitizedValue, "", options); err != nil {
				return err
			}
		}
	}
	for _, field := range []string{"createdBy", "lastModifiedBy"} {
		if err := recording.AddBodyKeySanitizer("$..systemData."+field, recording.SanitizedValue, "", options); err != nil {
			return err
		}
	}
	return nil
}
