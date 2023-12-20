// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGlobalEndpointManagerEmulator(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)
	emulatorRegionName := "South Central US"
	preferredRegions := []string{}
	emulatorRegion := accountRegion{Name: emulatorRegionName, Endpoint: "https://127.0.0.1:8081/"}

	gem, err := newGlobalEndpointManager(client, preferredRegions, 5*time.Minute)
	assert.NoError(t, err)

	accountProps, err := gem.GetAccountProperties(context.Background())
	assert.NoError(t, err)

	// Verify the expected account properties
	expectedAccountProps := accountProperties{
		ReadRegions:                  []accountRegion{emulatorRegion},
		WriteRegions:                 []accountRegion{emulatorRegion},
		EnableMultipleWriteLocations: false,
	}
	assert.Equal(t, expectedAccountProps, accountProps)

	emulatorEndpoint, err := url.Parse("https://localhost:8081/")
	assert.NoError(t, err)

	// Verify the read endpoints
	readEndpoints, err := gem.GetReadEndpoints()
	assert.NoError(t, err)

	expectedEndpoints := []url.URL{
		*emulatorEndpoint,
	}
	assert.Equal(t, expectedEndpoints, readEndpoints)

	// Verify the write endpoints
	writeEndpoints, err := gem.GetWriteEndpoints()
	assert.NoError(t, err)

	assert.Equal(t, expectedEndpoints, writeEndpoints)

	// Assert location cache is not populated until update() is called
	locationInfo := gem.locationCache.locationInfo
	availableLocation := []string{}
	availableEndpointsByLocation := map[string]url.URL{}

	assert.Equal(t, locationInfo.availReadLocations, availableLocation)
	assert.Equal(t, locationInfo.availWriteLocations, availableLocation)
	assert.Equal(t, locationInfo.availReadEndpointsByLocation, availableEndpointsByLocation)
	assert.Equal(t, locationInfo.availWriteEndpointsByLocation, availableEndpointsByLocation)

	//update and assert available locations are now populated in location cache
	err = gem.Update(context.Background())
	assert.NoError(t, err)
	locationInfo = gem.locationCache.locationInfo

	assert.Equal(t, len(locationInfo.availReadLocations), len(availableLocation)+1)
	assert.Equal(t, len(locationInfo.availWriteLocations), len(availableLocation)+1)
	assert.Equal(t, locationInfo.availWriteLocations[0], emulatorRegionName)
	assert.Equal(t, locationInfo.availReadLocations[0], emulatorRegionName)
	assert.Equal(t, len(locationInfo.availReadEndpointsByLocation), len(availableEndpointsByLocation)+1)
	assert.Equal(t, len(locationInfo.availWriteEndpointsByLocation), len(availableEndpointsByLocation)+1)
}
