// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake

import (
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azlogs"
)

func preferHeaderToQueryOptions(s string) *azlogs.QueryOptions {
	var opts azlogs.QueryOptions
	for _, option := range strings.Split(s, ",") {
		switch {
		case strings.HasPrefix(option, "include-statistics="):
			val := strings.TrimPrefix(option, "include-statistics=")
			if val == "true" {
				opts.Statistics = to.Ptr(true)
			}
			if val == "false" {
				opts.Statistics = to.Ptr(false)
			}
		case strings.HasPrefix(option, "include-render="):
			val := strings.TrimPrefix(option, "include-render=")
			if val == "true" {
				opts.Visualization = to.Ptr(true)
			}
			if val == "false" {
				opts.Visualization = to.Ptr(false)
			}
		case strings.HasPrefix(option, "wait="):
			val := strings.TrimPrefix(option, "wait=")
			wait, err := strconv.Atoi(val)
			if err == nil {
				opts.Wait = &wait
			}
		}
	}
	return &opts
}
