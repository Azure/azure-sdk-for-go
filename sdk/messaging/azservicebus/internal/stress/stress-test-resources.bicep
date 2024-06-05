// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

@description('The base resource name.')
param baseName string = resourceGroup().name

@description('The client OID to grant access to test resources.')
param testApplicationOid string

module sb '../../test-resources.bicep' = {
  name: 'test_servicebus'
  params: {
    baseName: baseName
    location: resourceGroup().location
    testApplicationOid: testApplicationOid
    enablePremium: false // we don't use/need Premium for our stress tests
  }
}

output SERVICEBUS_ENDPOINT string = sb.outputs.SERVICEBUS_ENDPOINT
