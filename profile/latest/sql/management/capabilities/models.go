package capabilities

import (
	 original "github.com/Azure/azure-sdk-for-go/service/sql/management/2014-04-01/capabilities"
)

type (
	 GroupClient = original.GroupClient
	 ManagementClient = original.ManagementClient
	 CapabilityStatus = original.CapabilityStatus
	 MaxSizeUnits = original.MaxSizeUnits
	 PerformanceLevelUnit = original.PerformanceLevelUnit
	 EditionCapability = original.EditionCapability
	 ElasticPoolDtuCapability = original.ElasticPoolDtuCapability
	 ElasticPoolEditionCapability = original.ElasticPoolEditionCapability
	 ElasticPoolPerDatabaseMaxDtuCapability = original.ElasticPoolPerDatabaseMaxDtuCapability
	 ElasticPoolPerDatabaseMinDtuCapability = original.ElasticPoolPerDatabaseMinDtuCapability
	 LocationCapabilities = original.LocationCapabilities
	 MaxSizeCapability = original.MaxSizeCapability
	 PerformanceLevel = original.PerformanceLevel
	 ServerVersionCapability = original.ServerVersionCapability
	 ServiceObjectiveCapability = original.ServiceObjectiveCapability
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Available = original.Available
	 Default = original.Default
	 Disabled = original.Disabled
	 Visible = original.Visible
	 Gigabytes = original.Gigabytes
	 Megabytes = original.Megabytes
	 Petabytes = original.Petabytes
	 Terabytes = original.Terabytes
	 DTU = original.DTU
)
