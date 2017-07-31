package servicefabric

import (
	 original "github.com/Azure/azure-sdk-for-go/service/servicefabric/management/2016-09-01/servicefabric"
)

type (
	 ClustersClient = original.ClustersClient
	 ClusterVersionsClient = original.ClusterVersionsClient
	 ClusterState = original.ClusterState
	 DurabilityLevel = original.DurabilityLevel
	 Environment = original.Environment
	 ProvisioningState = original.ProvisioningState
	 ReliabilityLevel = original.ReliabilityLevel
	 ReliabilityLevel1 = original.ReliabilityLevel1
	 UpgradeMode = original.UpgradeMode
	 UpgradeMode1 = original.UpgradeMode1
	 X509StoreName = original.X509StoreName
	 AvailableOperationDisplay = original.AvailableOperationDisplay
	 AzureActiveDirectory = original.AzureActiveDirectory
	 CertificateDescription = original.CertificateDescription
	 ClientCertificateCommonName = original.ClientCertificateCommonName
	 ClientCertificateThumbprint = original.ClientCertificateThumbprint
	 Cluster = original.Cluster
	 ClusterCodeVersionsListResult = original.ClusterCodeVersionsListResult
	 ClusterCodeVersionsResult = original.ClusterCodeVersionsResult
	 ClusterHealthPolicy = original.ClusterHealthPolicy
	 ClusterListResult = original.ClusterListResult
	 ClusterProperties = original.ClusterProperties
	 ClusterPropertiesUpdateParameters = original.ClusterPropertiesUpdateParameters
	 ClusterUpdateParameters = original.ClusterUpdateParameters
	 ClusterUpgradeDeltaHealthPolicy = original.ClusterUpgradeDeltaHealthPolicy
	 ClusterUpgradePolicy = original.ClusterUpgradePolicy
	 ClusterVersionDetails = original.ClusterVersionDetails
	 DiagnosticsStorageAccountConfig = original.DiagnosticsStorageAccountConfig
	 EndpointRangeDescription = original.EndpointRangeDescription
	 ErrorModel = original.ErrorModel
	 ErrorModelError = original.ErrorModelError
	 NodeTypeDescription = original.NodeTypeDescription
	 OperationListResult = original.OperationListResult
	 OperationResult = original.OperationResult
	 Resource = original.Resource
	 SettingsParameterDescription = original.SettingsParameterDescription
	 SettingsSectionDescription = original.SettingsSectionDescription
	 OperationsClient = original.OperationsClient
	 ManagementClient = original.ManagementClient
)

const (
	 AutoScale = original.AutoScale
	 BaselineUpgrade = original.BaselineUpgrade
	 Deploying = original.Deploying
	 EnforcingClusterVersion = original.EnforcingClusterVersion
	 Ready = original.Ready
	 UpdatingInfrastructure = original.UpdatingInfrastructure
	 UpdatingUserCertificate = original.UpdatingUserCertificate
	 UpdatingUserConfiguration = original.UpdatingUserConfiguration
	 UpgradeServiceUnreachable = original.UpgradeServiceUnreachable
	 WaitingForNodes = original.WaitingForNodes
	 Bronze = original.Bronze
	 Gold = original.Gold
	 Silver = original.Silver
	 Linux = original.Linux
	 Windows = original.Windows
	 Canceled = original.Canceled
	 Failed = original.Failed
	 Succeeded = original.Succeeded
	 Updating = original.Updating
	 ReliabilityLevelBronze = original.ReliabilityLevelBronze
	 ReliabilityLevelGold = original.ReliabilityLevelGold
	 ReliabilityLevelSilver = original.ReliabilityLevelSilver
	 ReliabilityLevel1Bronze = original.ReliabilityLevel1Bronze
	 ReliabilityLevel1Gold = original.ReliabilityLevel1Gold
	 ReliabilityLevel1Platinum = original.ReliabilityLevel1Platinum
	 ReliabilityLevel1Silver = original.ReliabilityLevel1Silver
	 Automatic = original.Automatic
	 Manual = original.Manual
	 UpgradeMode1Automatic = original.UpgradeMode1Automatic
	 UpgradeMode1Manual = original.UpgradeMode1Manual
	 AddressBook = original.AddressBook
	 AuthRoot = original.AuthRoot
	 CertificateAuthority = original.CertificateAuthority
	 Disallowed = original.Disallowed
	 My = original.My
	 Root = original.Root
	 TrustedPeople = original.TrustedPeople
	 TrustedPublisher = original.TrustedPublisher
	 DefaultBaseURI = original.DefaultBaseURI
)
