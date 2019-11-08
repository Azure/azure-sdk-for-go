// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"fmt"
	"io/ioutil"
)

// newBodyDownloadPolicy creates a policy object that downloads the response's body to a []byte.
func newBodyDownloadPolicy() Policy {
	return PolicyFunc(func(ctx context.Context, req *Request) (*Response, error) {
		response, err := req.Do(ctx)
		if err != nil {
			return response, err
		}
		var opValues bodyDownloadPolicyOpValues
		if req.OperationValue(&opValues); !opValues.skip {
			// Either bodyDownloadPolicyOpValues was not specified (so skip is false)
			// or it was specified and skip is false: don't skip downloading the body
			response.Payload, err = ioutil.ReadAll(response.Body)
			response.Body.Close()
			if err != nil {
				err = fmt.Errorf("body download policy: %w", err)
			}
		}
		return response, err
	})
}

// bodyDownloadPolicyOpValues is the struct containing the per-operation values
type bodyDownloadPolicyOpValues struct {
	skip bool
}
