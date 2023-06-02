param ns string = 'ripark-eg1'

resource ns_resource 'Microsoft.EventGrid/namespaces@2023-06-01-preview' = {
  name: ns
  location: 'eastus'
  sku: {
    name: 'Standard'
    capacity: 1
  }
  properties: {
    topicsConfiguration: {
    }
    isZoneRedundant: true
    publicNetworkAccess: 'Enabled'
    inboundIpRules: []
    minimumTlsVersionAllowed: '1.2'
  }
}

resource ns_testtopic1 'Microsoft.EventGrid/namespaces/topics@2023-06-01-preview' = {
  parent: ns_resource
  name: 'testtopic1'
  properties: {
    publisherType: 'Custom'
    inputSchema: 'CloudEventSchemaV1_0'
    eventRetentionInDays: 1
  }
}

resource ns_testtopic1_testsubscription1 'Microsoft.EventGrid/namespaces/topics/eventSubscriptions@2023-06-01-preview' = {
  parent: ns_testtopic1
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
  dependsOn: [
    ns_resource
  ]
}
