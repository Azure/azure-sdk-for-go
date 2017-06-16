package storage

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// GetSASURIWithSignedIP creates an URL to the specified queue which contains the Shared
// Access Signature with specified permissions and expiration time.
//
// See https://msdn.microsoft.com/en-us/library/azure/ee395415.aspx
func (q *Queue) GetSASURIWithSignedIP(expiry time.Time, permissions string, signedIPRange string) (string, error) {
	canonicalizedResource, err := q.qsc.client.buildCanonicalizedResource(q.buildPath(), q.qsc.auth)
	if err != nil {
		return "", err
	}

	// "The canonicalizedresouce portion of the string is a canonical path to the signed resource.
	// It must include the service name (blob, table, queue or file) for version 2015-02-21 or
	// later, the storage account name, and the resource name, and must be URL-decoded.
	// -- https://msdn.microsoft.com/en-us/library/azure/dn140255.aspx
	// We need to replace + with %2b first to avoid being treated as a space (which is correct for query strings, but not the path component).
	canonicalizedResource = strings.Replace(canonicalizedResource, "+", "%2b", -1)
	canonicalizedResource, err = url.QueryUnescape(canonicalizedResource)
	if err != nil {
		return "", err
	}

	// assumption that start time is now.
	//signedStart := time.Now().UTC().Format(time.RFC3339)
	signedExpiry := expiry.UTC().Format(time.RFC3339)

	// Cannot get this working yet. Any values entered generates bad URL.
	protocols := ""
	signedIdentifier := ""

	stringToSign, err := queueSASStringToSign(q.qsc.client.apiVersion, canonicalizedResource,
		signedExpiry, signedIPRange, permissions, protocols, signedIdentifier)
	if err != nil {
		return "", err
	}

	sig := q.qsc.client.computeHmac256(stringToSign)
	sasParams := url.Values{
		"sv":  {q.qsc.client.apiVersion},
		"se":  {signedExpiry},
		"sp":  {permissions},
		"sig": {sig},
	}

	uri := q.qsc.client.getEndpoint(queueServiceName, q.buildPath(), nil)
	sasURL, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	sasURL.RawQuery = sasParams.Encode()
	return sasURL.String(), nil
}

// GetSASURI creates an URL to the specified queue which contains the Shared
// Access Signature with specified permissions and expiration time.
//
// See https://msdn.microsoft.com/en-us/library/azure/ee395415.aspx
func (q *Queue) GetSASURI(expiry time.Time, permissions string) (string, error) {
	return q.GetSASURIWithSignedIP(expiry, permissions, "")
}

func queueSASStringToSign(signedVersion, canonicalizedResource, signedExpiry, signedIP string, signedPermissions string, protocols string, signedIdentifier string) (string, error) {
	var signedStart string

	if signedVersion >= "2015-02-21" {
		canonicalizedResource = "/queue" + canonicalizedResource
	}

	// https://msdn.microsoft.com/en-us/library/azure/dn140255.aspx#Anchor_12
	if signedVersion >= "2015-04-05" {
		return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
			signedPermissions,
			signedStart,
			signedExpiry,
			canonicalizedResource,
			signedIdentifier,
			signedIP,
			protocols,
			signedVersion), nil

	}

	// reference: http://msdn.microsoft.com/en-us/library/azure/dn140255.aspx
	if signedVersion >= "2013-08-15" {
		return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", signedPermissions, signedStart, signedExpiry, canonicalizedResource, signedIdentifier, signedVersion), nil
	}

	return "", errors.New("storage: not implemented SAS for versions earlier than 2013-08-15")
}
