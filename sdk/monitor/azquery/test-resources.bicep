param workspaceName string
param location string
param sku string = 'pergb2018'
param retentionInDays int = 30
param resourcePermissions bool = false

resource log_analytics 'Microsoft.OperationalInsights/workspaces@2020-08-01' = {
  name: workspaceName
  location: location
  properties: {
    sku: {
      name: sku
    }
    retentionInDays: retentionInDays
    features: {
      searchVersion: 1
      legacy: 0
      enableLogAccessUsingOnlyResourcePermissions: resourcePermissions
    }
  }
}

output WORKSPACE_ID string = log_analytics.id
