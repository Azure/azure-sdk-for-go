@description('The base resource name.')
param baseName string = resourceGroup().name

@description('The resource location')
param location string = resourceGroup().location

@description('The client OID to grant access to test resources.')
param testApplicationOid string

#disable-next-line no-hardcoded-location // resources for the Event Grid namespaces aren't available in all locations
var egnsLocation = 'eastus'

resource eventGridNamespace 'Microsoft.EventGrid/namespaces@2023-06-01-preview' = {
  name: 'egns${baseName}'
  location: egnsLocation
  sku: {
    name: 'Standard'
    capacity: 1
  }
  properties: {
    isZoneRedundant: true
    publicNetworkAccess: 'Enabled'
  }
}

resource eventGridNamespaceTopic 'Microsoft.EventGrid/namespaces/topics@2023-06-01-preview' = {
  parent: eventGridNamespace
  name: 'testtopic1'
  properties: {
    publisherType: 'Custom'
    inputSchema: 'CloudEventSchemaV1_0'
    eventRetentionInDays: 1
  }
}

resource eventGridNamespaceSubscription 'Microsoft.EventGrid/namespaces/topics/eventSubscriptions@2023-06-01-preview' = {
  parent: eventGridNamespaceTopic
  name: 'testsubscription1'
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

module egBasic './test-resources-eventgrid-basic.bicep' = {
  name: 'egBasic'
  params: {
    baseName: baseName
    location: location
    testApplicationOid: testApplicationOid
  }
}

output RESOURCE_GROUP string = resourceGroup().name
output AZURE_SUBSCRIPTION_ID string = subscription().subscriptionId

// from Event Grid namespaces
#disable-next-line outputs-should-not-contain-secrets // (this is just how our test deployments work)
output EVENTGRID_KEY string = eventGridNamespace.listKeys().key1
output EVENTGRID_ENDPOINT string = 'https://${eventGridNamespace.properties.topicsConfiguration.hostname}'
output EVENTGRID_TOPIC string = eventGridNamespaceTopic.name
output EVENTGRID_SUBSCRIPTION string = eventGridNamespaceSubscription.name

// from our Event Grid Basic SKU
output EVENTGRID_TOPIC_NAME string = egBasic.outputs.EVENTGRID_TOPIC_NAME
#disable-next-line outputs-should-not-contain-secrets // (this is just how our test deployments work)
output EVENTGRID_TOPIC_KEY string = egBasic.outputs.EVENTGRID_TOPIC_KEY
output EVENTGRID_TOPIC_ENDPOINT string = egBasic.outputs.EVENTGRID_TOPIC_ENDPOINT

output EVENTGRID_CE_TOPIC_NAME string = egBasic.outputs.EVENTGRID_CE_TOPIC_NAME
#disable-next-line outputs-should-not-contain-secrets // (this is just how our test deployments work)
output EVENTGRID_CE_TOPIC_KEY string = egBasic.outputs.EVENTGRID_CE_TOPIC_KEY
output EVENTGRID_CE_TOPIC_ENDPOINT string = egBasic.outputs.EVENTGRID_CE_TOPIC_ENDPOINT

output STORAGE_ACCOUNT_BLOB string = egBasic.outputs.STORAGE_ACCOUNT_BLOB
output STORAGE_ACCOUNT_QUEUE string = egBasic.outputs.STORAGE_ACCOUNT_QUEUE
output STORAGE_QUEUE_NAME string = egBasic.outputs.STORAGE_QUEUE_NAME
