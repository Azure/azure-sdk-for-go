param baseName string

@description('The principal to assign the role to. This is application object id.')
param testApplicationOid string

var mgmtApiVersion = '2019-04-01'
var authorizationApiVersion = '2018-09-01-preview'
var blobDataContributorRoleId = '/subscriptions/${subscription().subscriptionId}/providers/Microsoft.Authorization/roleDefinitions/0a9a7e1f-b9d0-4cc4-a60d-0319b160aaa3'
var location = resourceGroup().location
var primaryAccountName_var = '${baseName}prim'
var encryption = {
  services: {
    file: {
      enabled: true
    }
    blob: {
      enabled: true
    }
  }
  keySource: 'Microsoft.Storage'
}
var networkAcls = {
  bypass: 'AzureServices'
  virtualNetworkRules: []
  ipRules: []
  defaultAction: 'Allow'
}

resource tableDataContributorRoleId_id 'Microsoft.Authorization/roleAssignments@[variables(\'authorizationApiVersion\')]' = {
  name: guid('tableDataContributorRoleId${resourceGroup().id}')
  properties: {
    roleDefinitionId: blobDataContributorRoleId
    principalId: testApplicationOid
  }
}

resource primaryAccountName 'Microsoft.Storage/storageAccounts@[variables(\'mgmtApiVersion\')]' = {
  name: primaryAccountName_var
  location: location
  sku: {
    name: 'Standard_RAGRS'
    tier: 'Standard'
  }
  kind: 'StorageV2'
  properties: {
    networkAcls: networkAcls
    supportsHttpsTrafficOnly: true
    encryption: encryption
    accessTier: 'Cool'
  }
}

resource Microsoft_DocumentDB_databaseAccounts_primaryAccountName 'Microsoft.DocumentDB/databaseAccounts@2020-04-01' = {
  name: primaryAccountName_var
  location: location
  tags: {
    defaultExperience: 'Azure Table'
    'hidden-cosmos-mmspecial': ''
    CosmosAccountType: 'Non-Production'
  }
  kind: 'GlobalDocumentDB'
  properties: {
    publicNetworkAccess: 'Enabled'
    enableAutomaticFailover: false
    enableMultipleWriteLocations: false
    isVirtualNetworkFilterEnabled: false
    virtualNetworkRules: []
    disableKeyBasedMetadataWriteAccess: false
    enableFreeTier: false
    enableAnalyticalStorage: false
    databaseAccountOfferType: 'Standard'
    consistencyPolicy: {
      defaultConsistencyLevel: 'BoundedStaleness'
      maxIntervalInSeconds: 86400
      maxStalenessPrefix: 1000000
    }
    locations: [
      {
        locationName: location
        provisioningState: 'Succeeded'
        failoverPriority: 0
        isZoneRedundant: false
      }
    ]
    capabilities: [
      {
        name: 'EnableTable'
      }
    ]
    ipRules: []
  }
}

output TABLES_STORAGE_ACCOUNT_NAME string = primaryAccountName_var
output TABLES_PRIMARY_STORAGE_ACCOUNT_KEY string = listKeys(primaryAccountName.id, mgmtApiVersion).keys[0].value
output TABLES_COSMOS_ACCOUNT_NAME string = primaryAccountName_var
output TABLES_PRIMARY_COSMOS_ACCOUNT_KEY string = listKeys(Microsoft_DocumentDB_databaseAccounts_primaryAccountName.id, '2020-04-01').primaryMasterKey