// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package audience

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const aadAudienceErrorCode = "AADSTS500011"

type AudienceErrorHandlingPolicy struct {
	AudienceConfigured bool
}

func NewAudienceErrorHandlingPolicy(audienceConfigured bool) *AudienceErrorHandlingPolicy {
	return &AudienceErrorHandlingPolicy{
		AudienceConfigured: audienceConfigured,
	}
}

func (p *AudienceErrorHandlingPolicy) Do(req *policy.Request) (*http.Response, error) {
	resp, err := req.Next()
	if err != nil {
		if strings.Contains(err.Error(), aadAudienceErrorCode) {
			if p.AudienceConfigured {
				return nil, errors.New("unable to authenticate to Azure App Configuration. An incorrect token audience was provided. Please set ClientOptions.Cloud for the target cloud . For details on how to configure the authentication token audience visit https://aka.ms/appconfig/client-token-audience and examples in https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/data/azappconfig/cloud_config.go")
			} else {
				return nil, errors.New("unable to authenticate to Azure App Configuration. No authentication token audience was provided. Please set ClientOptions.Cloud for the target cloud. For details on how to configure the authentication token audience visit https://aka.ms/appconfig/client-token-audience and examples in https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/data/azappconfig/cloud_config.go")
			}
		}
		return nil, err
	}

	return resp, nil
}
