package virtualmachineimage

import (
	"encoding/xml"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
)

//ImageClient is used to manage operations on Azure Locations
type ImageClient struct {
	client management.Client
}

type ImageList struct {
	XMLName  xml.Name  `xml:"Images"`
	Xmlns    string    `xml:"xmlns,attr"`
	OSImages []OSImage `xml:"OSImage"`
}

type OSImage struct {
	Category        string
	Label           string
	LogicalSizeInGB string
	Name            string
	OS              string
	Eula            string
	Description     string
	Location        string
}
