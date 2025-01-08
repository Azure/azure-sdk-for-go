# frequently-asked-questions-for-azure-go-mgmt-sdk

## json unmarshal error
* when users reports bugs which contains "json unmarshal error",we need to do the followings to solve:
* 1、first ,we can  suggest user to open the logger button to show the response body,so that we can see the details the api returned to check whether it is a bug of sdk,[example like this](https://github.com/Azure/azure-sdk-for-go/issues/23883#event-15755949673)
*2、if the sdk does diff from the api info, then check the reported sdk version, if it is not the latest version, we can first try the latest version locally, then suggest user to upgrade sdk version to the latest if the latest version is ok [example like this](https://github.com/Azure/azure-sdk-for-go/issues/23883#event-15755949673)


## add service-attention label

