# Release History

## 0.2.0 (2026-08-03)
### Features Added

- New enum type `SKUMixPlacementAllocationStrategy` with values `SKUMixPlacementAllocationStrategyEvictionOptimized`, `SKUMixPlacementAllocationStrategyLowestPrice`, `SKUMixPlacementAllocationStrategyPrioritized`
- New enum type `SKUMixPlacementCapacityType` with values `SKUMixPlacementCapacityTypeVCPU`, `SKUMixPlacementCapacityTypeVM`
- New enum type `SKUMixPlacementOSType` with values `SKUMixPlacementOSTypeLinux`, `SKUMixPlacementOSTypeWindows`
- New enum type `SKUMixPlacementPartialFulfillmentReason` with values `SKUMixPlacementPartialFulfillmentReasonInsufficientCapacity`, `SKUMixPlacementPartialFulfillmentReasonInsufficientQuota`, `SKUMixPlacementPartialFulfillmentReasonNone`
- New enum type `SKUMixPlacementPriority` with values `SKUMixPlacementPriorityRegular`, `SKUMixPlacementPrioritySpot`
- New enum type `SKUMixPlacementZonalDistributionStrategy` with values `SKUMixPlacementZonalDistributionStrategyBestEffortBalanced`, `SKUMixPlacementZonalDistributionStrategyBestEffortSingleZone`, `SKUMixPlacementZonalDistributionStrategyPrioritized`
- New function `*ClientFactory.NewSKUMixPlacementScoresClient() *SKUMixPlacementScoresClient`
- New function `NewSKUMixPlacementScoresClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SKUMixPlacementScoresClient, error)`
- New function `*SKUMixPlacementScoresClient.Get(ctx context.Context, location string, options *SKUMixPlacementScoresClientGetOptions) (SKUMixPlacementScoresClientGetResponse, error)`
- New function `*SKUMixPlacementScoresClient.Post(ctx context.Context, location string, skuMixPlacementRequest SKUMixPlacementRequest, options *SKUMixPlacementScoresClientPostOptions) (SKUMixPlacementScoresClientPostResponse, error)`
- New struct `SKUMixPlacementBase`
- New struct `SKUMixPlacementCapacityProfile`
- New struct `SKUMixPlacementDeploymentChoice`
- New struct `SKUMixPlacementInstanceDescription`
- New struct `SKUMixPlacementItem`
- New struct `SKUMixPlacementProperties`
- New struct `SKUMixPlacementRequest`
- New struct `SKUMixPlacementResponse`
- New struct `SKUMixPlacementSpotPriorityProfile`
- New struct `SKUMixPlacementVMSize`
- New struct `SKUMixPlacementZoneAllocationPolicy`
- New struct `SKUMixPlacementZonePreference`


## 0.1.0 (2025-09-22)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armrecommender` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).