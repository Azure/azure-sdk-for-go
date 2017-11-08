# ARM Examples in Go

## Configuring your environment

Install the Azure CLI tool and login using the following commands.

```bash
$ curl -L https://aka.ms/InstallAzureCli | bash
$ exec -l $SHELL
$ az login
```

Create a `Service Principal` for the Azure Active Directory using the following command.

```bash
$ az ad sp create-for-rbac
{
  "appId": "1234567-1234-1234-1234-1234567890ab",
  "displayName": "azure-cli-2017-08-18-19-25-59",
  "name": "http://azure-cli-2017-08-18-19-25-59",
  "password": "1234567-1234-1234-be18-1234567890ab",
  "tenant": "1234567-1234-1234-be18-1234567890ab"
}
```

Translate the output from the previous command to newly exported environmental variables.

**Warning**: The names of the values might change from the previous command to the environmental variables.
Follow the chart closely.

Service Principal Variable Name | Environmental variable
--- | ---
appId | AZURE_CLIENT_ID
password | AZURE_CLIENT_SECRET
tenant | AZURE_TENANT_ID

Run the following command to get you Azure subscription ID.

```bash
$ az account show --query id
"1234567-1234-1234-1234567890ab"
```

Finally export that value as an environmental variable as well.

Command| Environmental variable
--- | ---
az account show --query id | AZURE_SUBSCRIPTION_ID

**At this point you should have the following 4 environmental variables set!**

```bash
export AZURE_CLIENT_ID = "1234567-1234-1234-1234567890ab"
export AZURE_CLIENT_SECRET = "1234567-1234-1234-1234567890ab"
export AZURE_TENANT_ID = "1234567-1234-1234-1234567890ab"
export AZURE_SUBSCRIPTION_ID = "1234567-1234-1234-1234567890ab"
```

## Running an example

Navigate to the example you would like to run and modify any configuration for the example before running.

To run the example use the following command

```bash
$ go run <example>.go
```
