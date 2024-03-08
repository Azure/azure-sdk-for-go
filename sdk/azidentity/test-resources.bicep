@minLength(6)
@maxLength(23)
@description('The base resource name.')
param baseName string = resourceGroup().name

@description('The location of the resource. By default, this is the same as the resource group.')
param location string = resourceGroup().location

// https://docs.microsoft.com/azure/role-based-access-control/built-in-roles
var blobContributor = subscriptionResourceId('Microsoft.Authorization/roleDefinitions', 'ba92f5b4-2d11-453d-a403-e96b0029c9fe') // Storage Blob Data Contributor

resource sa 'Microsoft.Storage/storageAccounts@2021-08-01' = {
  kind: 'StorageV2'
  location: location
  name: baseName
  properties: {
    accessTier: 'Hot'
  }
  sku: {
    name: 'Standard_LRS'
  }
}

resource saUserAssigned 'Microsoft.Storage/storageAccounts@2021-08-01' = {
  kind: 'StorageV2'
  location: location
  name: '${baseName}2'
  properties: {
    accessTier: 'Hot'
  }
  sku: {
    name: 'Standard_LRS'
  }
}

resource usermgdid 'Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30' = {
  location: location
  name: baseName
}

resource blobRoleUserAssigned 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  scope: saUserAssigned
  name: guid(resourceGroup().id, blobContributor, usermgdid.id)
  properties: {
    principalId: usermgdid.properties.principalId
    principalType: 'ServicePrincipal'
    roleDefinitionId: blobContributor
  }
}

resource blobRoleFunc 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(resourceGroup().id, blobContributor, 'azfunc')
  properties: {
    principalId: func.identity.principalId
    roleDefinitionId: blobContributor
    principalType: 'ServicePrincipal'
  }
  scope: sa
}

resource farm 'Microsoft.Web/serverfarms@2021-03-01' = {
  kind: 'app'
  location: location
  name: '${baseName}_asp'
  properties: {}
  sku: {
    capacity: 1
    family: 'B'
    name: 'B1'
    size: 'B1'
    tier: 'Basic'
  }
}

resource func 'Microsoft.Web/sites@2021-03-01' = {
  identity: {
    type: 'SystemAssigned, UserAssigned'
    userAssignedIdentities: {
      '${usermgdid.id}': {}
    }
  }
  kind: 'functionapp'
  location: location
  name: '${baseName}func'
  properties: {
    enabled: true
    httpsOnly: true
    keyVaultReferenceIdentity: 'SystemAssigned'
    serverFarmId: farm.id
    siteConfig: {
      alwaysOn: true
      appSettings: [
        {
          name: 'AZIDENTITY_STORAGE_NAME'
          value: sa.name
        }
        {
          name: 'AZIDENTITY_STORAGE_NAME_USER_ASSIGNED'
          value: saUserAssigned.name
        }
        {
          name: 'AZIDENTITY_USER_ASSIGNED_IDENTITY'
          value: usermgdid.id
        }
        {
          name: 'AzureWebJobsStorage'
          value: 'DefaultEndpointsProtocol=https;AccountName=${sa.name};EndpointSuffix=${environment().suffixes.storage};AccountKey=${sa.listKeys().keys[0].value}'
        }
        {
          name: 'FUNCTIONS_EXTENSION_VERSION'
          value: '~4'
        }
        {
          name: 'FUNCTIONS_WORKER_RUNTIME'
          value: 'custom'
        }
        {
          name: 'WEBSITE_CONTENTAZUREFILECONNECTIONSTRING'
          value: 'DefaultEndpointsProtocol=https;AccountName=${sa.name};EndpointSuffix=${environment().suffixes.storage};AccountKey=${sa.listKeys().keys[0].value}'
        }
        {
          name: 'WEBSITE_CONTENTSHARE'
          value: toLower('${baseName}-func')
        }
      ]
      http20Enabled: true
      minTlsVersion: '1.2'
    }
  }
}

output AZIDENTITY_FUNCTION_NAME string = func.name
output AZIDENTITY_STORAGE_NAME string = sa.name
output AZIDENTITY_STORAGE_NAME_USER_ASSIGNED string = saUserAssigned.name
output AZIDENTITY_USER_ASSIGNED_IDENTITY string = usermgdid.id
