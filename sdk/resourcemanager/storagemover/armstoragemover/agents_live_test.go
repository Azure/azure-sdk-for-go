// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// AgentScenarioSuite mirrors .NET AgentTests. Both methods require an agent VM that has been
// registered with the Storage Mover RP via the Azure Connected Machine onboarding flow. The RP does
// not allow callers to create agents directly, so the only way to record this test is against a
// pre-existing registered agent. We skip with reason instead of attempting `Agents_CreateOrUpdate`
// which would fail with "Agent not registered".
type AgentScenarioSuite struct {
	scenarioBaseSuite
}

func TestAgentScenarioSuite(t *testing.T) {
	suite.Run(t, new(AgentScenarioSuite))
}

func (s *AgentScenarioSuite) SetupSuite() {
	s.T().Skip("Agent tests require a registered Hybrid Compute machine with the Storage Mover agent installed; the RP does not allow creating agents directly. See cross-language playbook for details.")
}

func (s *AgentScenarioSuite) TearDownSuite() {
	// Nothing to do — SetupSuite skipped.
}

// TestAgentGetExist mirrors .NET AgentTests.GetExistTest. Skipped — see SetupSuite.
func (s *AgentScenarioSuite) TestAgentGetExist() {
	s.T().Skip("see SetupSuite")
}
