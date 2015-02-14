package hostedservice

import (
	"encoding/xml"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
)

//HostedServiceClient is used to manage operations on Azure Hosted Services
type HostedServiceClient struct {
	client management.Client
}

type CreateHostedService struct {
	XMLName        xml.Name
	Xmlns          string `xml:"xmlns,attr"`
	ServiceName    string
	Label          string
	Description    string
	Location       string
	ReverseDnsFqdn string `xml:"omitempty"`
}

type AvailabilityResponse struct {
	Xmlns  string `xml:"xmlns,attr"`
	Result bool
	Reason string
}

type HostedService struct {
	Url                               string
	ServiceName                       string
	Description                       string `xml:"HostedServiceProperties>Description"`
	AffinityGroup                     string `xml:"HostedServiceProperties>AffinityGroup"`
	Location                          string `xml:"HostedServiceProperties>Location"`
	LabelBase64                       string `xml:"HostedServiceProperties>Label"`
	Label                             string
	Status                            string `xml:"HostedServiceProperties>Status"`
	ReverseDnsFqdn                    string `xml:"HostedServiceProperties>ReverseDnsFqdn"`
	DefaultWinRmCertificateThumbprint string
}
