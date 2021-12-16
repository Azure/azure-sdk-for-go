//go:build go1.9
// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/eng/tools/profileBuilder

package dtl

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/preview/devtestlabs/mgmt/2015-05-21-preview/dtl"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type CostPropertyType = original.CostPropertyType

const (
	Projected   CostPropertyType = original.Projected
	Reported    CostPropertyType = original.Reported
	Unavailable CostPropertyType = original.Unavailable
)

type CustomImageOsType = original.CustomImageOsType

const (
	Linux   CustomImageOsType = original.Linux
	None    CustomImageOsType = original.None
	Windows CustomImageOsType = original.Windows
)

type EnableStatus = original.EnableStatus

const (
	Disabled EnableStatus = original.Disabled
	Enabled  EnableStatus = original.Enabled
)

type LabStorageType = original.LabStorageType

const (
	Premium  LabStorageType = original.Premium
	Standard LabStorageType = original.Standard
)

type LinuxOsState = original.LinuxOsState

const (
	DeprovisionApplied   LinuxOsState = original.DeprovisionApplied
	DeprovisionRequested LinuxOsState = original.DeprovisionRequested
	NonDeprovisioned     LinuxOsState = original.NonDeprovisioned
)

type PolicyEvaluatorType = original.PolicyEvaluatorType

const (
	AllowedValuesPolicy PolicyEvaluatorType = original.AllowedValuesPolicy
	MaxValuePolicy      PolicyEvaluatorType = original.MaxValuePolicy
)

type PolicyFactName = original.PolicyFactName

const (
	PolicyFactNameGalleryImage                PolicyFactName = original.PolicyFactNameGalleryImage
	PolicyFactNameLabVMCount                  PolicyFactName = original.PolicyFactNameLabVMCount
	PolicyFactNameLabVMSize                   PolicyFactName = original.PolicyFactNameLabVMSize
	PolicyFactNameUserOwnedLabVMCount         PolicyFactName = original.PolicyFactNameUserOwnedLabVMCount
	PolicyFactNameUserOwnedLabVMCountInSubnet PolicyFactName = original.PolicyFactNameUserOwnedLabVMCountInSubnet
)

type PolicyStatus = original.PolicyStatus

const (
	PolicyStatusDisabled PolicyStatus = original.PolicyStatusDisabled
	PolicyStatusEnabled  PolicyStatus = original.PolicyStatusEnabled
)

type SourceControlType = original.SourceControlType

const (
	GitHub SourceControlType = original.GitHub
	VsoGit SourceControlType = original.VsoGit
)

type SubscriptionNotificationState = original.SubscriptionNotificationState

const (
	Deleted      SubscriptionNotificationState = original.Deleted
	NotDefined   SubscriptionNotificationState = original.NotDefined
	Registered   SubscriptionNotificationState = original.Registered
	Suspended    SubscriptionNotificationState = original.Suspended
	Unregistered SubscriptionNotificationState = original.Unregistered
	Warned       SubscriptionNotificationState = original.Warned
)

type TaskType = original.TaskType

const (
	LabBillingTask     TaskType = original.LabBillingTask
	LabVmsShutdownTask TaskType = original.LabVmsShutdownTask
	LabVmsStartupTask  TaskType = original.LabVmsStartupTask
)

type UsagePermissionType = original.UsagePermissionType

const (
	Allow   UsagePermissionType = original.Allow
	Default UsagePermissionType = original.Default
	Deny    UsagePermissionType = original.Deny
)

type WindowsOsState = original.WindowsOsState

const (
	NonSysprepped    WindowsOsState = original.NonSysprepped
	SysprepApplied   WindowsOsState = original.SysprepApplied
	SysprepRequested WindowsOsState = original.SysprepRequested
)

