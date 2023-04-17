//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznotificationhubs

import (
	"strings"
)

type (
	// Represents a request to send a notification to Azure Notification Hubs.
	NotificationRequest struct {
		Message     string            // The message body to send.
		Headers     map[string]string // The HTTP headers to send.
		Platform    string            // The platform to send to such as "adm", "apple", "baidu", "browser", "gcm", "windows", "xiaomi".
		ContentType string            // The content type of the message such as "application/json;charset=utf-8", "application/xml;charset=utf-8", "application/octet-stream".
	}

	// Represents a response from Azure Notification Hubs
	NotificationResponse struct {
		TrackingId    string // The TrackingId for the request.
		CorrelationId string // The correlation ID for the request.
		Location      string // The location (if present) for the resource.
	}

	// Represents a send-based response from Azure Notification Hubs
	NotificationMessageResponse struct {
		TrackingId     string // The TrackingId for the request.
		CorrelationId  string // The correlation ID for the request.
		NotificationId string // The notification ID for the request.
	}
)

// Creates a tag expression
func CreateTagExpression(tags []string) string {
	return strings.Join(tags, "||")
}
