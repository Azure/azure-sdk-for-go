package management

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/MSOpenTech/azure-sdk-for-go/core/http"
	"github.com/MSOpenTech/azure-sdk-for-go/core/tls"
)

const (
	msVersionHeader           = "x-ms-version"
	msVersionHeaderValue      = "2014-10-01"
	contentHeader             = "Content-Type"
	defaultContentHeaderValue = "application/xml"
	requestIdHeader           = "X-Ms-Request-Id"
)

const (
	EnvEnableRequestLogs = "AZURESDK_LOGREQUESTS_MANAGEMENT"      // Set this environment variable to a non-empty value to enable request logging.
	EnvRequestLogPath    = "AZURESDK_LOGREQUESTS_MANAGEMENT_PATH" // Set this environment variable to the path of the log file. If not set, os.Stderr will be used.
)

//SendAzureGetRequest sends a request to the management API using the HTTP GET method
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

//SendAzurePostRequest sends a request to the management API using the HTTP POST method
//and returns the request ID or an error.
func (client *Client) SendAzurePostRequest(url string, data []byte) (string, error) {
	if url == "" {
		return "", fmt.Errorf(errParamNotSpecified, "url")
	}

	response, err := client.sendAzureRequest(url, "POST", "", data)
	if err != nil {
		return "", err
	}

	requestId := response.Header[requestIdHeader]
	return requestId[0], nil
}

//SendAzurePutRequest sends a request to the management API using the HTTP PUT method
//and returns the request ID or an error. The content type can be specified, however
//if an empty string is passed, the default of "application/xml" will be used.
func (client *Client) SendAzurePutRequest(url string, contentType string, data []byte) (string, error) {
	if url == "" {
		return "", fmt.Errorf(errParamNotSpecified, contentType, "url")
	}

	response, err := client.sendAzureRequest(url, "PUT", contentType, data)
	if err != nil {
		return "", err
	}

	requestId := response.Header[requestIdHeader]
	return requestId[0], nil
}

//SendAzureDeleteRequest sends a request to the management API using the HTTP DELETE method
//and returns the request ID or an error.
func (client *Client) SendAzureDeleteRequest(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf(errParamNotSpecified, "url")
	}

	response, err := client.sendAzureRequest(url, "DELETE", "", nil)
	if err != nil {
		return "", err
	}

	requestId := response.Header[requestIdHeader]
	return requestId[0], nil
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

	var requestLogger *log.Logger
	if os.Getenv(EnvEnableRequestLogs) != "" {
		out := os.Stderr
		if logfile := os.Getenv(EnvRequestLogPath); logfile != "" {
			var err error
			out, err = os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				panic(err)
			}
		}
		requestLogger = log.New(out, "azuresdk.management", log.LstdFlags)
	}

	httpClient := &http.Client{
		Transport: &loggedTransport{logger: requestLogger, inner: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: ssl,
		}},
	}

	return httpClient
}

type loggedTransport struct {
	logger *log.Logger
	inner  *http.Transport
}

func (l *loggedTransport) RoundTrip(request *http.Request) (response *http.Response, err error) {
	// log request
	if l.logger != nil {
		l.logger.Printf("[REQUEST] %s %s", request.Method, request.URL)
		l.logger.Printf("[REQUEST] headers: %+v", request.Header)
		if request.Body != nil {
			body, ok := request.Body.(*nopCloser)
			if !ok {
				body = &nopCloser{&bytes.Buffer{}}
				body.Buffer.ReadFrom(request.Body)
				request.Body = body
			}
			l.logger.Printf("[REQUEST] body: %s", string(body.Bytes()))
		} else {
			l.logger.Printf("[REQUEST] body: nil")
		}
	}

	response, err = l.inner.RoundTrip(request)

	//log response
	if l.logger != nil {
		l.logger.Printf("[RESPONSE] %s", response.Status)
		l.logger.Printf("[RESPONSE] headers: %+v", response.Header)
		body := &nopCloser{&bytes.Buffer{}}
		body.Buffer.ReadFrom(response.Body)
		response.Body = body
		l.logger.Printf("[RESPONSE] body: %s", string(body.Bytes()))
	}
	return
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
		request, err = http.NewRequest(requestType, url, &nopCloser{body})
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

// provide our own closer wrapper so that http.NewRequest doesn't need to wrap it with
// ioutil.NopCloser. This way, we can find our way back to buffer for logging the request body
type nopCloser struct {
	*bytes.Buffer
}

func (nopCloser) Close() error { return nil }

//getAzureError converts an error response body into an AzureError type.
func getAzureError(responseBody []byte) error {
	error := new(AzureError)
	err := xml.Unmarshal(responseBody, error)
	if err != nil {
		return err
	}

	return error
}
