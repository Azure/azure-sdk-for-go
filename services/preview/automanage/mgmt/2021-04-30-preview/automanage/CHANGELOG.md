# Unreleased

## Additive Changes

### New Funcs

1. *ServicePrincipal.UnmarshalJSON([]byte) error
1. ConfigurationProfileAssignmentsClient.ListByClusterName(context.Context, string, string) (ConfigurationProfileAssignmentList, error)
1. ConfigurationProfileAssignmentsClient.ListByClusterNamePreparer(context.Context, string, string) (*http.Request, error)
1. ConfigurationProfileAssignmentsClient.ListByClusterNameResponder(*http.Response) (ConfigurationProfileAssignmentList, error)
1. ConfigurationProfileAssignmentsClient.ListByClusterNameSender(*http.Request) (*http.Response, error)
1. ConfigurationProfileAssignmentsClient.ListByMachineName(context.Context, string, string) (ConfigurationProfileAssignmentList, error)
1. ConfigurationProfileAssignmentsClient.ListByMachineNamePreparer(context.Context, string, string) (*http.Request, error)
1. ConfigurationProfileAssignmentsClient.ListByMachineNameResponder(*http.Response) (ConfigurationProfileAssignmentList, error)
1. ConfigurationProfileAssignmentsClient.ListByMachineNameSender(*http.Request) (*http.Response, error)
1. ConfigurationProfileAssignmentsClient.ListByVirtualMachines(context.Context, string, string) (ConfigurationProfileAssignmentList, error)
1. ConfigurationProfileAssignmentsClient.ListByVirtualMachinesPreparer(context.Context, string, string) (*http.Request, error)
1. ConfigurationProfileAssignmentsClient.ListByVirtualMachinesResponder(*http.Response) (ConfigurationProfileAssignmentList, error)
1. ConfigurationProfileAssignmentsClient.ListByVirtualMachinesSender(*http.Request) (*http.Response, error)
1. ConfigurationProfileHCIAssignmentsClient.CreateOrUpdate(context.Context, ConfigurationProfileAssignment, string, string, string) (ConfigurationProfileAssignment, error)
1. ConfigurationProfileHCIAssignmentsClient.CreateOrUpdatePreparer(context.Context, ConfigurationProfileAssignment, string, string, string) (*http.Request, error)
1. ConfigurationProfileHCIAssignmentsClient.CreateOrUpdateResponder(*http.Response) (ConfigurationProfileAssignment, error)
1. ConfigurationProfileHCIAssignmentsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ConfigurationProfileHCIAssignmentsClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. ConfigurationProfileHCIAssignmentsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. ConfigurationProfileHCIAssignmentsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ConfigurationProfileHCIAssignmentsClient.DeleteSender(*http.Request) (*http.Response, error)
1. ConfigurationProfileHCIAssignmentsClient.Get(context.Context, string, string, string) (ConfigurationProfileAssignment, error)
1. ConfigurationProfileHCIAssignmentsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ConfigurationProfileHCIAssignmentsClient.GetResponder(*http.Response) (ConfigurationProfileAssignment, error)
1. ConfigurationProfileHCIAssignmentsClient.GetSender(*http.Request) (*http.Response, error)
1. ConfigurationProfileHCRPAssignmentsClient.CreateOrUpdate(context.Context, ConfigurationProfileAssignment, string, string, string) (ConfigurationProfileAssignment, error)
1. ConfigurationProfileHCRPAssignmentsClient.CreateOrUpdatePreparer(context.Context, ConfigurationProfileAssignment, string, string, string) (*http.Request, error)
1. ConfigurationProfileHCRPAssignmentsClient.CreateOrUpdateResponder(*http.Response) (ConfigurationProfileAssignment, error)
1. ConfigurationProfileHCRPAssignmentsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ConfigurationProfileHCRPAssignmentsClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. ConfigurationProfileHCRPAssignmentsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. ConfigurationProfileHCRPAssignmentsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ConfigurationProfileHCRPAssignmentsClient.DeleteSender(*http.Request) (*http.Response, error)
1. ConfigurationProfileHCRPAssignmentsClient.Get(context.Context, string, string, string) (ConfigurationProfileAssignment, error)
1. ConfigurationProfileHCRPAssignmentsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ConfigurationProfileHCRPAssignmentsClient.GetResponder(*http.Response) (ConfigurationProfileAssignment, error)
1. ConfigurationProfileHCRPAssignmentsClient.GetSender(*http.Request) (*http.Response, error)
1. HCIReportsClient.Get(context.Context, string, string, string, string) (Report, error)
1. HCIReportsClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. HCIReportsClient.GetResponder(*http.Response) (Report, error)
1. HCIReportsClient.GetSender(*http.Request) (*http.Response, error)
1. HCIReportsClient.ListByConfigurationProfileAssignments(context.Context, string, string, string) (ReportList, error)
1. HCIReportsClient.ListByConfigurationProfileAssignmentsPreparer(context.Context, string, string, string) (*http.Request, error)
1. HCIReportsClient.ListByConfigurationProfileAssignmentsResponder(*http.Response) (ReportList, error)
1. HCIReportsClient.ListByConfigurationProfileAssignmentsSender(*http.Request) (*http.Response, error)
1. HCRPReportsClient.Get(context.Context, string, string, string, string) (Report, error)
1. HCRPReportsClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. HCRPReportsClient.GetResponder(*http.Response) (Report, error)
1. HCRPReportsClient.GetSender(*http.Request) (*http.Response, error)
1. HCRPReportsClient.ListByConfigurationProfileAssignments(context.Context, string, string, string) (ReportList, error)
1. HCRPReportsClient.ListByConfigurationProfileAssignmentsPreparer(context.Context, string, string, string) (*http.Request, error)
1. HCRPReportsClient.ListByConfigurationProfileAssignmentsResponder(*http.Response) (ReportList, error)
1. HCRPReportsClient.ListByConfigurationProfileAssignmentsSender(*http.Request) (*http.Response, error)
1. NewConfigurationProfileHCIAssignmentsClient(string) ConfigurationProfileHCIAssignmentsClient
1. NewConfigurationProfileHCIAssignmentsClientWithBaseURI(string, string) ConfigurationProfileHCIAssignmentsClient
1. NewConfigurationProfileHCRPAssignmentsClient(string) ConfigurationProfileHCRPAssignmentsClient
1. NewConfigurationProfileHCRPAssignmentsClientWithBaseURI(string, string) ConfigurationProfileHCRPAssignmentsClient
1. NewHCIReportsClient(string) HCIReportsClient
1. NewHCIReportsClientWithBaseURI(string, string) HCIReportsClient
1. NewHCRPReportsClient(string) HCRPReportsClient
1. NewHCRPReportsClientWithBaseURI(string, string) HCRPReportsClient
1. NewServicePrincipalsClient(string) ServicePrincipalsClient
1. NewServicePrincipalsClientWithBaseURI(string, string) ServicePrincipalsClient
1. ServicePrincipal.MarshalJSON() ([]byte, error)
1. ServicePrincipalProperties.MarshalJSON() ([]byte, error)
1. ServicePrincipalsClient.Get(context.Context) (ServicePrincipal, error)
1. ServicePrincipalsClient.GetPreparer(context.Context) (*http.Request, error)
1. ServicePrincipalsClient.GetResponder(*http.Response) (ServicePrincipal, error)
1. ServicePrincipalsClient.GetSender(*http.Request) (*http.Response, error)
1. ServicePrincipalsClient.ListBySubscription(context.Context) (ServicePrincipalListResult, error)
1. ServicePrincipalsClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. ServicePrincipalsClient.ListBySubscriptionResponder(*http.Response) (ServicePrincipalListResult, error)
1. ServicePrincipalsClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. ConfigurationProfileHCIAssignmentsClient
1. ConfigurationProfileHCRPAssignmentsClient
1. HCIReportsClient
1. HCRPReportsClient
1. ServicePrincipal
1. ServicePrincipalListResult
1. ServicePrincipalProperties
1. ServicePrincipalsClient
