package azure

import (
	"encoding/xml"
	"errors"
	"fmt"
)

const (
	azureLocationListURL = "locations"
	invalidLocationError = "Invalid location: %s. Available locations: %s"
)

//LocationClient is used to manage operations on Azure Locations
type LocationClient struct {
	client *Client
}

//Location is used to return a handle to the Location API
func (client *Client) Location() *LocationClient {
	return &LocationClient{client: client}
}

func (self *LocationClient) ResolveLocation(location string) error {
	if len(location) == 0 {
		return fmt.Errorf(paramNotSpecifiedError, "location")
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

	return errors.New(fmt.Sprintf(invalidLocationError, location, locations))
}

func (self *LocationClient) GetLocationList() (LocationList, error) {
	locationList := LocationList{}

	response, err := self.client.sendAzureGetRequest(azureLocationListURL)
	if err != nil {
		return locationList, err
	}

	err = xml.Unmarshal(response, &locationList)
	if err != nil {
		return locationList, err
	}

	return locationList, nil
}

func (self *LocationClient) GetLocation(location string) (*Location, error) {
	if len(location) == 0 {
		return nil, fmt.Errorf(paramNotSpecifiedError, "location")
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

	return nil, errors.New(fmt.Sprintf(invalidLocationError, location, locations))
}
