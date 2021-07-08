// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"github.com/stretchr/testify/assert"
)

func failOnError(err error, s *tableClientLiveTests) {
	if err != nil {
		assert.FailNow(s.T(), err.Error())
	}
}
