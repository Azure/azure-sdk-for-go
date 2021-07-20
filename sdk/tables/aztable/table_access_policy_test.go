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
	client, delete := s.init(true)
	defer delete()

	id := "empty"

	signedIdentifiers := make([]*SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: &id,
	})

	permission2 := "r"
	id2 := "partial"

	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: &id2,
		AccessPolicy: &AccessPolicy{
			Permission: &permission2,
		},
	})

	id3 := "full"
	permission3 := "r"
	start := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)
	expiry := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)

	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: &id3,
		AccessPolicy: &AccessPolicy{
			Start:      &start,
			Expiry:     &expiry,
			Permission: &permission3,
		},
	})

	param := TableSetAccessPolicyOptions{
		TableACL: signedIdentifiers,
	}

	_, err := client.SetAccessPolicy(ctx, &param)
	if err != nil {
		assert.FailNow("Set access policy failed")
	}

	// Make a Get to assert two access policies
	resp, err := client.GetAccessPolicy(ctx)
	assert.Nil(err, "Get Access Policy failed")
	assert.Equal(len(resp.SignedIdentifiers), 3)
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
	assert.NotNil(err, "Set access policy succeeded but should have failed")
	assert.Contains(err.Error(), tooManyAccessPoliciesError.Error())
}

func (s *tableClientLiveTests) TestSetNullAccessPolicy() {
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip("TableAccessPolicies are not available on Cosmos Accounts")
	}

	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	id := "null"

	signedIdentifiers := make([]*SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: &id,
	})

	param := TableSetAccessPolicyOptions{
		TableACL: signedIdentifiers,
	}

	_, err := client.SetAccessPolicy(ctx, &param)
	if err != nil {
		assert.FailNow("Set access policy failed")
	}

	resp, err := client.GetAccessPolicy(ctx)
	assert.Nil(err, "Get Access Policy failed")
	assert.Equal(len(resp.SignedIdentifiers), 1)
}
