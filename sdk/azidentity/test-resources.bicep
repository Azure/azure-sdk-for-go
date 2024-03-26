// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

@description('Kubernetes cluster admin user name.')
param adminUser string = 'azureuser'

@minLength(6)
@maxLength(23)
@description('The base resource name.')
param baseName string = resourceGroup().name

@description('Whether to deploy resources. When set to false, this file deploys nothing.')
param deployResources bool = false

param sshPubKey string = ''

@description('The location of the resource. By default, this is the same as the resource group.')
param location string = resourceGroup().location

// https://learn.microsoft.com/azure/role-based-access-control/built-in-roles
var blobContributor = subscriptionResourceId('Microsoft.Authorization/roleDefinitions', 'ba92f5b4-2d11-453d-a403-e96b0029c9fe') // Storage Blob Data Contributor

resource sa 'Microsoft.Storage/storageAccounts@2021-08-01' = if (deployResources) {
  kind: 'StorageV2'
  location: location
  name: 'sa${uniqueString(baseName)}'
  properties: {
    accessTier: 'Hot'
  }
  sku: {
    name: 'Standard_LRS'
  }
}

resource saUserAssigned 'Microsoft.Storage/storageAccounts@2021-08-01' = if (deployResources) {
  kind: 'StorageV2'
  location: location
  name: 'sa2${uniqueString(baseName)}'
  properties: {
    accessTier: 'Hot'
  }
  sku: {
    name: 'Standard_LRS'
  }
}

resource usermgdid 'Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30' = if (deployResources) {
  location: location
  name: baseName
}

resource blobRoleUserAssigned 'Microsoft.Authorization/roleAssignments@2022-04-01' = if (deployResources) {
  scope: saUserAssigned
  name: guid(resourceGroup().id, blobContributor, usermgdid.id)
  properties: {
    principalId: deployResources ? usermgdid.properties.principalId : ''
    principalType: 'ServicePrincipal'
    roleDefinitionId: blobContributor
  }
}

resource blobRoleFunc 'Microsoft.Authorization/roleAssignments@2022-04-01' = if (deployResources) {
  name: guid(resourceGroup().id, blobContributor, 'azfunc')
  properties: {
    principalId: deployResources ? azfunc.identity.principalId : ''
    roleDefinitionId: blobContributor
    principalType: 'ServicePrincipal'
  }
  scope: sa
}

resource containerRegistry 'Microsoft.ContainerRegistry/registries@2023-01-01-preview' = if (deployResources) {
  location: location
  name: uniqueString(resourceGroup().id)
  properties: {
    adminUserEnabled: true
  }
  sku: {
    name: 'Basic'
  }
}

resource farm 'Microsoft.Web/serverfarms@2021-03-01' = if (deployResources) {
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

resource azfunc 'Microsoft.Web/sites@2021-03-01' = if (deployResources) {
  identity: {
    type: 'SystemAssigned, UserAssigned'
    userAssignedIdentities: {
      '${deployResources ? usermgdid.id : ''}': {}
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
          value: deployResources ? sa.name : null
        }
        {
          name: 'AZIDENTITY_STORAGE_NAME_USER_ASSIGNED'
          value: deployResources ? saUserAssigned.name : null
        }
        {
          name: 'AZIDENTITY_USER_ASSIGNED_IDENTITY'
          value: deployResources ? usermgdid.id : null
        }
        {
          name: 'AzureWebJobsStorage'
          value: 'DefaultEndpointsProtocol=https;AccountName=${deployResources ? sa.name : ''};EndpointSuffix=${deployResources ? environment().suffixes.storage : ''};AccountKey=${deployResources ? sa.listKeys().keys[0].value : ''}'
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
          value: 'DefaultEndpointsProtocol=https;AccountName=${deployResources ? sa.name : ''};EndpointSuffix=${deployResources ? environment().suffixes.storage : ''};AccountKey=${deployResources ? sa.listKeys().keys[0].value : ''}'
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

resource aks 'Microsoft.ContainerService/managedClusters@2023-06-01' = if (deployResources) {
  name: baseName
  location: location
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    agentPoolProfiles: [
      {
        count: 1
        enableAutoScaling: false
        kubeletDiskType: 'OS'
        mode: 'System'
        name: 'agentpool'
        osDiskSizeGB: 128
        osDiskType: 'Managed'
        osSKU: 'Ubuntu'
        osType: 'Linux'
        type: 'VirtualMachineScaleSets'
        vmSize: 'Standard_D2s_v3'
      }
    ]
    dnsPrefix: 'identitytest'
    enableRBAC: true
    linuxProfile: {
      adminUsername: adminUser
      ssh: {
        publicKeys: [
          {
            keyData: sshPubKey
          }
        ]
      }
    }
    oidcIssuerProfile: {
      enabled: true
    }
    securityProfile: {
      workloadIdentity: {
        enabled: true
      }
    }
  }
}

output AZIDENTITY_ACR_LOGIN_SERVER string = deployResources ? containerRegistry.properties.loginServer : ''
output AZIDENTITY_ACR_NAME string = deployResources ? containerRegistry.name : ''
output AZIDENTITY_AKS_NAME string = deployResources ? aks.name : ''
output AZIDENTITY_FUNCTION_NAME string = deployResources ? azfunc.name : ''
output AZIDENTITY_STORAGE_NAME string = deployResources ? sa.name : ''
output AZIDENTITY_STORAGE_NAME_USER_ASSIGNED string = deployResources ? saUserAssigned.name : ''
output AZIDENTITY_USER_ASSIGNED_IDENTITY string = deployResources ? usermgdid.id : ''
output AZIDENTITY_USER_ASSIGNED_IDENTITY_CLIENT_ID string = deployResources ? usermgdid.properties.clientId : ''
output AZIDENTITY_USER_ASSIGNED_IDENTITY_NAME string = deployResources ? usermgdid.name : ''
