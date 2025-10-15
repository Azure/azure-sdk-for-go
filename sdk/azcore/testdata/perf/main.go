// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

func main() {
	perf.Run(map[string]perf.PerfMethods{
		"ClientGET": {Register: clientTestRegister, New: newClientGETTest},
		"MockJson":  {Register: mockJsonTestRegister, New: NewMockJsonTest},
		"MockXml":   {Register: mockXmlTestRegister, New: NewMockXmlTest},
	})
}
