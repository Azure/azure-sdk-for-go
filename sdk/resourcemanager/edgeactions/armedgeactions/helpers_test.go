// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armedgeactions_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/edgeactions/armedgeactions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
