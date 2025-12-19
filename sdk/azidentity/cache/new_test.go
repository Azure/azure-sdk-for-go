//go:build darwin || linux || windows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"errors"
	"path/filepath"
	"sync"
	"testing"

	"github.com/AzureAD/microsoft-authentication-extensions-for-go/cache/accessor"
	"github.com/AzureAD/microsoft-authentication-extensions-for-go/cache/accessor/file"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	errBefore := storageError
	onceBefore := once
	storageBefore := storage
	tryStorageBefore := tryStorage
	t.Cleanup(func() {
		once = onceBefore
		storage = storageBefore
		storageError = errBefore
		tryStorage = tryStorageBefore
	})
	for _, expectedErr := range []error{nil, errors.New("it didn't work")} {
		name := "storage error"
		if expectedErr == nil {
			name = "no storage error"
		}
		t.Run(name, func(t *testing.T) {
			once = &sync.Once{}
			storage = func(string) (accessor.Accessor, error) {
				p := filepath.Join(t.TempDir(), t.Name())
				return file.New(p)
			}
			storageError = nil
			tries := 0
			tryStorage = func() {
				tries++
				storageError = expectedErr
			}
			wg := &sync.WaitGroup{}
			ch := make(chan error, 1)
			for i := 0; i < 50; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if _, err := New(&Options{Name: t.Name()}); err != nil {
						select {
						case ch <- err:
						default:
						}
					}
				}()
			}
			wg.Wait()
			select {
			case err := <-ch:
				if expectedErr == nil {
					t.Fatal(err)
				}
				require.EqualError(t, err, expectedErr.Error())
			default:
			}
			require.Equal(t, 1, tries, "tryStorage was called more than once")
		})
	}
}
