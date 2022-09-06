# testing azidentity in Azure Functions

# prerequisite tools
- Azure CLI
- Azure Functions Core Tools 4.x

# Azure resources
This test requires instances of these Azure resources:
- Azure Key Vault
- Azure Managed Identity
  - with secrets/set and secrets/delete permission for the Key Vault
- Azure Storage account
- Azure App Service Plan
- Azure Function App x2
  - one for system-assigned identity, one for user-assigned

The rest of this section is a walkthrough of deploying these resources.

## Set environment variables
- RESOURCE_GROUP
  - name of an Azure resource group
  - must be unique in the Azure subscription
  - e.g. 'identity-test-rg'
- APP_SERVICE_PLAN
  - name of an Azure App Service Plan
- FUNCTION_APP_SYSTEM_ASSIGNED
  - name of an Azure function app
  - must be globally unique
- FUNCTION_APP_USER_ASSIGNED
  - name of an Azure function app
  - must be globally unique
- MANAGED_IDENTITY_NAME
  - name of the user-assigned identity
  - 3-128 alphanumeric characters
  - must be unique in the resource group
- STORAGE_ACCOUNT_NAME
  - 3-24 alphanumeric characters
  - must be globally unique (check it with `az storage account check-name`)
- KEY_VAULT_NAME
  - 3-24 alphanumeric characters
  - must begin with a letter
  - must be globally unique

## resource group
```sh
az group create -n $RESOURCE_GROUP --location westus2
```

## Key Vault:
```sh
az keyvault create -g $RESOURCE_GROUP -n $KEY_VAULT_NAME --sku standard
```

## Storage account
```sh
az storage account create -g $RESOURCE_GROUP -n $STORAGE_ACCOUNT_NAME
```

## App Service Plan
```sh
az appservice plan create -g $RESOURCE_GROUP -n $APP_SERVICE_PLAN -l westus2 --sku B1 --is-linux
```

## Functions App: system-assigned identity
```sh
az functionapp create -g $RESOURCE_GROUP -n $FUNCTION_APP_SYSTEM_ASSIGNED -s $STORAGE_ACCOUNT_NAME -p $APP_SERVICE_PLAN --runtime custom
```

Set app configuration:
```sh
az functionapp config appsettings set -g $RESOURCE_GROUP -n $FUNCTION_APP_SYSTEM_ASSIGNED \
  --settings AZURE_IDENTITY_TEST_VAULT_URL=$(az keyvault show -g $RESOURCE_GROUP -n $KEY_VAULT_NAME --query properties.vaultUri -o tsv)
```

Assign a system-assigned identity:
```sh
az functionapp identity assign -g $RESOURCE_GROUP -n $FUNCTION_APP_SYSTEM_ASSIGNED
```

Allow the system-assigned identity to access the Key Vault:
```sh
az keyvault set-policy -n $KEY_VAULT_NAME \
    --object-id $(az functionapp identity show -g $RESOURCE_GROUP -n $FUNCTION_APP_SYSTEM_ASSIGNED --query principalId -o tsv) \
    --secret-permissions list
```

## managed identity
Create the identity:
```sh
az identity create -n $MANAGED_IDENTITY_NAME -g $RESOURCE_GROUP -l westus2
```

Allow it to access the Key Vault:
```sh
az keyvault set-policy -n $KEY_VAULT_NAME \
    --object-id $(az identity show -g $RESOURCE_GROUP -n $MANAGED_IDENTITY_NAME --query principalId -o tsv) \
    --secret-permissions list
```

## Functions App: user-assigned identity
```sh
az functionapp create -g $RESOURCE_GROUP -n $FUNCTION_APP_USER_ASSIGNED -s $STORAGE_ACCOUNT_NAME -p $APP_SERVICE_PLAN --runtime custom
```

Set app configuration:
```sh
az functionapp config appsettings set -g $RESOURCE_GROUP -n $FUNCTION_APP_USER_ASSIGNED \
  --settings AZURE_IDENTITY_TEST_VAULT_URL=$(az keyvault show -g $RESOURCE_GROUP -n $KEY_VAULT_NAME --query properties.vaultUri -o tsv) \
   AZURE_IDENTITY_TEST_MANAGED_IDENTITY_CLIENT_ID=$(az identity show -g $RESOURCE_GROUP -n $MANAGED_IDENTITY_NAME -o tsv --query clientId)
```

Assign the identity:
```sh
az functionapp identity assign -g $RESOURCE_GROUP -n $FUNCTION_APP_USER_ASSIGNED --identities $(az identity show -g $RESOURCE_GROUP -n $MANAGED_IDENTITY_NAME -o tsv --query id)
```

# deploy test code

Clone the repository:
```
git clone https://github.com/Azure/azure-sdk-for-go
```

Build the custom handler:
```
cd azure-sdk-for-go/sdk/samples/azidentity/manual-tests/managed-identity/functions
GOOS=linux GOARCH=amd64 go build handler.go
```

## publish to Azure
For example, publishing to the app using a system assigned identity:
```sh
func azure functionapp publish $FUNCTION_APP_SYSTEM_ASSIGNED
```
Do this again for the app using a user-assigned identity (replace `FUNCTION_APP_SYSTEM_ASSIGNED` with `FUNCTION_APP_USER_ASSIGNED`).


# run tests
For each Functions App, get the tests' invocation URLs, and browse to each. For example, for the app using system-assigned identity:
```sh
func azure functionapp list-functions $FUNCTION_APP_SYSTEM_ASSIGNED --show-keys
```
Do this again for the app using a user-assigned identity (replace `FUNCTION_APP_SYSTEM_ASSIGNED` with `FUNCTION_APP_USER_ASSIGNED`).

# Delete Azure resources
```sh
az group delete -n $RESOURCE_GROUP -y --no-wait
```