type ApplyArtifactsRequest = original.ApplyArtifactsRequest
type ArmTemplateInfo = original.ArmTemplateInfo
type Artifact = original.Artifact
type ArtifactClient = original.ArtifactClient
type ArtifactDeploymentStatusProperties = original.ArtifactDeploymentStatusProperties
type ArtifactInstallProperties = original.ArtifactInstallProperties
type ArtifactParameterProperties = original.ArtifactParameterProperties
type ArtifactProperties = original.ArtifactProperties
type ArtifactSource = original.ArtifactSource
type ArtifactSourceClient = original.ArtifactSourceClient
type ArtifactSourceProperties = original.ArtifactSourceProperties
type BaseClient = original.BaseClient
type CloudError = original.CloudError
type CloudErrorBody = original.CloudErrorBody
type Cost = original.Cost
type CostClient = original.CostClient
type CostInsight = original.CostInsight
type CostInsightClient = original.CostInsightClient
type CostInsightProperties = original.CostInsightProperties
type CostInsightRefreshDataFuture = original.CostInsightRefreshDataFuture
type CostPerDayProperties = original.CostPerDayProperties
type CostProperties = original.CostProperties
type CostRefreshDataFuture = original.CostRefreshDataFuture
type CustomImage = original.CustomImage
type CustomImageClient = original.CustomImageClient
type CustomImageCreateOrUpdateResourceFuture = original.CustomImageCreateOrUpdateResourceFuture
type CustomImageDeleteResourceFuture = original.CustomImageDeleteResourceFuture
type CustomImageProperties = original.CustomImageProperties
type CustomImagePropertiesCustom = original.CustomImagePropertiesCustom
type CustomImagePropertiesFromVM = original.CustomImagePropertiesFromVM
type DayDetails = original.DayDetails
type EvaluatePoliciesProperties = original.EvaluatePoliciesProperties
type EvaluatePoliciesRequest = original.EvaluatePoliciesRequest
type EvaluatePoliciesResponse = original.EvaluatePoliciesResponse
type Formula = original.Formula
type FormulaClient = original.FormulaClient
type FormulaCreateOrUpdateResourceFuture = original.FormulaCreateOrUpdateResourceFuture
type FormulaProperties = original.FormulaProperties
type FormulaPropertiesFromVM = original.FormulaPropertiesFromVM
type GalleryImage = original.GalleryImage
type GalleryImageClient = original.GalleryImageClient
type GalleryImageProperties = original.GalleryImageProperties
type GalleryImageReference = original.GalleryImageReference
type GenerateArmTemplateRequest = original.GenerateArmTemplateRequest
type GenerateUploadURIParameter = original.GenerateUploadURIParameter
type GenerateUploadURIResponse = original.GenerateUploadURIResponse
type HourDetails = original.HourDetails
type Lab = original.Lab
type LabClient = original.LabClient
type LabCreateEnvironmentFuture = original.LabCreateEnvironmentFuture
type LabCreateOrUpdateResourceFuture = original.LabCreateOrUpdateResourceFuture
type LabDeleteResourceFuture = original.LabDeleteResourceFuture
type LabProperties = original.LabProperties
type LabVhd = original.LabVhd
type LabVirtualMachine = original.LabVirtualMachine
type LabVirtualMachineProperties = original.LabVirtualMachineProperties
type LinuxOsInfo = original.LinuxOsInfo
type ParameterInfo = original.ParameterInfo
type Policy = original.Policy
type PolicyClient = original.PolicyClient
type PolicyProperties = original.PolicyProperties
type PolicySetClient = original.PolicySetClient
type PolicySetResult = original.PolicySetResult
type PolicyViolation = original.PolicyViolation
type ResponseWithContinuationArtifact = original.ResponseWithContinuationArtifact
type ResponseWithContinuationArtifactIterator = original.ResponseWithContinuationArtifactIterator
type ResponseWithContinuationArtifactPage = original.ResponseWithContinuationArtifactPage
type ResponseWithContinuationArtifactSource = original.ResponseWithContinuationArtifactSource
type ResponseWithContinuationArtifactSourceIterator = original.ResponseWithContinuationArtifactSourceIterator
type ResponseWithContinuationArtifactSourcePage = original.ResponseWithContinuationArtifactSourcePage
type ResponseWithContinuationCost = original.ResponseWithContinuationCost
type ResponseWithContinuationCostInsight = original.ResponseWithContinuationCostInsight
type ResponseWithContinuationCostInsightIterator = original.ResponseWithContinuationCostInsightIterator
type ResponseWithContinuationCostInsightPage = original.ResponseWithContinuationCostInsightPage
type ResponseWithContinuationCostIterator = original.ResponseWithContinuationCostIterator
type ResponseWithContinuationCostPage = original.ResponseWithContinuationCostPage
type ResponseWithContinuationCustomImage = original.ResponseWithContinuationCustomImage
type ResponseWithContinuationCustomImageIterator = original.ResponseWithContinuationCustomImageIterator
type ResponseWithContinuationCustomImagePage = original.ResponseWithContinuationCustomImagePage
type ResponseWithContinuationFormula = original.ResponseWithContinuationFormula
type ResponseWithContinuationFormulaIterator = original.ResponseWithContinuationFormulaIterator
type ResponseWithContinuationFormulaPage = original.ResponseWithContinuationFormulaPage
type ResponseWithContinuationGalleryImage = original.ResponseWithContinuationGalleryImage
type ResponseWithContinuationGalleryImageIterator = original.ResponseWithContinuationGalleryImageIterator
type ResponseWithContinuationGalleryImagePage = original.ResponseWithContinuationGalleryImagePage
type ResponseWithContinuationLab = original.ResponseWithContinuationLab
type ResponseWithContinuationLabIterator = original.ResponseWithContinuationLabIterator
type ResponseWithContinuationLabPage = original.ResponseWithContinuationLabPage
type ResponseWithContinuationLabVhd = original.ResponseWithContinuationLabVhd
type ResponseWithContinuationLabVhdIterator = original.ResponseWithContinuationLabVhdIterator
type ResponseWithContinuationLabVhdPage = original.ResponseWithContinuationLabVhdPage
type ResponseWithContinuationLabVirtualMachine = original.ResponseWithContinuationLabVirtualMachine
type ResponseWithContinuationLabVirtualMachineIterator = original.ResponseWithContinuationLabVirtualMachineIterator
type ResponseWithContinuationLabVirtualMachinePage = original.ResponseWithContinuationLabVirtualMachinePage
type ResponseWithContinuationPolicy = original.ResponseWithContinuationPolicy
type ResponseWithContinuationPolicyIterator = original.ResponseWithContinuationPolicyIterator
type ResponseWithContinuationPolicyPage = original.ResponseWithContinuationPolicyPage
type ResponseWithContinuationSchedule = original.ResponseWithContinuationSchedule
type ResponseWithContinuationScheduleIterator = original.ResponseWithContinuationScheduleIterator
type ResponseWithContinuationSchedulePage = original.ResponseWithContinuationSchedulePage
type ResponseWithContinuationVirtualNetwork = original.ResponseWithContinuationVirtualNetwork
type ResponseWithContinuationVirtualNetworkIterator = original.ResponseWithContinuationVirtualNetworkIterator
type ResponseWithContinuationVirtualNetworkPage = original.ResponseWithContinuationVirtualNetworkPage
type Schedule = original.Schedule
type ScheduleClient = original.ScheduleClient
type ScheduleCreateOrUpdateResourceFuture = original.ScheduleCreateOrUpdateResourceFuture
type ScheduleDeleteResourceFuture = original.ScheduleDeleteResourceFuture
type ScheduleExecuteFuture = original.ScheduleExecuteFuture
type ScheduleProperties = original.ScheduleProperties
type Subnet = original.Subnet
type SubnetOverride = original.SubnetOverride
type SubscriptionNotification = original.SubscriptionNotification
type SubscriptionNotificationProperties = original.SubscriptionNotificationProperties
type VMCostProperties = original.VMCostProperties
type VirtualMachineApplyArtifactsFuture = original.VirtualMachineApplyArtifactsFuture
type VirtualMachineClient = original.VirtualMachineClient
type VirtualMachineCreateOrUpdateResourceFuture = original.VirtualMachineCreateOrUpdateResourceFuture
type VirtualMachineDeleteResourceFuture = original.VirtualMachineDeleteResourceFuture
type VirtualMachineStartFuture = original.VirtualMachineStartFuture
type VirtualMachineStopFuture = original.VirtualMachineStopFuture
type VirtualNetwork = original.VirtualNetwork
type VirtualNetworkClient = original.VirtualNetworkClient
type VirtualNetworkCreateOrUpdateResourceFuture = original.VirtualNetworkCreateOrUpdateResourceFuture
type VirtualNetworkDeleteResourceFuture = original.VirtualNetworkDeleteResourceFuture
type VirtualNetworkProperties = original.VirtualNetworkProperties
type WeekDetails = original.WeekDetails
type WindowsOsInfo = original.WindowsOsInfo

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewArtifactClient(subscriptionID string) ArtifactClient {
	return original.NewArtifactClient(subscriptionID)
}
func NewArtifactClientWithBaseURI(baseURI string, subscriptionID string) ArtifactClient {
	return original.NewArtifactClientWithBaseURI(baseURI, subscriptionID)
}
func NewArtifactSourceClient(subscriptionID string) ArtifactSourceClient {
	return original.NewArtifactSourceClient(subscriptionID)
}
func NewArtifactSourceClientWithBaseURI(baseURI string, subscriptionID string) ArtifactSourceClient {
	return original.NewArtifactSourceClientWithBaseURI(baseURI, subscriptionID)
}
func NewCostClient(subscriptionID string) CostClient {
	return original.NewCostClient(subscriptionID)
}
func NewCostClientWithBaseURI(baseURI string, subscriptionID string) CostClient {
	return original.NewCostClientWithBaseURI(baseURI, subscriptionID)
}
func NewCostInsightClient(subscriptionID string) CostInsightClient {
	return original.NewCostInsightClient(subscriptionID)
}
func NewCostInsightClientWithBaseURI(baseURI string, subscriptionID string) CostInsightClient {
	return original.NewCostInsightClientWithBaseURI(baseURI, subscriptionID)
}
func NewCustomImageClient(subscriptionID string) CustomImageClient {
	return original.NewCustomImageClient(subscriptionID)
}
func NewCustomImageClientWithBaseURI(baseURI string, subscriptionID string) CustomImageClient {
	return original.NewCustomImageClientWithBaseURI(baseURI, subscriptionID)
}
func NewFormulaClient(subscriptionID string) FormulaClient {
	return original.NewFormulaClient(subscriptionID)
}
func NewFormulaClientWithBaseURI(baseURI string, subscriptionID string) FormulaClient {
	return original.NewFormulaClientWithBaseURI(baseURI, subscriptionID)
}
func NewGalleryImageClient(subscriptionID string) GalleryImageClient {
	return original.NewGalleryImageClient(subscriptionID)
}
func NewGalleryImageClientWithBaseURI(baseURI string, subscriptionID string) GalleryImageClient {
	return original.NewGalleryImageClientWithBaseURI(baseURI, subscriptionID)
}
func NewLabClient(subscriptionID string) LabClient {
	return original.NewLabClient(subscriptionID)
}
func NewLabClientWithBaseURI(baseURI string, subscriptionID string) LabClient {
	return original.NewLabClientWithBaseURI(baseURI, subscriptionID)
}
func NewPolicyClient(subscriptionID string) PolicyClient {
	return original.NewPolicyClient(subscriptionID)
}
func NewPolicyClientWithBaseURI(baseURI string, subscriptionID string) PolicyClient {
	return original.NewPolicyClientWithBaseURI(baseURI, subscriptionID)
}
func NewPolicySetClient(subscriptionID string) PolicySetClient {
	return original.NewPolicySetClient(subscriptionID)
}
func NewPolicySetClientWithBaseURI(baseURI string, subscriptionID string) PolicySetClient {
	return original.NewPolicySetClientWithBaseURI(baseURI, subscriptionID)
}
func NewResponseWithContinuationArtifactIterator(page ResponseWithContinuationArtifactPage) ResponseWithContinuationArtifactIterator {
	return original.NewResponseWithContinuationArtifactIterator(page)
}
func NewResponseWithContinuationArtifactPage(cur ResponseWithContinuationArtifact, getNextPage func(context.Context, ResponseWithContinuationArtifact) (ResponseWithContinuationArtifact, error)) ResponseWithContinuationArtifactPage {
	return original.NewResponseWithContinuationArtifactPage(cur, getNextPage)
}
func NewResponseWithContinuationArtifactSourceIterator(page ResponseWithContinuationArtifactSourcePage) ResponseWithContinuationArtifactSourceIterator {
	return original.NewResponseWithContinuationArtifactSourceIterator(page)
}
func NewResponseWithContinuationArtifactSourcePage(cur ResponseWithContinuationArtifactSource, getNextPage func(context.Context, ResponseWithContinuationArtifactSource) (ResponseWithContinuationArtifactSource, error)) ResponseWithContinuationArtifactSourcePage {
	return original.NewResponseWithContinuationArtifactSourcePage(cur, getNextPage)
}
func NewResponseWithContinuationCostInsightIterator(page ResponseWithContinuationCostInsightPage) ResponseWithContinuationCostInsightIterator {
	return original.NewResponseWithContinuationCostInsightIterator(page)
}
func NewResponseWithContinuationCostInsightPage(cur ResponseWithContinuationCostInsight, getNextPage func(context.Context, ResponseWithContinuationCostInsight) (ResponseWithContinuationCostInsight, error)) ResponseWithContinuationCostInsightPage {
	return original.NewResponseWithContinuationCostInsightPage(cur, getNextPage)
}
func NewResponseWithContinuationCostIterator(page ResponseWithContinuationCostPage) ResponseWithContinuationCostIterator {
	return original.NewResponseWithContinuationCostIterator(page)
}
func NewResponseWithContinuationCostPage(cur ResponseWithContinuationCost, getNextPage func(context.Context, ResponseWithContinuationCost) (ResponseWithContinuationCost, error)) ResponseWithContinuationCostPage {
	return original.NewResponseWithContinuationCostPage(cur, getNextPage)
}
func NewResponseWithContinuationCustomImageIterator(page ResponseWithContinuationCustomImagePage) ResponseWithContinuationCustomImageIterator {
	return original.NewResponseWithContinuationCustomImageIterator(page)
}
func NewResponseWithContinuationCustomImagePage(cur ResponseWithContinuationCustomImage, getNextPage func(context.Context, ResponseWithContinuationCustomImage) (ResponseWithContinuationCustomImage, error)) ResponseWithContinuationCustomImagePage {
	return original.NewResponseWithContinuationCustomImagePage(cur, getNextPage)
}
func NewResponseWithContinuationFormulaIterator(page ResponseWithContinuationFormulaPage) ResponseWithContinuationFormulaIterator {
	return original.NewResponseWithContinuationFormulaIterator(page)
}
func NewResponseWithContinuationFormulaPage(cur ResponseWithContinuationFormula, getNextPage func(context.Context, ResponseWithContinuationFormula) (ResponseWithContinuationFormula, error)) ResponseWithContinuationFormulaPage {
	return original.NewResponseWithContinuationFormulaPage(cur, getNextPage)
}
func NewResponseWithContinuationGalleryImageIterator(page ResponseWithContinuationGalleryImagePage) ResponseWithContinuationGalleryImageIterator {
	return original.NewResponseWithContinuationGalleryImageIterator(page)
}
func NewResponseWithContinuationGalleryImagePage(cur ResponseWithContinuationGalleryImage, getNextPage func(context.Context, ResponseWithContinuationGalleryImage) (ResponseWithContinuationGalleryImage, error)) ResponseWithContinuationGalleryImagePage {
	return original.NewResponseWithContinuationGalleryImagePage(cur, getNextPage)
}
func NewResponseWithContinuationLabIterator(page ResponseWithContinuationLabPage) ResponseWithContinuationLabIterator {
	return original.NewResponseWithContinuationLabIterator(page)
}
func NewResponseWithContinuationLabPage(cur ResponseWithContinuationLab, getNextPage func(context.Context, ResponseWithContinuationLab) (ResponseWithContinuationLab, error)) ResponseWithContinuationLabPage {
	return original.NewResponseWithContinuationLabPage(cur, getNextPage)
}
func NewResponseWithContinuationLabVhdIterator(page ResponseWithContinuationLabVhdPage) ResponseWithContinuationLabVhdIterator {
	return original.NewResponseWithContinuationLabVhdIterator(page)
}
func NewResponseWithContinuationLabVhdPage(cur ResponseWithContinuationLabVhd, getNextPage func(context.Context, ResponseWithContinuationLabVhd) (ResponseWithContinuationLabVhd, error)) ResponseWithContinuationLabVhdPage {
	return original.NewResponseWithContinuationLabVhdPage(cur, getNextPage)
}
func NewResponseWithContinuationLabVirtualMachineIterator(page ResponseWithContinuationLabVirtualMachinePage) ResponseWithContinuationLabVirtualMachineIterator {
	return original.NewResponseWithContinuationLabVirtualMachineIterator(page)
}
func NewResponseWithContinuationLabVirtualMachinePage(cur ResponseWithContinuationLabVirtualMachine, getNextPage func(context.Context, ResponseWithContinuationLabVirtualMachine) (ResponseWithContinuationLabVirtualMachine, error)) ResponseWithContinuationLabVirtualMachinePage {
	return original.NewResponseWithContinuationLabVirtualMachinePage(cur, getNextPage)
}
func NewResponseWithContinuationPolicyIterator(page ResponseWithContinuationPolicyPage) ResponseWithContinuationPolicyIterator {
	return original.NewResponseWithContinuationPolicyIterator(page)
}
func NewResponseWithContinuationPolicyPage(cur ResponseWithContinuationPolicy, getNextPage func(context.Context, ResponseWithContinuationPolicy) (ResponseWithContinuationPolicy, error)) ResponseWithContinuationPolicyPage {
	return original.NewResponseWithContinuationPolicyPage(cur, getNextPage)
}
func NewResponseWithContinuationScheduleIterator(page ResponseWithContinuationSchedulePage) ResponseWithContinuationScheduleIterator {
	return original.NewResponseWithContinuationScheduleIterator(page)
}
func NewResponseWithContinuationSchedulePage(cur ResponseWithContinuationSchedule, getNextPage func(context.Context, ResponseWithContinuationSchedule) (ResponseWithContinuationSchedule, error)) ResponseWithContinuationSchedulePage {
	return original.NewResponseWithContinuationSchedulePage(cur, getNextPage)
}
func NewResponseWithContinuationVirtualNetworkIterator(page ResponseWithContinuationVirtualNetworkPage) ResponseWithContinuationVirtualNetworkIterator {
	return original.NewResponseWithContinuationVirtualNetworkIterator(page)
}
func NewResponseWithContinuationVirtualNetworkPage(cur ResponseWithContinuationVirtualNetwork, getNextPage func(context.Context, ResponseWithContinuationVirtualNetwork) (ResponseWithContinuationVirtualNetwork, error)) ResponseWithContinuationVirtualNetworkPage {
	return original.NewResponseWithContinuationVirtualNetworkPage(cur, getNextPage)
}
func NewScheduleClient(subscriptionID string) ScheduleClient {
	return original.NewScheduleClient(subscriptionID)
}
func NewScheduleClientWithBaseURI(baseURI string, subscriptionID string) ScheduleClient {
	return original.NewScheduleClientWithBaseURI(baseURI, subscriptionID)
}
func NewVirtualMachineClient(subscriptionID string) VirtualMachineClient {
	return original.NewVirtualMachineClient(subscriptionID)
}
func NewVirtualMachineClientWithBaseURI(baseURI string, subscriptionID string) VirtualMachineClient {
	return original.NewVirtualMachineClientWithBaseURI(baseURI, subscriptionID)
}
func NewVirtualNetworkClient(subscriptionID string) VirtualNetworkClient {
	return original.NewVirtualNetworkClient(subscriptionID)
}
func NewVirtualNetworkClientWithBaseURI(baseURI string, subscriptionID string) VirtualNetworkClient {
	return original.NewVirtualNetworkClientWithBaseURI(baseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func PossibleCostPropertyTypeValues() []CostPropertyType {
	return original.PossibleCostPropertyTypeValues()
}
func PossibleCustomImageOsTypeValues() []CustomImageOsType {
	return original.PossibleCustomImageOsTypeValues()
}
func PossibleEnableStatusValues() []EnableStatus {
	return original.PossibleEnableStatusValues()
}
func PossibleLabStorageTypeValues() []LabStorageType {
	return original.PossibleLabStorageTypeValues()
}
func PossibleLinuxOsStateValues() []LinuxOsState {
	return original.PossibleLinuxOsStateValues()
}
func PossiblePolicyEvaluatorTypeValues() []PolicyEvaluatorType {
	return original.PossiblePolicyEvaluatorTypeValues()
}
func PossiblePolicyFactNameValues() []PolicyFactName {
	return original.PossiblePolicyFactNameValues()
}
func PossiblePolicyStatusValues() []PolicyStatus {
	return original.PossiblePolicyStatusValues()
}
func PossibleSourceControlTypeValues() []SourceControlType {
	return original.PossibleSourceControlTypeValues()
}
func PossibleSubscriptionNotificationStateValues() []SubscriptionNotificationState {
	return original.PossibleSubscriptionNotificationStateValues()
}
func PossibleTaskTypeValues() []TaskType {
	return original.PossibleTaskTypeValues()
}
func PossibleUsagePermissionTypeValues() []UsagePermissionType {
	return original.PossibleUsagePermissionTypeValues()
}
func PossibleWindowsOsStateValues() []WindowsOsState {
	return original.PossibleWindowsOsStateValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
