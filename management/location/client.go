package location

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/management"
)

const (
	azureLocationListURL = "locations"
	errInvalidLocation   = "Invalid location: %s. Available locations: %s"
	errParamNotSpecified = "Parameter %s is not specified."
)

//NewClient is used to instantiate a new LocationClient from an Azure client
func NewClient(client management.Client) LocationClient {
	return LocationClient{client: client}
}

func (self LocationClient) ResolveLocation(location string) error {
	if location == "" {
		return fmt.Errorf(errParamNotSpecified, "location")
	}

	locations, err := self.GetLocationList()
	if err != nil {
		return err
	}

	for _, existingLocation := range locations.Locations {
		if existingLocation.Name != location {
			continue
		}

		return nil
	}

	return errors.New(fmt.Sprintf(errInvalidLocation, location, locations))
}

func (self LocationClient) GetLocationList() (LocationList, error) {
	locationList := LocationList{}

	response, err := self.client.SendAzureGetRequest(azureLocationListURL)
	if err != nil {
		return locationList, err
	}

	err = xml.Unmarshal(response, &locationList)
	if err != nil {
		return locationList, err
	}

	return locationList, nil
}

func (self LocationClient) GetLocation(location string) (*Location, error) {
	if location == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "location")
	}

	locations, err := self.GetLocationList()
	if err != nil {
		return nil, err
	}

	for _, existingLocation := range locations.Locations {
		if existingLocation.Name != location {
			continue
		}

		return &existingLocation, nil
	}

	return nil, errors.New(fmt.Sprintf(errInvalidLocation, location, locations))
}
