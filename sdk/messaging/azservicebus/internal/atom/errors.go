// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func NewResponseError(inner error, resp *http.Response) error {
	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return &azcore.ResponseError{
			StatusCode:  resp.StatusCode,
			RawResponse: resp,
		}
	}

	var ec *managementError

	if err := xml.Unmarshal(bytes, &ec); err != nil {
		return &azcore.ResponseError{
			StatusCode:  resp.StatusCode,
			RawResponse: resp,
		}
	}

	return &azcore.ResponseError{
		ErrorCode:   ec.Detail,
		StatusCode:  resp.StatusCode,
		RawResponse: resp,
	}
}

// ResponseError conforms to the older azcore.HTTPResponse
// NOTE: after breaking changes have been incorporated we'll move
// to the newer azcore.HTTPResponseError
type ResponseError struct {
	inner error
	resp  *http.Response
}

func (e ResponseError) RawResponse() *http.Response {
	return e.resp
}

func (e ResponseError) Error() string {
	if e.inner == nil {
		return fmt.Sprintf("%s: %d", e.resp.Status, e.resp.StatusCode)
	} else {
		return e.inner.Error()
	}
}

func (e ResponseError) Unwrap() error {
	return e.inner
}
