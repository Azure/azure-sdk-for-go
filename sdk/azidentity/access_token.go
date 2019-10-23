// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"encoding/json"
	"hash/fnv"
	"time"
)

// AccessToken is used to set and maintain tokens for authentication
type AccessToken struct {
	Token     string      `json:"access_token"`
	ExpiresIn json.Number `json:"expires_in"`
	ExpiresOn time.Time
}

// SetToken assigns the given string into the Token field of the type
// CP: this isnt necessary if Token is exported?
func (c *AccessToken) SetToken(token string) {
	c.Token = token
}

// GetToken provides a way to check the unexported token variable
// CP: This also isnt necessary if this is exported
func (c *AccessToken) GetToken() string {
	return c.Token
}

// SetExpiresOn assigns the given integer to the ExpiresOn field of the AccessToken
// The int should be the time in seconds that the token expires in
// CP: check this implementation for type
func (c *AccessToken) SetExpiresOn() {
	t, err := c.ExpiresIn.Int64()
	if err != nil {
		return
	}
	c.ExpiresOn = time.Now().Add(time.Second * time.Duration(t)).UTC()

}

// GetExpiresOn provides a way to check the unexported ExpiresOn variable
// ExpiresOn is now exported for unmarshaling purposes, this func might not be necessary anymore
func (c *AccessToken) GetExpiresOn() time.Time {
	return c.ExpiresOn
}

// Equals determines whether the AccessToken that is being used is the same as the one that is currently being used
func (c *AccessToken) Equals(accessToken AccessToken) bool {
	if accessToken == *c {
		return accessToken.GetExpiresOn() == c.GetExpiresOn() && accessToken.GetToken() == c.GetToken()
	}
	return false
}

// NewAccessToken constructs the AccessToken type
func NewAccessToken(accessToken string, expiresOn json.Number) *AccessToken {
	c := &AccessToken{Token: accessToken, ExpiresIn: expiresOn}
	return c
}

// GetHashCode returns a uint32 hash of the access token
func (c *AccessToken) GetHashCode() uint32 {
	h := fnv.New32a()
	h.Write([]byte(c.Token))
	return h.Sum32()
}
