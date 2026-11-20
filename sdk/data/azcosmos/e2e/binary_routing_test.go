// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	"github.com/stretchr/testify/require"
)

type binaryWireAudit struct {
	mu       sync.Mutex
	requests []bool
	bodies   []bool
}

func rewriteAccountEndpoints(response *http.Response, endpoint string) (err error) {
	if response.Request.Method != http.MethodGet || response.Request.URL.Path != "/" || response.StatusCode != http.StatusOK {
		return nil
	}
	original := response.Body
	defer func() { err = errors.Join(err, original.Close()) }()
	var properties map[string]json.RawMessage
	if err := json.NewDecoder(original).Decode(&properties); err != nil {
		return err
	}
	for _, field := range []string{"readableLocations", "writableLocations"} {
		var locations []map[string]any
		if err := json.Unmarshal(properties[field], &locations); err != nil {
			return err
		}
		for _, location := range locations {
			location["databaseAccountEndpoint"] = endpoint + "/"
		}
		properties[field], err = json.Marshal(locations)
		if err != nil {
			return err
		}
	}
	body, err := json.Marshal(properties)
	if err != nil {
		return err
	}
	response.Body = io.NopCloser(bytes.NewReader(body))
	response.ContentLength = int64(len(body))
	response.Header.Set("Content-Length", strconv.Itoa(len(body)))
	return nil
}

func newBinaryAuditProxy(t *testing.T, endpoint string) (string, *binaryWireAudit) {
	t.Helper()
	target, err := url.Parse(endpoint)
	require.NoError(t, err)
	audit := &binaryWireAudit{}
	proxy := httputil.NewSingleHostReverseProxy(target)
	var proxyURL string
	proxy.ModifyResponse = func(response *http.Response) error {
		if err := rewriteAccountEndpoints(response, proxyURL); err != nil {
			return err
		}
		if response.Request.Method != http.MethodGet || !strings.Contains(response.Request.URL.Path, "/docs/") || response.StatusCode != http.StatusOK {
			return nil
		}
		body, err := io.ReadAll(response.Body)
		closeErr := response.Body.Close()
		if err != nil || closeErr != nil {
			return errors.Join(err, closeErr)
		}
		response.Body = io.NopCloser(bytes.NewReader(body))
		audit.mu.Lock()
		defer audit.mu.Unlock()
		audit.requests = append(audit.requests, strings.Contains(response.Request.Header.Get("x-ms-cosmos-supported-serialization-formats"), "CosmosBinary"))
		audit.bodies = append(audit.bodies, len(body) > 0 && body[0] == 0x80)
		return nil
	}
	server := httptest.NewUnstartedServer(proxy)
	proxyURL = "http://" + server.Listener.Addr().String()
	server.Config.Protocols = new(http.Protocols)
	server.Config.Protocols.SetHTTP1(true)
	server.Config.Protocols.SetUnencryptedHTTP2(true)
	server.Start()
	t.Cleanup(server.Close)
	return server.URL, audit
}

func TestBinaryEncodingAndRouting(t *testing.T) {
	profile := selectedProfile(t, "configuration.binary-routing")
	endpoint, managementEndpoint, databaseID, containerID := requireProfileEmulator(t, profile)
	require.NoError(t, validateEmulatorBackend(os.Getenv("AZURE_COSMOS_E2E_BACKEND")))
	probeEmulatorCapabilities(t, profile, endpoint, managementEndpoint)
	ctx := context.Background()

	proxyEndpoint, audit := newBinaryAuditProxy(t, endpoint)
	runtime, err := azcosmos.NewRuntime(nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	fx := newFixtureAt(t, proxyEndpoint, databaseID, containerID, "", &azcosmos.ClientOptions{
		Runtime: runtime, Operation: azcosmos.OperationOptions{BinaryEncoding: &azcosmos.BinaryEncodingOptions{Disabled: true}},
	})
	id := uniqueID(t)
	pk := azcosmos.NewPartitionKeyString(id)
	expected := item(id, id, 42)
	_, err = fx.Container.CreateItem(ctx, pk, id, mustMarshalItem(t, expected), nil)
	require.NoError(t, err)
	trackItem(t, fx.Container, pk, id)
	enabled := newFixtureAt(t, proxyEndpoint, databaseID, containerID, "", &azcosmos.ClientOptions{
		Runtime: runtime, Operation: azcosmos.OperationOptions{BinaryEncoding: &azcosmos.BinaryEncodingOptions{}},
	})
	var etag string
	for _, readCase := range []struct {
		container *azcosmos.ContainerClient
		override  *azcosmos.BinaryEncodingOptions
	}{
		{fx.Container, nil},
		{fx.Container, &azcosmos.BinaryEncodingOptions{}},
		{enabled.Container, nil},
		{enabled.Container, &azcosmos.BinaryEncodingOptions{Disabled: true}},
	} {
		read, err := readCase.container.ReadItem(ctx, pk, id, &azcosmos.ReadItemOptions{
			Operation: azcosmos.OperationOptions{BinaryEncoding: readCase.override},
		})
		require.NoError(t, err)
		require.True(t, json.Valid(read.Value))
		require.Equal(t, expected, mustUnmarshalItem(t, read.Value))
		if etag == "" {
			etag = string(read.ETag)
			require.NotEmpty(t, etag)
		}
		require.Equal(t, etag, string(read.ETag))
	}
	audit.mu.Lock()
	requests := append([]bool(nil), audit.requests...)
	bodies := append([]bool(nil), audit.bodies...)
	audit.mu.Unlock()
	require.Equal(t, []bool{false, true, true, false}, requests)
	require.Equal(t, []bool{false, true, true, false}, bodies)

	writer := newFixtureAt(t, endpoint, databaseID, containerID, "", profile.clientOptions)
	regionID := uniqueID(t) + "-routing"
	regionPK := azcosmos.NewPartitionKeyString(regionID)
	_, err = writer.Container.CreateItem(ctx, regionPK, regionID, mustMarshalItem(t, item(regionID, regionID, 99)), nil)
	require.NoError(t, err)
	trackItem(t, writer.Container, regionPK, regionID)

	for _, routing := range []struct {
		strategy azcosmos.RoutingStrategy
		region   string
	}{
		{azcosmos.PreferredRegions(azcosmos.RegionWestUS, azcosmos.RegionEastUS), string(azcosmos.RegionWestUS)},
		{azcosmos.RoutingStrategy{}, string(azcosmos.RegionEastUS)},
	} {
		regionFixture := newFixtureAt(t, endpoint, databaseID, containerID, "", &azcosmos.ClientOptions{Routing: routing.strategy})
		var read azcosmos.ItemResponse
		var readErr error
		require.Eventually(t, func() bool {
			read, readErr = regionFixture.Container.ReadItem(ctx, regionPK, regionID, &azcosmos.ReadItemOptions{
				Operation: azcosmos.OperationOptions{ConsistencyStrategy: azcosmos.ReadConsistencyStrategyEventual},
			})
			if readErr == nil {
				return true
			}
			var cosmosErr *azcosmos.Error
			return !errors.As(readErr, &cosmosErr) || cosmosErr.Code != azcosmos.CodeNotFound
		}, 10*time.Second, 100*time.Millisecond)
		require.NoError(t, readErr)
		require.Equal(t, item(regionID, regionID, 99), mustUnmarshalItem(t, read.Value))
		require.NotNil(t, read.Diagnostics)
		require.NotEmpty(t, read.Diagnostics.RegionsContacted)
		require.Equal(t, routing.region, read.Diagnostics.RegionsContacted[0])
	}
}
