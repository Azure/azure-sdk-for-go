// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import "fmt"

// implements `internal/errorinfo/NonRetriable`
type errClosed struct {
	link string
}

func (ec errClosed) NonRetriable() {}
func (ec errClosed) Error() string {
	return fmt.Sprintf("%s is closed and can no longer be used", ec.link)
}
