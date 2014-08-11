package azureSdkForGo

import (
	"fmt"
	"os"
	"github.com/MSOpenTech/azure-sdk-for-go/core/tls"
	"github.com/MSOpenTech/azure-sdk-for-go/core/http"
	"io"
	"bytes"
	"strings"
	"os/exec"
)


func PrintErrorAndExit(err error) {
	fmt.Println("Error: ")
	fmt.Println(err)
	os.Exit(2)
}

func SendAzureGetRequest(url string) ([]byte){
	response := SendAzureRequest(url, "GET", nil)
	return response
}

func SendAzurePostRequest(url string, data []byte) ([]byte){
	response := SendAzureRequest(url, "POST", data)
	return response
}

func SendAzureRequest(url string, requestType string,  data []byte) ([]byte){
	client := createHttpClient()

	request, reqErr := createAzureRequest(url, requestType, data)
	if reqErr != nil {
		PrintErrorAndExit(reqErr)
	}

	response, err := client.Do(request)
	if err != nil {
		PrintErrorAndExit(err)
	}

	responseContent := getResponseBody(response)
	return responseContent
}

func ExecuteCommand(command string) ([]byte) {
	parts := strings.Fields(command)
	head := parts[0]
	parts = parts[1:len(parts)]

	cmd := exec.Command(head, parts...)

	out, err := cmd.Output()

	if err != nil {
		PrintErrorAndExit(err)
	}

	return out
}

func createAzureRequest(url string, requestType string,  data []byte) (*http.Request, error){
	var request *http.Request
	var err error
	if data != nil {
		body := bytes.NewBuffer(data)
		request, err = http.NewRequest(requestType, url, body)
	} else {
		request, err = http.NewRequest(requestType, url, nil)
	}

	request.Header.Add("Content-Type", "application/xml")
	request.Header.Add("x-ms-version", "2014-05-01")

	return request, err
}

func createHttpClient() (*http.Client){
	cert, _ := tls.LoadX509KeyPair(GetPublishSettings().SubscriptionCert, GetPublishSettings().SubscriptionCert)

	ssl := &tls.Config{}
	ssl.Certificates = []tls.Certificate{cert}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: ssl,
		},
	}

	return client
}

func getResponseBody(response *http.Response) ([]byte){

	responseBody := make([]byte, response.ContentLength)
	io.ReadFull(response.Body, responseBody)
	return responseBody
}
