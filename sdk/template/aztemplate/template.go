//go:build go1.17
// +build go1.17

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztemplate

import (
	generated "github.com/Azure/azure-sdk-for-go/sdk/template/aztemplate/internal/packageinfo"
)

type TemplateClient struct{}

func NewTemplateClient() (TemplateClient, error) {

	return TemplateClient{}, nil

}

func ClientVersion() string {

	return packageinfo.Version

}
