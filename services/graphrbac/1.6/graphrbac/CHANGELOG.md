Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

### New Constants

1. ObjectType.ObjectTypeAppRoleAssignment

### New Funcs

1. *AppRoleAssignment.UnmarshalJSON([]byte) error
1. *AppRoleAssignmentListResultIterator.Next() error
1. *AppRoleAssignmentListResultIterator.NextWithContext(context.Context) error
1. *AppRoleAssignmentListResultPage.Next() error
1. *AppRoleAssignmentListResultPage.NextWithContext(context.Context) error
1. ADGroup.AsAppRoleAssignment() (*AppRoleAssignment, bool)
1. AppRoleAssignment.AsADGroup() (*ADGroup, bool)
1. AppRoleAssignment.AsAppRoleAssignment() (*AppRoleAssignment, bool)
1. AppRoleAssignment.AsApplication() (*Application, bool)
1. AppRoleAssignment.AsBasicDirectoryObject() (BasicDirectoryObject, bool)
1. AppRoleAssignment.AsDirectoryObject() (*DirectoryObject, bool)
1. AppRoleAssignment.AsServicePrincipal() (*ServicePrincipal, bool)
1. AppRoleAssignment.AsUser() (*User, bool)
1. AppRoleAssignment.MarshalJSON() ([]byte, error)
1. AppRoleAssignmentListResult.IsEmpty() bool
1. AppRoleAssignmentListResultIterator.NotDone() bool
1. AppRoleAssignmentListResultIterator.Response() AppRoleAssignmentListResult
1. AppRoleAssignmentListResultIterator.Value() AppRoleAssignment
1. AppRoleAssignmentListResultPage.NotDone() bool
1. AppRoleAssignmentListResultPage.Response() AppRoleAssignmentListResult
1. AppRoleAssignmentListResultPage.Values() []AppRoleAssignment
1. Application.AsAppRoleAssignment() (*AppRoleAssignment, bool)
1. DirectoryObject.AsAppRoleAssignment() (*AppRoleAssignment, bool)
1. NewAppRoleAssignmentListResultIterator(AppRoleAssignmentListResultPage) AppRoleAssignmentListResultIterator
1. NewAppRoleAssignmentListResultPage(AppRoleAssignmentListResult, func(context.Context, AppRoleAssignmentListResult) (AppRoleAssignmentListResult, error)) AppRoleAssignmentListResultPage
1. ServicePrincipal.AsAppRoleAssignment() (*AppRoleAssignment, bool)
1. ServicePrincipalsClient.ListAppRoleAssignedTo(context.Context, string) (AppRoleAssignmentListResultPage, error)
1. ServicePrincipalsClient.ListAppRoleAssignedToComplete(context.Context, string) (AppRoleAssignmentListResultIterator, error)
1. ServicePrincipalsClient.ListAppRoleAssignedToPreparer(context.Context, string) (*http.Request, error)
1. ServicePrincipalsClient.ListAppRoleAssignedToResponder(*http.Response) (AppRoleAssignmentListResult, error)
1. ServicePrincipalsClient.ListAppRoleAssignedToSender(*http.Request) (*http.Response, error)
1. ServicePrincipalsClient.ListAppRoleAssignments(context.Context, string) (AppRoleAssignmentListResultPage, error)
1. ServicePrincipalsClient.ListAppRoleAssignmentsComplete(context.Context, string) (AppRoleAssignmentListResultIterator, error)
1. ServicePrincipalsClient.ListAppRoleAssignmentsPreparer(context.Context, string) (*http.Request, error)
1. ServicePrincipalsClient.ListAppRoleAssignmentsResponder(*http.Response) (AppRoleAssignmentListResult, error)
1. ServicePrincipalsClient.ListAppRoleAssignmentsSender(*http.Request) (*http.Response, error)
1. User.AsAppRoleAssignment() (*AppRoleAssignment, bool)

## Struct Changes

### New Structs

1. AppRoleAssignment
1. AppRoleAssignmentListResult
1. AppRoleAssignmentListResultIterator
1. AppRoleAssignmentListResultPage
