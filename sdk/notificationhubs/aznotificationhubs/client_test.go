//go:build go1.20
// +build go1.20

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznotificationhubs

import (
	"context"
	"os"
	"testing"
	"time"
)

const (
	anhConnectionStringVariable = "NOTIFICATIONHUBS_CONNECTION_STRING"
	anhAppleDeviceVariable      = "NOTIFICATIONHUBS_APPLE_DEVICE_TOKEN"
	anhAppleDeviceToken         = "00fc13adff785122b4ad28809a3420982341241421348097878e577c991de8f0"
	anhHubVariable              = "NOTIFICATION_HUB_NAME"
	anhInstallationIdVariable   = "NOTIFICATIONHUBS_INSTALLATION_ID"
	anhInstallationId           = "3e2cdcaed0f7-4fe5-853f-40c2e2efe387"
	messageBody                 = `{"aps": { "alert": { "title": "My title", "body": "My body" } } }`
)

func getEnvWithFallback(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func TestDirectSend(t *testing.T) {
	connectionString := os.Getenv(anhConnectionStringVariable)
	hubName := os.Getenv(anhHubVariable)
	deviceToken := getEnvWithFallback(anhAppleDeviceVariable, anhAppleDeviceToken)

	if connectionString == "" || hubName == "" || deviceToken == "" {
		t.Skip(`Skipping test due to missing environment variables`)
	}

	client, err := NewNotificationHubClientFromConnectionString(connectionString, hubName, nil)
	if client == nil || err != nil {
		t.Fatalf(`NewNotificationHubClientFromConnectionString %v`, err)
	}

	headers := make(map[string]string)
	headers["apns-priority"] = "10"
	headers["apns-push-type"] = "alert"

	contentType := "application/json;charset=utf-8"
	platform := "apple"

	request := &NotificationRequest{
		Message:     messageBody,
		Headers:     headers,
		Platform:    platform,
		ContentType: contentType,
	}

	response, err := client.SendDirectNotification(context.TODO(), request, deviceToken)
	if response == nil || err != nil {
		t.Fatalf(`SendDirectNotification %v`, err)
	}
}

func TestTagExpression(t *testing.T) {
	connectionString := os.Getenv(anhConnectionStringVariable)
	hubName := os.Getenv(anhHubVariable)

	if connectionString == "" || hubName == "" {
		t.Skip(`Skipping test due to missing environment variables`)
	}

	client, err := NewNotificationHubClientFromConnectionString(connectionString, hubName, nil)
	if client == nil || err != nil {
		t.Fatalf(`NewNotificationHubClientFromConnectionString %v`, err)
	}

	headers := make(map[string]string)
	headers["apns-priority"] = "10"
	headers["apns-push-type"] = "alert"

	contentType := "application/json;charset=utf-8"
	platform := "apple"

	request := &NotificationRequest{
		Message:     messageBody,
		Headers:     headers,
		Platform:    platform,
		ContentType: contentType,
	}

	tagExpression := "language_en&&country_US"

	response, err := client.SendNotification(context.TODO(), request, tagExpression)
	if response == nil || err != nil {
		t.Fatalf(`SendNotificationWithTagExpression %v`, err)
	}
}

func TestScheduledSend(t *testing.T) {
	connectionString := os.Getenv(anhConnectionStringVariable)
	hubName := os.Getenv(anhHubVariable)

	if connectionString == "" || hubName == "" {
		t.Skip(`Skipping test due to missing environment variables`)
	}

	client, err := NewNotificationHubClientFromConnectionString(connectionString, hubName, nil)
	if client == nil || err != nil {
		t.Fatalf(`NewNotificationHubClientFromConnectionString %v`, err)
	}

	headers := make(map[string]string)
	headers["apns-priority"] = "10"
	headers["apns-push-type"] = "alert"

	contentType := "application/json;charset=utf-8"
	platform := "apple"

	request := &NotificationRequest{
		Message:     messageBody,
		Headers:     headers,
		Platform:    platform,
		ContentType: contentType,
	}

	tagExpression := "language_en&&country_US"
	scheduleTime := time.Now().Add(time.Hour * 8)

	// Schedule a notification
	response, err := client.SendScheduledNotification(context.TODO(), request, tagExpression, scheduleTime)
	if response == nil || err != nil {
		t.Fatalf(`SendScheduledNotification %v`, err)
	}

	// Cancel the scheduled notification
	notificationId := response.NotificationId
	res, err := client.CancelScheduledNotification(context.TODO(), notificationId)
	if res == nil || err != nil {
		t.Fatalf(`CancelScheduledNotification %v`, err)
	}
}

func TestCreateInstallation(t *testing.T) {
	connectionString := os.Getenv(anhConnectionStringVariable)
	hubName := os.Getenv(anhHubVariable)
	deviceToken := getEnvWithFallback(anhAppleDeviceVariable, anhAppleDeviceToken)
	installationID := getEnvWithFallback(anhInstallationIdVariable, anhInstallationId)

	if connectionString == "" || hubName == "" || deviceToken == "" {
		t.Skip(`Skipping test due to missing environment variables`)
	}

	client, err := NewNotificationHubClientFromConnectionString(connectionString, hubName, nil)
	if client == nil || err != nil {
		t.Fatalf(`NewNotificationHubClientFromConnectionString %v`, err)
	}

	installation := &Installation{
		InstallationID: installationID,
		PushChannel:    deviceToken,
		Platform:       "apns",
		Tags:           []string{"language_en", "country_US"},
	}

	response, err := client.CreateOrUpdateInstallation(context.TODO(), installation)
	if response == nil || err != nil {
		t.Fatalf(`CreateOrUpdateInstallation %v`, err)
	}
}

func TestUpdateInstallation(t *testing.T) {
	connectionString := os.Getenv(anhConnectionStringVariable)
	hubName := os.Getenv(anhHubVariable)
	installationId := getEnvWithFallback(anhInstallationIdVariable, anhInstallationId)

	if connectionString == "" || hubName == "" {
		t.Skip(`Skipping test due to missing environment variables`)
	}

	client, err := NewNotificationHubClientFromConnectionString(connectionString, hubName, nil)
	if client == nil || err != nil {
		t.Fatalf(`NewNotificationHubClientFromConnectionString %v`, err)
	}

	installation, err := client.GetInstallation(context.TODO(), installationId)
	if installation == nil || err != nil {
		t.Fatalf(`GetInstallation %v`, err)
	}

	updates := []InstallationPatch{
		{Op: "add", Path: "/tags", Value: "likes_dogs"},
	}

	response, err := client.UpdateInstallation(context.TODO(), installationId, updates)
	if response == nil || err != nil {
		t.Fatalf(`UpdateInstallation %v`, err)
	}
}
