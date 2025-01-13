// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

param baseName string

@description('The principal to assign the role to. This is application object id.')
param testApplicationOid string

var storageTableContributorRoleId = resourceId(
  'Microsoft.Authorization/roleDefinitions',
  '0a9a7e1f-b9d0-4cc4-a60d-0319b160aaa3'
)
var location = resourceGroup().location
var primaryAccountName = '${baseName}prim'
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

//
// Accounts
// 

resource stgAccount 'Microsoft.Storage/storageAccounts@2019-04-01' = {
  name: primaryAccountName
  location: location
  sku: {
    name: 'Standard_RAGRS'
  }
  kind: 'StorageV2'
  properties: {
    networkAcls: networkAcls
    supportsHttpsTrafficOnly: true
    encryption: encryption
    accessTier: 'Cool'
  }
}

resource cosmosdDBAccount 'Microsoft.DocumentDB/databaseAccounts@2020-04-01' = {
  name: primaryAccountName
  location: location
  tags: {
    defaultExperience: 'Azure Table'
    'hidden-cosmos-mmspecial': ''
    CosmosAccountType: 'Non-Production'
  }
  kind: 'GlobalDocumentDB'
  properties: {
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

//
// Roles and assignments
// 

resource tableDataContributorRoleId_id 'Microsoft.Authorization/roleAssignments@2018-09-01-preview' = {
  name: guid('tableDataContributorRoleId${resourceGroup().id}')
  properties: {
    roleDefinitionId: storageTableContributorRoleId
    principalId: testApplicationOid
  }
}

// CosmosDB has it's own data plane RBAC, so we need to set that up _slightly_ differently than our 
// blob storage account, for instance.

// we're missing _one_ permission that we need for our tests (reading throughput)
resource cosmosDBThroughputRoleDef 'Microsoft.DocumentDB/databaseAccounts/sqlRoleDefinitions@2024-05-15' = {
  name: guid('cosmosDBThroughputRoleDef${resourceGroup().id}')
  parent: cosmosdDBAccount
  properties: {
    assignableScopes: [
      cosmosdDBAccount.id
    ]
    permissions: [
      {
        dataActions: [
          'Microsoft.DocumentDB/databaseAccounts/throughputSettings/read'
        ]
      }
    ]
    roleName: 'Custom role to read throughput'
    type: 'CustomRole'
  }
}

resource cosmosDBThroughputRoleAssignment 'Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments@2024-05-15' = {
  name: guid('cosmosDBThroughputRoleAssignment${resourceGroup().id}')
  parent: cosmosdDBAccount
  properties: {
    principalId: testApplicationOid
    roleDefinitionId: cosmosDBThroughputRoleDef.id
    scope: cosmosdDBAccount.id
  }
}

resource cosmosDBThroughputRole 'Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments@2024-05-15' = {
  name: guid('customRoleAssignment${resourceGroup().id}')
  parent: cosmosdDBAccount
  properties: {
    principalId: testApplicationOid
    // Built-in role 'Azure Cosmos DB Built-in Data Contributor'
    roleDefinitionId: '/${subscription().id}/resourceGroups/${resourceGroup().name}/providers/Microsoft.DocumentDB/databaseAccounts/${cosmosdDBAccount.name}/sqlRoleDefinitions/00000000-0000-0000-0000-000000000002'

    // resourceId(
    //   subscription().id,
    //   resourceGroup().name,
    //   'providers/Microsoft.DocumentDB/databaseAccounts',
    //   cosmosdDBAccount.name,
    //   'sqlRoleDefinitions',
    //   '00000000-0000-0000-0000-000000000002' // Built-in role 'Azure Cosmos DB Built-in Data Contributor'
    // )
    scope: cosmosdDBAccount.id
  }
}

output TABLES_STORAGE_ACCOUNT_NAME string = primaryAccountName
#disable-next-line outputs-should-not-contain-secrets
output TABLES_PRIMARY_STORAGE_ACCOUNT_KEY string = stgAccount.listKeys().keys[0].value
output TABLES_COSMOS_ACCOUNT_NAME string = primaryAccountName
#disable-next-line outputs-should-not-contain-secrets
output TABLES_PRIMARY_COSMOS_ACCOUNT_KEY string = cosmosdDBAccount.listKeys().primaryMasterKey
