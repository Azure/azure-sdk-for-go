//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

//
//import "github.com/stretchr/testify/require"
//
////nolint
//func (s *azblobUnrecordedTestSuite) TestDeserializeORSPolicies() {
//	_require := require.New(s.T())
//
//	headers := map[string]string{
//		"x-ms-or-111_111":   "Completed",
//		"x-ms-or-111_222":   "Failed",
//		"x-ms-or-222_111":   "Completed",
//		"x-ms-or-222_222":   "Failed",
//		"x-ms-or-policy-id": "333",     // to be ignored
//		"x-ms-not-related":  "garbage", // to be ignored
//	}
//
//	result := deserializeORSPolicies(headers)
//	_require.NotNil(result)
//	rules0, rules1 := *result[0].Rules, *result[1].Rules
//	_require.Len(result, 2)
//	_require.Len(rules0, 2)
//	_require.Len(rules1, 2)
//
//	if rules0[0].RuleId == "111" {
//		_require.Equal(rules0[0].Status, "Completed")
//	} else {
//		_require.Equal(rules0[0].Status, "Failed")
//	}
//
//	if rules0[1].RuleId == "222" {
//		_require.Equal(rules0[1].Status, "Failed")
//	} else {
//		_require.Equal(rules0[1].Status, "Completed")
//	}
//
//	if rules1[0].RuleId == "111" {
//		_require.Equal(rules1[0].Status, "Completed")
//	} else {
//		_require.Equal(rules1[0].Status, "Failed")
//	}
//
//	if rules1[1].RuleId == "222" {
//		_require.Equal(rules1[1].Status, "Failed")
//	} else {
//		_require.Equal(rules1[1].Status, "Completed")
//	}
//}
//
////func (s * azblobUnrecordedTestSuite) TestORSSource() {
////	_require := require.New(s.T())
////	testName := s.T().Name()
////	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
////	if err != nil {
////		s.Fail("Unable to fetch service client because " + err.Error())
////	}
////
////	containerName := generateContainerName(testName)
////	containerClient := createNewContainer(_require, containerName, svcClient)
////	defer deleteContainer(_require, containerClient)
////
////	bbName := generateBlobName(testName)
////	bbClient := createNewBlockBlob(_require, bbName, containerClient)
////
////	getResp, err := bbClient.GetProperties(ctx, nil)
////	_require.Nil(err)
////	_require.Nil(getResp.ObjectReplicationRules)
////}
