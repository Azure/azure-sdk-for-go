# frequently-asked-questions-for-azure-go-mgmt-sdk


## json unmarshal error
* when user reports a bug which description contains "json unmarshal error",we need to do the followings to solve:
	- first ,we can  suggest user to open the logger button to show the response body,so that we can see the details the api returned to check whether it is a bug of sdk,[example like this](https://github.com/Azure/azure-sdk-for-go/issues/23578#event-15727039059)
    - if the sdk unmarshal result does different from the api info, then check the reported sdk version, if it is not the latest version, we can first try the latest version locally, then suggest user to upgrade sdk version to the latest if the latest version is ok [example like this](https://github.com/Azure/azure-sdk-for-go/issues/23883#event-15755949673).otherwise, add label `service-attention`


## add service-attention label
when user report something related to product experience or related to functions about the service, add label  `service-attention`, [example like this](https://github.com/Azure/azure-sdk-for-go/issues/23867)


## some cases is not belong mgmt
* some reports does not belong to sdk of resource mangement,e.g:https://github.com/Azure/azure-sdk-for-go/issues/23895, this revolved in `management group`,not resource
* some reports is related to namespace started with `az`, e.g: https://github.com/Azure/azure-sdk-for-go/issues/23889#issuecomment-2565844706, this is reported about `azcosmos` which is not under `sdk/resourcemanager` directory, so we will not deal issues like this

