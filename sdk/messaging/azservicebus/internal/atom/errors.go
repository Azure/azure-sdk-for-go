// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"fmt"
	"net/http"
)

// ResponseError conforms to the azcore.HTTPResponse
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
