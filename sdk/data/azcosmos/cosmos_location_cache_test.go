package azcosmos

import (
	"fmt"
	"net/url"
	"os"
	"testing"
)

var defaultEndpt url.URL
var loc1Endpt url.URL
var loc2Endpt url.URL
var loc3Endpt url.URL
var loc4Endpt url.URL
var writeEndpts []url.URL
var readEndpts []url.URL
var endptsByLoc map[string]url.URL
var lc LocationCache

func TestMain(m *testing.M) {
	defaultEndpt, err := url.Parse("https://default.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse default endpoint URI")
		os.Exit(1)
	}
	loc1Endpt, err := url.Parse("https://location1.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse location1 endpoint URI")
		os.Exit(1)
	}
	loc2Endpt, err := url.Parse("https://location2.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse location2 endpoint URI")
		os.Exit(1)
	}
	loc3Endpt, err := url.Parse("https://location3.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse location3 endpoint URI")
		os.Exit(1)
	}
	loc4Endpt, err := url.Parse("https://location4.documents.azure.com")
	if err != nil {
		fmt.Println("Unable to parse location4 endpoint URI")
		os.Exit(1)
	}

	writeEndpts = []url.URL{*loc1Endpt, *loc2Endpt, *loc3Endpt}
	readEndpts = []url.URL{*loc1Endpt, *loc2Endpt, *loc4Endpt}
	endptsByLoc = map[string]url.URL{"location1": *loc1Endpt, "location2": *loc2Endpt, "location3": *loc3Endpt, "location4": *loc4Endpt}

	prefLocs := make([]string, 0)
	lc := NewLocationCache(prefLocs, *defaultEndpt)
	lc.enableEndptDiscovery = true
	lc.connLimit = 10

	status := m.Run()
	os.Exit(status)
}

func CreateDbAcct(useMultipleWriteLocations bool, enforceSingleMasterWriteLoc bool) AcctProperties {
	loc1 := AcctRegion{name: "location1", endpoint: loc1Endpt.String()}
	loc2 := AcctRegion{name: "location2", endpoint: loc2Endpt.String()}
	loc3 := AcctRegion{name: "location3", endpoint: loc3Endpt.String()}
	loc4 := AcctRegion{name: "location4", endpoint: loc4Endpt.String()}
	writeRegions := []AcctRegion{loc1, loc2, loc3}
	if !useMultipleWriteLocations && enforceSingleMasterWriteLoc {
		writeRegions = []AcctRegion{loc1}
	}
	readRegions := []AcctRegion{loc1, loc2, loc4}
	return AcctProperties{writeRegions: writeRegions, readRegions: readRegions, enableMultipleWriteLocations: useMultipleWriteLocations}
}

func TestGetLocation(t *testing.T) {
	dbAcct := CreateDbAcct(lc.enableMultipleWriteLocations, false)
	lc.DbAcctRead(dbAcct) // requires unit test of update
	if dbAcct.writeRegions == nil || len(dbAcct.writeRegions) == 0 {
		t.Fatal("Write Regions are empty")
	}
	expected, actual := dbAcct.writeRegions[0].name, lc.GetLocation(defaultEndpt)
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
