package backup

import (
	 original "github.com/Azure/azure-sdk-for-go/service/recoveryservices/management/2016-12-01/backup"
)

type (
	 VaultConfigsClient = original.VaultConfigsClient
	 ManagementClient = original.ManagementClient
	 EnhancedSecurityState = original.EnhancedSecurityState
	 SkuName = original.SkuName
	 StorageModelType = original.StorageModelType
	 StorageType = original.StorageType
	 StorageTypeState = original.StorageTypeState
	 TriggerType = original.TriggerType
	 VaultUpgradeState = original.VaultUpgradeState
	 Resource = original.Resource
	 Sku = original.Sku
	 StorageConfig = original.StorageConfig
	 StorageConfigProperties = original.StorageConfigProperties
	 TrackedResource = original.TrackedResource
	 UpgradeDetails = original.UpgradeDetails
	 Vault = original.Vault
	 VaultConfig = original.VaultConfig
	 VaultConfigProperties = original.VaultConfigProperties
	 VaultExtendedInfo = original.VaultExtendedInfo
	 VaultExtendedInfoResource = original.VaultExtendedInfoResource
	 VaultProperties = original.VaultProperties
	 StorageConfigsClient = original.StorageConfigsClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Disabled = original.Disabled
	 Enabled = original.Enabled
	 Invalid = original.Invalid
	 RS0 = original.RS0
	 Standard = original.Standard
	 StorageModelTypeGeoRedundant = original.StorageModelTypeGeoRedundant
	 StorageModelTypeInvalid = original.StorageModelTypeInvalid
	 StorageModelTypeLocallyRedundant = original.StorageModelTypeLocallyRedundant
	 StorageTypeGeoRedundant = original.StorageTypeGeoRedundant
	 StorageTypeInvalid = original.StorageTypeInvalid
	 StorageTypeLocallyRedundant = original.StorageTypeLocallyRedundant
	 StorageTypeStateInvalid = original.StorageTypeStateInvalid
	 StorageTypeStateLocked = original.StorageTypeStateLocked
	 StorageTypeStateUnlocked = original.StorageTypeStateUnlocked
	 ForcedUpgrade = original.ForcedUpgrade
	 UserTriggered = original.UserTriggered
	 Failed = original.Failed
	 InProgress = original.InProgress
	 Unknown = original.Unknown
	 Upgraded = original.Upgraded
)
