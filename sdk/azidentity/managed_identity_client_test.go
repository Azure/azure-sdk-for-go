// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type userAgentValidatingPolicy struct {
	t     *testing.T
	appID string
}

func (p userAgentValidatingPolicy) Do(req *policy.Request) (*http.Response, error) {
	expected := "azsdk-go-" + component + "/" + version
	if p.appID != "" {
		expected = p.appID + " " + expected
	}
	if ua := req.Raw().Header.Get("User-Agent"); !strings.HasPrefix(ua, expected) {
		p.t.Fatalf("unexpected User-Agent %s", ua)
	}
	return req.Next()
}

func TestManagedIdentityCredential_UserAgent(t *testing.T) {
	for _, appID := range []string{"", "customvalue"} {
		options := ManagedIdentityCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: &mockSTS{}, PerCallPolicies: []policy.Policy{userAgentValidatingPolicy{t: t, appID: appID}},
			},
		}
		c, err := NewManagedIdentityCredential(&options)
		if err != nil {
			t.Fatal(err)
		}
		_, err = c.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
		if err != nil {
			t.Fatal(err)
		}
	}
}
