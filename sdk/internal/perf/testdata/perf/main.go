// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import "github.com/Azure/azure-sdk-for-go/sdk/internal/perf"

func main() {
	perf.Run(map[string]perf.PerfMethods{
		"NoOpTest":  {Register: nil, New: NewNoOpTest},
		"SleepTest": {Register: sleepTestRegister, New: NewSleepTest},
	})
}
