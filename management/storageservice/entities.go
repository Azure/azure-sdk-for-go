package storageservice

import (
	"encoding/xml"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
)

//StorageServiceClient is used to manage operations on Azure Storage
type StorageServiceClient struct {
	client management.Client
}

type StorageServiceList struct {
	XMLName         xml.Name         `xml:"StorageServices"`
	Xmlns           string           `xml:"xmlns,attr"`
	StorageServices []StorageService `xml:"StorageService"`
}

type StorageService struct {
	Url                      string
	ServiceName              string
	StorageServiceProperties StorageServiceProperties
}

type StorageServiceProperties struct {
	Description           string
	Location              string
	Label                 string
	Status                string
	Endpoints             []string `xml:"Endpoints>Endpoint"`
	GeoReplicationEnabled string
	GeoPrimaryRegion      string
}

type StorageServiceDeployment struct {
	XMLName               xml.Name `xml:"CreateStorageServiceInput"`
	Xmlns                 string   `xml:"xmlns,attr"`
	ServiceName           string
	Description           string
	Label                 string
	AffinityGroup         string `xml:",omitempty"`
	Location              string `xml:",omitempty"`
	GeoReplicationEnabled bool
	ExtendedProperties    ExtendedPropertyList
	SecondaryReadEnabled  bool
}

type ExtendedPropertyList struct {
	ExtendedProperty []ExtendedProperty
}

type ExtendedProperty struct {
	Name  string
	Value string
}
