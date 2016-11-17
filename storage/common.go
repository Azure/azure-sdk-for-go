package storage

import (
	"io/ioutil"
	"time"
)

// AccessPolicyDetails are used for SETTING policies
type AccessPolicyDetails struct {
	ID         string
	StartTime  time.Time
	ExpiryTime time.Time
	CanRead    bool
	CanWrite   bool
	CanDelete  bool
}

// AccessPolicyDetailsXML has specifics about an access policy
// annotated with XML details.
type AccessPolicyDetailsXML struct {
	StartTime  time.Time `xml:"Start"`
	ExpiryTime time.Time `xml:"Expiry"`
	Permission string    `xml:"Permission"`
}

// SignedIdentifier is a wrapper for a specific policy
type SignedIdentifier struct {
	ID           string                 `xml:"Id"`
	AccessPolicy AccessPolicyDetailsXML `xml:"AccessPolicy"`
}

// SignedIdentifiers part of the response from GetPermissions call.
type SignedIdentifiers struct {
	SignedIdentifiers []SignedIdentifier `xml:"SignedIdentifier"`
}

// AccessPolicy is the response type from the GetPermissions call.
type AccessPolicy struct {
	SignedIdentifiersList SignedIdentifiers `xml:"SignedIdentifiers"`
}

func generatePermissions(accessPolicy AccessPolicyDetails) (permissions string) {
	// generate the permissions string (rwd).
	// still want the end user API to have bool flags.
	permissions = ""

	if accessPolicy.CanRead {
		permissions += "r"
	}

	if accessPolicy.CanWrite {
		permissions += "w"
	}

	if accessPolicy.CanDelete {
		permissions += "d"
	}

	return permissions
}

// convertAccessPolicyToXMLStructs converts between AccessPolicyDetails which is a struct better for API usage to the
// AccessPolicy struct which will get converted to XML.
func convertAccessPolicyToXMLStructs(accessPolicy AccessPolicyDetails) SignedIdentifiers {
	return SignedIdentifiers{
		SignedIdentifiers: []SignedIdentifier{
			{
				ID: accessPolicy.ID,
				AccessPolicy: AccessPolicyDetailsXML{
					StartTime:  accessPolicy.StartTime.UTC().Round(time.Second),
					ExpiryTime: accessPolicy.ExpiryTime.UTC().Round(time.Second),
					Permission: generatePermissions(accessPolicy),
				},
			},
		},
	}
}

// generateAccessPolicy generates the XML access policy used as the payload for SetContainerPermissions.
func generateAccessPolicy(accessPolicy AccessPolicyDetails) (accessPolicyXML string, err error) {

	if accessPolicy.ID != "" {
		signedIdentifiers := convertAccessPolicyToXMLStructs(accessPolicy)
		body, _, err := xmlMarshal(signedIdentifiers)
		if err != nil {
			return "", err
		}

		xmlByteArray, err := ioutil.ReadAll(body)
		if err != nil {
			return "", err
		}
		accessPolicyXML = string(xmlByteArray)
		return accessPolicyXML, nil
	}

	return "", nil
}
