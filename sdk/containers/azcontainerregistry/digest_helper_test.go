// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseDigestValidator(t *testing.T) {
	tests := []struct {
		name    string
		digest  string
		want    digestValidator
		wantErr error
	}{
		{"sha256", "sha256:test", newSha256Validator(), nil},
		{"not supported", "sha512:test", nil, ErrDigestAlgNotSupported},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDigestValidator(tt.digest)
			if err != nil || tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("parseDigestValidator() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDigestValidator() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlobDigestCalculator_saveAndRestoreState(t *testing.T) {
	calculator := NewBlobDigestCalculator()
	calculator.restoreState()
	calculator.saveState()
	calculator.restoreState()
	calculator.h.Write([]byte("test1"))
	sum := calculator.h.Sum(nil)
	calculator.saveState()
	calculator.h.Write([]byte("test2"))
	require.NotEqual(t, sum, calculator.h.Sum(nil))
	calculator.restoreState()
	require.Equal(t, sum, calculator.h.Sum(nil))
}
