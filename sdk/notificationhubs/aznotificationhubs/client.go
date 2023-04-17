//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznotificationhubs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/notificationhubs/aznotificationhubs/auth"
)

const (
	AZNHApiVersion = "2020-06"
)

type (
	// Represents a client for interacting with the Azure Notification Hubs service.
	NotificationHubClient struct {
		hubName       string
		hostName      string
		endpointUrl   string
		tokenProvider auth.TokenProvider
		pl            runtime.Pipeline
	}

	// NotificationHubsClientOptions contains optional settings for NotificationHubClient.
	NotificationHubsClientOptions struct {
		azcore.ClientOptions
	}
)

// Creates a new Notification Hubs client with Access Policy connection string and hub name.
func NewNotificationHubClientFromConnectionString(connectionString string, hubName string, options *NotificationHubsClientOptions) (*NotificationHubClient, error) {
	if options == nil {
		options = &NotificationHubsClientOptions{}
	}

	parsedConnection, err := auth.FromConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &options.ClientOptions)
	tokenProvider := auth.NewNotificationHubsTokenProvider(parsedConnection.KeyName, parsedConnection.KeyValue)

	return &NotificationHubClient{
		hubName:       hubName,
		hostName:      parsedConnection.Endpoint,
		endpointUrl:   strings.Replace(parsedConnection.Endpoint, "sb://", "https://", -1),
		tokenProvider: tokenProvider,
		pl:            pl,
	}, nil
}

// Cancels the scheduled notification
func (n *NotificationHubClient) CancelScheduledNotification(ctx context.Context, notificationId string) (*NotificationResponse, error) {
	requestUri := fmt.Sprintf("%v%v/schedulednotifications/%v?api-version=%v", n.endpointUrl, n.hubName, notificationId, AZNHApiVersion)

	req, err := runtime.NewRequest(ctx, http.MethodDelete, requestUri)
	if err != nil {
		return nil, err
	}

	n.setRequestHeaders(req)

	res, err := n.pl.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response from Azure Notification Hubs: %v", res.StatusCode)
	}

	return createNotificationResponse(res)
}

// Schedules a notification to be sent at a specified time.
func (n *NotificationHubClient) SendScheduledNotification(ctx context.Context, notificationRequest *NotificationRequest, tagExpression string, scheduledTime time.Time) (*NotificationMessageResponse, error) {
	requestUri := fmt.Sprintf("%v%v/schedulednotifications/?api-version=%v", n.endpointUrl, n.hubName, AZNHApiVersion)

	req, err := runtime.NewRequest(ctx, http.MethodPost, requestUri)
	if err != nil {
		return nil, err
	}

	body := streaming.NopCloser(strings.NewReader(notificationRequest.Message))
	req.SetBody(body, notificationRequest.ContentType)

	for headerName, headerValue := range notificationRequest.Headers {
		req.Raw().Header.Add(headerName, headerValue)
	}

	if tagExpression != "" {
		req.Raw().Header.Add("ServiceBusNotification-Tags", tagExpression)
	}

	n.setRequestHeaders(req)
	req.Raw().Header.Add("Content-Type", notificationRequest.ContentType)
	req.Raw().Header.Add("ServiceBusNotification-Format", notificationRequest.Platform)
	req.Raw().Header.Add("ServiceBusNotification-ScheduleTime", scheduledTime.Format(time.RFC3339))

	res, err := n.pl.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 201 {
		return nil, fmt.Errorf("invalid response from Azure Notification Hubs: %v", res.StatusCode)
	}

	return createNotificationMessageResponse(res)
}

// Sends a direct notification to a device.
func (n *NotificationHubClient) SendDirectNotification(ctx context.Context, notificationRequest *NotificationRequest, deviceToken interface{}) (*NotificationMessageResponse, error) {
	return n.sendNotification(ctx, notificationRequest, deviceToken, nil)
}

// Sends a notification to a tag expression or if not specified, to all devices.
func (n *NotificationHubClient) SendNotification(ctx context.Context, notificationRequest *NotificationRequest, tagExpression string) (*NotificationMessageResponse, error) {
	return n.sendNotification(ctx, notificationRequest, nil, &tagExpression)
}

func (n *NotificationHubClient) sendNotification(ctx context.Context, notificationRequest *NotificationRequest, deviceHandle interface{}, tagExpression *string) (*NotificationMessageResponse, error) {
	requestUri := fmt.Sprintf("%v%v/messages/?api-version=%v", n.endpointUrl, n.hubName, AZNHApiVersion)
	if deviceHandle != nil {
		requestUri += "&direct=true"
	}

	req, err := runtime.NewRequest(ctx, http.MethodPost, requestUri)
	if err != nil {
		return nil, err
	}

	body := streaming.NopCloser(strings.NewReader(notificationRequest.Message))
	req.SetBody(body, notificationRequest.ContentType)

	for headerName, headerValue := range notificationRequest.Headers {
		req.Raw().Header.Add(headerName, headerValue)
	}

	if deviceHandle != nil {
		err := setDeviceHandle(deviceHandle, req)
		if err != nil {
			return nil, err
		}
	}

	if tagExpression != nil {
		req.Raw().Header.Add("ServiceBusNotification-Tags", *tagExpression)
	}

	n.setRequestHeaders(req)
	req.Raw().Header.Add("Content-Type", notificationRequest.ContentType)
	req.Raw().Header.Add("ServiceBusNotification-Format", notificationRequest.Platform)

	res, err := n.pl.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 201 {
		return nil, fmt.Errorf("invalid response from Azure Notification Hubs: %v", res.StatusCode)
	}

	return createNotificationMessageResponse(res)
}

