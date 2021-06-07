package azblob

import (
	"errors"
	"fmt"
	"strings"
)

const (
	AccountName              = "ACCOUNTNAME"
	AccountKey               = "ACCOUNTKEY"
	BlobEndpoint             = "BLOBENDPOINT"
	BlobSecondaryEndpoint    = "BLOBSECONDARYENDPOINT"
	DefaultEndpointsProtocol = "DEFAULTENDPOINTSPROTOCOL"
	EndpointSuffix           = "ENDPOINTSUFFIX"
	Service                  = "SERVICE"
)

// ParseConnectionString accepts connection string to an Azure Storage account and returns primary account url, secondary account url, and Shared key credential
// primaryURL, secondaryURL, cred, err := parseConnectionString(connectionString, "blob")
func ParseConnectionString(connectionString string, relativePath string) (string, string, *SharedKeyCredential, error) {
	if connectionString == "" {
		return "", "", nil, errors.New("connection string cannot be empty")
	}
	connStrParts := strings.Split(connectionString, ";")
	connSettings := make(map[string]string)
	connSettings[Service] = "blob"
	for _, connStrPart := range connStrParts {
		param := strings.SplitN(connStrPart, "=", 2)
		if len(param) != 2 {
			return "", "", nil, errors.New("connection string is either blank or malformed")
		}
		connSettings[strings.ToUpper(param[0])] = param[1]
	}

	credentials, err := NewSharedKeyCredential(connSettings[AccountName], connSettings[AccountKey])
	if err != nil {
		return "", "", credentials, err
	}

	primary, secondary := "", ""
	if _, ok := connSettings[BlobEndpoint]; ok {
		primary = connSettings[BlobEndpoint]
		if _, ok := connSettings[BlobSecondaryEndpoint]; ok {
			secondary = connSettings[BlobSecondaryEndpoint]
		}
	} else {
		if _, ok := connSettings[BlobSecondaryEndpoint]; ok {
			return "", "", nil, errors.New("connection string specifies only secondary endpoint")
		}
		primary = fmt.Sprintf("%s://%s.%s.%s", connSettings[DefaultEndpointsProtocol],
			connSettings[AccountName], connSettings[Service], connSettings[EndpointSuffix])
		secondary = fmt.Sprintf("%s-secondary.%s.%s", connSettings[AccountName],
			connSettings[Service], connSettings[EndpointSuffix])
	}

	if primary == "" {
		primary = fmt.Sprintf("https://%s.%s.%s", connSettings[AccountName],
			connSettings[Service], connSettings[EndpointSuffix])
	}

	return primary + relativePath, secondary + relativePath, credentials, nil
}
