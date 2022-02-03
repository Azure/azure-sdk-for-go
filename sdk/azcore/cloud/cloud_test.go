//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cloud

import (
	"fmt"
	"testing"
)

func TestPreconfiguredClouds(t *testing.T) {
	for _, config := range []Configuration{AzureChina, AzurePublicCloud, AzureGovernment} {
		if arm, ok := config.Services[ResourceManager]; !ok || arm.Audience == "" || arm.Endpoint == "" {
			t.Fatalf("%s is missing ARM configuration", config.Name)
		}
	}
}

func TestGetConfigurations_Stack(t *testing.T) {
	audience := "https://management.adfs.selfhost.local/d48c9ef1-e46d-4148-b41e"
	loginEndpoint := "https://adfs.azurestack.contoso.com/adfs"
	name := "AzureStack-User-d48c9ef1-e46d-4148-b41e"
	// real responses look like this (no "resourceManager")
	metadata := []byte(fmt.Sprintf(`[
    {
        "portal": "https://portal.azurestack.contoso.com/",
        "authentication": {
            "loginEndpoint": "%s",
            "audiences": ["%s"]
        },
        "graphAudience": "https://graph.azurestack.contoso.com/",
        "graph": "https://graph.azurestack.contoso.com/",
        "name": "%s",
        "suffixes": {
            "keyVaultDns": "vault.azurestack.contoso.com",
            "storage": "azurestack.contoso.com"
        },
        "gallery": "https://providers.selfhost.local:30016/"
    }
]`, loginEndpoint, audience, name))
	clouds, err := getConfigurationsFromMetadata(metadata)
	if err != nil {
		t.Fatal(err)
	}
	if len(clouds) != 1 {
		t.Fatalf("unmarshaled configuration for %d clouds; expected 1", len(clouds))
	}
	if v := clouds[0].Name; v != name {
		t.Fatalf(`unexpected Name "%s"`, v)
	}
	if v := clouds[0].LoginEndpoint; v != loginEndpoint {
		t.Fatalf(`unexpected LoginEndpoint "%s"`, v)
	}
	if len(clouds[0].Services) != 1 {
		t.Fatal("expected one service configuration")
	}
	if v, ok := clouds[0].Services[ResourceManager]; !ok {
		t.Fatal("no Resource Manager configuration")
	} else if v.Audience != audience {
		t.Fatalf(`unexpected audience "%v"`, v.Audience)
	}
}
