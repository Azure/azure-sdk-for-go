package keyvault

import (
	 original "github.com/Azure/azure-sdk-for-go/service/keyvault/management/2016-10-01/keyvault"
)

type (
	 ManagementClient = original.ManagementClient
	 CertificatePermissions = original.CertificatePermissions
	 CreateMode = original.CreateMode
	 KeyPermissions = original.KeyPermissions
	 SecretPermissions = original.SecretPermissions
	 SkuName = original.SkuName
	 StoragePermissions = original.StoragePermissions
	 AccessPolicyEntry = original.AccessPolicyEntry
	 DeletedVault = original.DeletedVault
	 DeletedVaultListResult = original.DeletedVaultListResult
	 DeletedVaultProperties = original.DeletedVaultProperties
	 Permissions = original.Permissions
	 Resource = original.Resource
	 ResourceListResult = original.ResourceListResult
	 Sku = original.Sku
	 Vault = original.Vault
	 VaultCreateOrUpdateParameters = original.VaultCreateOrUpdateParameters
	 VaultListResult = original.VaultListResult
	 VaultProperties = original.VaultProperties
	 VaultsClient = original.VaultsClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Create = original.Create
	 Delete = original.Delete
	 Deleteissuers = original.Deleteissuers
	 Get = original.Get
	 Getissuers = original.Getissuers
	 Import = original.Import
	 List = original.List
	 Listissuers = original.Listissuers
	 Managecontacts = original.Managecontacts
	 Manageissuers = original.Manageissuers
	 Purge = original.Purge
	 Recover = original.Recover
	 Setissuers = original.Setissuers
	 Update = original.Update
	 CreateModeDefault = original.CreateModeDefault
	 CreateModeRecover = original.CreateModeRecover
	 KeyPermissionsBackup = original.KeyPermissionsBackup
	 KeyPermissionsCreate = original.KeyPermissionsCreate
	 KeyPermissionsDecrypt = original.KeyPermissionsDecrypt
	 KeyPermissionsDelete = original.KeyPermissionsDelete
	 KeyPermissionsEncrypt = original.KeyPermissionsEncrypt
	 KeyPermissionsGet = original.KeyPermissionsGet
	 KeyPermissionsImport = original.KeyPermissionsImport
	 KeyPermissionsList = original.KeyPermissionsList
	 KeyPermissionsPurge = original.KeyPermissionsPurge
	 KeyPermissionsRecover = original.KeyPermissionsRecover
	 KeyPermissionsRestore = original.KeyPermissionsRestore
	 KeyPermissionsSign = original.KeyPermissionsSign
	 KeyPermissionsUnwrapKey = original.KeyPermissionsUnwrapKey
	 KeyPermissionsUpdate = original.KeyPermissionsUpdate
	 KeyPermissionsVerify = original.KeyPermissionsVerify
	 KeyPermissionsWrapKey = original.KeyPermissionsWrapKey
	 SecretPermissionsBackup = original.SecretPermissionsBackup
	 SecretPermissionsDelete = original.SecretPermissionsDelete
	 SecretPermissionsGet = original.SecretPermissionsGet
	 SecretPermissionsList = original.SecretPermissionsList
	 SecretPermissionsPurge = original.SecretPermissionsPurge
	 SecretPermissionsRecover = original.SecretPermissionsRecover
	 SecretPermissionsRestore = original.SecretPermissionsRestore
	 SecretPermissionsSet = original.SecretPermissionsSet
	 Premium = original.Premium
	 Standard = original.Standard
	 StoragePermissionsDelete = original.StoragePermissionsDelete
	 StoragePermissionsDeletesas = original.StoragePermissionsDeletesas
	 StoragePermissionsGet = original.StoragePermissionsGet
	 StoragePermissionsGetsas = original.StoragePermissionsGetsas
	 StoragePermissionsList = original.StoragePermissionsList
	 StoragePermissionsListsas = original.StoragePermissionsListsas
	 StoragePermissionsRegeneratekey = original.StoragePermissionsRegeneratekey
	 StoragePermissionsSet = original.StoragePermissionsSet
	 StoragePermissionsSetsas = original.StoragePermissionsSetsas
	 StoragePermissionsUpdate = original.StoragePermissionsUpdate
)
