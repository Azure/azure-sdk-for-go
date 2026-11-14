// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

type ruleNameTransport struct {
	t        *testing.T
	expected string
	updating bool
}

func (transport ruleNameTransport) Do(req *http.Request) (*http.Response, error) {
	require.Equal(transport.t, http.MethodPut, req.Method)
	require.Equal(transport.t, "/topic/Subscriptions/subscription/Rules/"+transport.expected, req.URL.Path)
	if transport.updating {
		require.Equal(transport.t, "*", req.Header.Get("If-Match"))
	} else {
		require.Empty(transport.t, req.Header.Get("If-Match"))
	}
	body, err := io.ReadAll(req.Body)
	require.NoError(transport.t, err)
	var envelope struct {
		Content struct {
			Rule struct {
				Name string `xml:"Name"`
			} `xml:"RuleDescription"`
		} `xml:"content"`
	}
	require.NoError(transport.t, xml.Unmarshal(body, &envelope))
	require.Equal(transport.t, transport.expected, envelope.Content.Rule.Name)

	reply := `<entry xmlns="http://www.w3.org/2005/Atom"><content type="application/xml"><RuleDescription xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><Filter i:type="TrueFilter"><SqlExpression>1=1</SqlExpression></Filter><Action i:type="EmptyRuleAction"/><Name>` + transport.expected + `</Name></RuleDescription></content></entry>`
	return &http.Response{
		StatusCode: http.StatusCreated,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(reply)),
		Request:    req,
	}, nil
}

func TestAdminClient_RuleNameRequest(t *testing.T) {
	for _, tc := range []struct {
		name     string
		value    *string
		expected string
	}{
		{name: "omitted", expected: "$Default"},
		{name: "empty", value: to.Ptr(""), expected: "$Default"},
		{name: "explicitDefault", value: to.Ptr("$Default"), expected: "$Default"},
		{name: "custom", value: to.Ptr("custom-rule"), expected: "custom-rule"},
	} {
		for _, updating := range []bool{false, true} {
			operation := "create"
			if updating {
				operation = "update"
			}
			t.Run(operation+"/"+tc.name, func(t *testing.T) {
				client, err := NewClient("fake.servicebus.windows.net", fakeAPIVersionCredential{}, &ClientOptions{
					ClientOptions: azcore.ClientOptions{
						Transport: ruleNameTransport{t: t, expected: tc.expected, updating: updating},
					},
				})
				require.NoError(t, err)
				if updating {
					properties := RuleProperties{}
					if tc.value != nil {
						properties.Name = *tc.value
					}
					response, err := client.UpdateRule(context.Background(), "topic", "subscription", properties)
					require.NoError(t, err)
					require.Equal(t, tc.expected, response.Name)
				} else {
					response, err := client.CreateRule(context.Background(), "topic", "subscription", &CreateRuleOptions{Name: tc.value})
					require.NoError(t, err)
					require.Equal(t, tc.expected, response.Name)
				}
			})
		}
	}
}
