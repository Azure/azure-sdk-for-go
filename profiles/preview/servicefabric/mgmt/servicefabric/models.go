// +build go1.9

// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package servicefabric

import original "github.com/Azure/azure-sdk-for-go/services/servicefabric/mgmt/2017-07-01-preview/servicefabric"

type ClusterState = original.ClusterState

const (
	AutoScale                 ClusterState = original.AutoScale
	BaselineUpgrade           ClusterState = original.BaselineUpgrade
	Deploying                 ClusterState = original.Deploying
	EnforcingClusterVersion   ClusterState = original.EnforcingClusterVersion
	Ready                     ClusterState = original.Ready
	UpdatingInfrastructure    ClusterState = original.UpdatingInfrastructure
	UpdatingUserCertificate   ClusterState = original.UpdatingUserCertificate
	UpdatingUserConfiguration ClusterState = original.UpdatingUserConfiguration
	UpgradeServiceUnreachable ClusterState = original.UpgradeServiceUnreachable
	WaitingForNodes           ClusterState = original.WaitingForNodes
)

func PossibleClusterStateValues() []ClusterState {
	return original.PossibleClusterStateValues()
}

type DefaultMoveCost = original.DefaultMoveCost

const (
	High   DefaultMoveCost = original.High
	Low    DefaultMoveCost = original.Low
	Medium DefaultMoveCost = original.Medium
	Zero   DefaultMoveCost = original.Zero
)

func PossibleDefaultMoveCostValues() []DefaultMoveCost {
	return original.PossibleDefaultMoveCostValues()
}

type DurabilityLevel = original.DurabilityLevel

const (
	Bronze DurabilityLevel = original.Bronze
	Gold   DurabilityLevel = original.Gold
	Silver DurabilityLevel = original.Silver
)

func PossibleDurabilityLevelValues() []DurabilityLevel {
	return original.PossibleDurabilityLevelValues()
}

type Environment = original.Environment

const (
	Linux   Environment = original.Linux
	Windows Environment = original.Windows
)

func PossibleEnvironmentValues() []Environment {
	return original.PossibleEnvironmentValues()
}

type PartitionScheme = original.PartitionScheme

const (
	PartitionSchemeNamed                      PartitionScheme = original.PartitionSchemeNamed
	PartitionSchemePartitionSchemeDescription PartitionScheme = original.PartitionSchemePartitionSchemeDescription
	PartitionSchemeSingleton                  PartitionScheme = original.PartitionSchemeSingleton
	PartitionSchemeUniformInt64Range          PartitionScheme = original.PartitionSchemeUniformInt64Range
)

func PossiblePartitionSchemeValues() []PartitionScheme {
	return original.PossiblePartitionSchemeValues()
}

type ProvisioningState = original.ProvisioningState

const (
	Canceled  ProvisioningState = original.Canceled
	Failed    ProvisioningState = original.Failed
	Succeeded ProvisioningState = original.Succeeded
	Updating  ProvisioningState = original.Updating
)

func PossibleProvisioningStateValues() []ProvisioningState {
	return original.PossibleProvisioningStateValues()
}

type ReliabilityLevel = original.ReliabilityLevel

const (
	ReliabilityLevelBronze   ReliabilityLevel = original.ReliabilityLevelBronze
	ReliabilityLevelGold     ReliabilityLevel = original.ReliabilityLevelGold
	ReliabilityLevelNone     ReliabilityLevel = original.ReliabilityLevelNone
	ReliabilityLevelPlatinum ReliabilityLevel = original.ReliabilityLevelPlatinum
	ReliabilityLevelSilver   ReliabilityLevel = original.ReliabilityLevelSilver
)

func PossibleReliabilityLevelValues() []ReliabilityLevel {
	return original.PossibleReliabilityLevelValues()
}

type ReliabilityLevel1 = original.ReliabilityLevel1

const (
	ReliabilityLevel1Bronze ReliabilityLevel1 = original.ReliabilityLevel1Bronze
	ReliabilityLevel1Gold   ReliabilityLevel1 = original.ReliabilityLevel1Gold
	ReliabilityLevel1Silver ReliabilityLevel1 = original.ReliabilityLevel1Silver
)

func PossibleReliabilityLevel1Values() []ReliabilityLevel1 {
	return original.PossibleReliabilityLevel1Values()
}

type Scheme = original.Scheme

const (
	Affinity           Scheme = original.Affinity
	AlignedAffinity    Scheme = original.AlignedAffinity
	Invalid            Scheme = original.Invalid
	NonAlignedAffinity Scheme = original.NonAlignedAffinity
)

func PossibleSchemeValues() []Scheme {
	return original.PossibleSchemeValues()
}

type ServiceKind = original.ServiceKind

const (
	ServiceKindServiceProperties ServiceKind = original.ServiceKindServiceProperties
	ServiceKindStateful          ServiceKind = original.ServiceKindStateful
	ServiceKindStateless         ServiceKind = original.ServiceKindStateless
)

func PossibleServiceKindValues() []ServiceKind {
	return original.PossibleServiceKindValues()
}

