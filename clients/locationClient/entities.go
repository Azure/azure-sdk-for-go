package locationClient

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

type LocationList struct {
	XMLName   xml.Name   `xml:"Locations"`
	Xmlns     string     `xml:"xmlns,attr"`
	Locations []Location `xml:"Location"`
}

type Location struct {
	Name                    string
	DisplayName             string
	AvailableServices       []string `xml:"AvailableServices>AvailableService"`
	WebWorkerRoleSizes      []string `xml:"ComputeCapabilities>WebWorkerRoleSizes>RoleSize"`
	VirtualMachineRoleSizes []string `xml:"ComputeCapabilities>VirtualMachinesRoleSizes>RoleSize"`
}

func (locationList LocationList) String() string {
	var buf bytes.Buffer

	for _, location := range locationList.Locations {
		buf.WriteString(fmt.Sprintf("%s, ", location.Name))
	}

	return strings.Trim(buf.String(), ", ")
}
