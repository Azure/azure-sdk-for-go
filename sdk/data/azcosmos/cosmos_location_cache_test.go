// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azcosmos

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

var defaultEndpoint *url.URL
var loc1Endpoint *url.URL
var loc2Endpoint *url.URL
var loc3Endpoint *url.URL
var loc4Endpoint *url.URL
var writeEndpoints []url.URL
var readEndpoints []url.URL
var endpointsByLoc map[regionId]url.URL
var loc1 accountRegion
var loc2 accountRegion
var loc3 accountRegion
var loc4 accountRegion
var prefLocs []string

func TestMain(m *testing.M) {
	var err error
	defaultEndpoint, err = url.Parse("https://default.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse default endpoint URI")
		os.Exit(1)
	}
	loc1Endpoint, err = url.Parse("https://location1.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse location1 endpoint URI")
		os.Exit(1)
	}
	loc2Endpoint, err = url.Parse("https://location2.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse location2 endpoint URI")
		os.Exit(1)
	}
	loc3Endpoint, err = url.Parse("https://location3.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse location3 endpoint URI")
		os.Exit(1)
	}
	loc4Endpoint, err = url.Parse("https://location4.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse location4 endpoint URI")
		os.Exit(1)
	}

	loc1 = accountRegion{Name: newRegionId("location1"), Endpoint: loc1Endpoint.String()}
	loc2 = accountRegion{Name: newRegionId("location2"), Endpoint: loc2Endpoint.String()}
	loc3 = accountRegion{Name: newRegionId("location3"), Endpoint: loc3Endpoint.String()}
	loc4 = accountRegion{Name: newRegionId("location4"), Endpoint: loc4Endpoint.String()}

	writeEndpoints = []url.URL{*loc1Endpoint, *loc2Endpoint, *loc3Endpoint}
	readEndpoints = []url.URL{*loc1Endpoint, *loc2Endpoint, *loc4Endpoint}
	endpointsByLoc = map[regionId]url.URL{newRegionId("location1"): *loc1Endpoint, newRegionId("location2"): *loc2Endpoint, newRegionId("location3"): *loc3Endpoint, newRegionId("location4"): *loc4Endpoint}

	prefLocs = make([]string, 0)

	status := m.Run()
	os.Exit(status)
}

func endpointListToString(endpoints []url.URL) string {
	var endpointStr strings.Builder
	for _, endpoint := range endpoints {
		endpointStr.WriteString(endpoint.String() + ",")
	}
	return endpointStr.String()
}

func CreateDatabaseAccount(useMultipleWriteLocations bool, enforceSingleMasterWriteLoc bool) accountProperties {
	writeRegions := []accountRegion{loc1, loc2, loc3}
	if !useMultipleWriteLocations && enforceSingleMasterWriteLoc {
		writeRegions = []accountRegion{loc1}
	}
	readRegions := []accountRegion{loc1, loc2, loc4}
	return accountProperties{WriteRegions: writeRegions, ReadRegions: readRegions, EnableMultipleWriteLocations: useMultipleWriteLocations}
}

func ResetLocationCache() *locationCache {
	lc := newLocationCache(prefLocs, *defaultEndpoint, true)
	lc.enableCrossRegionRetries = true
	return lc
}

func TestMarkEndpointUnavailable(t *testing.T) {
	lc := ResetLocationCache()
	var firstCheckTime time.Time
	// mark endpoint unavailable for first time
	err := lc.markEndpointUnavailableForRead(*loc1Endpoint)
	if err != nil {
		t.Fatalf("Received error marking endpoint unavailable: %s", err.Error())
	}
	if info, ok := lc.locationUnavailabilityInfoMap[*loc1Endpoint]; ok {
		var zeroTime time.Time
		if firstCheckTime = info.lastCheckTime; firstCheckTime == zeroTime {
			t.Errorf("Expected lastCheckTime to be set, but was not")
		}
		if info.unavailableOps != read {
			t.Errorf("Expected unavailableOps to be 1 (read-only), but was %d", info.unavailableOps)
		}
	} else {
		t.Errorf("Expected locationUnavailabilityInfoMap to contain %s, but it did not", loc1Endpoint.String())
	}
	// mark endpoint unavailable for second time
	time.Sleep(100 * time.Millisecond)
	err = lc.markEndpointUnavailableForWrite(*loc1Endpoint)
	if err != nil {
		t.Fatalf("Received error marking endpoint unavailable: %s", err.Error())
	}
	if info, ok := lc.locationUnavailabilityInfoMap[*loc1Endpoint]; ok {
		var zeroTime time.Time
		if info.lastCheckTime == zeroTime || info.lastCheckTime == firstCheckTime {
			t.Errorf("Expected lastCheckTime to be updated, but was not. First check time: %s, last check time: %s", firstCheckTime, info.lastCheckTime)
		}
		if info.unavailableOps != all {
			t.Errorf("Expected unavailableOps to be 3 (read+write), but was %d", info.unavailableOps)
		}
	} else {
		t.Errorf("Expected locationUnavailabilityInfoMap to contain %s, but it did not", loc1Endpoint.String())
	}
}

