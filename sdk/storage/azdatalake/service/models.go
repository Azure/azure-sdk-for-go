//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
)

type CreateFilesystemOptions = filesystem.CreateOptions

type DeleteFilesystemOptions = filesystem.DeleteOptions

// CORSRule - CORS is an HTTP feature that enables a web application running under one domain to access resources in another
// domain. Web browsers implement a security restriction known as same-origin policy that
// prevents a web page from calling APIs in a different domain; CORS provides a secure way to allow one domain (the origin
// domain) to call APIs in another domain.
type CORSRule = service.CORSRule

// RetentionPolicy - the retention policy which determines how long the associated data should persist.
type RetentionPolicy = service.RetentionPolicy

// Metrics - a summary of request statistics grouped by API in hour or minute aggregates for blobs.
type Metrics = service.Metrics

// Logging - Azure Analytics Logging settings.
type Logging = service.Logging

// StaticWebsite - The properties that enable an account to host a static website.
type StaticWebsite = service.StaticWebsite

// GetPropertiesOptions contains the optional parameters for the Client.GetProperties method.
type GetPropertiesOptions = service.GetPropertiesOptions

// SetPropertiesOptions provides set of options for Client.SetProperties
type SetPropertiesOptions = service.SetPropertiesOptions

// ListFilesystemsOptions contains the optional parameters for the Client.List method.
type ListFilesystemsOptions struct {
	// The number of filesystem names to retrieve. If the request does not specify the
	//  maximum number of filesystem names to retrieve, or specifies a value greater than 5,000, the server will return up to
	//  5,000 items.
	MaxResults *int32

	// A string value that identifies the portion of the list of filesystems to be
	//  returned with the next listing operation. The operation returns a marker value within the response body if the listing
	//  operation did not return all filesystem names and the value can be used with a subsequent recursive call to request
	//  the next set of filesystem names.
	Marker *string
}
