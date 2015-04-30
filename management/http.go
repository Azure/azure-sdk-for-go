package management

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/core/http"
	"github.com/Azure/azure-sdk-for-go/core/tls"
)

const (
	msVersionHeader           = "x-ms-version"
	msVersionHeaderValue      = "2014-10-01"
	contentHeader             = "Content-Type"
	defaultContentHeaderValue = "application/xml"
	requestIdHeader           = "X-Ms-Request-Id"
)

//sendAzureGetRequest sends a request to the management API using the HTTP GET method
//and returns the response body or an error.
func (client *Client) SendAzureGetRequest(url string) ([]byte, error) {
	if url == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "url")
	}

	response, err := client.sendAzureRequest(url, "GET", "", nil)
	if err != nil {
		return nil, err
	}

	responseContent := getResponseBody(response)
	return responseContent, nil
}

//sendAzurePostRequest sends a request to the management API using the HTTP POST method
//and returns the request ID or an error.
func (client *Client) SendAzurePostRequest(url string, data []byte) (OperationId, error) {
	if url == "" {
		return "", fmt.Errorf(errParamNotSpecified, "url")
	}

	response, err := client.sendAzureRequest(url, "POST", "", data)
	if err != nil {
		return "", err
	}

	return getOperationId(response)
}

//sendAzurePutRequest sends a request to the management API using the HTTP PUT method
//and returns the request ID or an error. The content type can be specified, however
//if an empty string is passed, the default of "application/xml" will be used.
func (client *Client) SendAzurePutRequest(url string, contentType string, data []byte) (OperationId, error) {
	if url == "" {
		return "", fmt.Errorf(errParamNotSpecified, contentType, "url")
	}

	response, err := client.sendAzureRequest(url, "PUT", contentType, data)
	if err != nil {
		return "", err
	}

	return getOperationId(response)
}

//sendAzureDeleteRequest sends a request to the management API using the HTTP DELETE method
//and returns the request ID or an error.
func (client *Client) SendAzureDeleteRequest(url string) (OperationId, error) {
	if url == "" {
		return "", fmt.Errorf(errParamNotSpecified, "url")
	}

	response, err := client.sendAzureRequest(url, "DELETE", "", nil)
	if err != nil {
		return "", err
	}

	return getOperationId(response)
}

func getOperationId(response *http.Response) (OperationId, error) {
	requestId := response.Header[requestIdHeader]
	if len(requestId) == 0 {
		return "", fmt.Errorf("Could not retrieve operation id from %q header", requestIdHeader)
	}
	return OperationId(requestId[0]), nil
}

//sendAzureRequest constructs an HTTP client for the request, sends it to the
//management API and returns the response or an error.
func (client *Client) sendAzureRequest(url string, requestType string, contentType string, data []byte) (*http.Response, error) {
	if url == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "url")
	}
	if requestType == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "requestType")
	}

	httpClient := client.createHttpClient()

	response, err := client.sendRequest(httpClient, url, requestType, contentType, data, 7)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//createHttpClient creates an HTTP Client configured with the key pair for
//the subscription for this client.
func (client *Client) createHttpClient() *http.Client {
	cert, _ := tls.X509KeyPair(client.publishSettings.SubscriptionCert, client.publishSettings.SubscriptionKey)

	ssl := &tls.Config{}
	ssl.Certificates = []tls.Certificate{cert}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: ssl,
		},
	}

	return httpClient
}

//sendRequest sends a request to the Azure management API using the given
//HTTP client and parameters. It returns the response from the call or an
//error.
func (client *Client) sendRequest(httpClient *http.Client, url string, requestType string, contentType string, data []byte, numberOfRetries int) (*http.Response, error) {
	request, reqErr := client.createAzureRequest(url, requestType, contentType, data)
	if reqErr != nil {
		return nil, reqErr
	}

	response, err := httpClient.Do(request)
	if err != nil {
		if numberOfRetries == 0 {
			return nil, err
		}

		return client.sendRequest(httpClient, url, requestType, contentType, data, numberOfRetries-1)
	}

	if response.StatusCode >= http.StatusBadRequest {
		responseContent := getResponseBody(response)
		azureErr := getAzureError(responseContent)
		if azureErr != nil {
			if numberOfRetries == 0 {
				return nil, azureErr
			}

			return client.sendRequest(httpClient, url, requestType, contentType, data, numberOfRetries-1)
		}
	}

	return response, nil
}

//createAzureRequest packages up the request with the correct set of headers and returns
//the request object or an error.
func (client *Client) createAzureRequest(url string, requestType string, contentType string, data []byte) (*http.Request, error) {
	var request *http.Request
	var err error

	url = fmt.Sprintf("%s/%s/%s", client.managementURL, client.publishSettings.SubscriptionID, url)
	if data != nil {
		body := bytes.NewBuffer(data)
		request, err = http.NewRequest(requestType, url, body)
	} else {
		request, err = http.NewRequest(requestType, url, nil)
	}

	if err != nil {
		return nil, err
	}

	request.Header.Add(msVersionHeader, msVersionHeaderValue)
	if len(contentType) > 0 {
		request.Header.Add(contentHeader, contentType)
	} else {
		request.Header.Add(contentHeader, defaultContentHeaderValue)
	}

	return request, nil
}

//getAzureError converts an error response body into an AzureError type.
func getAzureError(responseBody []byte) error {
	error := new(AzureError)
	err := xml.Unmarshal(responseBody, error)
	if err != nil {
		return err
	}

	return error
}
