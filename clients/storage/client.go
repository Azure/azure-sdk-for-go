package storage

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

const (
	DefaultBaseUrl    = "core.windows.net"
	DefaultApiVersion = "2014-02-14"
	defaultUseHttps   = true

	blobServiceName  = "blob"
	tableServiceName = "table"
	queueServiceName = "queue"
)

type StorageClient struct {
	accountName string
	accountKey  string
	useHttps    bool
	baseUrl     string
	apiVersion  string
}

func NewBasicClient(accountName, accountKey string) (*StorageClient, error) {
	return NewClient(accountName, accountKey, DefaultBaseUrl, DefaultApiVersion, defaultUseHttps)
}

func NewClient(accountName, accountKey, blobServiceBaseUrl, apiVersion string, useHttps bool) (*StorageClient, error) {
	if accountName == "" {
		return nil, fmt.Errorf("azure: account name required")
	} else if accountKey == "" {
		return nil, fmt.Errorf("azure: account key required")
	} else if blobServiceBaseUrl == "" {
		return nil, fmt.Errorf("azure: base storage service url required")
	}

	return &StorageClient{
		accountName: accountName,
		accountKey:  accountKey,
		useHttps:    useHttps,
		baseUrl:     blobServiceBaseUrl,
		apiVersion:  apiVersion}, nil
}

func (c StorageClient) getBaseUrl(service string) string {
	scheme := "http"
	if c.useHttps {
		scheme = "https"
	}

	host := fmt.Sprintf("%s.%s.%s", c.accountName, service, c.baseUrl)

	u := &url.URL{
		Scheme: scheme,
		Host:   host}
	return u.String()
}

func (c StorageClient) getEndpoint(service, path string, params url.Values) string {
	u, err := url.Parse(c.getBaseUrl(service))
	if err != nil {
		// really should not happen
		panic(err)
	}

	if path == "" {
		path = "/" // API doesn't accept path segments not starting with '/''
	}

	u.Path = path
	u.RawQuery = params.Encode()
	return u.String()
}

func (c StorageClient) GetBlobService() *BlobStorageClient {
	return &BlobStorageClient{c}
}

func (c StorageClient) createAuthorizationHeader(canonicalizedString string) (string, error) {
	signature, err := c.computeHmac256(canonicalizedString)
	if err != nil {
		return "", err
	}
	authorizationHeader := fmt.Sprintf("%s %s:%s", "SharedKey", c.accountName, signature)
	return authorizationHeader, nil
}

func (c StorageClient) getAuthorizationHeader(verb, url string, headers map[string]string) (string, error) {
	canonicalizedResource, err := c.buildCanonicalizedResource(url)
	if err != nil {
		return "", err
	}

	canonicalizedString := c.buildCanonicalizedString(verb, headers, canonicalizedResource)
	return c.createAuthorizationHeader(canonicalizedString)
}

func (c StorageClient) getStandardHeaders() map[string]string {
	// TODO (ahmetalpbalkan) test
	return map[string]string{
		"x-ms-version": c.apiVersion,
		"x-ms-date":    currentTimeRfc1123Formatted(),
	}
}

func (c StorageClient) buildCanonicalizedHeader(headers map[string]string) string {
	// TODO (ahmetalpbalkan) write test case for imported code

	cm := make(map[string]string)

	for k, v := range headers {
		headerName := strings.TrimSpace(strings.ToLower(k))
		match, _ := regexp.MatchString("x-ms-", headerName)
		if match {
			cm[headerName] = v
		}
	}

	if len(cm) == 0 {
		return ""
	}

	keys := make([]string, 0, len(cm))
	for key, _ := range cm {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	ch := ""

	for i, key := range keys {
		if i == len(keys)-1 {
			ch += fmt.Sprintf("%s:%s", key, cm[key])
		} else {
			ch += fmt.Sprintf("%s:%s\n", key, cm[key])
		}
	}
	return ch
}

func (c StorageClient) buildCanonicalizedResource(uri string) (string, error) {
	errMsg := "buildCanonicalizedResource error: %s"
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf(errMsg, err.Error())
	}

	cr := "/" + c.accountName
	if len(u.Path) > 0 {
		cr += u.Path
	}

	params, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", fmt.Errorf(errMsg, err.Error())
	}

	if len(params) > 0 {
		cr += "\n"
		keys := make([]string, 0, len(params))
		for key, _ := range params {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for i, key := range keys {
			if len(params[key]) > 1 {
				sort.Strings(params[key])
			}

			if i == len(keys)-1 {
				cr += fmt.Sprintf("%s:%s", key, strings.Join(params[key], ","))
			} else {
				cr += fmt.Sprintf("%s:%s\n", key, strings.Join(params[key], ","))
			}
		}
	}
	return cr, nil
}

func (c StorageClient) buildCanonicalizedString(verb string, headers map[string]string, canonicalizedResource string) string {
	canonicalizedString := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
		verb,
		headers["Content-Encoding"],
		headers["Content-Language"],
		headers["Content-Length"],
		headers["Content-MD5"],
		headers["Content-Type"],
		headers["Date"],
		headers["If-Modified-Singe"],
		headers["If-Match"],
		headers["If-None-Match"],
		headers["If-Unmodified-Singe"],
		headers["Range"],
		c.buildCanonicalizedHeader(headers),
		canonicalizedResource)

	return canonicalizedString
}

func (c StorageClient) exec(verb, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	authHeader, err := c.getAuthorizationHeader(verb, url, headers)
	if err != nil {
		return nil, err
	}
	headers["Authorization"] = authHeader

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(verb, url, body)
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	statusCode := resp.StatusCode
	if statusCode >= 400 && statusCode <= 505 {
		defer resp.Body.Close()
		errXml := new(ErrorXml)

		var respBody []byte
		respBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// TODO (ahmetalpbalkan) write test case for imported deserialization code
		// TODO (ahmetalpbalkan) extract
		err = xml.Unmarshal(respBody, errXml)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s %s", "Remote server returned error:", errXml.Message)
	}

	return readResponseBody(resp)
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// TODO (ahmetalpbalkan) refactor
type ErrorXml struct {
	Code                      string
	Message                   string
	AuthenticationErrorDetail string
}
