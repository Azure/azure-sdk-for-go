package containerregistry

import (
	 original "github.com/Azure/azure-sdk-for-go/service/containerregistry/management/2017-06-01-preview/containerregistry"
)

type (
	 ManagementClient = original.ManagementClient
	 PasswordName = original.PasswordName
	 ProvisioningState = original.ProvisioningState
	 RegistryUsageUnit = original.RegistryUsageUnit
	 SkuName = original.SkuName
	 SkuTier = original.SkuTier
	 WebhookAction = original.WebhookAction
	 WebhookStatus = original.WebhookStatus
	 Actor = original.Actor
	 CallbackConfig = original.CallbackConfig
	 Event = original.Event
	 EventContent = original.EventContent
	 EventInfo = original.EventInfo
	 EventListResult = original.EventListResult
	 EventRequestMessage = original.EventRequestMessage
	 EventResponseMessage = original.EventResponseMessage
	 OperationDefinition = original.OperationDefinition
	 OperationDisplayDefinition = original.OperationDisplayDefinition
	 OperationListResult = original.OperationListResult
	 RegenerateCredentialParameters = original.RegenerateCredentialParameters
	 Registry = original.Registry
	 RegistryListCredentialsResult = original.RegistryListCredentialsResult
	 RegistryListResult = original.RegistryListResult
	 RegistryNameCheckRequest = original.RegistryNameCheckRequest
	 RegistryNameStatus = original.RegistryNameStatus
	 RegistryPassword = original.RegistryPassword
	 RegistryProperties = original.RegistryProperties
	 RegistryPropertiesUpdateParameters = original.RegistryPropertiesUpdateParameters
	 RegistryUpdateParameters = original.RegistryUpdateParameters
	 RegistryUsage = original.RegistryUsage
	 RegistryUsageListResult = original.RegistryUsageListResult
	 Replication = original.Replication
	 ReplicationListResult = original.ReplicationListResult
	 ReplicationProperties = original.ReplicationProperties
	 Request = original.Request
	 Resource = original.Resource
	 Sku = original.Sku
	 Source = original.Source
	 Status = original.Status
	 StorageAccountProperties = original.StorageAccountProperties
	 Target = original.Target
	 Webhook = original.Webhook
	 WebhookCreateParameters = original.WebhookCreateParameters
	 WebhookListResult = original.WebhookListResult
	 WebhookProperties = original.WebhookProperties
	 WebhookPropertiesCreateParameters = original.WebhookPropertiesCreateParameters
	 WebhookPropertiesUpdateParameters = original.WebhookPropertiesUpdateParameters
	 WebhookUpdateParameters = original.WebhookUpdateParameters
	 OperationsClient = original.OperationsClient
	 RegistriesClient = original.RegistriesClient
	 ReplicationsClient = original.ReplicationsClient
	 WebhooksClient = original.WebhooksClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Password = original.Password
	 Password2 = original.Password2
	 Canceled = original.Canceled
	 Creating = original.Creating
	 Deleting = original.Deleting
	 Failed = original.Failed
	 Succeeded = original.Succeeded
	 Updating = original.Updating
	 Bytes = original.Bytes
	 Count = original.Count
	 Basic = original.Basic
	 ManagedBasic = original.ManagedBasic
	 ManagedPremium = original.ManagedPremium
	 ManagedStandard = original.ManagedStandard
	 SkuTierBasic = original.SkuTierBasic
	 SkuTierManaged = original.SkuTierManaged
	 Delete = original.Delete
	 Push = original.Push
	 Disabled = original.Disabled
	 Enabled = original.Enabled
)