func TestRefreshStaleEndpoints(t *testing.T) {
	lc := ResetLocationCache()
	// mark endpoint unavailable for first time
	err := lc.markEndpointUnavailableForRead(*loc1Endpoint)
	if err != nil {
		t.Fatalf("Received error marking endpoint unavailable: %s", err.Error())
	}
	if info, ok := lc.locationUnavailabilityInfoMap[*loc1Endpoint]; ok {
		info.lastCheckTime = time.Now().Add(-1*defaultExpirationTime - 1*time.Second)
		lc.locationUnavailabilityInfoMap[*loc1Endpoint] = info
	} else {
		t.Errorf("Expected locationUnavailabilityInfoMap to contain %s, but it did not", loc1Endpoint.String())
	}
	// refresh stale endpoints, since time since last check is greater default expiration time
	lc.refreshStaleEndpoints()
	if len(lc.locationUnavailabilityInfoMap) != 0 {
		t.Errorf("Expected locationUnavailabilityInfoMap to be empty, but it was not")
	}
}

func TestIsEndpointUnavailable(t *testing.T) {
	lc := ResetLocationCache()
	err := lc.markEndpointUnavailableForRead(*loc1Endpoint)
	if err != nil {
		t.Fatalf("Received error marking endpoint unavailable: %s", err.Error())
	}
	err = lc.markEndpointUnavailableForWrite(*loc2Endpoint)
	if err != nil {
		t.Fatalf("Received error marking endpoint unavailable: %s", err.Error())
	}

	if lc.isEndpointUnavailable(*loc1Endpoint, none) {
		t.Errorf("Expected IsEndpointUnavailable to return false, but it returned true for ops = none")
	}
	if lc.isEndpointUnavailable(*loc1Endpoint, write) {
		t.Errorf("Expected IsEndpointUnavailable to return false, but it returned true for ops = write when region is unavailable for read")
	}
	if lc.isEndpointUnavailable(*loc3Endpoint, all) {
		t.Errorf("Expected IsEndpointUnavailable to return false, but it returned true for an available region")
	}
	if !lc.isEndpointUnavailable(*loc1Endpoint, read) {
		t.Errorf("Expected IsEndpointUnavailable to return true, but it returned false for ops = read when region is unavailable for read")
	}

	if info, ok := lc.locationUnavailabilityInfoMap[*loc1Endpoint]; ok {
		info.lastCheckTime = time.Now().Add(-1*defaultExpirationTime - 1*time.Second)
		lc.locationUnavailabilityInfoMap[*loc1Endpoint] = info
	} else {
		t.Errorf("Expected locationUnavailabilityInfoMap to contain %s, but it did not", loc1Endpoint.String())
	}

	if lc.isEndpointUnavailable(*loc1Endpoint, read) {
		t.Errorf("Expected IsEndpointUnavailable to return false, but it returned true stale unavailability")
	}
}

