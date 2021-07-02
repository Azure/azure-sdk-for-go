// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *tableClientLiveTests) TestSetAccessPolicy() {
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip("TableAccessPolicies are not available on Cosmos Accounts")
	}

	s.T().Skip("Skipping for now, Body of request needs to include ID")

	assert := assert.New(s.T())
	// context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	expiration := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	permission := "r"
	id := "1"

	signedIdentifiers := make([]*SignedIdentifier, 0)

	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		AccessPolicy: &AccessPolicy{
			Expiry:     &expiration,
			Start:      &start,
			Permission: &permission,
		},
		ID: &id,
	})

	param := TableSetAccessPolicyOptions{
		TableACL: signedIdentifiers,
	}

	_, err := client.SetAccessPolicy(ctx, &param)
	assert.Nil(err, "Set access policy failed")
}

func (s *tableClientLiveTests) TestSetEmptyAccessPolicy() {
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip("TableAccessPolicies are not available on Cosmos Accounts")
	}

	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	_, err := client.SetAccessPolicy(ctx, &TableSetAccessPolicyOptions{})
	assert.Nil(err, "Set access policy failed")
}