// Gets the installation for the specified installation ID.
func (n *NotificationHubClient) GetInstallation(ctx context.Context, installationId string) (*Installation, error) {
	requestUri := fmt.Sprintf("%v%v/installations/%v?api-version=%v", n.endpointUrl, n.hubName, installationId, AZNHApiVersion)

	req, err := runtime.NewRequest(ctx, http.MethodGet, requestUri)
	if err != nil {
		return nil, err
	}

	n.setRequestHeaders(req)

	res, err := n.pl.Do(req)
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
func (n *NotificationHubClient) CreateOrUpdateInstallation(ctx context.Context, installation *Installation) (*NotificationResponse, error) {
	requestUri := fmt.Sprintf("%v%v/installations/%v?api-version=%v", n.endpointUrl, n.hubName, installation.InstallationID, AZNHApiVersion)

	installationJSON, err := json.Marshal(installation)
	if err != nil {
		return nil, err
	}

	req, err := runtime.NewRequest(ctx, http.MethodPut, requestUri)
	if err != nil {
		return nil, err
	}

	body := streaming.NopCloser(bytes.NewReader(installationJSON))
	req.SetBody(body, "application/json")

	n.setRequestHeaders(req)
	req.Raw().Header.Add("Content-Type", "application/json")

	res, err := n.pl.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response from Azure Notification Hubs: %v", res.StatusCode)
	}

	return createNotificationResponse(res)
}

// Updates an installation with the specified patches.
func (n *NotificationHubClient) UpdateInstallation(ctx context.Context, installationId string, patches []InstallationPatch) (*NotificationResponse, error) {
	requestUri := fmt.Sprintf("%v%v/installations/%v?api-version=%v", n.endpointUrl, n.hubName, installationId, AZNHApiVersion)

	patchesJSON, err := json.Marshal(patches)
	if err != nil {
		return nil, err
	}

	req, err := runtime.NewRequest(ctx, http.MethodPatch, requestUri)
	if err != nil {
		return nil, err
	}

	body := streaming.NopCloser(bytes.NewReader(patchesJSON))
	req.SetBody(body, "application/json-patch+json")

	n.setRequestHeaders(req)
	req.Raw().Header.Add("Content-Type", "application/json-patch+json")

	res, err := n.pl.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response from Azure Notification Hubs: %v", res.StatusCode)
	}

	return createNotificationResponse(res)
}

// Deletes an installation from Azure Notification Hubs.
func (n *NotificationHubClient) DeleteInstallation(ctx context.Context, installationId string) (*NotificationResponse, error) {
	requestUri := fmt.Sprintf("%v%v/installations/%v?api-version=%v", n.endpointUrl, n.hubName, installationId, AZNHApiVersion)

	req, err := runtime.NewRequest(ctx, http.MethodDelete, requestUri)
	if err != nil {
		return nil, err
	}

	n.setRequestHeaders(req)

	res, err := n.pl.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 204 {
		return nil, fmt.Errorf("invalid response from Azure Notification Hubs: %v", res.StatusCode)
	}

	return createNotificationResponse(res)
}

func (n *NotificationHubClient) setRequestHeaders(req *policy.Request) {
	sasToken, _ := n.tokenProvider.GetToken(n.hostName)

	req.Raw().Header.Add("Authorization", sasToken.Token)
}

func setDeviceHandle(deviceHandle interface{}, req *policy.Request) error {
	var (
		endpoint, p256dh, auth string
		ok                     bool
	)

	switch v := deviceHandle.(type) {
	case string:
		req.Raw().Header.Add("ServiceBusNotification-DeviceHandle", v)
	case map[string]interface{}:
		if endpoint, ok = v["endpoint"].(string); !ok {
			return fmt.Errorf("missing endpoint")
		}
		if p256dh, ok = v["p256dh"].(string); !ok {
			return fmt.Errorf("missing p256dh")
		}
		if auth, ok = v["auth"].(string); !ok {
			return fmt.Errorf("missing auth")
		}
		req.Raw().Header.Add("ServiceBusNotification-DeviceHandle", endpoint)
		req.Raw().Header.Add("p256", p256dh)
		req.Raw().Header.Add("auth", auth)
	default:
		return fmt.Errorf("invalid deviceHandle type")
	}

	return nil
}

func createNotificationMessageResponse(res *http.Response) (*NotificationMessageResponse, error) {
	var notificationId string
	correlationId := res.Header.Get("x-ms-correlation-request-id")
	trackingId := res.Header.Get("TrackingId")
	location := res.Header.Get("Location")

	if location != "" {
		parsedLocation, err := url.Parse(location)
		if err != nil {
			return nil, err
		}

		path := parsedLocation.Path
		parts := strings.Split(path, "/")
		notificationId = parts[len(parts)-1]
	}

	return &NotificationMessageResponse{
		CorrelationId:  correlationId,
		TrackingId:     trackingId,
		NotificationId: notificationId,
	}, nil
}

func createNotificationResponse(res *http.Response) (*NotificationResponse, error) {
	correlationId := res.Header.Get("x-ms-correlation-request-id")
	trackingId := res.Header.Get("TrackingId")
	location := res.Header.Get("Location")

	return &NotificationResponse{
		CorrelationId: correlationId,
		TrackingId:    trackingId,
		Location:      location,
	}, nil
}
