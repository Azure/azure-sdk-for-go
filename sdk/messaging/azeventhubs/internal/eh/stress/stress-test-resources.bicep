// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

@description('The base resource name.')
param baseName string = resourceGroup().name

@description('The resource location')
param location string = resourceGroup().location

module eh '../../../test-resources.bicep' = {
  name: 'test_eventhub'
  params: {
    baseName: baseName
    location: location
    tenantIsTME: true
    partitions: 32
  }
}

output EVENTHUB_NAMESPACE string = eh.outputs.EVENTHUB_NAMESPACE
output EVENTHUB_NAME_STRESS string = eh.outputs.EVENTHUB_NAME
output CHECKPOINTSTORE_STORAGE_ENDPOINT string = eh.outputs.CHECKPOINTSTORE_STORAGE_ENDPOINT
