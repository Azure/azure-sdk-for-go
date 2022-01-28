//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cloud

import "encoding/json"

// WellKnownClouds contains configuration settings for public Azure clouds.
var WellKnownClouds = map[Name]Configuration{}

type Name string

const (
	AzureChina       Name = "AzureChina"
	AzureGovernment  Name = "AzureGovernment"
	AzurePublicCloud Name = "AzurePublicCloud"
)

func init() {
	clouds, _ := getConfigurationsFromMetadata(armMetadata)
	for _, cloud := range clouds {
		switch cloud.Name {
		case "AzureChinaCloud":
			WellKnownClouds[AzureChina] = cloud
		case "AzureCloud":
			WellKnownClouds[AzurePublicCloud] = cloud
		case "AzureUSGovernment":
			WellKnownClouds[AzureGovernment] = cloud
		}
	}
}

type ServiceName string

const ResourceManager ServiceName = "resourceManager"

type ServiceConfiguration struct {
	Audiences []string
	Endpoint  string
	Suffix    string
}

type Configuration struct {
	LoginEndpoint string
	Name          string
	Services      map[ServiceName]ServiceConfiguration
}

// getConfigurationsFromMetadata unmarshals Configuration objects from ARM endpoint metadata
func getConfigurationsFromMetadata(b []byte) ([]Configuration, error) {
	var raw []cloudConfiguration
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}

	ret := make([]Configuration, 0, len(raw))
	for _, r := range raw {
		c := Configuration{
			LoginEndpoint: r.Authentication.LoginEndpoint,
			Name:          r.Name,
			Services: map[ServiceName]ServiceConfiguration{
				ResourceManager: {
					Audiences: r.Authentication.Audiences, Endpoint: r.ResourceManager,
				},
			},
		}
		ret = append(ret, c)
	}

	return ret, nil
}

// types for unmarshaling ARM metadata
type authentication struct {
	Audiences        []string `json:"audiences"`
	IdentityProvider string   `json:"identityProvider,omitempty"`
	LoginEndpoint    string   `json:"loginEndpoint"`
	Tenant           string   `json:"tenant,omitempty"`
}

type cloudConfiguration struct {
	Authentication  authentication    `json:"authentication,omitempty"`
	Name            string            `json:"name,omitempty"`
	ResourceManager string            `json:"resourceManager,omitempty"`
	Suffixes        map[string]string `json:"suffixes,omitempty"`
}

// https://management.azure.com/metadata/endpoints?api-version=2019-05-01 (minus Azure Germany)
var armMetadata = []byte(`[
    {
        "portal": "https://portal.azure.com",
        "authentication": {
            "loginEndpoint": "https://login.microsoftonline.com/",
            "audiences": [
                "https://management.core.windows.net/",
                "https://management.azure.com/"
            ],
            "tenant": "common",
            "identityProvider": "AAD"
        },
        "media": "https://rest.media.azure.net",
        "graphAudience": "https://graph.windows.net/",
        "graph": "https://graph.windows.net/",
        "name": "AzureCloud",
        "suffixes": {
            "azureDataLakeStoreFileSystem": "azuredatalakestore.net",
            "acrLoginServer": "azurecr.io",
            "sqlServerHostname": "database.windows.net",
            "azureDataLakeAnalyticsCatalogAndJob": "azuredatalakeanalytics.net",
            "keyVaultDns": "vault.azure.net",
            "storage": "core.windows.net",
            "azureFrontDoorEndpointSuffix": "azurefd.net"
        },
        "batch": "https://batch.core.windows.net/",
        "resourceManager": "https://management.azure.com/",
        "vmImageAliasDoc": "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/arm-compute/quickstart-templates/aliases.json",
        "activeDirectoryDataLake": "https://datalake.azure.net/",
        "sqlManagement": "https://management.core.windows.net:8443/",
        "gallery": "https://gallery.azure.com/"
    },
    {
        "portal": "https://portal.azure.cn",
        "authentication": {
            "loginEndpoint": "https://login.chinacloudapi.cn",
            "audiences": [
                "https://management.core.chinacloudapi.cn",
                "https://management.chinacloudapi.cn"
            ],
            "tenant": "common",
            "identityProvider": "AAD"
        },
        "media": "https://rest.media.chinacloudapi.cn",
        "graphAudience": "https://graph.chinacloudapi.cn",
        "graph": "https://graph.chinacloudapi.cn",
        "name": "AzureChinaCloud",
        "suffixes": {
            "acrLoginServer": "azurecr.cn",
            "sqlServerHostname": "database.chinacloudapi.cn",
            "keyVaultDns": "vault.azure.cn",
            "storage": "core.chinacloudapi.cn",
            "azureFrontDoorEndpointSuffix": ""
        },
        "batch": "https://batch.chinacloudapi.cn",
        "resourceManager": "https://management.chinacloudapi.cn",
        "vmImageAliasDoc": "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/arm-compute/quickstart-templates/aliases.json",
        "sqlManagement": "https://management.core.chinacloudapi.cn:8443",
        "gallery": "https://gallery.chinacloudapi.cn"
    },
    {
        "portal": "https://portal.azure.us",
        "authentication": {
            "loginEndpoint": "https://login.microsoftonline.us",
            "audiences": [
                "https://management.core.usgovcloudapi.net",
                "https://management.usgovcloudapi.net"
            ],
            "tenant": "common",
            "identityProvider": "AAD"
        },
        "media": "https://rest.media.usgovcloudapi.net",
        "graphAudience": "https://graph.windows.net",
        "graph": "https://graph.windows.net",
        "name": "AzureUSGovernment",
        "suffixes": {
            "acrLoginServer": "azurecr.us",
            "sqlServerHostname": "database.usgovcloudapi.net",
            "keyVaultDns": "vault.usgovcloudapi.net",
            "storage": "core.usgovcloudapi.net",
            "azureFrontDoorEndpointSuffix": ""
        },
        "batch": "https://batch.core.usgovcloudapi.net",
        "resourceManager": "https://management.usgovcloudapi.net",
        "vmImageAliasDoc": "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/arm-compute/quickstart-templates/aliases.json",
        "sqlManagement": "https://management.core.usgovcloudapi.net:8443",
        "gallery": "https://gallery.usgovcloudapi.net"
    }
]`)
