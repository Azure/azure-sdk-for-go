package locationClient

import (
	"fmt"
	"strings"
	"encoding/xml"
	"errors"
	"bytes"
	azure "github.com/MSOpenTech/azure-sdk-for-go"
)


func ResolveLocation(specifiedLocation string) (error) {
	locations, err := GetLocationList()
	if err != nil {
		return err
	}

	for _, location := range locations.Locations {
		if location.Name != specifiedLocation {
			continue
		}

		return nil
	}

	var availableLocations bytes.Buffer
	for _, location := range locations.Locations {
		availableLocations.WriteString(location.Name + ", ")
	}

	return errors.New(fmt.Sprintf("Invalid location. Available locations: %s", strings.Trim(availableLocations.String(), ", ")))
}

func GetLocationList() (LocationList, error) {
	locationList := LocationList{}

	requestURL :=  "locations"
	response, err := azure.SendAzureGetRequest(requestURL)
	if err != nil {
		return locationList, err
	}

	err = xml.Unmarshal(response, &locationList)
	if err != nil {
		return locationList, err
	}

	return locationList, nil
}
