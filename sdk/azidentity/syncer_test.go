//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func TestSyncer(t *testing.T) {
	silentAuths, tokenRequests := 0, 0
	s := newSyncer("", "tenant",
		func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
			tokenRequests++
			return azcore.AccessToken{}, nil
		},
		func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
			var err error
			if tokenRequests == 0 {
				err = errors.New("cache miss")
			}
			silentAuths++
			return azcore.AccessToken{}, err
		},
		syncerOptions{},
	)
	goroutines := 50
	wg := sync.WaitGroup{}
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			_, err := s.GetToken(context.Background(), testTRO)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if tokenRequests != 1 {
		t.Errorf("expected 1 token request, got %d", tokenRequests)
	}
	if silentAuths != goroutines {
		t.Errorf("expected %d silent auth attempts, got %d", goroutines, silentAuths)
	}
}