type ServiceKindBasicServiceUpdateProperties = original.ServiceKindBasicServiceUpdateProperties

const (
	ServiceKindBasicServiceUpdatePropertiesServiceKindServiceUpdateProperties ServiceKindBasicServiceUpdateProperties = original.ServiceKindBasicServiceUpdatePropertiesServiceKindServiceUpdateProperties
	ServiceKindBasicServiceUpdatePropertiesServiceKindStateful                ServiceKindBasicServiceUpdateProperties = original.ServiceKindBasicServiceUpdatePropertiesServiceKindStateful
	ServiceKindBasicServiceUpdatePropertiesServiceKindStateless               ServiceKindBasicServiceUpdateProperties = original.ServiceKindBasicServiceUpdatePropertiesServiceKindStateless
)

func PossibleServiceKindBasicServiceUpdatePropertiesValues() []ServiceKindBasicServiceUpdateProperties {
	return original.PossibleServiceKindBasicServiceUpdatePropertiesValues()
}

type Type = original.Type

const (
	TypeServicePlacementPolicyDescription Type = original.TypeServicePlacementPolicyDescription
)

func PossibleTypeValues() []Type {
	return original.PossibleTypeValues()
}

type UpgradeMode = original.UpgradeMode

const (
	Automatic UpgradeMode = original.Automatic
	Manual    UpgradeMode = original.Manual
)

func PossibleUpgradeModeValues() []UpgradeMode {
	return original.PossibleUpgradeModeValues()
}

type UpgradeMode1 = original.UpgradeMode1

const (
	UpgradeMode1Automatic UpgradeMode1 = original.UpgradeMode1Automatic
	UpgradeMode1Manual    UpgradeMode1 = original.UpgradeMode1Manual
)

func PossibleUpgradeMode1Values() []UpgradeMode1 {
	return original.PossibleUpgradeMode1Values()
}

type Weight = original.Weight

const (
	WeightHigh   Weight = original.WeightHigh
	WeightLow    Weight = original.WeightLow
	WeightMedium Weight = original.WeightMedium
	WeightZero   Weight = original.WeightZero
)

func PossibleWeightValues() []Weight {
	return original.PossibleWeightValues()
}

type X509StoreName = original.X509StoreName

const (
	AddressBook          X509StoreName = original.AddressBook
	AuthRoot             X509StoreName = original.AuthRoot
	CertificateAuthority X509StoreName = original.CertificateAuthority
	Disallowed           X509StoreName = original.Disallowed
	My                   X509StoreName = original.My
	Root                 X509StoreName = original.Root
	TrustedPeople        X509StoreName = original.TrustedPeople
	TrustedPublisher     X509StoreName = original.TrustedPublisher
)

func PossibleX509StoreNameValues() []X509StoreName {
	return original.PossibleX509StoreNameValues()
}

