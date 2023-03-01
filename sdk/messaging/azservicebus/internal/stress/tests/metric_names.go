// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

type MetricName string

const (
	MetricNameSessionTimeoutMS MetricName = "SessionTimeoutMS"
	MetricNameConnectionLost   MetricName = "ConnectionLost"
	MetricNameMessageSent      MetricName = "MessageSent"
)
