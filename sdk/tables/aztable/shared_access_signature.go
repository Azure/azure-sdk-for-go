// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type AzureSasCredential struct {
	signature string
}

func NewAzureSasCredential(signature string) (*AzureSasCredential, error) {
	return &AzureSasCredential{
		signature: signature,
	}, nil
}

func (a *AzureSasCredential) Update(signature string) {
	a.signature = signature
}

// AuthenticationPolicy implements the Credential interface on SharedKeyCredential.
func (a *AzureSasCredential) AuthenticationPolicy(azcore.AuthenticationPolicyOptions) azcore.Policy {
	return azcore.PolicyFunc(func(req *azcore.Request) (*azcore.Response, error) {
		currentUrl := req.URL.String()
		query := req.URL.Query()

		signature := strings.TrimPrefix(a.signature, "?")

		if query.Encode() != "" {
			if !strings.Contains(currentUrl, signature) {
				currentUrl = currentUrl + "?" + signature
			}
		} else {
			if strings.HasSuffix(currentUrl, "?") {
				currentUrl = currentUrl + signature
			} else {
				currentUrl = currentUrl + "?" + signature
			}
		}

		newUrl, err := url.Parse(currentUrl)
		if err != nil {
			return nil, err
		}
		req.URL = newUrl

		return req.Next()
	})
}

type SharedAccessSignature struct {
	accountName string
	accountKey  string
}

func NewSharedAccessSignature(accountName, accountKey string) (*SharedAccessSignature, error) {
	return &SharedAccessSignature{
		accountName: accountName,
		accountKey:  accountKey,
	}, nil
}

type ResourceType struct {
	Service   bool
	Object    bool
	Container bool
}

func (r ResourceType) String() string {
	s := ""
	if r.Service {
		s += "s"
	}
	if r.Container {
		s += "c"
	}
	if r.Object {
		s += "o"
	}
	return s
}

type AccountSasPermissions struct {
	Read    bool
	Write   bool
	Delete  bool
	List    bool
	Add     bool
	Create  bool
	Update  bool
	Process bool
}

// TODO: Build bitmap
func (a *AccountSasPermissions) String() string {
	var ret string
	if a.Read {
		ret += "r"
	}
	if a.Write {
		ret += "w"
	}
	if a.Delete {
		ret += "d"
	}
	if a.List {
		ret += "l"
	}
	if a.Add {
		ret += "a"
	}
	if a.Create {
		ret += "c"
	}
	if a.Update {
		ret += "u"
	}
	if a.Process {
		ret += "p"
	}
	return ret
}

type SasProtocol string

const (
	SasProtocolHttps = "https"
	SasProtocolHttp  = "http"
)

type AccountSignatureProperties struct {
	// Required: Specifies the services as a bitmap accessible with the account SAS. Default is Tables (t)
	Services string
	// Required: Specifies the resource types that are accessible with the account SAS. Options are Contain
	ResourceTypes ResourceType
	// Required: The permissions associated with the shared access signature, the user is restricted to operations
	// allowed by the permissions
	Permissions AccountSasPermissions
	// Required: The time the shared access signature becomes invalid
	Expiry *time.Time
	// The time when the shared access signature becomes valid
	Start *time.Time
	// Specifies an IP address or range of IP addresses from which to accept requests
	IpAddress string
	// Specifies the protocol permitted for a request made.
	Protocol SasProtocol
}

var X_MS_VERSION = "2020-08-04"

// GenerateAccountSignature creates a signature that delegates service-level operations.
func GenerateAccountSignature(cred SharedKeyCredential, properties AccountSignatureProperties) (string, error) {
	sas := newSharedAccessSignature(cred, X_MS_VERSION)
	return sas.generateAccount("t", properties)
}

type TableSignatureProperties struct {
	TableName string
}

func GenerateTableSignature(cred SharedKeyCredential, properties TableSignatureProperties) (string, error) {
	sas := newSharedAccessSignature(cred, X_MS_VERSION)
	return sas.generateTable(properties)
}

type sharedAccessSignature struct {
	cred     SharedKeyCredential
	version  string
	queryMap map[string]string
}

