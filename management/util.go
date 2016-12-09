package management

import (
	"io/ioutil"
	"net/http"
)

func getResponseBody(response *http.Response) ([]byte, error) {
	defer func() {
		_ = response.Body.Close()
	}()

	return ioutil.ReadAll(response.Body)
}
