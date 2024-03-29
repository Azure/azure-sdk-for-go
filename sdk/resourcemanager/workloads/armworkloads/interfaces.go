//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armworkloads

// FileShareConfigurationClassification provides polymorphic access to related types.
// Call the interface's GetFileShareConfiguration() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *CreateAndMountFileShareConfiguration, *FileShareConfiguration, *MountFileShareConfiguration, *SkipFileShareConfiguration
type FileShareConfigurationClassification interface {
	// GetFileShareConfiguration returns the FileShareConfiguration content of the underlying type.
	GetFileShareConfiguration() *FileShareConfiguration
}

// InfrastructureConfigurationClassification provides polymorphic access to related types.
// Call the interface's GetInfrastructureConfiguration() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *InfrastructureConfiguration, *SingleServerConfiguration, *ThreeTierConfiguration
type InfrastructureConfigurationClassification interface {
	// GetInfrastructureConfiguration returns the InfrastructureConfiguration content of the underlying type.
	GetInfrastructureConfiguration() *InfrastructureConfiguration
}

// OSConfigurationClassification provides polymorphic access to related types.
// Call the interface's GetOSConfiguration() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *LinuxConfiguration, *OSConfiguration, *WindowsConfiguration
type OSConfigurationClassification interface {
	// GetOSConfiguration returns the OSConfiguration content of the underlying type.
	GetOSConfiguration() *OSConfiguration
}

// ProviderSpecificPropertiesClassification provides polymorphic access to related types.
// Call the interface's GetProviderSpecificProperties() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *DB2ProviderInstanceProperties, *HanaDbProviderInstanceProperties, *MsSQLServerProviderInstanceProperties, *PrometheusHaClusterProviderInstanceProperties,
// - *PrometheusOSProviderInstanceProperties, *ProviderSpecificProperties, *SapNetWeaverProviderInstanceProperties
type ProviderSpecificPropertiesClassification interface {
	// GetProviderSpecificProperties returns the ProviderSpecificProperties content of the underlying type.
	GetProviderSpecificProperties() *ProviderSpecificProperties
}

// SAPConfigurationClassification provides polymorphic access to related types.
// Call the interface's GetSAPConfiguration() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *DeploymentConfiguration, *DeploymentWithOSConfiguration, *DiscoveryConfiguration, *SAPConfiguration
type SAPConfigurationClassification interface {
	// GetSAPConfiguration returns the SAPConfiguration content of the underlying type.
	GetSAPConfiguration() *SAPConfiguration
}

// SAPSizingRecommendationResultClassification provides polymorphic access to related types.
// Call the interface's GetSAPSizingRecommendationResult() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *SAPSizingRecommendationResult, *SingleServerRecommendationResult, *ThreeTierRecommendationResult
type SAPSizingRecommendationResultClassification interface {
	// GetSAPSizingRecommendationResult returns the SAPSizingRecommendationResult content of the underlying type.
	GetSAPSizingRecommendationResult() *SAPSizingRecommendationResult
}

// SingleServerCustomResourceNamesClassification provides polymorphic access to related types.
// Call the interface's GetSingleServerCustomResourceNames() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *SingleServerCustomResourceNames, *SingleServerFullResourceNames
type SingleServerCustomResourceNamesClassification interface {
	// GetSingleServerCustomResourceNames returns the SingleServerCustomResourceNames content of the underlying type.
	GetSingleServerCustomResourceNames() *SingleServerCustomResourceNames
}

// SoftwareConfigurationClassification provides polymorphic access to related types.
// Call the interface's GetSoftwareConfiguration() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *ExternalInstallationSoftwareConfiguration, *SAPInstallWithoutOSConfigSoftwareConfiguration, *ServiceInitiatedSoftwareConfiguration,
// - *SoftwareConfiguration
type SoftwareConfigurationClassification interface {
	// GetSoftwareConfiguration returns the SoftwareConfiguration content of the underlying type.
	GetSoftwareConfiguration() *SoftwareConfiguration
}

// ThreeTierCustomResourceNamesClassification provides polymorphic access to related types.
// Call the interface's GetThreeTierCustomResourceNames() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *ThreeTierCustomResourceNames, *ThreeTierFullResourceNames
type ThreeTierCustomResourceNamesClassification interface {
	// GetThreeTierCustomResourceNames returns the ThreeTierCustomResourceNames content of the underlying type.
	GetThreeTierCustomResourceNames() *ThreeTierCustomResourceNames
}
