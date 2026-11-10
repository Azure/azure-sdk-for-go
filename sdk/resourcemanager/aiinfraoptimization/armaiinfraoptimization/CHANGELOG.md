# Release History

## 0.1.1 (2026-11-10)
### Bugs Fixed

- Fixed `RecommendationsClient.RecommendVMSKU` so that errors returned while handling the response are recorded on the tracing span.

### Other Changes

- Updated the documentation of `RecommendVMSKUTokenBudget.MaxOutputTokens` and `RecommendVMSKUTokenBudget.PromptTokens` to reflect the revised token budget limits.
- Upgraded `github.com/Azure/azure-sdk-for-go/sdk/azcore` dependency to `v1.24.0`.

## 0.1.0 (2026-11-02)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/aiinfraoptimization/armaiinfraoptimization` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).