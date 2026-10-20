// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armedgeactions_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/edgeactions/armedgeactions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelperDefaultParallelism(t *testing.T) {
	for _, tc := range []struct {
		name     string
		args     []string
		expected int
	}{
		{"default", nil, 7},
		{"sequential", []string{"-test.parallel=1"}, 1},
		{"bounded", []string{"-test.parallel=2"}, 2},
		{"higher", []string{"-test.parallel=12"}, 12},
	} {
		t.Run(tc.name, func(t *testing.T) {
			flags := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
			parallelism := flags.Int("test.parallel", 1, "")
			require.NoError(t, flags.Parse(tc.args))
			require.NoError(t, configureTestParallelism(flags))
			assert.Equal(t, tc.expected, *parallelism)
		})
	}
}

func TestHelperLoadTestEnv(t *testing.T) {
	for _, mode := range []string{recording.RecordingMode, recording.LiveMode} {
		t.Run(mode, func(t *testing.T) {
			t.Chdir(t.TempDir())
			t.Setenv("AZURE_SUBSCRIPTION_ID", "shell-subscription")
			t.Setenv("LOCATION", "")
			require.NoError(t, os.Unsetenv("LOCATION"))
			require.NoError(t, os.WriteFile(".env", []byte("AZURE_SUBSCRIPTION_ID=file-subscription\nLOCATION=westus\n"), 0600))
			require.NoError(t, loadTestEnv(mode))
			assert.Equal(t, "shell-subscription", os.Getenv("AZURE_SUBSCRIPTION_ID"))
			assert.Equal(t, "westus", os.Getenv("LOCATION"))
		})
	}
}

func TestHelperLoadTestEnvFiles(t *testing.T) {
	for _, tc := range []struct {
		name, mode, contents string
		wantError            bool
	}{
		{"missing", recording.RecordingMode, "", false},
		{"malformed", recording.RecordingMode, "LOCATION='unterminated", true},
		{"playback-malformed", recording.PlaybackMode, "LOCATION='unterminated", false},
		{"playback-valid", recording.PlaybackMode, "LOCATION=westus\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			t.Setenv("LOCATION", "")
			require.NoError(t, os.Unsetenv("LOCATION"))
			if tc.contents != "" {
				require.NoError(t, os.WriteFile(".env", []byte(tc.contents), 0600))
			}
			err := loadTestEnv(tc.mode)
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			_, exists := os.LookupEnv("LOCATION")
			assert.False(t, exists)
		})
	}
}

func TestHelperResourceGroupName(t *testing.T) {
	for _, mode := range []string{recording.RecordingMode, recording.LiveMode} {
		for _, suffix := range []uint64{0, 123456789, ^uint64(0)} {
			name := scenarioResourceGroupName(mode, suffix)
			assert.Equal(t, fmt.Sprintf("go-sdk-test-%d", suffix), name)
			assert.Regexp(t, `^go-sdk-test-\d+$`, name)
		}
	}
	assert.Equal(t, recording.SanitizedValue, scenarioResourceGroupName(recording.PlaybackMode, 123456789))
}

func TestHelperScenarioTimeout(t *testing.T) {
	for _, tc := range []struct {
		deployments int
		expected    time.Duration
	}{
		{0, 20 * time.Minute},
		{1, 50 * time.Minute},
		{2, 90 * time.Minute},
		{3, 130 * time.Minute},
	} {
		timeout, err := scenarioTimeout(tc.deployments)
		require.NoError(t, err)
		assert.Equal(t, tc.expected, timeout)
	}
	for _, deployments := range []int{-1, 4} {
		_, err := scenarioTimeout(deployments)
		require.ErrorContains(t, err, "unsupported scenario deployment count")
	}
}

