// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import "github.com/stretchr/testify/assert"

//nolint
func (s *azblobUnrecordedTestSuite) TestDeserializeORSPolicies() {
	_assert := assert.New(s.T())

	headers := map[string]string{
		"x-ms-or-111_111":   "Completed",
		"x-ms-or-111_222":   "Failed",
		"x-ms-or-222_111":   "Completed",
		"x-ms-or-222_222":   "Failed",
		"x-ms-or-policy-id": "333",     // to be ignored
		"x-ms-not-related":  "garbage", // to be ignored
	}

	result := deserializeORSPolicies(headers)
	_assert.NotNil(result)
	rules0, rules1 := *result[0].Rules, *result[1].Rules
	_assert.Len(result, 2)
	_assert.Len(rules0, 2)
	_assert.Len(rules1, 2)

	if rules0[0].RuleId == "111" {
		_assert.Equal(rules0[0].Status, "Completed")
	} else {
		_assert.Equal(rules0[0].Status, "Failed")
	}

	if rules0[1].RuleId == "222" {
		_assert.Equal(rules0[1].Status, "Failed")
	} else {
		_assert.Equal(rules0[1].Status, "Completed")
	}

	if rules1[0].RuleId == "111" {
		_assert.Equal(rules1[0].Status, "Completed")
	} else {
		_assert.Equal(rules1[0].Status, "Failed")
	}

	if rules1[1].RuleId == "222" {
		_assert.Equal(rules1[1].Status, "Failed")
	} else {
		_assert.Equal(rules1[1].Status, "Completed")
	}
}

//func (s * azblobUnrecordedTestSuite) TestORSSource() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	bbName := generateBlobName(testName)
//	bbClient := createNewBlockBlob(_assert, bbName, containerClient)
//
//	getResp, err := bbClient.GetProperties(ctx, nil)
//	_assert.Nil(err)
//	_assert.Nil(getResp.ObjectReplicationRules)
//}
