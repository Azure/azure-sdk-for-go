// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

@description('Unique base name supplied by the Azure SDK test resource provisioner.')
@minLength(3)
@maxLength(44)
param baseName string

@description('Object ID of the service principal used by the live test task.')
param testApplicationOid string

@description('Azure region for the single-region live smoke account.')
param location string = resourceGroup().location

var accountName = toLower(baseName)
var databaseName = 'go-live'
var containerName = 'items'
var hierarchicalContainerName = 'items-hpk'
var dataContributorRole = '00000000-0000-0000-0000-000000000002'

resource account 'Microsoft.DocumentDB/databaseAccounts@2024-05-15' = {
  name: accountName
  location: location
  kind: 'GlobalDocumentDB'
  properties: {
    databaseAccountOfferType: 'Standard'
    publicNetworkAccess: 'Enabled'
    // Entra ID only: reject account-key auth so live tests cannot regress to keys.
    disableLocalAuth: true
    enableMultipleWriteLocations: false
    enableAutomaticFailover: false
    consistencyPolicy: {
      defaultConsistencyLevel: 'Session'
    }
    locations: [
      {
        locationName: location
        failoverPriority: 0
        isZoneRedundant: false
      }
    ]
  }
}

resource database 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases@2024-05-15' = {
  parent: account
  name: databaseName
  properties: {
    resource: {
      id: databaseName
    }
  }
}

resource container 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers@2024-05-15' = {
  parent: database
  name: containerName
  properties: {
    resource: {
      id: containerName
      partitionKey: {
        paths: [
          '/pk'
        ]
        kind: 'Hash'
        version: 2
      }
      // Fallback for interrupted local runs; CI also removes the owned resource group.
      defaultTtl: 86400
    }
    options: {
      throughput: 400
    }
  }
}

resource role 'Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments@2024-05-15' = {
  parent: account
  name: guid(account.id, testApplicationOid, database.id)
  properties: {
    principalId: testApplicationOid
    roleDefinitionId: '${account.id}/sqlRoleDefinitions/${dataContributorRole}'
    scope: '${account.id}/dbs/${databaseName}'
  }
}

// A separate, minimum-throughput container keeps prefix coverage bounded. The test
// reads routing metadata and fails unless this still has exactly one physical range.
resource hierarchicalContainer 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers@2024-05-15' = {
  parent: database
  name: hierarchicalContainerName
  properties: {
    resource: {
      id: hierarchicalContainerName
      partitionKey: {
        paths: [
          '/tenant'
          '/user'
        ]
        kind: 'MultiHash'
        version: 2
      }
      defaultTtl: 86400
    }
    options: {
      throughput: 400
    }
  }
}

// No account keys are emitted or supplied to the test process.
output AZCOSMOS_ENDPOINT string = account.properties.documentEndpoint
output AZCOSMOS_DATABASE string = database.name
output AZCOSMOS_CONTAINER string = container.name
output AZCOSMOS_HPK_CONTAINER string = hierarchicalContainer.name
output AZCOSMOS_REGION string = location
