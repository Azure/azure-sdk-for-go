# Troubleshooting Azure Identity Authentication Issues

`azidentity` credential types return errors when authentication fails or they are otherwise unable to authenticate.
This troubleshooting guide describes steps you can take to resolve such errors.

## Table of contents

- [Troubleshooting Default Azure Credential Authentication Issues](#troubleshooting-default-azure-credential-authentication-issues)
- [Troubleshooting Environment Credential Authentication Issues](#troubleshooting-environment-credential-authentication-issues)
- [Troubleshooting Service Principal Authentication Issues](#troubleshooting-service-principal-authentication-issues)
- [Troubleshooting User Password Authentication Issues](#troubleshooting-user-password-authentication-issues)
- [Troubleshooting Managed Identity Authentication Issues](#troubleshooting-managed-identity-authentication-issues)
- [Troubleshooting Azure CLI Authentication Issues](#troubleshooting-azure-cli-authentication-issues)

## Troubleshooting Default Azure Credential Issues

`DefaultAzureCredential` attempts to retrieve an access token by sequentially invoking a chain of credentials. An error from this credential signifies that every credential in the chain failed to acquire a token. To address this, follow the configuration instructions for the respective credential you intend to use. This document contains more guidance for handling errors from each credential type:

-  [Environment Credential](#troubleshooting-environment-credential-authentication-issues) |
-  [Managed Identity](#troubleshooting-managed-identity-authentication-issues) |
-  [Azure CLI](#troubleshooting-azure-cli-authentication-issues) |

## Troubleshooting Environment Credential Issues

#### Environment variables not configured

`EnvironmentCredential` supports service principal and user password authentication. `NewEnvironmentCredential()` returns an error when environment configuration is incomplete. To fix this, set the appropriate environment variables for the identity you want to authenticate:

##### Service principal with secret

| Variable Name | Value |
| --- | --- |
AZURE_CLIENT_ID | ID of an Azure Active Directory application. |
AZURE_TENANT_ID | ID of the application's Azure Active Directory tenant. |
AZURE_CLIENT_SECRET | One of the application's client secrets. |

##### Service principal with certificate

| Variable name | Value |
| --- | --- |
AZURE_CLIENT_ID | ID of an Azure Active Directory application. |
AZURE_TENANT_ID | ID of the application's Azure Active Directory tenant. |
AZURE_CLIENT_CERTIFICATE_PATH | Path to a PEM-encoded or PKCS12 certificate file including private key (without password protection). |

##### User password

| Variable name | Value |
| --- | --- |
AZURE_CLIENT_ID | ID of an Azure Active Directory application. |
AZURE_USERNAME | A username (usually an email address). |
AZURE_PASSWORD | The associated password for the given username. |

### Authentication failures

`EnvironmentCredential` supports service principal and user password authentication.
Please follow the troubleshooting guidelines below for the respective authentication method.

- [service principal](#troubleshooting-service-principal-authentication-issues)
- [user password](#troubleshooting-user-password-authentication-issues)

## Troubleshooting user password authentication issues

### Two factor authentication required

`UsernamePasswordCredential` isn't compatible with any kind of multifactor authentication.

## Troubleshooting service principal authentication issues

### Create a new service principal

Please follow the instructions [here](https://docs.microsoft.com/cli/azure/create-an-azure-service-principal-azure-cli)
to create a new service principal.

### Invalid arguments

#### Client ID

Authenticating a service principal requires a client or "application" ID and tenant ID. These are
required parameters for `NewClientSecretCredential()` and `NewClientCertificateCredential()`. If you
have already created your service principal, you can retrieve these IDs by following the instructions
[here](https://docs.microsoft.com/azure/active-directory/develop/howto-create-service-principal-portal#get-tenant-and-app-id-values-for-signing-in).

#### Client secret

A client secret, also called an application password, is a secret string that an application uses to prove its identity.
Azure Active Directory doesn't expose the values of existing secrets. If you have already created a service principal you can follow the instructions
[here](https://docs.microsoft.com/azure/active-directory/develop/howto-create-service-principal-portal#option-2-create-a-new-application-secret)
to create a new secret.

### Client certificate credential issues

`ClientCertificateCredential` authenticates with a certificate in PKCS12 (PFX) or PEM format. The certificate must first be registered for your service principal.
Follow the instructions [here](https://docs.microsoft.com/azure/active-directory/develop/howto-create-service-principal-portal#option-1-upload-a-certificate).
to register a certificate.

## Troubleshooting managed identity authentication issues

### Managed identity unavailable

[Managed identity authentication](https://docs.microsoft.com/azure/active-directory/managed-identities-azure-resources/overview) requires support from the hosting environment, which may require configuration
to expose a managed identity to your application. `azidentity` has been tested with managed
identities on these Azure services:

- [Azure Virtual Machines](https://docs.microsoft.com/azure/active-directory/managed-identities-azure-resources/qs-configure-portal-windows-vm)
- [Azure App Service](https://docs.microsoft.com/azure/app-service/overview-managed-identity)
- [Azure Kubernetes Service](https://docs.microsoft.com/azure/aks/use-managed-identity)
- [Azure Cloud Shell](https://docs.microsoft.com/azure/cloud-shell/overview)
- [Azure Arc](https://docs.microsoft.com/azure/azure-arc/servers/managed-identity-authentication)
- [Azure Service Fabric](https://docs.microsoft.com/azure/service-fabric/configure-existing-cluster-enable-managed-identity-token-service)


## Troubleshooting Azure CLI Authentication Issues

### Azure CLI Not Installed

The Azure CLI must be installed and on the application's path. Follow the instructions
[here](https://docs.microsoft.com/cli/azure/install-azure-cli) to install it and then try authenticating again.

### Azure account not logged in

`AzureCLICredential` authenticates as the identity currently logged in to Azure CLI.
You need to login to your account in Azure CLI via `az login` command. You can further read instructions to [Sign in with Azure CLI](https://docs.microsoft.com/cli/azure/authenticate-azure-cli).
Once logged in try running the credential again.
