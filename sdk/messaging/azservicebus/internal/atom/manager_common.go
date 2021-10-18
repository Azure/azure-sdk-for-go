// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/devigned/tab"
)

// constructAtomPath adds the proper parameters for skip and top
// This is common for the list operations for queues, topics and subscriptions.
func constructAtomPath(basePath string, skip int, top int) string {
	values := url.Values{}

	if skip > 0 {
		values.Add("$skip", fmt.Sprintf("%d", skip))
	}

	if top > 0 {
		values.Add("$top", fmt.Sprintf("%d", top))
	}

	if len(values) == 0 {
		return basePath
	}

	return fmt.Sprintf("%s?%s", basePath, values.Encode())
}

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
