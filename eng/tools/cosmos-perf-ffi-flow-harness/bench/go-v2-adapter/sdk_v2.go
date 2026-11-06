// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build go_v2

package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync/atomic"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
)

const experimentalReactorShards = 1

type v2Client struct {
	runtime    *azcosmos.Runtime
	clients    []*azcosmos.Client
	containers []*azcosmos.ContainerClient
	next       atomic.Uint64
}

func newSDKClient(ctx context.Context, options options) (sdkClient, error) {
	credential, err := azcosmos.NewKeyCredential(options.key)
	if err != nil {
		return nil, err
	}
	runtime, err := azcosmos.NewRuntime(nil)
	if err != nil {
		return nil, err
	}
	result := &v2Client{runtime: runtime}
	for range experimentalReactorShards {
		client, err := azcosmos.NewClientWithKey(ctx, options.endpoint, credential, &azcosmos.ClientOptions{
			Runtime: runtime,
			DiagnosticsVerbosity: azcosmos.DiagnosticsVerbositySummary,
		})
		if err != nil {
			_ = result.Close()
			return nil, err
		}
		container, err := client.NewContainer(ctx, options.database, options.container)
		if err != nil {
			_ = client.Close()
			_ = result.Close()
			return nil, err
		}
		result.clients = append(result.clients, client)
		result.containers = append(result.containers, container)
	}
	return result, nil
}

func (client *v2Client) Read(ctx context.Context, pk string, id string) (operationResponse, error) {
	index := int((client.next.Add(1) - 1) % uint64(len(client.containers)))
	response, err := client.containers[index].ReadItem(ctx, azcosmos.NewPartitionKeyString(pk), id, nil)
	if err != nil {
		return operationResponse{}, err
	}
	return operationResponse{
		Body:       response.Value,
		StatusCode: response.StatusCode,
		Attempts:   response.AttemptCount,
	}, nil
}

func (client *v2Client) Create(ctx context.Context, pk string, id string, body []byte) (operationResponse, error) {
	response, err := client.containers[0].CreateItem(ctx, azcosmos.NewPartitionKeyString(pk), id, json.RawMessage(body), nil)
	if err != nil {
		return operationResponse{}, err
	}
	return operationResponse{StatusCode: response.StatusCode, Attempts: response.AttemptCount}, nil
}

func (client *v2Client) Upsert(ctx context.Context, pk string, id string, body []byte) (operationResponse, error) {
	response, err := client.containers[0].UpsertItem(ctx, azcosmos.NewPartitionKeyString(pk), id, json.RawMessage(body), nil)
	if err != nil {
		return operationResponse{}, err
	}
	return operationResponse{StatusCode: response.StatusCode, Attempts: response.AttemptCount}, nil
}

func (client *v2Client) Delete(ctx context.Context, pk string, id string) (operationResponse, error) {
	response, err := client.containers[0].DeleteItem(ctx, azcosmos.NewPartitionKeyString(pk), id, nil)
	if err != nil {
		return operationResponse{}, err
	}
	return operationResponse{StatusCode: response.StatusCode, Attempts: response.AttemptCount}, nil
}

func (client *v2Client) Close() error {
	var firstErr error
	for _, sdkClient := range client.clients {
		if err := sdkClient.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if client.runtime != nil {
		if err := client.runtime.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func isNotFound(err error) bool {
	var cosmosError *azcosmos.Error
	return errors.As(err, &cosmosError) && cosmosError.StatusCode == http.StatusNotFound
}
