// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azcosmos

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"
)

var defaultEndpt *url.URL
var loc1Endpt *url.URL
var loc2Endpt *url.URL
var loc3Endpt *url.URL
var loc4Endpt *url.URL
var writeEndpts []url.URL
var readEndpts []url.URL
var endptsByLoc map[string]url.URL
var lc *LocationCache
var loc1 AcctRegion
var loc2 AcctRegion
var loc3 AcctRegion
var loc4 AcctRegion

func TestMain(m *testing.M) {
	var err error
	defaultEndpt, err = url.Parse("https://default.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse default endpoint URI")
		os.Exit(1)
	}
	loc1Endpt, err = url.Parse("https://location1.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse location1 endpoint URI")
		os.Exit(1)
	}
	loc2Endpt, err = url.Parse("https://location2.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse location2 endpoint URI")
		os.Exit(1)
	}
	loc3Endpt, err = url.Parse("https://location3.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse location3 endpoint URI")
		os.Exit(1)
	}
	loc4Endpt, err = url.Parse("https://location4.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse location4 endpoint URI")
		os.Exit(1)
	}

	loc1 = AcctRegion{name: "location1", endpoint: loc1Endpt.String()}
	loc2 = AcctRegion{name: "location2", endpoint: loc2Endpt.String()}
	loc3 = AcctRegion{name: "location3", endpoint: loc3Endpt.String()}
	loc4 = AcctRegion{name: "location4", endpoint: loc4Endpt.String()}

	writeEndpts = []url.URL{*loc1Endpt, *loc2Endpt, *loc3Endpt}
	readEndpts = []url.URL{*loc1Endpt, *loc2Endpt, *loc4Endpt}
	endptsByLoc = map[string]url.URL{"location1": *loc1Endpt, "location2": *loc2Endpt, "location3": *loc3Endpt, "location4": *loc4Endpt}

	prefLocs := make([]string, 0)
	lc = NewLocationCache(prefLocs, *defaultEndpt)
	lc.enableEndptDiscovery = true
	lc.connLimit = 10

	status := m.Run()
	os.Exit(status)
}

func CreateDbAcct(useMultipleWriteLocations bool, enforceSingleMasterWriteLoc bool) AcctProperties {
	writeRegions := []AcctRegion{loc1, loc2, loc3}
	if !useMultipleWriteLocations && enforceSingleMasterWriteLoc {
		writeRegions = []AcctRegion{loc1}
	}
	readRegions := []AcctRegion{loc1, loc2, loc4}
	return AcctProperties{writeRegions: writeRegions, readRegions: readRegions, enableMultipleWriteLocations: useMultipleWriteLocations}
}

func TestMarkEndptUnavailable(t *testing.T) {
	var firstCheckTime time.Time
	// mark endpoint unavailable for first time
	lc.MarkEndptUnavailableForRead(*loc1Endpt)
	if info, ok := lc.locationUnavailabilityInfoMap[*loc1Endpt]; ok {
		var zeroTime time.Time
		if firstCheckTime = info.lastCheckTime; firstCheckTime == zeroTime {
			t.Errorf("Expected lastCheckTime to be set, but was not")
		}
		if info.unavailableOps != read {
			t.Errorf("Expected unavailableOps to be 1 (read-only), but was %d", info.unavailableOps)
		}
	} else {
		t.Errorf("Expected locationUnavailabilityInfoMap to contain %s, but it did not", loc1Endpt.String())
	}
	// mark endpoint unavailable for second time
	time.Sleep(100 * time.Millisecond)
	lc.MarkEndptUnavailableForWrite(*loc1Endpt)
	if info, ok := lc.locationUnavailabilityInfoMap[*loc1Endpt]; ok {
		var zeroTime time.Time
		if info.lastCheckTime == zeroTime || info.lastCheckTime == firstCheckTime {
			t.Errorf("Expected lastCheckTime to be updated, but was not. First check time: %s, last check time: %s", firstCheckTime, info.lastCheckTime)
		}
		if info.unavailableOps != all {
			t.Errorf("Expected unavailableOps to be 3 (read+write), but was %d", info.unavailableOps)
		}
	} else {
		t.Errorf("Expected locationUnavailabilityInfoMap to contain %s, but it did not", loc1Endpt.String())
	}
}

