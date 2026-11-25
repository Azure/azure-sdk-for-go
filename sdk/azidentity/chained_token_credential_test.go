// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/stretchr/testify/require"
)

func TestChainedTokenCredential_AuthenticationFailure(t *testing.T) {
	t.Run("nested chains", func(t *testing.T) {
		unavailable := newCredentialUnavailableError("unavailable", "not configured")
		terminal := newAuthenticationFailedErrorFromMSAL("terminal", context.Canceled)
		first, last, fallback := NewFakeCredential(), NewFakeCredential(), NewFakeCredential()
		first.SetResponse(azcore.AccessToken{}, unavailable)
		last.SetResponse(azcore.AccessToken{}, terminal)
		fallback.SetResponse(azcore.AccessToken{Token: tokenValue}, nil)
		inner, err := NewChainedTokenCredential([]azcore.TokenCredential{first, last}, nil)
		require.NoError(t, err)
		outer, err := NewChainedTokenCredential([]azcore.TokenCredential{inner, fallback}, nil)
		require.NoError(t, err)
		_, err = outer.GetToken(context.Background(), testTRO)
		require.Zero(t, fallback.getTokenCalls, "an inner unavailable failure must not hide a terminal chain failure")

		var af errorinfo.AuthenticationFailure
		require.ErrorAs(t, err, &af)
		info := af.AuthenticationFailure()
		require.True(t, info.Chained)
		require.False(t, info.Unavailable)
		require.Len(t, info.Causes, 1)
		nested, ok := info.Causes[0].(errorinfo.AuthenticationFailure)
		require.True(t, ok)
		require.Equal(t, errorinfo.AuthenticationFailureInfo{Chained: true, Causes: []error{unavailable, terminal}}, nested.AuthenticationFailure())
		require.NotErrorIs(t, err, terminal, "chain errors must not expose their causes to errors.Is")
	})
	t.Run("all unavailable", func(t *testing.T) {
		first := newCredentialUnavailableError("first", "not configured")
		second := newCredentialUnavailableError("second", "not configured")
		c1, c2 := NewFakeCredential(), NewFakeCredential()
		c1.SetResponse(azcore.AccessToken{}, first)
		c2.SetResponse(azcore.AccessToken{}, second)
		cred, err := NewChainedTokenCredential([]azcore.TokenCredential{c1, c2}, nil)
		require.NoError(t, err)
		_, err = cred.GetToken(context.Background(), testTRO)
		var af errorinfo.AuthenticationFailure
		require.ErrorAs(t, err, &af)
		require.Equal(t, errorinfo.AuthenticationFailureInfo{Unavailable: true, Chained: true, Causes: []error{first, second}}, af.AuthenticationFailure())
	})
}

type fakeCredentialResponse struct {
	token azcore.AccessToken
	err   error
}

type fakeCredential struct {
	getTokenCalls int
	mut           *sync.Mutex
	responses     []fakeCredentialResponse
	static        *fakeCredentialResponse
}

func NewFakeCredential() *fakeCredential {
	return &fakeCredential{mut: &sync.Mutex{}}
}

func (c *fakeCredential) SetResponse(tk azcore.AccessToken, err error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.static = &fakeCredentialResponse{tk, err}
}

func (c *fakeCredential) AppendResponse(tk azcore.AccessToken, err error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.responses = append(c.responses, fakeCredentialResponse{tk, err})
}

func (c *fakeCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.getTokenCalls += 1
	if c.static != nil {
		return c.static.token, c.static.err
	}
	response := c.responses[0]
	c.responses = c.responses[1:]
	return response.token, response.err
}

func testGoodGetTokenResponse(t *testing.T, token azcore.AccessToken, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if v := token.Token; v != tokenValue {
		t.Fatalf(`unexpected token "%s"`, v)
	}
	if token.ExpiresOn.IsZero() {
		t.Fatal("token's ExpiresOn is zero value")
	}
}

func TestChainedTokenCredential_NilSource(t *testing.T) {
	_, err := NewChainedTokenCredential([]azcore.TokenCredential{NewFakeCredential(), nil}, nil)
	if err == nil {
		t.Fatalf("Expected an error for sending a nil credential in the chain")
	}
	_, err = NewChainedTokenCredential([]azcore.TokenCredential{}, nil)
	if err == nil {
		t.Fatalf("Expected an error for not sending any credential sources")
	}
}