type ApplicationDeleteFuture = original.ApplicationDeleteFuture
type ApplicationHealthPolicy = original.ApplicationHealthPolicy
type ApplicationMetricDescription = original.ApplicationMetricDescription
type ApplicationParameter = original.ApplicationParameter
type ApplicationPatchFuture = original.ApplicationPatchFuture
type ApplicationProperties = original.ApplicationProperties
type ApplicationPutFuture = original.ApplicationPutFuture
type ApplicationResource = original.ApplicationResource
type ApplicationResourceList = original.ApplicationResourceList
type ApplicationResourceUpdate = original.ApplicationResourceUpdate
type ApplicationTypeDeleteFuture = original.ApplicationTypeDeleteFuture
type ApplicationTypeProperties = original.ApplicationTypeProperties
type ApplicationTypeResource = original.ApplicationTypeResource
type ApplicationTypeResourceList = original.ApplicationTypeResourceList
type ApplicationUpdateProperties = original.ApplicationUpdateProperties
type ApplicationUpgradePolicy = original.ApplicationUpgradePolicy
type AvailableOperationDisplay = original.AvailableOperationDisplay
type AzureActiveDirectory = original.AzureActiveDirectory
type CertificateDescription = original.CertificateDescription
type ClientCertificateCommonName = original.ClientCertificateCommonName
type ClientCertificateThumbprint = original.ClientCertificateThumbprint
type Cluster = original.Cluster
type ClusterCodeVersionsListResult = original.ClusterCodeVersionsListResult
type ClusterCodeVersionsResult = original.ClusterCodeVersionsResult
type ClusterHealthPolicy = original.ClusterHealthPolicy
type ClusterListResult = original.ClusterListResult
type ClusterProperties = original.ClusterProperties
type ClusterPropertiesUpdateParameters = original.ClusterPropertiesUpdateParameters
type ClustersCreateFuture = original.ClustersCreateFuture
type ClustersUpdateFuture = original.ClustersUpdateFuture
type ClusterUpdateParameters = original.ClusterUpdateParameters
type ClusterUpgradeDeltaHealthPolicy = original.ClusterUpgradeDeltaHealthPolicy
type ClusterUpgradePolicy = original.ClusterUpgradePolicy
type ClusterVersionDetails = original.ClusterVersionDetails
type DiagnosticsStorageAccountConfig = original.DiagnosticsStorageAccountConfig
type EndpointRangeDescription = original.EndpointRangeDescription
type ErrorModel = original.ErrorModel
type NamedPartitionSchemeDescription = original.NamedPartitionSchemeDescription
type NodeTypeDescription = original.NodeTypeDescription
type OperationListResult = original.OperationListResult
type OperationListResultIterator = original.OperationListResultIterator
type OperationListResultPage = original.OperationListResultPage
type OperationResult = original.OperationResult
type BasicPartitionSchemeDescription = original.BasicPartitionSchemeDescription
type PartitionSchemeDescription = original.PartitionSchemeDescription
type ProxyResource = original.ProxyResource
type Resource = original.Resource
type RollingUpgradeMonitoringPolicy = original.RollingUpgradeMonitoringPolicy
type ServiceCorrelationDescription = original.ServiceCorrelationDescription
type ServiceDeleteFuture = original.ServiceDeleteFuture
type ServiceLoadMetricDescription = original.ServiceLoadMetricDescription
type ServicePatchFuture = original.ServicePatchFuture
type BasicServicePlacementPolicyDescription = original.BasicServicePlacementPolicyDescription
type ServicePlacementPolicyDescription = original.ServicePlacementPolicyDescription
type BasicServiceProperties = original.BasicServiceProperties
type ServiceProperties = original.ServiceProperties
type ServicePropertiesBase = original.ServicePropertiesBase
type ServicePutFuture = original.ServicePutFuture
type ServiceResource = original.ServiceResource
type ServiceResourceList = original.ServiceResourceList
type ServiceResourceUpdate = original.ServiceResourceUpdate
type ServiceTypeDeltaHealthPolicy = original.ServiceTypeDeltaHealthPolicy
type ServiceTypeHealthPolicy = original.ServiceTypeHealthPolicy
type ServiceTypeHealthPolicyMapItem = original.ServiceTypeHealthPolicyMapItem
type BasicServiceUpdateProperties = original.BasicServiceUpdateProperties
type ServiceUpdateProperties = original.ServiceUpdateProperties
type SettingsParameterDescription = original.SettingsParameterDescription
type SettingsSectionDescription = original.SettingsSectionDescription
type SingletonPartitionSchemeDescription = original.SingletonPartitionSchemeDescription
type StatefulServiceProperties = original.StatefulServiceProperties
type StatefulServiceUpdateProperties = original.StatefulServiceUpdateProperties
type StatelessServiceProperties = original.StatelessServiceProperties
type StatelessServiceUpdateProperties = original.StatelessServiceUpdateProperties
type UniformInt64RangePartitionSchemeDescription = original.UniformInt64RangePartitionSchemeDescription
type VersionDeleteFuture = original.VersionDeleteFuture
type VersionProperties = original.VersionProperties
type VersionPutFuture = original.VersionPutFuture
type VersionResource = original.VersionResource
type VersionResourceList = original.VersionResourceList

func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}

type ServiceClient = original.ServiceClient

func NewServiceClient() ServiceClient {
	return original.NewServiceClient()
}
func NewServiceClientWithBaseURI(baseURI string) ServiceClient {
	return original.NewServiceClientWithBaseURI(baseURI)
}

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient

func New() BaseClient {
	return original.New()
}
func NewWithBaseURI(baseURI string) BaseClient {
	return original.NewWithBaseURI(baseURI)
}

type ClustersClient = original.ClustersClient

func NewClustersClient() ClustersClient {
	return original.NewClustersClient()
}
func NewClustersClientWithBaseURI(baseURI string) ClustersClient {
	return original.NewClustersClientWithBaseURI(baseURI)
}

type OperationsClient = original.OperationsClient

func NewOperationsClient() OperationsClient {
	return original.NewOperationsClient()
}
func NewOperationsClientWithBaseURI(baseURI string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI)
}

type ApplicationClient = original.ApplicationClient

func NewApplicationClient() ApplicationClient {
	return original.NewApplicationClient()
}
func NewApplicationClientWithBaseURI(baseURI string) ApplicationClient {
	return original.NewApplicationClientWithBaseURI(baseURI)
}

type VersionClient = original.VersionClient

func NewVersionClient() VersionClient {
	return original.NewVersionClient()
}
func NewVersionClientWithBaseURI(baseURI string) VersionClient {
	return original.NewVersionClientWithBaseURI(baseURI)
}

type ClusterVersionsClient = original.ClusterVersionsClient

func NewClusterVersionsClient() ClusterVersionsClient {
	return original.NewClusterVersionsClient()
}
func NewClusterVersionsClientWithBaseURI(baseURI string) ClusterVersionsClient {
	return original.NewClusterVersionsClientWithBaseURI(baseURI)
}

type ApplicationTypeClient = original.ApplicationTypeClient

func NewApplicationTypeClient() ApplicationTypeClient {
	return original.NewApplicationTypeClient()
}
func NewApplicationTypeClientWithBaseURI(baseURI string) ApplicationTypeClient {
	return original.NewApplicationTypeClientWithBaseURI(baseURI)
}
