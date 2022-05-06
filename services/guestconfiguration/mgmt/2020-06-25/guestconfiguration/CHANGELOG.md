# Unreleased

## Additive Changes

### New Funcs

1. AssignmentsClient.RGList(context.Context, string) (AssignmentList, error)
1. AssignmentsClient.RGListPreparer(context.Context, string) (*http.Request, error)
1. AssignmentsClient.RGListResponder(*http.Response) (AssignmentList, error)
1. AssignmentsClient.RGListSender(*http.Request) (*http.Response, error)
1. AssignmentsClient.SubscriptionList(context.Context) (AssignmentList, error)
1. AssignmentsClient.SubscriptionListPreparer(context.Context) (*http.Request, error)
1. AssignmentsClient.SubscriptionListResponder(*http.Response) (AssignmentList, error)
1. AssignmentsClient.SubscriptionListSender(*http.Request) (*http.Response, error)
1. Navigation.MarshalJSON() ([]byte, error)
1. VMSSVMInfo.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. VMSSVMInfo

#### New Struct Fields

1. AssignmentProperties.ParameterHash
1. AssignmentProperties.ResourceType
1. AssignmentProperties.VmssVMList
1. AssignmentReportProperties.VmssResourceID
1. Navigation.ConfigurationProtectedParameter
1. Navigation.ContentType