func TestGetLocation(t *testing.T) {
	lc := ResetLocationCache()
	dbAcct := CreateDatabaseAccount(lc.enableMultipleWriteLocations, false)
	err := lc.databaseAccountRead(dbAcct)
	if err != nil {
		t.Fatalf("Received error Reading DB account: %s", err.Error())
	}
	if len(dbAcct.WriteRegions) == 0 {
		t.Fatal("Write Regions are empty")
	}
	actual := lc.getLocation(*defaultEndpoint)
	if actual == "" {
		t.Errorf("Expected GetLocation to return a valid location when provided the default endpoint, but it did not")
	}
	for _, region := range dbAcct.WriteRegions {
		url, err := url.Parse(region.Endpoint)
		if err != nil {
			t.Errorf("Failed to parse endpoint %s, %s", region.Endpoint, err)
			continue
		}
		expected, actual := region.Name, lc.getLocation(*url)
		if !expected.Equal(actual) {
			t.Errorf("Expected GetLocation to return Write Region %s, but was %s", expected, actual)
		}
	}

	for _, region := range dbAcct.ReadRegions {
		url, err := url.Parse(region.Endpoint)
		if err != nil {
			t.Errorf("Failed to parse endpoint %s, %s", region.Endpoint, err)
			continue
		}
		expected, actual := region.Name, lc.getLocation(*url)
		if !expected.Equal(actual) {
			t.Errorf("Expected GetLocation to return Read Region %s, but was %s", expected, actual)
		}
	}
}

func TestGetEndpointsByLocation(t *testing.T) {
	locs := []accountRegion{loc1, loc2, loc3, loc4}
	newEndpointsByLoc, parsedLocs, err := getEndpointsByLocation(locs)
	if err != nil {
		t.Fatalf("Received error getting endpoints by location: %s", err.Error())
	}
	if len(newEndpointsByLoc) != len(endpointsByLoc) {
		t.Errorf("Expected %d endpoints, but got %d", len(endpointsByLoc), len(newEndpointsByLoc))
	}
	for loc, endpoint := range endpointsByLoc {
		if newEndpoint, ok := newEndpointsByLoc[loc]; ok {
			if newEndpoint != endpoint {
				t.Errorf("Expected endpoint %s for location %s, but was %s", endpoint.String(), loc, newEndpoint.String())
			}
		} else {
			t.Errorf("Expected newEndpointsByLoc to contain location %s, but it did not", loc)
		}
	}

	if len(parsedLocs) != len(locs) {
		t.Errorf("Expected parsedLocs to contain %d locations, but it contained %d", len(locs), len(parsedLocs))
	}
	for i, loc := range locs {
		if parsedLocs[i] != loc.Name {
			t.Errorf("Expected parsedLocs to contain location %s, but it did not", loc.Name)
		}
	}
}

func TestGetPrefAvailableEndpoints(t *testing.T) {
	lc := ResetLocationCache()
	lc.enableMultipleWriteLocations = true
	dbAcct := CreateDatabaseAccount(lc.enableMultipleWriteLocations, false)
	// will set write locations to loc1, loc2, loc3
	err := lc.databaseAccountRead(dbAcct)
	if err != nil {
		t.Fatalf("Received error Reading DB account: %s", err.Error())
	}
	// marks loc1 unavailable, which will put it last in the preferred available endpoint list
	err = lc.markEndpointUnavailableForWrite(*loc1Endpoint)
	if err != nil {
		t.Fatalf("Received error marking endpoint unavailable: %s", err.Error())
	}
	// loc1: unavailable, loc2: available, loc5: non-existent
	lc.locationInfo.prefLocations = []regionId{loc1.Name, loc2.Name, newRegionId("location5")}
	prefWriteEndpoints := lc.getPrefAvailableEndpoints(lc.locationInfo.availWriteEndpointsByLocation, lc.locationInfo.availWriteLocations, write, lc.defaultEndpoint)
	prefWriteEndpointsString := endpointListToString(prefWriteEndpoints)

	// loc2: preferred + available, default: fallback endpoint, loc1: unavailable + preferred
	expectedWriteEndpoints := endpointListToString([]url.URL{*loc2Endpoint, *defaultEndpoint, *loc1Endpoint})

	if prefWriteEndpointsString != expectedWriteEndpoints {
		t.Errorf("Expected preferred available write endpoints to be %s, but was %s", expectedWriteEndpoints, prefWriteEndpointsString)
	}
}

