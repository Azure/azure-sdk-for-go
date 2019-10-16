// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"fmt"
)

// TODO: capture frame info for marshal, unmarshal, and parsing errors
// built in frame in xerror?  %w
type frameInfo struct {
	file string
	line int
}

func (f frameInfo) String() string {
	if f.zero() {
		return ""
	}
	return fmt.Sprintf("file: %s, line: %d", f.file, f.line)
}

func (f frameInfo) zero() bool {
	return f.file == "" && f.line == 0
}

// RequestError is returned when the service returns an unsuccessful resopnse code (4xx, 5xx).
type RequestError struct {
	msg  string
	resp *Response
}

func newRequestError(message string, response *Response) error {
	return RequestError{msg: message, resp: response}
}

func (re RequestError) Error() string {
	return re.msg
}

// Response returns the underlying response.
func (re RequestError) Response() *Response {
	return re.resp
}
