// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ClientOptions defines the options for the Cosmos client.
type ClientOptions struct {
	azcore.ClientOptions
	// When EnableContentResponseOnWrite is false will cause the response to have a null resource. This reduces networking and CPU load by not sending the resource back over the network and serializing it on the client.
	// The default is false.
	EnableContentResponseOnWrite bool
	// PreferredRegions is a list of regions to be used when initializing the client in case the default region fails.
	PreferredRegions ClientRegions
}

// PreferredRegionsAsStringArray returns the preferred regions as a string array.
func (co *ClientOptions) PreferredRegionsAsStringArray() []string {
	regions := make([]string, len(co.PreferredRegions))
	for i, region := range co.PreferredRegions {
		regions[i] = string(region)
	}
	return regions
}

type ClientRegions []ClientRegion
type ClientRegion string

const (
	ClientRegionEastUS              ClientRegion = "East US"
	ClientRegionEastUS2             ClientRegion = "East US 2"
	ClientRegionSouthCentralUS      ClientRegion = "South Central US"
	ClientRegionWestUS              ClientRegion = "West US"
	ClientRegionWestUS2             ClientRegion = "West US 2"
	ClientRegionWestUS3             ClientRegion = "West US 3"
	ClientRegionAustraliaEast       ClientRegion = "Australia East"
	ClientRegionSoutheastAsia       ClientRegion = "Southeast Asia"
	ClientRegionNorthEurope         ClientRegion = "North Europe"
	ClientRegionSwedenCentral       ClientRegion = "Sweden Central"
	ClientRegionUKSouth             ClientRegion = "UK South"
	ClientRegionWestEurope          ClientRegion = "West Europe"
	ClientRegionCentralUS           ClientRegion = "Central US"
	ClientRegionSouthAfricaNorth    ClientRegion = "South Africa North"
	ClientRegionCentralIndia        ClientRegion = "Central India"
	ClientRegionEastAsia            ClientRegion = "East Asia"
	ClientRegionJapanEast           ClientRegion = "Japan East"
	ClientRegionKoreaCentral        ClientRegion = "Korea Central"
	ClientRegionCanadaCentral       ClientRegion = "Canada Central"
	ClientRegionFranceCentral       ClientRegion = "France Central"
	ClientRegionGermanyWestCentral  ClientRegion = "Germany West Central"
	ClientRegionNorwayEast          ClientRegion = "Norway East"
	ClientRegionPolandCentral       ClientRegion = "Poland Central"
	ClientRegionSwitzerlandNorth    ClientRegion = "Switzerland North"
	ClientRegionUAENorth            ClientRegion = "UAE North"
	ClientRegionBrazilSouth         ClientRegion = "Brazil South"
	ClientRegionCentralUSEUAP       ClientRegion = "Central US EUAP"
	ClientRegionQatarCentral        ClientRegion = "Qatar Central"
	ClientRegionCentralUSStage      ClientRegion = "Central US (Stage)"
	ClientRegionEastUSStage         ClientRegion = "East US (Stage)"
	ClientRegionEastUS2Stage        ClientRegion = "East US 2 (Stage)"
	ClientRegionNorthCentralUSStage ClientRegion = "North Central US (Stage)"
	ClientRegionSouthCentralUSStage ClientRegion = "South Central US (Stage)"
	ClientRegionWestUSStage         ClientRegion = "West US (Stage)"
	ClientRegionWestUS2Stage        ClientRegion = "West US 2 (Stage)"
	ClientRegionAsia                ClientRegion = "Asia"
	ClientRegionAsiaPacific         ClientRegion = "Asia Pacific"
	ClientRegionAustralia           ClientRegion = "Australia"
	ClientRegionBrazil              ClientRegion = "Brazil"
	ClientRegionCanada              ClientRegion = "Canada"
	ClientRegionEurope              ClientRegion = "Europe"
	ClientRegionFrance              ClientRegion = "France"
	ClientRegionGermany             ClientRegion = "Germany"
	ClientRegionGlobal              ClientRegion = "Global"
	ClientRegionIndia               ClientRegion = "India"
	ClientRegionJapan               ClientRegion = "Japan"
	ClientRegionKorea               ClientRegion = "Korea"
	ClientRegionNorway              ClientRegion = "Norway"
	ClientRegionSingapore           ClientRegion = "Singapore"
	ClientRegionSouthAfrica         ClientRegion = "South Africa"
	ClientRegionSwitzerland         ClientRegion = "Switzerland"
	ClientRegionUnitedArabEmirates  ClientRegion = "United Arab Emirates"
	ClientRegionUnitedKingdom       ClientRegion = "United Kingdom"
	ClientRegionUnitedStates        ClientRegion = "United States"
	ClientRegionUnitedStatesEUAP    ClientRegion = "United States EUAP"
	ClientRegionEastAsiaStage       ClientRegion = "East Asia (Stage)"
	ClientRegionSoutheastAsiaStage  ClientRegion = "Southeast Asia (Stage)"
	ClientRegionBrazilUS            ClientRegion = "Brazil US"
	ClientRegionEastUSSTG           ClientRegion = "East US STG"
	ClientRegionNorthCentralUS      ClientRegion = "North Central US"
	ClientRegionJioIndiaWest        ClientRegion = "Jio India West"
	ClientRegionEastUS2EUAP         ClientRegion = "East US 2 EUAP"
	ClientRegionSouthCentralUSSTG   ClientRegion = "South Central US STG"
	ClientRegionWestCentralUS       ClientRegion = "West Central US"
	ClientRegionSouthAfricaWest     ClientRegion = "South Africa West"
	ClientRegionAustraliaCentral    ClientRegion = "Australia Central"
	ClientRegionAustraliaCentral2   ClientRegion = "Australia Central 2"
	ClientRegionAustraliaSoutheast  ClientRegion = "Australia Southeast"
	ClientRegionJapanWest           ClientRegion = "Japan West"
	ClientRegionJioIndiaCentral     ClientRegion = "Jio India Central"
	ClientRegionKoreaSouth          ClientRegion = "Korea South"
	ClientRegionSouthIndia          ClientRegion = "South India"
	ClientRegionWestIndia           ClientRegion = "West India"
	ClientRegionCanadaEast          ClientRegion = "Canada East"
	ClientRegionFranceSouth         ClientRegion = "France South"
	ClientRegionGermanyNorth        ClientRegion = "Germany North"
	ClientRegionNorwayWest          ClientRegion = "Norway West"
	ClientRegionSwitzerlandWest     ClientRegion = "Switzerland West"
	ClientRegionUKWest              ClientRegion = "UK West"
	ClientRegionUAECentral          ClientRegion = "UAE Central"
	ClientRegionBrazilSoutheast     ClientRegion = "Brazil Southeast"
)
