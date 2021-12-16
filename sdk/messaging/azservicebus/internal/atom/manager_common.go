// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/devigned/tab"
)

// CloseRes closes the response (or if it's nil just no-ops)
func CloseRes(ctx context.Context, res *http.Response) {
	if res == nil {
		return
	}

	_, _ = io.Copy(ioutil.Discard, res.Body)

	if err := res.Body.Close(); err != nil {
		tab.For(ctx).Error(err)
	}
}