func TestHelperVersionCode(t *testing.T) {
	for _, deploymentType := range []armedgeactions.EdgeActionVersionDeploymentType{
		armedgeactions.EdgeActionVersionDeploymentTypeFile, armedgeactions.EdgeActionVersionDeploymentTypeZip,
	} {
		t.Run(string(deploymentType), func(t *testing.T) {
			first, err := versionCode(deploymentType)
			require.NoError(t, err)
			second, err := versionCode(deploymentType)
			require.NoError(t, err)
			assert.Equal(t, first, second)
			assertHandler(t, first)
			content, err := base64.StdEncoding.DecodeString(first)
			require.NoError(t, err)
			if deploymentType == armedgeactions.EdgeActionVersionDeploymentTypeZip {
				archive, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
				require.NoError(t, err)
				require.Len(t, archive.File, 1)
				assert.Equal(t, "handler.js", archive.File[0].Name)
				assert.Equal(t, uint16(zip.Store), archive.File[0].Method)
				assert.Equal(t, uint16(0), archive.File[0].CreatorVersion>>8)
				assert.Equal(t, uint16(33), archive.File[0].ModifiedDate)
				assert.Zero(t, archive.File[0].ModifiedTime)
			}
		})
	}
	_, err := versionCode(armedgeactions.EdgeActionVersionDeploymentTypeOthers)
	require.ErrorContains(t, err, "unsupported deployment type")
}

func TestHelperDeployScenarioVersion(t *testing.T) {
	for _, deploymentType := range []armedgeactions.EdgeActionVersionDeploymentType{
		armedgeactions.EdgeActionVersionDeploymentTypeFile, armedgeactions.EdgeActionVersionDeploymentTypeZip,
	} {
		t.Run(string(deploymentType), func(t *testing.T) {
			factory, server := newMockFactory(t)
			content, err := versionCode(deploymentType)
			require.NoError(t, err)
			response := strings.Replace(versionJSON, `"file"`, fmt.Sprintf("%q", deploymentType), 1)
			expectRequest(t, server, "PUT", versionPath,
				fmt.Sprintf(`{"location":"global","properties":{"deploymentType":%q,"isDefaultVersion":"False"}}`, deploymentType),
				201, response)
			expectRequest(t, server, "POST", versionPath+"/deployVersionCode",
				fmt.Sprintf(`{"name":"v1","content":%q}`, content), 202, "",
				mock.WithHeader("Azure-AsyncOperation", "https://management.azure.com/operations/one?api-version=2026-10-01"),
				mock.WithHeader("Location", "https://management.azure.com/results/one?api-version=2026-10-01"))
			expectRequest(t, server, "GET", "/operations/one", "", 200, `{"status":"Succeeded"}`)
			expectRequest(t, server, "GET", "/results/one", "", 200, response)
			expectRequest(t, server, "GET", versionPath, "", 200, response)
			deployed := deployScenarioVersion(t, t.Context(), factory.NewEdgeActionVersionsClient(),
				"test-rg", "test-action", "v1", deploymentType)
			assert.Equal(t, to.Ptr(versionPath), deployed.ID)
			assert.Equal(t, 5, server.Requests(), "create, deploy, poll, fetch the result, and get the deployed version")
		})
	}
}

func TestHelperDeleteScenarioVersions(t *testing.T) {
	for _, tc := range []struct {
		name     string
		versions []string
	}{
		{"deployment", []string{"v1"}},
		{"swap", []string{"v1", "v2"}},
		{"filters", []string{"v2", "v1"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			factory, server := newMockFactory(t)
			for _, version := range tc.versions {
				expectRequest(t, server, "DELETE", actionPath+"/versions/"+version, "", 204, "")
			}
			expectRequest(t, server, "GET", actionPath+"/versions", "", 200, `{"value":[]}`)
			deleteScenarioVersions(t, t.Context(), factory.NewEdgeActionVersionsClient(),
				"test-rg", "test-action", tc.versions...)
			assert.Equal(t, len(tc.versions)+1, server.Requests(), "delete in scenario order and verify the version list is empty")
		})
	}
}

func TestHelperCleanupOrder(t *testing.T) {
	for _, defaultVersion := range []string{"v1", "v2"} {
		t.Run(defaultVersion, func(t *testing.T) {
			factory, server := newMockFactory(t)
			expectRequest(t, server, "GET", actionPath, "", 200, actionJSON)
			expectRequest(t, server, "GET", actionPath+"/executionFilters", "", 200, `{"value":[`+filterJSON+`]}`)
			states := map[string]string{"v1": "False", "v2": "False"}
			states[defaultVersion] = "True"
			next := "https://management.azure.com" + actionPath + "/versions?api-version=2026-10-01&$skiptoken=next"
			expectRequest(t, server, "GET", actionPath+"/versions", "", 200,
				fmt.Sprintf(`{"value":[{"name":"v1","properties":{"isDefaultVersion":%q,"provisioningState":"Succeeded"}}],"nextLink":%q}`, states["v1"], next))
			expectRequest(t, server, "GET", actionPath+"/versions", "", 200,
				fmt.Sprintf(`{"value":[{"name":"v2","properties":{"isDefaultVersion":%q,"provisioningState":"Succeeded"}}]}`, states["v2"]))
			expectRequest(t, server, "DELETE", filterPath, "", 204, "")
			other := "v1"
			if defaultVersion == other {
				other = "v2"
			}
			for _, version := range []string{other, defaultVersion} {
				expectRequest(t, server, "DELETE", actionPath+"/versions/"+version, "", 204, "")
			}
			expectRequest(t, server, "DELETE", actionPath, "", 204, "")
			require.NoError(t, cleanupAction(context.Background(), factory, "test-rg", "test-action"))
			assert.Equal(t, 8, server.Requests())
		})
	}
}

