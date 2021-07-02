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

	assert := assert.New(s.T())
	// context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	now := time.Now()
	expiration := now.Add(time.Duration(time.Hour * 1))
	permission := "rw"

	accessPolicies := make([]*AccessPolicy, 0)
	accessPolicies = append(accessPolicies, &AccessPolicy{
		Expiry:     &expiration,
		Start:      &now,
		Permission: &permission,
	})

	_, err := client.SetAccessPolicy(ctx, accessPolicies)
	assert.Nil(err, "Set access policy failed")
}

func (s *tableClientLiveTests) TestSetEmptyAccessPolicy() {
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip("TableAccessPolicies are not available on Cosmos Accounts")
	}

	assert := assert.New(s.T())
	// context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	accessPolicies := make([]*AccessPolicy, 0)

	_, err := client.SetAccessPolicy(ctx, accessPolicies)
	assert.Nil(err, "Set access policy failed")
}
