package batchmanagement

import (
	 original "github.com/Azure/azure-sdk-for-go/service/batch/management/2017-05-01/batchManagement"
)

type (
	 ManagementClient = original.ManagementClient
	 LocationClient = original.LocationClient
	 AccountKeyType = original.AccountKeyType
	 NameAvailabilityReason = original.NameAvailabilityReason
	 PackageState = original.PackageState
	 PoolAllocationMode = original.PoolAllocationMode
	 ProvisioningState = original.ProvisioningState
	 ActivateApplicationPackageParameters = original.ActivateApplicationPackageParameters
	 Application = original.Application
	 ApplicationCreateParameters = original.ApplicationCreateParameters
	 ApplicationPackage = original.ApplicationPackage
	 ApplicationUpdateParameters = original.ApplicationUpdateParameters
	 AutoStorageBaseProperties = original.AutoStorageBaseProperties
	 AutoStorageProperties = original.AutoStorageProperties
	 BatchAccount = original.BatchAccount
	 BatchAccountCreateParameters = original.BatchAccountCreateParameters
	 BatchAccountCreateProperties = original.BatchAccountCreateProperties
	 BatchAccountKeys = original.BatchAccountKeys
	 BatchAccountListResult = original.BatchAccountListResult
	 BatchAccountProperties = original.BatchAccountProperties
	 BatchAccountRegenerateKeyParameters = original.BatchAccountRegenerateKeyParameters
	 BatchAccountUpdateParameters = original.BatchAccountUpdateParameters
	 BatchAccountUpdateProperties = original.BatchAccountUpdateProperties
	 BatchLocationQuota = original.BatchLocationQuota
	 CheckNameAvailabilityParameters = original.CheckNameAvailabilityParameters
	 CheckNameAvailabilityResult = original.CheckNameAvailabilityResult
	 CloudError = original.CloudError
	 CloudErrorBody = original.CloudErrorBody
	 KeyVaultReference = original.KeyVaultReference
	 ListApplicationsResult = original.ListApplicationsResult
	 Operation = original.Operation
	 OperationDisplay = original.OperationDisplay
	 OperationListResult = original.OperationListResult
	 Resource = original.Resource
	 OperationsClient = original.OperationsClient
	 ApplicationClient = original.ApplicationClient
	 ApplicationPackageClient = original.ApplicationPackageClient
	 BatchAccountClient = original.BatchAccountClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Primary = original.Primary
	 Secondary = original.Secondary
	 AlreadyExists = original.AlreadyExists
	 Invalid = original.Invalid
	 Active = original.Active
	 Pending = original.Pending
	 Unmapped = original.Unmapped
	 BatchService = original.BatchService
	 UserSubscription = original.UserSubscription
	 ProvisioningStateCancelled = original.ProvisioningStateCancelled
	 ProvisioningStateCreating = original.ProvisioningStateCreating
	 ProvisioningStateDeleting = original.ProvisioningStateDeleting
	 ProvisioningStateFailed = original.ProvisioningStateFailed
	 ProvisioningStateInvalid = original.ProvisioningStateInvalid
	 ProvisioningStateSucceeded = original.ProvisioningStateSucceeded
)
