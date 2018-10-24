// Package links implements the Azure ARM Links service API version 2016-09-01.
//
// Azure resources can be linked together to form logical relationships. You can establish links between resources
// belonging to different resource groups. However, all the linked resources must belong to the same subscription. Each
// resource can be linked to 50 other resources. If any of the linked resources are deleted or moved, the link owner
// must clean up the remaining link.
package links

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
    "github.com/Azure/go-autorest/autorest"
)

const (
// DefaultBaseURI is the default URI used for the service Links
DefaultBaseURI = "https://management.azure.com")

// BaseClient is the base client for Links.
type BaseClient struct {
    autorest.Client
    BaseURI string
            SubscriptionID string
}

// New creates an instance of the BaseClient client.
func New(subscriptionID string)BaseClient {
    return NewWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewWithBaseURI creates an instance of the BaseClient client.
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
    return BaseClient{
        Client: autorest.NewClientWithUserAgent(UserAgent()),
        BaseURI: baseURI,
                SubscriptionID: subscriptionID,
    }
}