func TestChainedTokenCredential_GetTokenSuccess(t *testing.T) {
	// ChainedTokenCredential should continue iterating when a source returns credentialUnavailableError, wrapped or not
	c1 := NewFakeCredential()
	c1.SetResponse(azcore.AccessToken{}, newCredentialUnavailableError("test", "something went wrong"))
	c2 := NewFakeCredential()
	c2.SetResponse(azcore.AccessToken{}, fmt.Errorf("%w", newCredentialUnavailableError("...", "...")))
	c3 := NewFakeCredential()
	c3.SetResponse(azcore.AccessToken{Token: tokenValue, ExpiresOn: time.Now().Add(time.Hour)}, nil)
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{c1, c2, c3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cred.GetToken(context.Background(), testTRO)
	if err != nil {
		t.Fatal(err)
	}
	if v := tk.Token; v != tokenValue {
		t.Fatalf(`unexpected token "%s"`, v)
	}
	if tk.ExpiresOn.IsZero() {
		t.Fatal("Received an incorrect time in the response")
	}
	if count := c1.getTokenCalls; count != 1 {
		t.Fatalf("expected 1 GetToken call, got %d", count)
	}
	if count := c2.getTokenCalls; count != 1 {
		t.Fatalf("expected 1 GetToken call, got %d", count)
	}
}

func TestChainedTokenCredential_GetTokenFail(t *testing.T) {
	c := NewFakeCredential()
	c.SetResponse(azcore.AccessToken{}, newAuthenticationFailedError("test", "something went wrong", nil))
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{c}, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if _, ok := err.(*AuthenticationFailedError); !ok {
		t.Fatalf("expected AuthenticationFailedError, received %T", err)
	}
	if !strings.Contains(err.Error(), "something went wrong") {
		t.Fatalf(`unexpected error message "%s"`, err.Error())
	}
}

func TestChainedTokenCredential_MultipleCredentialsGetTokenUnavailable(t *testing.T) {
	c1 := NewFakeCredential()
	c1.SetResponse(azcore.AccessToken{}, newCredentialUnavailableError("unavailableCredential1", "Unavailable expected error"))
	c2 := NewFakeCredential()
	c2.SetResponse(azcore.AccessToken{}, newCredentialUnavailableError("unavailableCredential2", "Unavailable expected error"))
	c3 := NewFakeCredential()
	c3.SetResponse(azcore.AccessToken{}, newCredentialUnavailableError("unavailableCredential3", "Unavailable expected error"))
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{c1, c2, c3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if _, ok := err.(credentialUnavailable); !ok {
		t.Fatalf("expected credentialUnavailable, received %T", err)
	}
	expectedError := `ChainedTokenCredential: failed to acquire a token.
Attempted credentials:
	unavailableCredential1: Unavailable expected error
	unavailableCredential2: Unavailable expected error
	unavailableCredential3: Unavailable expected error`
	if err.Error() != expectedError {
		t.Fatalf("Did not create an appropriate error message.\n\nReceived:\n%s\n\nExpected:\n%s", err.Error(), expectedError)
	}
}

func TestChainedTokenCredential_MultipleCredentialsGetTokenAuthenticationFailed(t *testing.T) {
	c1 := NewFakeCredential()
	c1.SetResponse(azcore.AccessToken{}, newCredentialUnavailableError("unavailableCredential1", "Unavailable expected error"))
	c2 := NewFakeCredential()
	c2.SetResponse(azcore.AccessToken{}, newCredentialUnavailableError("unavailableCredential2", "Unavailable expected error"))
	c3 := NewFakeCredential()
	c3.SetResponse(azcore.AccessToken{}, newAuthenticationFailedError("authenticationFailedCredential3", "Authentication failed expected error", nil))
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{c1, c2, c3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if _, ok := err.(*AuthenticationFailedError); !ok {
		t.Fatalf("expected AuthenticationFailedError, received %T", err)
	}
	expectedError := `ChainedTokenCredential: failed to acquire a token.
Attempted credentials:
	unavailableCredential1: Unavailable expected error
	unavailableCredential2: Unavailable expected error
	authenticationFailedCredential3: Authentication failed expected error`
	if err.Error() != expectedError {
		t.Fatalf("Did not create an appropriate error message.\n\nReceived:\n%s\n\nExpected:\n%s", err.Error(), expectedError)
	}
}

func TestChainedTokenCredential_MultipleCredentialsGetTokenCustomName(t *testing.T) {
	c := NewFakeCredential()
	c.SetResponse(azcore.AccessToken{}, newCredentialUnavailableError("unavailableCredential1", "Unavailable expected error"))
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{c}, nil)
	if err != nil {
		t.Fatal(err)
	}
	cred.name = "CustomNameCredential"
	_, err = cred.GetToken(context.Background(), testTRO)
	if _, ok := err.(credentialUnavailable); !ok {
		t.Fatalf("expected credentialUnavailable, received %T", err)
	}
	expectedError := `CustomNameCredential: failed to acquire a token.
Attempted credentials:
	unavailableCredential1: Unavailable expected error`
	if err.Error() != expectedError {
		t.Fatalf("Did not create an appropriate error message.\n\nReceived:\n%s\n\nExpected:\n%s", err.Error(), expectedError)
	}
}

func TestChainedTokenCredential_RepeatedGetTokenWithSuccessfulCredential(t *testing.T) {
	failedCredential := NewFakeCredential()
	failedCredential.SetResponse(azcore.AccessToken{}, newCredentialUnavailableError("MockCredential", "Mocking a credential unavailable error"))
	successfulCredential := NewFakeCredential()
	successfulCredential.SetResponse(azcore.AccessToken{Token: tokenValue, ExpiresOn: time.Now()}, nil)

	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{failedCredential, successfulCredential}, nil)
	if err != nil {
		t.Fatal(err)
	}

	tk, err := cred.GetToken(context.Background(), testTRO)
	testGoodGetTokenResponse(t, tk, err)
	if failedCredential.getTokenCalls != 1 {
		t.Fatal("The failed credential getToken should have been called once")
	}
	if successfulCredential.getTokenCalls != 1 {
		t.Fatalf("The successful credential getToken should have been called once")
	}
	tk2, err2 := cred.GetToken(context.Background(), testTRO)
	testGoodGetTokenResponse(t, tk2, err2)
	if failedCredential.getTokenCalls != 1 {
		t.Fatalf("The failed credential getToken should not have been called again")
	}
	if successfulCredential.getTokenCalls != 2 {
		t.Fatalf("The successful credential getToken should have been called twice")
	}
}

func TestChainedTokenCredential_RepeatedGetTokenWithSuccessfulCredentialWithRetrySources(t *testing.T) {
	failedCredential := NewFakeCredential()
	failedCredential.SetResponse(azcore.AccessToken{}, newCredentialUnavailableError("MockCredential", "Mocking a credential unavailable error"))
	successfulCredential := NewFakeCredential()
	successfulCredential.SetResponse(azcore.AccessToken{Token: tokenValue, ExpiresOn: time.Now()}, nil)

	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{failedCredential, successfulCredential}, &ChainedTokenCredentialOptions{RetrySources: true})
	if err != nil {
		t.Fatal(err)
	}

	tk, err := cred.GetToken(context.Background(), testTRO)
	testGoodGetTokenResponse(t, tk, err)
	if failedCredential.getTokenCalls != 1 {
		t.Fatalf("The failed credential getToken should have been called once")
	}
	if successfulCredential.getTokenCalls != 1 {
		t.Fatalf("The successful credential getToken should have been called once")
	}
	tk2, err2 := cred.GetToken(context.Background(), testTRO)
	testGoodGetTokenResponse(t, tk2, err2)
	if failedCredential.getTokenCalls != 2 {
		t.Fatalf("The failed credential getToken should have been called twice")
	}
	if successfulCredential.getTokenCalls != 2 {
		t.Fatalf("The successful credential getToken should have been called twice")
	}
}

