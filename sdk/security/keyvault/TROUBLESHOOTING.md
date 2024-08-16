# Troubleshoot Azure Key Vault Client Module Issues

The Azure Key Vault SDKs for Go use a common HTTP pipeline and authentication to create, update, and delete secrets,
keys, and certificates in Key Vault and Managed HSM. This troubleshooting guide contains steps for diagnosing issues
common to these SDKs.

## Table of Contents

* [Authentication Errors](#authentication-errors)
  * [HTTP 401 Errors](#http-401-errors)
    * [Frequent HTTP 401 Errors in Logs](#frequent-http-401-errors-in-logs)
  * [HTTP 403 Errors](#http-403-errors)
    * [Operation Not Permitted](#operation-not-permitted)
    * [Access Denied to First Party Service](#access-denied-to-first-party-service)
  * [Other Authentication Issues](#other-authentication-issues)
    * [Incorrect Challenge Resource](#incorrect-challenge-resource)
* [Other Service Errors](#other-service-errors)
  * [HTTP 429: Too many requests](#http-429-too-many-requests)

## Authentication Errors

### HTTP 401 Errors

HTTP 401 errors may indicate authentication problems.

#### Frequent HTTP 401 Errors in Logs

Most often, this is expected. A Key Vault client sends its first request without authorization to discover authentication parameters. This can cause a 401 response to appear in logs without a corresponding error.

### HTTP 403 Errors

HTTP 403 errors indicate the user isn't authorized to perform a specific operation in Key Vault or Managed HSM.

#### Operation Not Permitted

You may see an error similar to:

```text
--------------------------------------------------------------------------------
RESPONSE 403: 403 Forbidden
ERROR CODE: Forbidden
--------------------------------------------------------------------------------
{
  "error": {
    "code": "Forbidden",
    "message": "Operation decrypt is not permitted on this key.",
    "innererror": {
      "code": "KeyOperationForbidden"
    }
  }
}
```

The operation and inner `code` may vary, but the rest of the text will indicate which operation isn't permitted.
This error indicates that the authenticated application or user doesn't have permission to perform that operation.

1. Check that the application or user has the appropriate permission:
   * [Access policies](https://learn.microsoft.com/azure/key-vault/general/assign-access-policy) (Key Vault)
   * [Role-Based Access Control (RBAC)](https://learn.microsoft.com/azure/key-vault/general/rbac-guide) (Key Vault and Managed HSM)
2. If the appropriate permission is assigned to your application or user, make sure you are authenticating that application or user.
   If using the [DefaultAzureCredential], a different credential might've been used than one you expected.
   [Enable logging](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/README.md#logging)
   and you will see which credential the [DefaultAzureCredential] used as shown below, and why previously-attempted credentials
   were rejected.

   ```text
   AzureCLICredential.GetToken() acquired a token for scope https://vault.azure.net/.default
   ```

#### Access Denied to First Party Service

You may see an error similar to:

```text
--------------------------------------------------------------------------------
RESPONSE 403: 403 Forbidden
ERROR CODE: Forbidden
--------------------------------------------------------------------------------
{
  "error": {
    "code": "Forbidden",
    "message": "Access denied to first party service...",
    "innererror": {
      "code": "AccessDenied"
    }
  }
}
```

The error `message` may also contain the tenant ID (`tid`) and application ID (`appid`). This error may occur because:

1. You have the **Allow trust services** option enabled and are trying to access the Key Vault from a service not on
   [this list](https://learn.microsoft.com/azure/key-vault/general/overview-vnet-service-endpoints#trusted-services) of
   trusted services.
2. You logged into a Microsoft Account (MSA). See [above](#operation-not-permitted) for troubleshooting steps.

### Other Authentication Issues

If you are using the `azidentity` module to authenticate Azure Key Vault clients, please see its
[troubleshooting guide](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/TROUBLESHOOTING.md).

#### Incorrect challenge resource

If an error is thrown with a message similar to:

```text
challenge resource 'myvault.vault.azure.net' doesn't match the requested domain. Set DisableChallengeResourceVerification to true in your client options to disable. See https://aka.ms/azsdk/blog/vault-uri for more information
```

Check that the resources is expected - that you're not receiving a challenge from an unknown host which may indicate an incorrect request URI. If it is correct but you are using a mock service or non-transparent proxy for testing, set the DisableChallengeResourceVerification to true in your client options:

```go
vaultURL := "https://myvault.vault.azure.net"
credential, err := azidentity.NewDefaultAzureCredential(nil)
options := azsecrets.ClientOptions{
    DisableChallengeResourceVerification: true,
}
client := azsecrets.NewClient(vaultURI, credential, &options)
```

Read our [release notes][release_notes_resource] for more information about this change.

## Other Service Errors

To troubleshoot Key Vault errors not described in this guide,
see [Azure Key Vault REST API Error Codes](https://learn.microsoft.com/azure/key-vault/general/rest-error-codes).

### HTTP 429: Too Many Requests

If you get an error or see logs describing an HTTP 429 response, you may be making too many requests to Key Vault too quickly.

Possible solutions include:

1. Use a single instance of any client in your application for a single Key Vault.
2. Use a single credential instance for all clients.
3. Cache Key Vault resources (certificates, keys, secrets) in memory to reduce calls to retrieve them.

See the [Azure Key Vault throttling guide](https://learn.microsoft.com/azure/key-vault/general/overview-throttling)
for more information.

[azidentity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[DefaultAzureCredential]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/README.md#defaultazurecredential
[release_notes_resource]: https://devblogs.microsoft.com/azure-sdk/guidance-for-applications-using-the-key-vault-libraries/