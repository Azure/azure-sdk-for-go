# armweightsandbiases - Service Retirement

The Microsoft.WeightsAndBiases service will be permanently retired on September 30, 2026. There is no replacement service or SDK module.

Stop using this module and remove it from your application:

```sh
go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/weightsandbiases/armweightsandbiases@none
go mod tidy
```

For more information, see the [Azure SDK support policy](https://aka.ms/azsdk/support-policies). For questions about this package, open an issue in the [Azure SDK for Go repository](https://github.com/Azure/azure-sdk-for-go/issues).