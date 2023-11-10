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

// ClientOptions contains optional settings for [Client]
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient creates a client that manages Web PubSub service
func NewClient(endpoint string, hub string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{internal.TokenScope}, nil)
	azcoreClient, err := azcore.NewClient(internal.ModuleName+".Client", internal.ModuleVersion,
		runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{
		internal: azcoreClient,
		endpoint: endpoint,
		hub:      hub,
	}, nil
}

// NewClientFromConnectionString creates a Client from a connection string
func NewClientFromConnectionString(connectionString string, hub string, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	props, err := internal.ParseConnectionString(connectionString)

	if err != nil {
		return nil, err
	}

	azcoreClient, err := azcore.NewClient(internal.ModuleName+".Client", internal.ModuleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{internal.NewWebPubSubKeyCredentialPolicy(props.AccessKey)},
	}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{
		internal: azcoreClient,
		endpoint: props.Endpoint,
		hub:      hub,
		key:      &props.AccessKey,
	}, nil
}

// GenerateClientAccessUrlOptions represents the options for generating a client access url
type GenerateClientAccessUrlOptions struct {
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

	// ExpirationTimeInMinutes is the number of minutes until the token expires.
	ExpirationTimeInMinutes *int32

	// Groups are the groups to join when the client connects.
	Groups []string
}

// GenerateClientAccessUrlResponse represents the response type for the generated client access url
type GenerateClientAccessUrlResponse struct {
	// The client token
	Token string
	// The base URL for the client to connect to
	BaseURL string
	// The URL client connects to with access_token query string
	URL string
}

// GenerateClientAccessUrl - generate URL for the WebSocket clients
//   - options - GenerateClientAccessUrlOptions contains the optional parameters for the Client.GenerateClientAccessUrl method.
func (c *Client) GenerateClientAccessUrl(ctx context.Context, options *GenerateClientAccessUrlOptions) (*GenerateClientAccessUrlResponse, error) {
	endpoint := c.endpoint
	hubName := c.hub
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.New("endpoint is not a valid URL")
	}

	parsedURL.Scheme = strings.Replace(strings.ToLower(parsedURL.Scheme), "http", "ws", 1)

	clientEndpoint := parsedURL.String()
	baseURL := fmt.Sprintf("%s/client/hubs/%s", clientEndpoint, hubName)

	var token string
	if c.key != nil {
		token, err = c.signJwtToken(baseURL, options)
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
		resp, err := c.generateClientToken(ctx, hubName, &ClientGenerateClientTokenOptions{UserID: userId, Role: options.Roles, Group: options.Groups, MinutesToExpire: options.ExpirationTimeInMinutes})
		if err != nil {
			return nil, err
		}

		token = *resp.Token
	}

	return &GenerateClientAccessUrlResponse{
		Token:   token,
		BaseURL: baseURL,
		URL:     fmt.Sprintf("%s?access_token=%s", baseURL, url.QueryEscape(token)),
	}, nil
}

const DefaultExpirationTime = time.Hour

func (c *Client) signJwtToken(baseURL string, options *GenerateClientAccessUrlOptions) (string, error) {
	if c.key == nil {
		return "", errors.New("key is nil")
	}
	key := []byte(*c.key)
	var exp int64
	if options == nil || options.ExpirationTimeInMinutes == nil {
		exp = time.Now().Add(DefaultExpirationTime).Unix()
	} else {
		exp = time.Now().Add(time.Minute * time.Duration(*options.ExpirationTimeInMinutes)).Unix()
	}
	claims := jwt.MapClaims{
		"aud": baseURL,
		"exp": exp,
	}

	if options != nil && options.UserID != "" {
		claims["sub"] = options.UserID
	}

	if options != nil && options.Groups != nil && len(options.Groups) > 0 {
		claims["webpubsub.group"] = options.Groups
	}

	if options != nil && options.Roles != nil && len(options.Roles) > 0 {
		claims["role"] = options.Roles
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}
