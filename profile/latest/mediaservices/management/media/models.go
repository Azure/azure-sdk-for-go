package media

import (
	 original "github.com/Azure/azure-sdk-for-go/service/mediaservices/management/2015-10-01/media"
)

type (
	 ManagementClient = original.ManagementClient
	 EntityNameUnavailabilityReason = original.EntityNameUnavailabilityReason
	 KeyType = original.KeyType
	 ResourceType = original.ResourceType
	 APIEndpoint = original.APIEndpoint
	 APIError = original.APIError
	 CheckNameAvailabilityInput = original.CheckNameAvailabilityInput
	 CheckNameAvailabilityOutput = original.CheckNameAvailabilityOutput
	 RegenerateKeyInput = original.RegenerateKeyInput
	 RegenerateKeyOutput = original.RegenerateKeyOutput
	 Resource = original.Resource
	 Service = original.Service
	 ServiceCollection = original.ServiceCollection
	 ServiceKeys = original.ServiceKeys
	 ServiceProperties = original.ServiceProperties
	 StorageAccount = original.StorageAccount
	 SyncStorageKeysInput = original.SyncStorageKeysInput
	 ServiceClient = original.ServiceClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 AlreadyExists = original.AlreadyExists
	 Invalid = original.Invalid
	 None = original.None
	 Primary = original.Primary
	 Secondary = original.Secondary
	 Mediaservices = original.Mediaservices
)
