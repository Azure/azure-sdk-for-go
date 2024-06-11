// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

@description('The base resource name.')
param baseName string = resourceGroup().name

@description('The resource location')
param location string = resourceGroup().location

@description('The client OID to grant access to test resources.')
param testApplicationOid string

output RESOURCE_GROUP string = resourceGroup().name
output AZURE_SUBSCRIPTION_ID string = subscription().subscriptionId

// 
// [BEGIN] Event Grid namespace
//

var namespaceName = '${baseName}-2'
var nsTopicName = 'testtopic1'
var nsSubscriptionName = 'testsubscription1'

resource ns_resource 'Microsoft.EventGrid/namespaces@2023-06-01-preview' = {
  name: namespaceName
  location: location
  sku: {
    name: 'Standard'
    capacity: 1
  }
  properties: {
    isZoneRedundant: true
    publicNetworkAccess: 'Enabled'
  }
}

resource ns_testtopic1 'Microsoft.EventGrid/namespaces/topics@2023-06-01-preview' = {
  parent: ns_resource
  name: nsTopicName
  properties: {
    publisherType: 'Custom'
    inputSchema: 'CloudEventSchemaV1_0'
    eventRetentionInDays: 1
  }
}

resource ns_testtopic1_testsubscription1 'Microsoft.EventGrid/namespaces/topics/eventSubscriptions@2023-06-01-preview' = {
  parent: ns_testtopic1
  name: nsSubscriptionName
  properties: {
    deliveryConfiguration: {
      deliveryMode: 'Queue'
      queue: {
        receiveLockDurationInSeconds: 60
        maxDeliveryCount: 10
        eventTimeToLive: 'P1D'
      }
    }
    eventDeliverySchema: 'CloudEventSchemaV1_0'
    filtersConfiguration: {
      includedEventTypes: []
    }
  }
}

// https://learn.microsoft.com/en-us/rest/api/eventgrid/controlplane-version2023-06-01-preview/namespaces/list-shared-access-keys?tabs=HTTP
#disable-next-line outputs-should-not-contain-secrets // (this is just how our test deployments work)
output EVENTGRID_KEY string = listKeys(
  resourceId('Microsoft.EventGrid/namespaces', namespaceName),
  '2023-06-01-preview'
).key1
output EVENTGRID_ENDPOINT string = 'https://${ns_resource.properties.topicsConfiguration.hostname}'

output EVENTGRID_TOPIC string = nsTopicName
output EVENTGRID_SUBSCRIPTION string = nsSubscriptionName

// [END] Event Grid namespace

//
// [BEGIN] Event Grid topics (publisher)
// 

resource egTopic 'Microsoft.EventGrid/topics@2023-06-01-preview' = {
  name: '${baseName}-eg'
  location: location
  kind: 'Azure'
  properties: {
    inputSchema: 'EventGridSchema'
  }
}

resource ceTopic 'Microsoft.EventGrid/topics@2023-06-01-preview' = {
  name: '${baseName}-ce'
  location: location
  kind: 'Azure'
  properties: {
    inputSchema: 'CloudEventSchemaV1_0'
  }
}

// https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles#eventgrid-contributor
resource egContributorRole 'Microsoft.Authorization/roleAssignments@2018-01-01-preview' = {
  name: guid('egContributorRoleId${baseName}')
  scope: resourceGroup()
  properties: {
    roleDefinitionId: subscriptionResourceId(
      'Microsoft.Authorization/roleDefinitions',
      '1e241071-0855-49ea-94dc-649edcd759de'
    )
    //    roleDefinitionId: '/subscriptions/${subscription().subscriptionId}/providers/Microsoft.Authorization/roleDefinitions/1e241071-0855-49ea-94dc-649edcd759de'
    principalId: testApplicationOid
  }
}

resource egDataContributorRole 'Microsoft.Authorization/roleAssignments@2018-01-01-preview' = {
  name: guid('egDataContributorRoleId${baseName}')
  scope: resourceGroup()
  properties: {
    roleDefinitionId: subscriptionResourceId(
      'Microsoft.Authorization/roleDefinitions',
      '1d8c3fe3-8864-474b-8749-01e3783e8157'
    )
    principalId: testApplicationOid
  }
}

resource egDataSenderRole 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid('egSenderRoleId${baseName}')
  scope: resourceGroup()
  properties: {
    roleDefinitionId: subscriptionResourceId(
      'Microsoft.Authorization/roleDefinitions',
      'd5a91429-5739-47e2-a06b-3470a27159e7'
    )
    principalId: testApplicationOid
  }
}

output EVENTGRID_TOPIC_NAME string = egTopic.name
#disable-next-line outputs-should-not-contain-secrets // (this is just how our test deployments work)
output EVENTGRID_TOPIC_KEY string = egTopic.listKeys().key1
output EVENTGRID_TOPIC_ENDPOINT string = egTopic.properties.endpoint

output EVENTGRID_CE_TOPIC_NAME string = ceTopic.name
#disable-next-line outputs-should-not-contain-secrets // (this is just how our test deployments work)
output EVENTGRID_CE_TOPIC_KEY string = ceTopic.listKeys().key1
output EVENTGRID_CE_TOPIC_ENDPOINT string = ceTopic.properties.endpoint

//
// [END] Event Grid topics (publisher)
// 
