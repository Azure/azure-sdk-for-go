param baseName string
param location string = resourceGroup().location

resource registry 'Microsoft.ContainerRegistry/registries@2022-02-01-preview' = {
  name: '${baseName}'
  location: location
  sku: {
    name: 'Standard'
  }
  properties: {
    publicNetworkAccess: 'Enabled'
    zoneRedundancy: 'Disabled'
    anonymousPullEnabled: true
  }
}

resource import1 'Microsoft.ContainerRegistry/registries/${baseName}/importImage@2022-02-01-preview' = {
  source: {
    sourceImage: 'library/hello-world:latest'
    registryUri: 'docker.io'
  }
  targetTags: ['hello-world:latest']
  mode: 'Force'
}

resource import2 'Microsoft.ContainerRegistry/registries/${baseName}/importImage@2022-02-01-preview' = {
  source: {
    sourceImage: 'library/alpine:3.17.1'
    registryUri: 'docker.io'
  }
  targetTags: ['alpine:3.17.1']
  mode: 'Force'
}

resource import3 'Microsoft.ContainerRegistry/registries/${baseName}/importImage@2022-02-01-preview' = {
  source: {
    sourceImage: 'library/alpine:3.16.3'
    registryUri: 'docker.io'
  }
  targetTags: ['alpine:3.16.3']
  mode: 'Force'
}

resource import4 'Microsoft.ContainerRegistry/registries/${baseName}/importImage@2022-02-01-preview' = {
  source: {
    sourceImage: 'library/alpine:3.15.6'
    registryUri: 'docker.io'
  }
  targetTags: ['alpine:3.15.6']
  mode: 'Force'
}

resource import5 'Microsoft.ContainerRegistry/registries/${baseName}/importImage@2022-02-01-preview' = {
  source: {
    sourceImage: 'library/alpine:3.14.8'
    registryUri: 'docker.io'
  }
  targetTags: ['alpine:3.14.8']
  mode: 'Force'
}

resource import6 'Microsoft.ContainerRegistry/registries/${baseName}/importImage@2022-02-01-preview' = {
  source: {
    sourceImage: 'library/ubuntu:20.04'
    registryUri: 'docker.io'
  }
  targetTags: ['alpine:ubuntu:20.04']
  mode: 'Force'
}

resource import7 'Microsoft.ContainerRegistry/registries/${baseName}/importImage@2022-02-01-preview' = {
  source: {
    sourceImage: 'library/nginx:latest'
    registryUri: 'docker.io'
  }
  targetTags: ['alpine:nginx:latest']
  mode: 'Force'
}

output LOGIN_SERVER string = registry.properties.loginServer
