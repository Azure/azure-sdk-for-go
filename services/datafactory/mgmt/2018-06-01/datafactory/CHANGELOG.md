# Change History

## Breaking Changes

### Removed Constants

1. AuthenticationType.AuthenticationTypeAuthenticationTypeAnonymous
1. AuthenticationType.AuthenticationTypeAuthenticationTypeBasic
1. AuthenticationType.AuthenticationTypeAuthenticationTypeClientCertificate
1. AuthenticationType.AuthenticationTypeAuthenticationTypeWebLinkedServiceTypeProperties
1. AuthorizationType.AuthorizationTypeAuthorizationTypeKey
1. AuthorizationType.AuthorizationTypeAuthorizationTypeLinkedIntegrationRuntimeType
1. AuthorizationType.AuthorizationTypeAuthorizationTypeRBAC
1. Type.TypeTypeAzureKeyVaultSecret
1. Type.TypeTypeSecretBase
1. Type.TypeTypeSecureString

## Additive Changes

### New Constants

1. AuthenticationType.AuthenticationTypeAnonymous
1. AuthenticationType.AuthenticationTypeBasic
1. AuthenticationType.AuthenticationTypeClientCertificate
1. AuthenticationType.AuthenticationTypeWebLinkedServiceTypeProperties
1. AuthorizationType.AuthorizationTypeKey
1. AuthorizationType.AuthorizationTypeLinkedIntegrationRuntimeType
1. AuthorizationType.AuthorizationTypeRBAC
1. Type.TypeAzureKeyVaultSecret
1. Type.TypeSecretBase
1. Type.TypeSecureString

### New Funcs

1. ArmIDWrapper.MarshalJSON() ([]byte, error)
1. ConnectionStateProperties.MarshalJSON() ([]byte, error)
1. ExposureControlResponse.MarshalJSON() ([]byte, error)
1. IntegrationRuntimeNodeIPAddress.MarshalJSON() ([]byte, error)
1. LinkedIntegrationRuntime.MarshalJSON() ([]byte, error)
1. ManagedIntegrationRuntimeError.MarshalJSON() ([]byte, error)
1. ManagedIntegrationRuntimeOperationResult.MarshalJSON() ([]byte, error)
1. ManagedIntegrationRuntimeStatusTypeProperties.MarshalJSON() ([]byte, error)
1. PipelineRunInvokedBy.MarshalJSON() ([]byte, error)
1. PrivateLinkResourceProperties.MarshalJSON() ([]byte, error)
1. SubResource.MarshalJSON() ([]byte, error)
1. TriggerSubscriptionOperationStatus.MarshalJSON() ([]byte, error)
