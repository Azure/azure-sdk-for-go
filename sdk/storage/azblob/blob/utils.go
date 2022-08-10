//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"strings"
	"time"
)

// ObjectReplicationRules struct
type ObjectReplicationRules struct {
	RuleId string
	Status string
}

// ObjectReplicationPolicy are deserialized attributes
type ObjectReplicationPolicy struct {
	PolicyId *string
	Rules    *[]ObjectReplicationRules
}

// deserializeORSPolicies is utility function to deserialize ORS Policies
func deserializeORSPolicies(policies map[string]string) (objectReplicationPolicies []ObjectReplicationPolicy) {
	if policies == nil {
		return nil
	}
	// For source blobs (blobs that have policy ids and rule ids applied to them),
	// the header will be formatted as "x-ms-or-<policy_id>_<rule_id>: {Complete, Failed}".
	// The value of this header is the status of the replication.
	orPolicyStatusHeader := make(map[string]string)
	for key, value := range policies {
		if strings.Contains(key, "or-") && key != "x-ms-or-policy-id" {
			orPolicyStatusHeader[key] = value
		}
	}

	parsedResult := make(map[string][]ObjectReplicationRules)
	for key, value := range orPolicyStatusHeader {
		policyAndRuleIDs := strings.Split(strings.Split(key, "or-")[1], "_")
		policyId, ruleId := policyAndRuleIDs[0], policyAndRuleIDs[1]

		parsedResult[policyId] = append(parsedResult[policyId], ObjectReplicationRules{RuleId: ruleId, Status: value})
	}

	for policyId, rules := range parsedResult {
		objectReplicationPolicies = append(objectReplicationPolicies, ObjectReplicationPolicy{
			PolicyId: &policyId,
			Rules:    &rules,
		})
	}
	return
}

// ParseHTTPHeaders parses GetPropertiesResponse and returns HTTPHeaders
func ParseHTTPHeaders(resp GetPropertiesResponse) HTTPHeaders {
	return generated.ParseHTTPHeaders(resp)
}

// ParseSASTimeString try to parse sas time string.
func ParseSASTimeString(val string) (t time.Time, timeFormat string, err error) {
	return exported.ParseSASTimeString(val)
}