func TestHelperCleanupRefreshesDefaultState(t *testing.T) {
	factory, server := newMockFactory(t)
	expectRequest(t, server, "GET", actionPath, "", 200, actionJSON)
	expectRequest(t, server, "GET", actionPath+"/executionFilters", "", 200, `{"value":[]}`)
	expectRequest(t, server, "GET", actionPath+"/versions", "", 200,
		`{"value":[{"name":"v1","properties":{"isDefaultVersion":"True","provisioningState":"Updating"}},{"name":"v2","properties":{"isDefaultVersion":"True","provisioningState":"Succeeded"}}]}`)
	expectRequest(t, server, "GET", versionPath, "", 200, versionJSON)
	expectRequest(t, server, "DELETE", versionPath, "", 200, "")
	expectRequest(t, server, "DELETE", actionPath+"/versions/v2", "", 200, "")
	expectRequest(t, server, "DELETE", actionPath, "", 200, "")
	require.NoError(t, cleanupAction(t.Context(), factory, "test-rg", "test-action"))
	assert.Equal(t, 7, server.Requests(), "use refreshed default flags after an in-flight swap")
}

func TestHelperCleanupInvalidState(t *testing.T) {
	for _, value := range []string{
		`[{"name":"v1","properties":{"isDefaultVersion":"True","provisioningState":"Succeeded"}},{"name":"v2","properties":{"isDefaultVersion":"True","provisioningState":"Succeeded"}}]`,
		`[{"name":"v1","properties":{"isDefaultVersion":"unknown","provisioningState":"Succeeded"}}]`,
		`[{"name":"v1","properties":{}}]`, `[{"properties":{"isDefaultVersion":"False"}}]`, `[null]`,
	} {
		t.Run(value, func(t *testing.T) {
			factory, server := newMockFactory(t)
			expectRequest(t, server, "GET", actionPath, "", 200, actionJSON)
			expectRequest(t, server, "GET", actionPath+"/executionFilters", "", 200, `{"value":[]}`)
			expectRequest(t, server, "GET", actionPath+"/versions", "", 200, `{"value":`+value+`}`)
			require.Error(t, cleanupAction(context.Background(), factory, "test-rg", "test-action"))
			assert.Equal(t, 3, server.Requests(), "ambiguous state must not issue deletes")
		})
	}
}

func TestHelperCleanupErrors(t *testing.T) {
	for _, status := range []int{http.StatusNotFound, http.StatusForbidden} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			factory, server := newMockFactory(t)
			expectRequest(t, server, "GET", actionPath, "", status, `{"error":{"code":"TestError","message":"lookup failed"}}`)
			err := cleanupAction(context.Background(), factory, "test-rg", "test-action")
			if status == http.StatusNotFound {
				require.NoError(t, err)
			} else {
				requireStatus(t, err, status)
			}
			assert.Equal(t, 1, server.Requests())
		})
	}
	t.Run("child-delete-failure", func(t *testing.T) {
		factory, server := newMockFactory(t)
		expectRequest(t, server, "GET", actionPath, "", 200, actionJSON)
		expectRequest(t, server, "GET", actionPath+"/executionFilters", "", 200, `{"value":[`+filterJSON+`]}`)
		expectRequest(t, server, "GET", actionPath+"/versions", "", 200, `{"value":[]}`)
		expectRequest(t, server, "DELETE", filterPath, "", 409, `{"error":{"code":"Conflict","message":"still in use"}}`)
		requireStatus(t, cleanupAction(context.Background(), factory, "test-rg", "test-action"), 409)
		assert.Equal(t, 4, server.Requests(), "do not delete the parent when child deletion fails")
	})
}