func TestRefreshStaleEndpts(t *testing.T) {
	// mark endpoint unavailable for first time
	lc.MarkEndptUnavailableForRead(*loc1Endpt)
	if info, ok := lc.locationUnavailabilityInfoMap[*loc1Endpt]; ok {
		info.lastCheckTime = time.Now().Add(-1 * DefaultExpirationTime)
		lc.locationUnavailabilityInfoMap[*loc1Endpt] = info
	} else {
		t.Errorf("Expected locationUnavailabilityInfoMap to contain %s, but it did not", loc1Endpt.String())
	}
	// refresh stale endpoints, since time since last check is greater default expiration time
	lc.RefreshStaleEndpts()
	if len(lc.locationUnavailabilityInfoMap) != 0 {
		t.Errorf("Expected locationUnavailabilityInfoMap to be empty, but it was not")
	}
}

func TestIsEndptUnavailable(t *testing.T) {
	lc.MarkEndptUnavailableForRead(*loc1Endpt)
	lc.MarkEndptUnavailableForWrite(*loc2Endpt)
	if lc.IsEndptUnavailable(*loc1Endpt, none) {
		t.Errorf("Expected IsEndptUnavailable to return false, but it returned true for ops = none")
	}
	if lc.IsEndptUnavailable(*loc1Endpt, write) {
		t.Errorf("Expected IsEndptUnavailable to return false, but it returned true for ops = write when region is unavailable for read")
	}
	if lc.IsEndptUnavailable(*loc3Endpt, all) {
		t.Errorf("Expected IsEndptUnavailable to return false, but it returned true for an available region")
	}
	if !lc.IsEndptUnavailable(*loc1Endpt, read) {
		t.Errorf("Expected IsEndptUnavailable to return true, but it returned false for ops = read when region is unavailable for read")
	}

	if info, ok := lc.locationUnavailabilityInfoMap[*loc1Endpt]; ok {
		info.lastCheckTime = time.Now().Add(-1 * DefaultExpirationTime)
		lc.locationUnavailabilityInfoMap[*loc1Endpt] = info
	} else {
		t.Errorf("Expected locationUnavailabilityInfoMap to contain %s, but it did not", loc1Endpt.String())
	}

	if lc.IsEndptUnavailable(*loc1Endpt, read) {
		t.Errorf("Expected IsEndptUnavailable to return false, but it returned true stale unavailablilty")
	}
}

func TestGetLocation(t *testing.T) {
	dbAcct := CreateDbAcct(lc.enableMultipleWriteLocations, false)
	lc.DbAcctRead(dbAcct) // requires unit test of update
	if dbAcct.writeRegions == nil || len(dbAcct.writeRegions) == 0 {
		t.Fatal("Write Regions are empty")
	}
	expected, actual := dbAcct.writeRegions[0].name, lc.GetLocation(*defaultEndpt)
	if expected != actual {
		t.Errorf("Expected GetLocation to return First Write Region %s, but was %s", expected, actual)
	}
	for _, region := range dbAcct.writeRegions {
		url, err := url.Parse(region.endpoint)
		if err != nil {
			t.Errorf("Failed to parse endpoint %s, %s", region.endpoint, err)
			continue
		}
		expected, actual = region.name, lc.GetLocation(*url)
		if expected != actual {
			t.Errorf("Expected GetLocation to return Write Region %s, but was %s", expected, actual)
		}
	}

	for _, region := range dbAcct.readRegions {
		url, err := url.Parse(region.endpoint)
		if err != nil {
			t.Errorf("Failed to parse endpoint %s, %s", region.endpoint, err)
			continue
		}
		expected, actual = region.name, lc.GetLocation(*url)
		if expected != actual {
			t.Errorf("Expected GetLocation to return Read Region %s, but was %s", expected, actual)
		}
	}
}
