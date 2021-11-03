# Release History

## 0.2.1 (2021-11-03)
### Breaking Changes

### New Content

- New const `HostingModelConnectedContainer`
- New const `HostingModelDisconnectedContainer`
- New const `DeploymentProvisioningStateMoving`
- New const `DeploymentScaleTypeManual`
- New const `HostingModelWeb`
- New const `DeploymentProvisioningStateFailed`
- New const `DeploymentProvisioningStateAccepted`
- New const `DeploymentProvisioningStateSucceeded`
- New const `DeploymentProvisioningStateDeleting`
- New const `DeploymentProvisioningStateCreating`
- New function `*DeploymentsListPager.NextPage(context.Context) bool`
- New function `*CommitmentPlansListPager.Err() error`
- New function `PossibleDeploymentProvisioningStateValues() []DeploymentProvisioningState`
- New function `DeploymentProvisioningState.ToPtr() *DeploymentProvisioningState`
- New function `*DeploymentsListPager.PageResponse() DeploymentsListResponse`
- New function `*DeploymentsDeletePollerResponse.Resume(context.Context, *DeploymentsClient, string) error`
- New function `*CommitmentTiersListPager.NextPage(context.Context) bool`
- New function `*DeploymentsListPager.Err() error`
- New function `*CommitmentPlansClient.BeginDelete(context.Context, string, string, string, *CommitmentPlansBeginDeleteOptions) (CommitmentPlansDeletePollerResponse, error)`
- New function `*DeploymentsDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `DeploymentsDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (DeploymentsDeleteResponse, error)`
- New function `HostingModel.ToPtr() *HostingModel`
- New function `CommitmentPlansDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (CommitmentPlansDeleteResponse, error)`
- New function `*CommitmentTiersListPager.Err() error`
- New function `*DeploymentsCreateOrUpdatePoller.FinalResponse(context.Context) (DeploymentsCreateOrUpdateResponse, error)`
- New function `NewCommitmentTiersClient(string, azcore.TokenCredential, *arm.ClientOptions) *CommitmentTiersClient`
- New function `*CommitmentPlansDeletePoller.FinalResponse(context.Context) (CommitmentPlansDeleteResponse, error)`
- New function `DeploymentsCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (DeploymentsCreateOrUpdateResponse, error)`
- New function `*CommitmentPlansClient.CreateOrUpdate(context.Context, string, string, string, CommitmentPlan, *CommitmentPlansCreateOrUpdateOptions) (CommitmentPlansCreateOrUpdateResponse, error)`
- New function `*CommitmentPlansListPager.PageResponse() CommitmentPlansListResponse`
- New function `*DeploymentsClient.List(string, string, *DeploymentsListOptions) *DeploymentsListPager`
- New function `*CommitmentPlansListPager.NextPage(context.Context) bool`
- New function `*DeploymentsCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `PossibleDeploymentScaleTypeValues() []DeploymentScaleType`
- New function `*CommitmentPlansDeletePoller.ResumeToken() (string, error)`
- New function `*CommitmentPlansDeletePollerResponse.Resume(context.Context, *CommitmentPlansClient, string) error`
- New function `Deployment.MarshalJSON() ([]byte, error)`
- New function `CommitmentPlan.MarshalJSON() ([]byte, error)`
- New function `DeploymentScaleType.ToPtr() *DeploymentScaleType`
- New function `CommitmentPlanListResult.MarshalJSON() ([]byte, error)`
- New function `CommitmentTierListResult.MarshalJSON() ([]byte, error)`
- New function `NewCommitmentPlansClient(string, azcore.TokenCredential, *arm.ClientOptions) *CommitmentPlansClient`
- New function `*DeploymentsDeletePoller.ResumeToken() (string, error)`
- New function `*DeploymentsClient.BeginCreateOrUpdate(context.Context, string, string, string, Deployment, *DeploymentsBeginCreateOrUpdateOptions) (DeploymentsCreateOrUpdatePollerResponse, error)`
- New function `*DeploymentsCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*DeploymentsCreateOrUpdatePoller.Done() bool`
- New function `*CommitmentPlansClient.List(string, string, *CommitmentPlansListOptions) *CommitmentPlansListPager`
- New function `*CommitmentPlansDeletePoller.Done() bool`
- New function `*CommitmentPlansDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*CommitmentTiersListPager.PageResponse() CommitmentTiersListResponse`
- New function `*DeploymentsClient.Get(context.Context, string, string, string, *DeploymentsGetOptions) (DeploymentsGetResponse, error)`
- New function `*DeploymentsDeletePoller.FinalResponse(context.Context) (DeploymentsDeleteResponse, error)`
- New function `NewDeploymentsClient(string, azcore.TokenCredential, *arm.ClientOptions) *DeploymentsClient`
- New function `*CommitmentPlansClient.Get(context.Context, string, string, string, *CommitmentPlansGetOptions) (CommitmentPlansGetResponse, error)`
- New function `*CommitmentTiersClient.List(string, *CommitmentTiersListOptions) *CommitmentTiersListPager`
- New function `*DeploymentsCreateOrUpdatePollerResponse.Resume(context.Context, *DeploymentsClient, string) error`
- New function `DeploymentListResult.MarshalJSON() ([]byte, error)`
- New function `*DeploymentsClient.BeginDelete(context.Context, string, string, string, *DeploymentsBeginDeleteOptions) (DeploymentsDeletePollerResponse, error)`
- New function `*DeploymentsDeletePoller.Done() bool`
- New function `PossibleHostingModelValues() []HostingModel`
- New struct `CommitmentCost`
- New struct `CommitmentPeriod`
- New struct `CommitmentPlan`
- New struct `CommitmentPlanListResult`
- New struct `CommitmentPlanProperties`
- New struct `CommitmentPlansBeginDeleteOptions`
- New struct `CommitmentPlansClient`
- New struct `CommitmentPlansCreateOrUpdateOptions`
- New struct `CommitmentPlansCreateOrUpdateResponse`
- New struct `CommitmentPlansCreateOrUpdateResult`
- New struct `CommitmentPlansDeletePoller`
- New struct `CommitmentPlansDeletePollerResponse`
- New struct `CommitmentPlansDeleteResponse`
- New struct `CommitmentPlansGetOptions`
- New struct `CommitmentPlansGetResponse`
- New struct `CommitmentPlansGetResult`
- New struct `CommitmentPlansListOptions`
- New struct `CommitmentPlansListPager`
- New struct `CommitmentPlansListResponse`
- New struct `CommitmentPlansListResult`
- New struct `CommitmentQuota`
- New struct `CommitmentTier`
- New struct `CommitmentTierListResult`
- New struct `CommitmentTiersClient`
- New struct `CommitmentTiersListOptions`
- New struct `CommitmentTiersListPager`
- New struct `CommitmentTiersListResponse`
- New struct `CommitmentTiersListResult`
- New struct `Deployment`
- New struct `DeploymentListResult`
- New struct `DeploymentModel`
- New struct `DeploymentProperties`
- New struct `DeploymentScaleSettings`
- New struct `DeploymentsBeginCreateOrUpdateOptions`
- New struct `DeploymentsBeginDeleteOptions`
- New struct `DeploymentsClient`
- New struct `DeploymentsCreateOrUpdatePoller`
- New struct `DeploymentsCreateOrUpdatePollerResponse`
- New struct `DeploymentsCreateOrUpdateResponse`
- New struct `DeploymentsCreateOrUpdateResult`
- New struct `DeploymentsDeletePoller`
- New struct `DeploymentsDeletePollerResponse`
- New struct `DeploymentsDeleteResponse`
- New struct `DeploymentsGetOptions`
- New struct `DeploymentsGetResponse`
- New struct `DeploymentsGetResult`
- New struct `DeploymentsListOptions`
- New struct `DeploymentsListPager`
- New struct `DeploymentsListResponse`
- New struct `DeploymentsListResult`
- New struct `ProxyResource`
- New field `Kind` in struct `DomainAvailability`
- New field `Kind` in struct `CheckDomainAvailabilityParameter`

Total 0 breaking change(s), 164 additive change(s).


## 0.2.0 (2021-10-29)

### Breaking Changes

- `arm.Connection` has been removed in `github.com/Azure/azure-sdk-for-go/sdk/azcore/v0.20.0`
- The parameters of `NewXXXClient` has been changed from `(con *arm.Connection, subscriptionID string)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`

## 0.1.0 (2021-10-26)

- Initial preview release.
