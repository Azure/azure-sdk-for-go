package management

import (
	"io/ioutil"
	"net/http"
)

func getResponseBody(response *http.Response) (out []byte, err error) {
	defer func() {
		err = response.Body.Close()
	}()
	return ioutil.ReadAll(response.Body)
}