func TestReadEndpoints(t *testing.T) {
	lc := ResetLocationCache()
	lc.locationInfo.prefLocations = []regionId{loc1.Name, loc2.Name, loc3.Name, loc4.Name}
	dbAcct := CreateDatabaseAccount(lc.enableMultipleWriteLocations, false)
	err := lc.databaseAccountRead(dbAcct)
	if err != nil {
		t.Fatalf("Received error Reading DB account: %s", err.Error())
	}

	lc.lastUpdateTime = time.Now().Add(-1*defaultExpirationTime - 1*time.Second)
	expectedReadEndpoints := endpointListToString([]url.URL{*loc1Endpoint, *loc2Endpoint, *loc4Endpoint})
	actualReadEndpoints, err := lc.readEndpoints()
	if err != nil {
		t.Fatalf("Received error getting read endpoints: %s", err.Error())
	}
	actualReadEndpointsOrder := endpointListToString(actualReadEndpoints)
	if actualReadEndpointsOrder != expectedReadEndpoints {
		t.Errorf("Expected read endpoints %s, but was %s", expectedReadEndpoints, actualReadEndpointsOrder)
	}

	lc.lastUpdateTime = time.Now().Add(-1*defaultExpirationTime - 1*time.Second)
	err = lc.markEndpointUnavailableForRead(*loc2Endpoint)
	if err != nil {
		t.Fatalf("Received error marking endpoint unavailable: %s", err.Error())
	}
	// loc2 is now unavailable, so it gets moved to the end of the list
	expectedReadEndpoints = endpointListToString([]url.URL{*loc1Endpoint, *loc4Endpoint, *loc2Endpoint})
	actualReadEndpoints, err = lc.readEndpoints()
	if err != nil {
		t.Fatalf("Received error getting read endpoints: %s", err.Error())
	}
	actualReadEndpointsOrder = endpointListToString(actualReadEndpoints)
	if actualReadEndpointsOrder != expectedReadEndpoints {
		t.Errorf("Expected read endpoints %s, but was %s", expectedReadEndpoints, actualReadEndpointsOrder)
	}
}

func TestWriteEndpoints(t *testing.T) {
	lc := ResetLocationCache()
	lc.enableMultipleWriteLocations = true
	lc.locationInfo.prefLocations = []regionId{loc1.Name, loc2.Name, loc3.Name, loc4.Name}
	dbAcct := CreateDatabaseAccount(lc.enableMultipleWriteLocations, false)
	err := lc.databaseAccountRead(dbAcct)
	if err != nil {
		t.Fatalf("Received error Reading DB account: %s", err.Error())
	}

	lc.lastUpdateTime = time.Now().Add(-1*defaultExpirationTime - 1*time.Second)
	expectedWriteEndpoints := []*url.URL{loc1Endpoint, loc2Endpoint, loc3Endpoint, defaultEndpoint}
	actualWriteEndpoints, err := lc.writeEndpoints()
	if err != nil {
		t.Fatalf("Received error getting write endpoints: %s", err.Error())
	}
	if len(expectedWriteEndpoints) != len(actualWriteEndpoints) {
		t.Errorf("Expected %d write endpoints, but got %d", len(expectedWriteEndpoints), len(actualWriteEndpoints))
	} else {
		for i, endpoint := range expectedWriteEndpoints {
			if endpoint.String() != actualWriteEndpoints[i].String() {
				t.Errorf("Expected endpoint %s, but was %s", endpoint.String(), actualWriteEndpoints[i].String())
			}
		}
	}

	lc.lastUpdateTime = time.Now().Add(-1*defaultExpirationTime - 1*time.Second)
	err = lc.markEndpointUnavailableForWrite(*loc1Endpoint)
	if err != nil {
		t.Fatalf("Received error marking endpoint unavailable: %s", err.Error())
	}
	expectedWriteEndpoints = []*url.URL{loc2Endpoint, loc3Endpoint, defaultEndpoint, loc1Endpoint}
	actualWriteEndpoints, err = lc.writeEndpoints()
	if err != nil {
		t.Fatalf("Received error getting write endpoints: %s", err.Error())
	}
	if len(expectedWriteEndpoints) != len(actualWriteEndpoints) {
		t.Errorf("Expected %d write endpoints, but got %d", len(expectedWriteEndpoints), len(actualWriteEndpoints))
	} else {
		for i, endpoint := range expectedWriteEndpoints {
			if endpoint.String() != actualWriteEndpoints[i].String() {
				t.Errorf("Expected endpoint %s, but was %s", endpoint.String(), actualWriteEndpoints[i].String())
			}
		}
	}
}

