// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/stretchr/testify/require"
)

type testCache []byte

func (testCache) Export(context.Context, cache.Marshaler, cache.ExportHints) error {
	return nil
}

func (testCache) Replace(context.Context, cache.Unmarshaler, cache.ReplaceHints) error {
	return nil
}

func TestExportReplace(t *testing.T) {
	countCAE, countNoCAE := 0, 0
	c := NewCache(func(cae bool) (cache.ExportReplace, error) {
		if cae {
			countCAE++
		} else {
			countNoCAE++
		}
		return (testCache)([]byte(fmt.Sprint(cae))), nil
	})
	wg := &sync.WaitGroup{}
	ch := make(chan error, 1)
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(cae bool) {
			defer wg.Done()
			if _, err := ExportReplace(c, cae); err != nil {
				select {
				case ch <- err:
					// set error
				default:
					// already set
				}
			}
		}(i%2 == 0)
	}
	wg.Wait()
	select {
	case err := <-ch:
		t.Fatal(err)
	default:
	}
	require.Equal(t, 1, countCAE)
	require.Equal(t, 1, countNoCAE)
	for _, b := range []bool{false, true} {
		xr, err := ExportReplace(c, b)
		require.NoError(t, err)
		require.EqualValues(t, []byte(fmt.Sprint(b)), xr.(testCache))
	}
}
