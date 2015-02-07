package hostedservice

import (
	"encoding/xml"
)

type HostedServiceDeployment struct {
	XMLName        xml.Name `xml:"CreateHostedService"`
	Xmlns          string   `xml:"xmlns,attr"`
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