func TestHelperCleanupWaitsForTerminalState(t *testing.T) {
	ctx := t.Context()
	for _, state := range []armedgeactions.ProvisioningState{
		armedgeactions.ProvisioningStateSucceeded, armedgeactions.ProvisioningStateFailed, armedgeactions.ProvisioningStateCanceled,
	} {
		calls := 0
		require.NoError(t, waitForTerminal(ctx, to.Ptr(armedgeactions.ProvisioningState("Creating")), func(context.Context) (*armedgeactions.ProvisioningState, error) {
			calls++
			return &state, nil
		}))
		assert.Equal(t, 1, calls)
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	require.ErrorIs(t, waitForTerminal(canceled, to.Ptr(armedgeactions.ProvisioningState("Creating")), nil), context.Canceled)
	require.ErrorContains(t, waitForTerminal(ctx, nil, nil), "requires a provisioning state")
	sentinel := fmt.Errorf("read failed")
	require.ErrorIs(t, waitForTerminal(ctx, to.Ptr(armedgeactions.ProvisioningState("Creating")), func(context.Context) (*armedgeactions.ProvisioningState, error) {
		return nil, sentinel
	}), sentinel)
}

func TestHelperCleanupPollsBeforeDelete(t *testing.T) {
	factory, server := newMockFactory(t)
	expectRequest(t, server, "GET", actionPath, "", 200, actionJSON)
	expectRequest(t, server, "GET", actionPath+"/executionFilters", "", 200, `{"value":[{"name":"sdkfilter","properties":{"provisioningState":"Creating"}}]}`)
	expectRequest(t, server, "GET", actionPath+"/versions", "", 200, `{"value":[]}`)
	expectRequest(t, server, "GET", filterPath, "", 200, filterJSON)
	expectRequest(t, server, "DELETE", filterPath, "", 200, "")
	expectRequest(t, server, "DELETE", actionPath, "", 200, "")
	require.NoError(t, cleanupAction(t.Context(), factory, "test-rg", "test-action"))
	assert.Equal(t, 6, server.Requests())
}

func TestHelperOwnedAction(t *testing.T) {
	ctx := context.Background()
	for _, status := range []int{http.StatusOK, http.StatusForbidden} {
		t.Run("preexisting-"+http.StatusText(status), func(t *testing.T) {
			factory, server := newMockFactory(t)
			expectRequest(t, server, "GET", actionPath, "", status, actionJSON)
			cleanup, err := createOwnedAction(ctx, factory, "test-rg", "test-action")
			require.Error(t, err)
			assert.Nil(t, cleanup)
			assert.Equal(t, 1, server.Requests(), "never create or delete an unowned action")
		})
	}
	for _, failure := range []string{"none", "create", "poll"} {
		t.Run(failure, func(t *testing.T) {
			factory, server := newMockFactory(t)
			expectRequest(t, server, "GET", actionPath, "", 404, `{"error":{"code":"ResourceNotFound"}}`)
			switch failure {
			case "none":
				expectRequest(t, server, "PUT", actionPath, "", 200, actionJSON)
			case "create":
				expectRequest(t, server, "PUT", actionPath, "", 400, `{"error":{"code":"BadRequest","message":"create failed"}}`)
			case "poll":
				expectRequest(t, server, "PUT", actionPath, "", 201, `{"properties":{"provisioningState":"Provisioning"}}`,
					mock.WithHeader("Azure-AsyncOperation", "https://management.azure.com/operations/create?api-version=2026-10-01"))
				expectRequest(t, server, "GET", "/operations/create", "", 200, `{"status":"Failed","error":{"code":"Failed","message":"create failed"}}`)
			}
			cleanup, err := createOwnedAction(ctx, factory, "test-rg", "test-action")
			require.NotNil(t, cleanup, "a partially created action still needs cleanup")
			if failure == "none" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			expectRequest(t, server, "GET", actionPath, "", 200, actionJSON)
			expectRequest(t, server, "GET", actionPath+"/executionFilters", "", 200, `{"value":[]}`)
			expectRequest(t, server, "GET", actionPath+"/versions", "", 200, `{"value":[]}`)
			expectRequest(t, server, "DELETE", actionPath, "", 204, "")
			require.NoError(t, cleanup(ctx))
			expected := 6
			if failure == "poll" {
				expected++
			}
			assert.Equal(t, expected, server.Requests())
		})
	}
}
