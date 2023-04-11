//go:build go1.20
// +build go1.20

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznotificationhubs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/notificationhubs/aznotificationhubs/auth"
)

const (
	AZNHApiVersion = "2020-06"
)

// Creates a tag expression
func CreateTagExpression(tags []string) string {
	return strings.Join(tags, "||")
}

// Creates a new Notification Hubs client with Access Policy connection string and hub name.
func NewNotificationHubClientWithConnectionString(connectionString string, hubName string) (*NotificationHubClient, error) {
	parsedConnection, err := auth.FromConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	tokenProvider := auth.NewNotificationHubsTokenProvider(parsedConnection.KeyName, parsedConnection.KeyValue)

	return &NotificationHubClient{
		hubName:       hubName,
		hostName:      parsedConnection.Endpoint,
		endpointUrl:   strings.Replace(parsedConnection.Endpoint, "sb://", "https://", -1),
		tokenProvider: tokenProvider,
	}, nil
}

// Cancels the scheduled notification
func (n *NotificationHubClient) CancelScheduledNotification(notificationId string) (*NotificationResponse, error) {
	requestUri := fmt.Sprintf("%v%v/schedulednotifications/%v?api-version=%v", n.endpointUrl, n.hubName, notificationId, AZNHApiVersion)

	client := &http.Client{Timeout: time.Second * 15}
	req, err := http.NewRequest(http.MethodDelete, requestUri, nil)
	if err != nil {
		return nil, err
	}

	n.addRequestHeaders(req)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return n.createNotificationResponse(res)
}

// Schedules a notification to be sent at a specified time.
func (n *NotificationHubClient) SendScheduledNotification(notificationRequest *NotificationRequest, tagExpression string, scheduledTime time.Time) (*NotificationMessageResponse, error) {
	requestUri := fmt.Sprintf("%v%v/schedulednotifications/?api-version=%v", n.endpointUrl, n.hubName, AZNHApiVersion)

	messageBody := []byte(notificationRequest.Message)

	client := &http.Client{Timeout: time.Second * 15}
	req, err := http.NewRequest(http.MethodPost, requestUri, bytes.NewBuffer(messageBody))
	if err != nil {
		return nil, err
	}

	for headerName, headerValue := range notificationRequest.Headers {
		req.Header.Add(headerName, headerValue)
	}

	if tagExpression != "" {
		req.Header.Add("ServiceBusNotification-Tags", tagExpression)
	}

	n.addRequestHeaders(req)
	req.Header.Add("Content-Type", notificationRequest.ContentType)
	req.Header.Add("ServiceBusNotification-Format", notificationRequest.Platform)
	req.Header.Add("ServiceBusNotification-ScheduleTime", scheduledTime.Format(time.RFC3339))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 201 {
		return nil, fmt.Errorf("invalid response from Azure Notification Hubs: %v", res.StatusCode)
	}

	return n.createNotificationMessageResponse(res)
}

// Sends a direct notification to a device.
func (n *NotificationHubClient) SendDirectNotification(notificationRequest *NotificationRequest, deviceToken string) (*NotificationMessageResponse, error) {
	return n.sendNotification(notificationRequest, &deviceToken, nil)
}

// Sends a notification to a tag expression or if not specified, to all devices.
func (n *NotificationHubClient) SendNotification(notificationRequest *NotificationRequest, tagExpression string) (*NotificationMessageResponse, error) {
	return n.sendNotification(notificationRequest, nil, &tagExpression)
}

func (n *NotificationHubClient) sendNotification(notificationRequest *NotificationRequest, deviceToken *string, tagExpression *string) (*NotificationMessageResponse, error) {
	requestUri := fmt.Sprintf("%v%v/messages/?api-version=%v", n.endpointUrl, n.hubName, AZNHApiVersion)
	if deviceToken != nil {
		requestUri += "&direct=true"
	}

	messageBody := []byte(notificationRequest.Message)

	client := &http.Client{Timeout: time.Second * 15}
	req, err := http.NewRequest(http.MethodPost, requestUri, bytes.NewBuffer(messageBody))
	if err != nil {
		return nil, err
	}

	for headerName, headerValue := range notificationRequest.Headers {
		req.Header.Add(headerName, headerValue)
	}

	if deviceToken != nil {
		req.Header.Add("ServiceBusNotification-DeviceHandle", *deviceToken)
	}

	if tagExpression != nil {
		req.Header.Add("ServiceBusNotification-Tags", *tagExpression)
	}

	n.addRequestHeaders(req)
	req.Header.Add("Content-Type", notificationRequest.ContentType)
	req.Header.Add("ServiceBusNotification-Format", notificationRequest.Platform)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 201 {
		return nil, fmt.Errorf("invalid response from Azure Notification Hubs: %v", res.StatusCode)
	}

	return n.createNotificationMessageResponse(res)
}