func TestReadEndpointsUsesDefaultEndpointAsFallback(t *testing.T) {
	lc := ResetLocationCache()
	lc.enableMultipleWriteLocations = true
	lc.locationInfo.prefLocations = []regionId{loc1.Name, loc2.Name, loc3.Name}
	dbAcct := CreateDatabaseAccount(lc.enableMultipleWriteLocations, false)
	err := lc.databaseAccountRead(dbAcct)
	if err != nil {
		t.Fatalf("Received error Reading DB account: %s", err.Error())
	}

	lc.lastUpdateTime = time.Now().Add(-1*defaultExpirationTime - 1*time.Second)

	// Get write endpoints - loc1 should be first (most preferred)
	actualWriteEndpoints, err := lc.writeEndpoints()
	if err != nil {
		t.Fatalf("Received error getting write endpoints: %s", err.Error())
	}

	// Get read endpoints
	actualReadEndpoints, err := lc.readEndpoints()
	if err != nil {
		t.Fatalf("Received error getting read endpoints: %s", err.Error())
	}

	// Assert the expected ordering of endpoints
	actualWriteEndpointsOrder := endpointListToString(actualWriteEndpoints)
	actualReadEndpointsOrder := endpointListToString(actualReadEndpoints)

	// The correct order is:
	// For write: loc1, loc2, loc3, default
	// For read: loc1, loc2
	// loc3 is not enabled for read, and loc4 (which IS enabled for read) is not in the preferred locations list
	expectedWriteEndpointOrder := endpointListToString([]url.URL{*loc1Endpoint, *loc2Endpoint, *loc3Endpoint, *defaultEndpoint})
	expectedReadEndpointOrder := endpointListToString([]url.URL{*loc1Endpoint, *loc2Endpoint})

	if actualWriteEndpointsOrder != expectedWriteEndpointOrder {
		t.Errorf("Expected write endpoint order %s, but was %s", expectedWriteEndpointOrder, actualWriteEndpointsOrder)
	}

	if actualReadEndpointsOrder != expectedReadEndpointOrder {
		t.Errorf("Expected read endpoint order %s, but was %s", expectedReadEndpointOrder, actualReadEndpointsOrder)
	}
}

func TestWriteEndpointUsedAsFallbackForRead(t *testing.T) {
	lc := ResetLocationCache()
	lc.enableMultipleWriteLocations = true
	lc.locationInfo.prefLocations = []regionId{loc3.Name, loc2.Name, loc1.Name}
	dbAcct := CreateDatabaseAccount(lc.enableMultipleWriteLocations, false)
	err := lc.databaseAccountRead(dbAcct)
	if err != nil {
		t.Fatalf("Received error Reading DB account: %s", err.Error())
	}

	lc.lastUpdateTime = time.Now().Add(-1*defaultExpirationTime - 1*time.Second)

	// Get write endpoints - loc1 should be first (most preferred)
	actualWriteEndpoints, err := lc.writeEndpoints()
	if err != nil {
		t.Fatalf("Received error getting write endpoints: %s", err.Error())
	}

	// Get read endpoints
	actualReadEndpoints, err := lc.readEndpoints()
	if err != nil {
		t.Fatalf("Received error getting read endpoints: %s", err.Error())
	}

	// Assert the expected ordering of endpoints
	actualWriteEndpointsOrder := endpointListToString(actualWriteEndpoints)
	actualReadEndpointsOrder := endpointListToString(actualReadEndpoints)

	// The correct order is:
	// For write: loc3, loc2, loc1, default
	// For read: loc2, loc1, loc3
	// loc3 is not enabled for read, but since it's the most-preferred write location, it should be used as the fallback for read
	expectedWriteEndpointOrder := endpointListToString([]url.URL{*loc3Endpoint, *loc2Endpoint, *loc1Endpoint, *defaultEndpoint})
	expectedReadEndpointOrder := endpointListToString([]url.URL{*loc2Endpoint, *loc1Endpoint, *loc3Endpoint})

	if actualWriteEndpointsOrder != expectedWriteEndpointOrder {
		t.Errorf("Expected write endpoint order %s, but was %s", expectedWriteEndpointOrder, actualWriteEndpointsOrder)
	}

	if actualReadEndpointsOrder != expectedReadEndpointOrder {
		t.Errorf("Expected read endpoint order %s, but was %s", expectedReadEndpointOrder, actualReadEndpointsOrder)
	}
}
