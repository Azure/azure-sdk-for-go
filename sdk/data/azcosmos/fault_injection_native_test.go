// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFaultInjectionNativeHeaderPresence(t *testing.T) {
	for _, tt := range []struct {
		name    string
		result  FaultInjectionResult
		wantErr bool
	}{
		{"predefined nil", FaultInjectionResult{Error: FaultInjectionInternalServerError}, false},
		{"predefined empty", FaultInjectionResult{Error: FaultInjectionInternalServerError, Headers: map[string]string{}}, false},
		{"custom nil", FaultInjectionResult{CustomStatusCode: 429}, false},
		{"custom empty", FaultInjectionResult{CustomStatusCode: 429, Headers: map[string]string{}}, false},
		{"custom populated", FaultInjectionResult{CustomStatusCode: 429, Headers: map[string]string{"x-test": "value"}}, false},
		{"predefined populated requires custom status", FaultInjectionResult{Error: FaultInjectionInternalServerError, Headers: map[string]string{"x-test": "value"}}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			options := ClientOptions{FaultInjectionRules: []FaultInjectionRule{{
				ID:        "header-presence",
				Condition: FaultInjectionCondition{Operation: FaultInjectionReadItem},
				Result:    tt.result,
			}}}
			require.NoError(t, options.validate())
			driver, err := openDriver(driverConfig{
				endpoint:   "https://myaccount.documents.azure.com",
				accountKey: emulatorKey,
				options:    options,
			})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, driver.close()) })
			err = driver.inspectFaultInjectionOptionsBuild()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
