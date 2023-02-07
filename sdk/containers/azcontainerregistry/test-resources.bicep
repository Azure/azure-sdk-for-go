param baseName string
param location string = resourceGroup().location

resource registry 'Microsoft.ContainerRegistry/registries@2021-09-01' = {
  name: '${baseName}'
  location: location
  sku: {
    name: 'Standard'
  }
  properties: {
    publicNetworkAccess: 'Enabled'
    zoneRedundancy: 'Disabled'
  }
}

output LOGIN_SERVER string = registry.properties.loginServer
