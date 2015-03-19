package management

import (
	"github.com/MSOpenTech/azure-sdk-for-go/core/http"
	"io"
)

func getResponseBody(response *http.Response) []byte {
	responseBody := make([]byte, response.ContentLength)
	io.ReadFull(response.Body, responseBody)
	return responseBody
}
