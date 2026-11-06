// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build go_v1

package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type v1Client struct {
	client    *azcosmos.Client
	container *azcosmos.ContainerClient
}

func newSDKClient(_ context.Context, options options) (sdkClient, error) {
	credential, err := azcosmos.NewKeyCredential(options.key)
	if err != nil {
		return nil, err
	}
	client, err := azcosmos.NewClientWithKey(options.endpoint, credential, nil)
	if err != nil {
		return nil, err
	}
	container, err := client.NewContainer(options.database, options.container)
	if err != nil {
		client.Close()
		return nil, err
	}
	return &v1Client{
		client:    client,
		container: container,
	}, nil
}

func (client *v1Client) Read(ctx context.Context, pk string, id string) (operationResponse, error) {
	response, err := client.container.ReadItem(ctx, azcosmos.NewPartitionKeyString(pk), id, nil)
	if err != nil {
		return operationResponse{}, err
	}
	result := operationResponse{Body: response.Value, StatusCode: response.RawResponse.StatusCode, Attempts: 1}
	if response.RawResponse.Header.Get("x-ms-gatewayversion") != "" {
		result.TransportKind = "gateway-v1"
	}
	result.HTTPVersion = response.RawResponse.Proto
	return result, nil
}

func (client *v1Client) Create(ctx context.Context, pk string, _ string, body []byte) (operationResponse, error) {
	response, err := client.container.CreateItem(ctx, azcosmos.NewPartitionKeyString(pk), body, nil)
	if err != nil {
		return operationResponse{}, err
	}
	return operationResponse{StatusCode: response.RawResponse.StatusCode, Attempts: 1}, nil
}

func (client *v1Client) Upsert(ctx context.Context, pk string, _ string, body []byte) (operationResponse, error) {
	response, err := client.container.UpsertItem(ctx, azcosmos.NewPartitionKeyString(pk), body, nil)
	if err != nil {
		return operationResponse{}, err
	}
	return operationResponse{StatusCode: response.RawResponse.StatusCode, Attempts: 1}, nil
}

func (client *v1Client) Delete(ctx context.Context, pk string, id string) (operationResponse, error) {
	response, err := client.container.DeleteItem(ctx, azcosmos.NewPartitionKeyString(pk), id, nil)
	if err != nil {
		return operationResponse{}, err
	}
	return operationResponse{StatusCode: response.RawResponse.StatusCode, Attempts: 1}, nil
}

func (client *v1Client) Close() error {
	client.client.Close()
	return nil
}

func isNotFound(err error) bool {
	var responseError *azcore.ResponseError
	return errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound
}
