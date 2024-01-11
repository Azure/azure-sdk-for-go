@description('The base resource name.')
param baseName string = resourceGroup().name

@description('The resource location')
param location string = resourceGroup().location

@description('The client OID to grant access to test resources.')
param testApplicationOid string

var storageQueueMessageSenderRoleID = 'c6a89b2d-59bc-44d0-9896-0f6e12d7b80a'
var storageQueueDataContributorRoleID = '974c5e8b-45b9-4653-ba55-5f855dd0fb88'
var storageBlobDataOwnerRoleID = 'b7e6dc6d-f1e8-4753-8033-0f276bb0955b'

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

resource egContributorRole 'Microsoft.Authorization/roleAssignments@2018-01-01-preview' = {
  name: guid('egContributorRoleId${baseName}')
  scope: resourceGroup()
  properties: {
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '1e241071-0855-49ea-94dc-649edcd759de')
    principalId: testApplicationOid
  }
}

resource egDataSenderRole 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid('egSenderRoleId${baseName}')
  scope: resourceGroup()
  properties: {
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', 'd5a91429-5739-47e2-a06b-3470a27159e7')
    principalId: testApplicationOid
  }
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2019-04-01' = {
  name: 'stgeg${baseName}'
  kind: 'StorageV2'
  location: location
  sku: {
    name: 'Standard_LRS'
  }
  properties: {
    networkAcls: {
      bypass: 'AzureServices'
      virtualNetworkRules: []
      ipRules: []
      defaultAction: 'Allow'
    }
    supportsHttpsTrafficOnly: true
    encryption: {
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
    accessTier: 'Hot'
  }
}

resource storageAccountQueueService 'Microsoft.Storage/storageAccounts/queueServices@2022-09-01' = {
  name: 'default'
  parent: storageAccount
  properties: {}
}

resource storageAccountQueue 'Microsoft.Storage/storageAccounts/queueServices/queues@2022-09-01' = {
  name: 'systemeventsdest'
  parent: storageAccountQueueService
  properties: {
    metadata: {}
  }
}

resource egSystemTopic 'Microsoft.EventGrid/systemTopics@2023-06-01-preview' = {
  name: 'egSystemTopic${baseName}'
  //location: 'eastus'
  location: location
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    source: storageAccount.id
    topicType: 'Microsoft.Storage.StorageAccounts'
  }
}

// Setting up the proper role for the Event Grid topic can send messages to our Azure storage queue.
// https://learn.microsoft.com/azure/event-grid/add-identity-roles#supported-destinations-and-azure-roles

// make it so the system topic has the rights to publish events to the storage queue.
resource storageQueueMessageSenderRole 'Microsoft.Authorization/roleAssignments@2018-01-01-preview' = {
  name: guid('storageQueueMessageSenderRole${baseName}')
  scope: resourceGroup()
  properties: {
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', storageQueueMessageSenderRoleID)
    principalId: egSystemTopic.identity.principalId
  }
}

// setup the roles needed for our test service principal to create blobs and receive messages
// from our storage queue.
resource storageQueueDataContributorRole 'Microsoft.Authorization/roleAssignments@2018-01-01-preview' = {
  name: guid('testAppSendRole${baseName}')
  scope: resourceGroup()
  properties: {
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', storageQueueDataContributorRoleID)
    principalId: testApplicationOid
  }
}

resource storageBlobDataOwnerRole 'Microsoft.Authorization/roleAssignments@2018-01-01-preview' = {
  name: guid('testAppSendRole2${baseName}')
  scope: resourceGroup()
  properties: {
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', storageBlobDataOwnerRoleID)
    principalId: testApplicationOid
  }
}

resource egSystemSubscription 'Microsoft.EventGrid/systemTopics/eventSubscriptions@2023-06-01-preview' = {
  parent: egSystemTopic
  name: 'egsyssub${baseName}'
  properties: {
    deliveryWithResourceIdentity: {
      identity: {
        type: 'SystemAssigned'
      }
      destination: {
        properties: {
          resourceId: storageAccount.id
          queueName: 'systemeventsdest'
        }
        endpointType: 'StorageQueue'
      }
    }
    filter: {
      includedEventTypes: [
        'Microsoft.Storage.BlobCreated'
        'Microsoft.Storage.BlobDeleted'
      ]
      enableAdvancedFilteringOnArrays: true
    }
    labels: []
    eventDeliverySchema: 'CloudEventSchemaV1_0'
    retryPolicy: {
      maxDeliveryAttempts: 30
      eventTimeToLiveInMinutes: 1440
    }
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

output STORAGE_ACCOUNT_BLOB string = storageAccount.properties.primaryEndpoints.blob
output STORAGE_ACCOUNT_QUEUE string = storageAccount.properties.primaryEndpoints.queue
output STORAGE_QUEUE_NAME string = storageAccountQueue.name
