package arm

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	EnvAuthorizationToken = "AzureAuthorizationToken"

	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"

	MimeTypeJson = "application/json"
)

var HeaderFormats map[string]string = map[string]string{
	HeaderAuthorization: "Bearer %s",
}

// Constructs a URL given the collection of base, path, and parameters. Path parameters are brace-enclosed strings (e.g., {replace-this}).
// The baseUrl must be absolute. Paths that begin with a forward-slash (i.e., '/') will cause removal of any path components within the baseUrl.
func BuildUrl(baseUrl string, path string, pathParameters map[string]interface{}, queryParameters map[string]interface{}) (urlFull string, err error) {

	// Shred the baseUrl into parts
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("%s is not a valid base URL (%v)", baseUrl, err)
	}

	// Ensure the path, if any, does not end with a forward-slash (one is added below)
	if strings.HasSuffix(u.Path, "/") {
		u.Path = strings.TrimRight(u.Path, "/")
	}

	// Replace path parameters and construct the path
	for key, value := range ensureValueStrings(pathParameters) {
		path = strings.Replace(path, "{"+key+"}", value, -1)
	}
	if strings.HasPrefix(path, "/") {
		u.Path = path
	} else {
		u.Path += "/" + path
	}

	// Encode query parameters (if any)
	v := u.Query()
	for key, value := range ensureValueStrings(queryParameters) {
		v.Add(key, value)
	}
	u.RawQuery = v.Encode()

	return u.String(), nil
}

func SendRequest(verb string, urlResource string, headers map[string]interface{}, body string) (response *http.Response, err error) {

	// Create the HTTP client
	client := &http.Client{}

	// Create the request
	request, err := http.NewRequest(verb, urlResource, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("Failed sending request to %s (%v)", urlResource, err)
	}

	// Set headers
	authToken := os.Getenv(EnvAuthorizationToken)
	if len(authToken) > 0 {
		request.Header.Add(HeaderAuthorization, fmt.Sprintf(HeaderFormats[HeaderAuthorization], authToken))
	}
	if len(body) > 0 {
		request.Header.Add(HeaderContentType, MimeTypeJson)
	}

	// Send the request
	return client.Do(request)
}

func ensureValueStrings(mapOfInterface map[string]interface{}) map[string]string {
	mapOfStrings := make(map[string]string)
	for key, value := range mapOfInterface {
		mapOfStrings[key] = ensureValueString(value)
	}
	return mapOfStrings
}

func ensureValueString(value interface{}) string {
	switch t := value.(type) {
	case nil:
		return ""
	case string:
		return t
	case []byte:
		return string(t)
	case []rune:
		return string(t)
	default:
		return fmt.Sprintf("%v", value)
	}
}
