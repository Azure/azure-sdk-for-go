package disk

import (
	 original "github.com/Azure/azure-sdk-for-go/service/compute/management/2017-03-30/disk"
)

type (
	 DisksClient = original.DisksClient
	 AccessLevel = original.AccessLevel
	 CreateOption = original.CreateOption
	 OperatingSystemTypes = original.OperatingSystemTypes
	 StorageAccountTypes = original.StorageAccountTypes
	 AccessURI = original.AccessURI
	 AccessURIOutput = original.AccessURIOutput
	 AccessURIRaw = original.AccessURIRaw
	 APIError = original.APIError
	 APIErrorBase = original.APIErrorBase
	 CreationData = original.CreationData
	 EncryptionSettings = original.EncryptionSettings
	 GrantAccessData = original.GrantAccessData
	 ImageDiskReference = original.ImageDiskReference
	 InnerError = original.InnerError
	 KeyVaultAndKeyReference = original.KeyVaultAndKeyReference
	 KeyVaultAndSecretReference = original.KeyVaultAndSecretReference
	 ListType = original.ListType
	 Model = original.Model
	 OperationStatusResponse = original.OperationStatusResponse
	 Properties = original.Properties
	 Resource = original.Resource
	 ResourceUpdate = original.ResourceUpdate
	 Sku = original.Sku
	 Snapshot = original.Snapshot
	 SnapshotList = original.SnapshotList
	 SnapshotUpdate = original.SnapshotUpdate
	 SourceVault = original.SourceVault
	 UpdateProperties = original.UpdateProperties
	 UpdateType = original.UpdateType
	 SnapshotsClient = original.SnapshotsClient
	 ManagementClient = original.ManagementClient
)

const (
	 None = original.None
	 Read = original.Read
	 Attach = original.Attach
	 Copy = original.Copy
	 Empty = original.Empty
	 FromImage = original.FromImage
	 Import = original.Import
	 Linux = original.Linux
	 Windows = original.Windows
	 PremiumLRS = original.PremiumLRS
	 StandardLRS = original.StandardLRS
	 DefaultBaseURI = original.DefaultBaseURI
)
