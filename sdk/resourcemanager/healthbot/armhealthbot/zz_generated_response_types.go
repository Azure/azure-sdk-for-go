//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhealthbot

// BotsClientCreateResponse contains the response from method BotsClient.Create.
type BotsClientCreateResponse struct {
	HealthBot
}

// BotsClientDeleteResponse contains the response from method BotsClient.Delete.
type BotsClientDeleteResponse struct {
	// placeholder for future response values
}

// BotsClientGetResponse contains the response from method BotsClient.Get.
type BotsClientGetResponse struct {
	HealthBot
}

// BotsClientListByResourceGroupResponse contains the response from method BotsClient.ListByResourceGroup.
type BotsClientListByResourceGroupResponse struct {
	BotResponseList
}

// BotsClientListResponse contains the response from method BotsClient.List.
type BotsClientListResponse struct {
	BotResponseList
}

// BotsClientUpdateResponse contains the response from method BotsClient.Update.
type BotsClientUpdateResponse struct {
	HealthBot
}

// OperationsClientListResponse contains the response from method OperationsClient.List.
type OperationsClientListResponse struct {
	AvailableOperations
}
