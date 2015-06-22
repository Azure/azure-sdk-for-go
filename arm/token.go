package arm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	oauthUrl = "https://login.microsoftonline.com/%s/oauth2/%s?api-version=1.0"
)

var expirationBase time.Time

func init() {
	expirationBase, _ = time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
}

type Token struct {
	AccessToken string `json:"access_token"`

	ExpiresIn string `json:"expires_in"`
	ExpiresOn string `json:"expires_on"`
	NotBefore string `json:"not_before"`

	Resource string `json:"resource"`
	Type     string `json:"token_type"`
}

func NewServicePrincipalToken(id string, secret string, tenentId string) (*Token, error) {
	v := url.Values{}
	v.Set("client_id", id)
	v.Set("client_secret", secret)
	v.Set("grant_type", "client_credentials")
	v.Set("resource", "https://management.azure.com/")

	u := fmt.Sprintf(oauthUrl, tenentId, "token")

	response, err := http.PostForm(u, v)
	if err != nil {
		return nil, err
	}
	defer (func() { response.Body.Close() })()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("Token request failed with status code %v", response.StatusCode)
	}

	token := &Token{}
	err = (json.NewDecoder(response.Body)).Decode(token)
	if err != nil {
		return nil, fmt.Errorf("Failed to deserialize Token (%v)", err)
	}

	return token, nil
}

func (t Token) Expires() time.Time {
	s, _ := strconv.Atoi(t.ExpiresOn)
	return expirationBase.Add(time.Duration(s) * time.Second)
}

func (t Token) IsExpired() bool {
	return t.WillExpireIn(0)
}

func (t Token) WillExpireIn(seconds int) bool {
	return !t.Expires().After(time.Now().Add(time.Duration(seconds) * time.Second))
}
