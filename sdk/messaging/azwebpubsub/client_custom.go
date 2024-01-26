//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azwebpubsub

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub/internal"
	"github.com/golang-jwt/jwt"
)

const defaultExpirationTime = time.Hour

// ClientOptions contains optional settings for [Client]
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient creates a client that manages Web PubSub service
func NewClient(endpoint string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{internal.TokenScope}, nil)
	azcoreClient, err := azcore.NewClient(internal.ModuleName, internal.ModuleVersion,
		runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{
		internal: azcoreClient,
		endpoint: endpoint,
	}, nil
}

// NewClientFromConnectionString creates a Client from a connection string
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	props, err := internal.ParseConnectionString(connectionString)

	if err != nil {
		return nil, err
	}

	azcoreClient, err := azcore.NewClient(internal.ModuleName, internal.ModuleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{internal.NewWebPubSubKeyCredentialPolicy(props.AccessKey)},
	}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{
		internal: azcoreClient,
		endpoint: props.Endpoint,
		key:      &props.AccessKey,
	}, nil
}

// GenerateClientAccessURLOptions represents the options for generating a client access url
type GenerateClientAccessURLOptions struct {
	// UserID is the user ID for the client.
	UserID string

	// Roles are the roles that the connection with the generated token will have.
	// Roles give the client initial permissions to leave, join, or publish to groups when using PubSub subprotocol.
	// Possible role values:
	// - webpubsub.joinLeaveGroup: the client can join or leave any group.
	// - webpubsub.sendToGroup: the client can send messages to any group.
	// - webpubsub.joinLeaveGroup.<group>: the client can join or leave group <group>.
	// - webpubsub.sendToGroup.<group>: the client can send messages to group <group>.
	// More info: https://azure.github.io/azure-webpubsub/references/pubsub-websocket-subprotocol#permissions
	Roles []string

	// ExpirationTimeInMinutes is the number of minutes until the token expires. Default value(60 minutes) is used if the value is 0.
	ExpirationTimeInMinutes int32

	// Groups are the groups to join when the client connects.
	Groups []string
}

// GenerateClientAccessURLResponse represents the response type for the generated client access url
type GenerateClientAccessURLResponse struct {
	// The client token
	Token string
	// The base URL for the client to connect to
	BaseURL string
	// The URL client connects to with access_token query string
	URL string
}

// GenerateClientAccessURL - generate URL for the WebSocket clients
//   - hub - The hub name.
//   - options - GenerateClientAccessUrlOptions contains the optional parameters for the Client.GenerateClientAccessURL method.
func (c *Client) GenerateClientAccessURL(ctx context.Context, hub string, options *GenerateClientAccessURLOptions) (*GenerateClientAccessURLResponse, error) {
	endpoint := c.endpoint
	if hub == "" {
		return nil, errors.New("empty hub name is not allowed")
	}
	hubPath := url.PathEscape(hub)
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.New("endpoint is not a valid URL")
	}

	audience := fmt.Sprintf("%sclient/hubs/%s", parsedURL.String(), hubPath)

	parsedURL.Scheme = strings.Replace(strings.ToLower(parsedURL.Scheme), "http", "ws", 1)
	baseURL := fmt.Sprintf("%sclient/hubs/%s", parsedURL.String(), hubPath)

	var token string
	if c.key != nil {
		token, err = c.signJwtToken(audience, options)
		if err != nil {
			return nil, err
		}
	} else {
		var userId *string
		if options.UserID == "" {
			userId = nil
		} else {
			userId = &options.UserID
		}
		// Replace with your logic to generate the token using a webPubSub method
		resp, err := c.generateClientToken(ctx, hub, &GenerateClientTokenOptions{UserID: userId, Role: options.Roles, Group: options.Groups, MinutesToExpire: &options.ExpirationTimeInMinutes})
		if err != nil {
			return nil, err
		}

		token = *resp.Token
	}

	return &GenerateClientAccessURLResponse{
		Token:   token,
		BaseURL: baseURL,
		URL:     fmt.Sprintf("%s?access_token=%s", baseURL, url.QueryEscape(token)),
	}, nil
}

func (c *Client) signJwtToken(audience string, options *GenerateClientAccessURLOptions) (string, error) {
	if c.key == nil {
		return "", errors.New("key is nil")
	}
	key := []byte(*c.key)
	var exp int64

	if options == nil || options.ExpirationTimeInMinutes == 0 {
		exp = time.Now().Add(defaultExpirationTime).Unix()
	} else {
		if options.ExpirationTimeInMinutes < 0 {
			return "", errors.New("the value of ExpirationTimeInMinutes is out of range")
		}
		exp = time.Now().Add(time.Minute * time.Duration(options.ExpirationTimeInMinutes)).Unix()
	}
	claims := jwt.MapClaims{
		"aud": audience,
		"exp": exp,
	}

	if options != nil && options.UserID != "" {
		claims["sub"] = options.UserID
	}

	if options != nil && len(options.Groups) > 0 {
		claims["webpubsub.group"] = options.Groups
	}

	if options != nil && options.Roles != nil && len(options.Roles) > 0 {
		claims["role"] = options.Roles
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}
