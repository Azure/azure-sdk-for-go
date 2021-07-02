// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"strconv"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func (s *tableClientLiveTests) TestSetAccessPolicy() {
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip("TableAccessPolicies are not available on Cosmos Accounts")
	}

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

func (s *tableClientLiveTests) TestSetMultipleAccessPolicies() {
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip("TableAccessPolicies are not available on Cosmos Accounts")
	}

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

	expiration = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	permission = "rw"
	id = "2"

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

	// Make a Get to assert two access policies
	resp, err := client.GetAccessPolicy(ctx)
	assert.Nil(err, "Get Access Policy failed")
	assert.Equal(len(resp.SignedIdentifiers), 2)
}

func (s *tableClientLiveTests) TestSetTooManyAccessPolicies() {
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip("TableAccessPolicies are not available on Cosmos Accounts")
	}

	assert := assert.New(s.T())
	// context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	expiration := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	permission := "r"
	id := "1"
	signedIdentifiers := make([]*SignedIdentifier, 0)

	for i := 0; i < 6; i++ {
		expiration = time.Date(2024+i, 1, 1, 0, 0, 0, 0, time.UTC)
		id = strconv.Itoa(i)

		signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
			AccessPolicy: &AccessPolicy{
				Expiry:     &expiration,
				Start:      &start,
				Permission: &permission,
			},
			ID: &id,
		})

	}

	param := TableSetAccessPolicyOptions{TableACL: signedIdentifiers}

	_, err := client.SetAccessPolicy(ctx, &param)
	assert.Nil(err, "Set access policy failed")

	// Make a Get to assert two access policies
	_, err = client.GetAccessPolicy(ctx)
	assert.NotNil(err)
	// assert.Nil(err, "Get Access Policy failed")
	// assert.Equal(len(resp.SignedIdentifiers), 2)
}
