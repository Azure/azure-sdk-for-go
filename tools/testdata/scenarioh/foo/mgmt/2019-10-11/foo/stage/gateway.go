// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package foo

// NatGatewaysClient is the network Client
type GatewaysClient struct {
	BaseClient
}

// NewGatewaysClient creates an instance of the NatGatewaysClient client.
func NewGatewaysClient(subscriptionID string) GatewaysClient {
	return NewGatewaysClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewGatewaysClientWithBaseURI creates an instance of the NatGatewaysClient client.
func NewGatewaysClientWithBaseURI(baseURI string, subscriptionID string) GatewaysClient {
	return GatewaysClient{NewWithBaseURI(baseURI, subscriptionID)}
}

func (client *GatewaysClient) CreateOrUpdate(resGroup string, parameters Gateway) error {
	return nil
}

func (client *GatewaysClient) DoSomething(something Something) error {
	return nil
}
