// Dummy parameter to handle defaults the script passes in
param testApplicationOid string = ''

resource config 'Microsoft.AppConfiguration/configurationStores@2020-07-01-preview' = {
  name: 'config-${resourceGroup().name}'
  location: resourceGroup().location
  sku: {
    name: 'Standard'
  }
}

output RESOURCE_GROUP string = resourceGroup().name
output AZURE_CLIENT_OID string = testApplicationOid