func TestChainedTokenCredential_Race(t *testing.T) {
	successFake := NewFakeCredential()
	successFake.SetResponse(azcore.AccessToken{Token: "*", ExpiresOn: time.Now().Add(time.Hour)}, nil)
	authFailFake := NewFakeCredential()
	authFailFake.SetResponse(azcore.AccessToken{}, newAuthenticationFailedError("", "", nil))
	unavailableFake := NewFakeCredential()
	unavailableFake.SetResponse(azcore.AccessToken{}, newCredentialUnavailableError("", ""))

	for _, b := range []bool{true, false} {
		t.Run(fmt.Sprintf("RetrySources_%v", b), func(t *testing.T) {
			success, _ := NewChainedTokenCredential(
				[]azcore.TokenCredential{successFake}, &ChainedTokenCredentialOptions{RetrySources: b},
			)
			failure, _ := NewChainedTokenCredential(
				[]azcore.TokenCredential{authFailFake}, &ChainedTokenCredentialOptions{RetrySources: b},
			)
			unavailable, _ := NewChainedTokenCredential(
				[]azcore.TokenCredential{unavailableFake}, &ChainedTokenCredentialOptions{RetrySources: b},
			)
			for i := 0; i < 5; i++ {
				go func() {
					_, _ = success.GetToken(context.Background(), testTRO)
					_, _ = failure.GetToken(context.Background(), testTRO)
					_, _ = unavailable.GetToken(context.Background(), testTRO)
				}()
			}
		})
	}
}
