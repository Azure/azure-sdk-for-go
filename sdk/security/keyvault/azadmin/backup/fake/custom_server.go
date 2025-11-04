// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake

import (
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal"
)

func init() {
	serverTransportInterceptor = &internal.FakeChallenge{}
}
