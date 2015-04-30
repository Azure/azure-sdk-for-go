package storageservice

import (
	"encoding/xml"

	"github.com/Azure/azure-sdk-for-go/management"
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

// Receiver type for Get Storage Account Keys operation
// See https://msdn.microsoft.com/en-us/library/azure/ee460785.aspx
type GetStorageServiceKeysResponse struct {
	Url          string
	PrimaryKey   string `xml:"StorageServiceKeys>Primary"`
	SecondaryKey string `xml:"StorageServiceKeys>Secondary"`
}

type CreateStorageServiceInput struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure CreateStorageServiceInput"`
	StorageAccountCreateParameters
}

type StorageAccountCreateParameters struct {
	ServiceName        string
	Description        string `xml:",omitempty"`
	Label              string
	AffinityGroup      string `xml:",omitempty"`
	Location           string `xml:",omitempty"`
	ExtendedProperties ExtendedPropertyList
	AccountType        AccountType
}

type AccountType string

const (
	AccountTypeStandardLRS   AccountType = "Standard_LRS"
	AccountTypeStandardZRS   AccountType = "Standard_ZRS"
	AccountTypeStandardGRS   AccountType = "Standard_GRS"
	AccountTypeStandardRAGRS AccountType = "Standard_RAGRS"
	AccountTypePremiumLRS    AccountType = "Premium_LRS"
)

type ExtendedPropertyList struct {
	ExtendedProperty []ExtendedProperty
}

type ExtendedProperty struct {
	Name  string
	Value string
}

type AvailabilityResponse struct {
	XMLName xml.Name `xml:"AvailabilityResponse"`
	Xmlns   string   `xml:"xmlns,attr"`
	Result  bool
	Reason  string
}