func newSharedAccessSignature(cred SharedKeyCredential, version string) *sharedAccessSignature {
	return &sharedAccessSignature{
		cred:     cred,
		version:  X_MS_VERSION,
		queryMap: make(map[string]string),
	}
}

func (s *sharedAccessSignature) generateAccount(service string, properties AccountSignatureProperties) (string, error) {
	s.addBase(properties.Permissions, properties.Expiry, properties.Start, properties.IpAddress, properties.Protocol, s.version)
	s.addAccount("t", properties.ResourceTypes)
	err := s.addAccountSignature()
	if err != nil {
		return "", err
	}
	return s.getToken()
}

func (s *sharedAccessSignature) generateTable(properties TableSignatureProperties) (string, error) {
	return "", nil
}

func toUtcDatetime(t *time.Time) string {
	return t.Format(utcLayout) + "Z"
}

func (s *sharedAccessSignature) addBase(permissions AccountSasPermissions, expiry *time.Time, start *time.Time, ip string, protocol SasProtocol, version string) {
	s.addQuery(queryStringSignedExpiry, toUtcDatetime(expiry))
	s.addQuery(queryStringSignedStart, toUtcDatetime(start))
	s.addQuery(queryStringSignedPermission, permissions.String())
	s.addQuery(queryStringSignedIp, ip)
	s.addQuery(querySignedProtocol, string(protocol))
	s.addQuery(querySignedVersion, version)
}

func (s *sharedAccessSignature) addAccount(services string, resources ResourceType) {
	s.addQuery(querySignedServices, services)
	s.addQuery(querySignedResourceTypes, resources.String())
}

func (s *sharedAccessSignature) addAccountSignature() error {
	// Get string to sign
	signedString := s.buildStringToSign()

	signed, err := s.cred.ComputeHMACSHA256(signedString)
	if err != nil {
		return err
	}
	s.addQuery(queryStringSignedSignature, signed)
	return nil
}

func (s *sharedAccessSignature) buildStringToSign() string {
	var ret string
	ret += s.cred.accountName + "\n"
	ret += s.getValueToAppend(queryStringSignedPermission)
	ret += s.getValueToAppend(querySignedServices)
	ret += s.getValueToAppend(querySignedResourceTypes)
	ret += s.getValueToAppend(queryStringSignedStart)
	ret += s.getValueToAppend(queryStringSignedExpiry)
	ret += s.getValueToAppend(queryStringSignedIp)
	ret += s.getValueToAppend(querySignedProtocol)
	ret += s.getValueToAppend(querySignedVersion)
	fmt.Println(ret)
	return ret
}

func (s *sharedAccessSignature) getValueToAppend(queryString string) string {
	val, ok := s.queryMap[queryString]
	if !ok {
		return "" + "\n"
	}
	return val + "\n"
}

func (s *sharedAccessSignature) addQuery(name, val string) {
	if val != "" {
		s.queryMap[name] = val
	}
}

func (s *sharedAccessSignature) getToken() (string, error) {
	fmt.Println(s.queryMap)
	pairsList := make([]string, 0)
	for k, v := range s.queryMap {
		if v != "" {
			pairsList = append(pairsList, fmt.Sprintf("%v=%v", k, v))
		}
	}
	return "?" + strings.Join(pairsList, "&"), nil
}

const (
	queryStringSignedSignature    = "sig"
	queryStringSignedPermission   = "sp"
	queryStringSignedStart        = "st"
	queryStringSignedExpiry       = "se"
	queryStringSignedResource     = "sr"
	queryStringSignedIdentifier   = "si"
	queryStringSignedIp           = "sip"
	querySignedProtocol           = "spr"
	querySignedVersion            = "sv"
	querySignedCacheControl       = "rscc"
	querySignedContentDisposition = "rscd"
	querySignedContentEncoding    = "rsce"
	querySignedContentLanguage    = "rscl"
	querySignedContentType        = "rsct"
	querySignedStartPk            = "spk"
	querySignedStartRk            = "srk"
	querySignedEndPk              = "epk"
	querySignedEndRk              = "erk"
	querySignedResourceTypes      = "srt"
	querySignedServices           = "ss"
	querySignedTableName          = "tn"
)
