// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

@minLength(6)
@maxLength(23)
@description('The base resource name.')
param baseName string = resourceGroup().name

@description('The location of the resource. By default, this is the same as the resource group.')
param location string = resourceGroup().location

resource batchAccount 'Microsoft.Batch/batchAccounts@2023-11-01' = {
  identity: {
    type: 'None'
  }
  location: location
  name: 'batch${uniqueString(baseName)}'
  properties: {
    allowedAuthenticationModes: [
      'AAD'
      'SharedKey'
      'TaskAuthenticationToken'
    ]
    networkProfile: {
      accountAccess: {
        defaultAction: 'Allow'
      }
    }
    poolAllocationMode: 'BatchService'
    publicNetworkAccess: 'Enabled'
  }
}

output AZBATCH_ENDPOINT string = batchAccount.properties.accountEndpoint
