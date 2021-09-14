# Unreleased

## Additive Changes

### New Constants

1. RegistrationStatus.Activated
1. RegistrationStatus.ActivationAttemptsExhausted
1. RegistrationStatus.ActivationPending

### New Funcs

1. EmailRegistration.MarshalJSON() ([]byte, error)
1. EmailRegistrationsClient.ActivateEmail(context.Context, string, EmailRegistration) (EmailRegistration, error)
1. EmailRegistrationsClient.ActivateEmailPreparer(context.Context, string, EmailRegistration) (*http.Request, error)
1. EmailRegistrationsClient.ActivateEmailResponder(*http.Response) (EmailRegistration, error)
1. EmailRegistrationsClient.ActivateEmailSender(*http.Request) (*http.Response, error)
1. EmailRegistrationsClient.RegisterEmail(context.Context, string) (EmailRegistration, error)
1. EmailRegistrationsClient.RegisterEmailPreparer(context.Context, string) (*http.Request, error)
1. EmailRegistrationsClient.RegisterEmailResponder(*http.Response) (EmailRegistration, error)
1. EmailRegistrationsClient.RegisterEmailSender(*http.Request) (*http.Response, error)
1. NewEmailRegistrationsClient(string) EmailRegistrationsClient
1. NewEmailRegistrationsClientWithBaseURI(string, string) EmailRegistrationsClient
1. PossibleRegistrationStatusValues() []RegistrationStatus

### Struct Changes

#### New Structs

1. EmailRegistration
1. EmailRegistrationsClient

#### New Struct Fields

1. OperationMetaMetricSpecification.FillGapWithZero
