// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ChainedTokenCredential provides a TokenCredential implementation that chains multiple TokenCredential sources to be tried in order
// and returns the token from the first successful call to GetToken().
type ChainedTokenCredential struct {
	sources []azcore.TokenCredential
}

// NewChainedTokenCredential creates an instance of ChainedTokenCredential with the specified TokenCredential sources.
func NewChainedTokenCredential(sources ...azcore.TokenCredential) (*ChainedTokenCredential, error) {
	if len(sources) == 0 {
		credErr := &CredentialUnavailableError{CredentialType: "Chained Token Credential", Message: "Length of sources cannot be 0"}
		azcore.Log().Write(azcore.LogError, logCredentialError(credErr.CredentialType, credErr))
		return nil, credErr
	}
	for _, source := range sources {
		if source == nil { // cannot have a nil credential in the chain or else the application will panic when GetToken() is called on nil
			credErr := &CredentialUnavailableError{CredentialType: "Chained Token Credential", Message: "Sources cannot contain a nil TokenCredential"}
			azcore.Log().Write(azcore.LogError, logCredentialError(credErr.CredentialType, credErr))
			return nil, credErr
		}
	}
	return &ChainedTokenCredential{sources: sources}, nil
}

// GetToken sequentially calls TokenCredential.GetToken on all the specified sources, returning the token from the first successful call to GetToken().
func (c *ChainedTokenCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (token *azcore.AccessToken, err error) {
	var errList []*CredentialUnavailableError
	for _, cred := range c.sources { // loop through all of the credentials provided in sources
		token, err = cred.GetToken(ctx, opts) // make a GetToken request for the current credential in the loop
		var credErr *CredentialUnavailableError
		if errors.As(err, &credErr) { // check if we received a CredentialUnavailableError
			errList = append(errList, credErr) // if we did receive a CredentialUnavailableError then we append it to our error slice and continue looping for a good credential
		} else if err != nil { // if we receive some other type of error then we must stop looping and process the error accordingly
			var authenticationFailed *AuthenticationFailedError
			if errors.As(err, &authenticationFailed) { // if the error is an AuthenticationFailedError we return the error related to the invalid credential and append all of the other error messages received prior to this point
				authErr := &AuthenticationFailedError{msg: "Received an AuthenticationFailedError, there is an invalid credential in the chain. " + createChainedErrorMessage(errList), inner: err}
				addGetTokenFailureLogs("Chained Token Credential", authErr)
				return nil, authErr
			}
			addGetTokenFailureLogs("Chained Token Credential", err)
			return nil, err // if we receive some other error type this is unexpected and we simple return the unexpected error
		} else {
			azcore.Log().Write(LogCredential, logGetTokenSuccess(c, opts))
			return token, nil // if we did not receive an error then we return the token
		}
	}
	// if we reach this point it means that all of the credentials in the chain returned CredentialUnavailableErrors
	credErr := &CredentialUnavailableError{CredentialType: "Chained Token Credential", Message: createChainedErrorMessage(errList)}
	addGetTokenFailureLogs("Chained Token Credential", credErr)
	return nil, credErr
}

// AuthenticationPolicy implements the azcore.Credential interface on ChainedTokenCredential and sets the bearer token
func (c *ChainedTokenCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}

// helper function used to chain the error messages of the CredentialUnavailableError slice
func createChainedErrorMessage(errList []*CredentialUnavailableError) string {
	msg := ""
	for _, err := range errList {
		msg += err.Error()
	}

	return msg
}