// Gets the installation for the specified installation ID.
func (n *NotificationHubClient) GetInstallation(installationId string) (*Installation, error) {
	requestUri := fmt.Sprintf("%v%v/installations/%v?api-version=%v", n.endpointUrl, n.hubName, installationId, AZNHApiVersion)

	client := &http.Client{Timeout: time.Second * 15}
	req, err := http.NewRequest(http.MethodGet, requestUri, nil)
	if err != nil {
		return nil, err
	}

	n.addRequestHeaders(req)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response from Azure Notification Hubs: %v", res.StatusCode)
	}

	defer res.Body.Close()

	var installation Installation
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&installation)
	if err != nil {
		return nil, err
	}

	return &installation, nil
}

// Creates or updates an installation.
func (n *NotificationHubClient) CreateOrUpdateInstallation(installation *Installation) (*NotificationResponse, error) {
	requestUri := fmt.Sprintf("%v%v/installations/%v?api-version=%v", n.endpointUrl, n.hubName, installation.InstallationId, AZNHApiVersion)

	installationJSON, err := json.Marshal(installation)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: time.Second * 15}
	req, err := http.NewRequest(http.MethodPut, requestUri, bytes.NewBuffer(installationJSON))
	if err != nil {
		return nil, err
	}

	n.addRequestHeaders(req)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response from Azure Notification Hubs: %v", res.StatusCode)
	}

	return n.createNotificationResponse(res)
}

// Updates an installation with the specified patches.
func (n *NotificationHubClient) UpdateInstallation(installationId string, patches []InstallationPatch) (*NotificationResponse, error) {
	requestUri := fmt.Sprintf("%v%v/installations/%v?api-version=%v", n.endpointUrl, n.hubName, installationId, AZNHApiVersion)

	patchesJSON, err := json.Marshal(patches)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: time.Second * 15}
	req, err := http.NewRequest(http.MethodPatch, requestUri, bytes.NewBuffer(patchesJSON))
	if err != nil {
		return nil, err
	}

	n.addRequestHeaders(req)
	req.Header.Add("Content-Type", "application/json-patch+json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response from Azure Notification Hubs: %v", res.StatusCode)
	}

	return n.createNotificationResponse(res)
}

func generateUserAgent() string {
	return fmt.Sprintf("NHub/%v} (api-origin=GoSDK;)", AZNHApiVersion)
}

func (n *NotificationHubClient) addRequestHeaders(req *http.Request) {
	sasToken, _ := n.tokenProvider.GetToken(n.hostName)

	req.Header.Add("Authorization", sasToken.Token)
	req.Header.Set("User-Agent", generateUserAgent())
}

func (*NotificationHubClient) createNotificationMessageResponse(res *http.Response) (*NotificationMessageResponse, error) {
	var notificationId string
	correlationId := res.Header.Get("x-ms-correlation-request-id")
	trackingId := res.Header.Get("TrackingId")
	location := res.Header.Get("Location")

	if location != "" {
		re := regexp.MustCompile(`/messages/(\w+-\d+-\w+-\d{2})\?`)
		match := re.FindStringSubmatch(location)
		if len(match) == 2 {
			notificationId = match[1]
		}
	}

	return &NotificationMessageResponse{
		CorrelationId:  correlationId,
		TrackingId:     trackingId,
		NotificationId: notificationId,
	}, nil
}

func (*NotificationHubClient) createNotificationResponse(res *http.Response) (*NotificationResponse, error) {
	correlationId := res.Header.Get("x-ms-correlation-request-id")
	trackingId := res.Header.Get("TrackingId")
	location := res.Header.Get("Location")

	return &NotificationResponse{
		CorrelationId: correlationId,
		TrackingId:    trackingId,
		Location:      location,
	}, nil
}
