// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const azureAuthorityHost = "AZURE_AUTHORITY_HOST"

// AuthorityHost is the base URL for Azure Active Directory
type AuthorityHost string

const (
	// AzureChina is a global constant to use in order to access the Azure China cloud.
	AzureChina AuthorityHost = "https://login.chinacloudapi.cn/"
	// AzureGovernment is a global constant to use in order to access the Azure Government cloud.
	AzureGovernment AuthorityHost = "https://login.microsoftonline.us/"
	// AzurePublicCloud is a global constant to use in order to access the Azure public cloud.
	AzurePublicCloud AuthorityHost = "https://login.microsoftonline.com/"
)

// defaultSuffix is the default AADv2 scope
const defaultSuffix = "/.default"

const (
	headerUserAgent   = "User-Agent"
	headerURLEncoded  = "application/x-www-form-urlencoded"
	headerMetadata    = "Metadata"
	headerContentType = "Content-Type"
)

const tenantIDValidationErr = "Invalid tenantID provided. You can locate your tenantID by following the instructions listed here: https://docs.microsoft.com/partner-center/find-ids-and-domain-names."

var (
	successStatusCodes = [2]int{
		http.StatusOK,      // 200
		http.StatusCreated, // 201
	}
)

type tokenResponse struct {
	token        *azcore.AccessToken
	refreshToken string
}

// setAuthorityHost initializes the authority host for credentials.
func setAuthorityHost(authorityHost AuthorityHost) (string, error) {
	host := string(authorityHost)
	if host == "" {
		host = string(AzurePublicCloud)
		if envAuthorityHost := os.Getenv(azureAuthorityHost); envAuthorityHost != "" {
			host = envAuthorityHost
		}
	}
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}
	if u.Scheme != "https" {
		return "", errors.New("cannot use an authority host without https")
	}
	return host, nil
}

// validTenantID return true is it receives a valid tenantID, returns false otherwise
func validTenantID(tenantID string) bool {
	match, err := regexp.MatchString("^[0-9a-zA-Z-.]+$", tenantID)
	if err != nil {
		return false
	}
	return match
}

// tokenEndpoint takes a given path and appends "/token" to the end of the path
func tokenEndpoint(p string) string {
	return path.Join(p, "/token")
}

// oauthPath returns the oauth path for AAD or ADFS
func oauthPath(tenantID string) string {
	if tenantID == "adfs" {
		return "/oauth2"
	}
	return "/oauth2/v2.0"
}
