// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/delta"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/report"
	"github.com/stretchr/testify/assert"
)

func TestCalculateNewVersion(t *testing.T) {
	fixChange := &model.Changelog{Modified: &report.Package{}}
	breakingChange := &model.Changelog{RemovedPackage: true, Modified: &report.Package{}}
	additiveChange := &model.Changelog{Modified: &report.Package{AdditiveChanges: &delta.Content{Content: exports.Content{Consts: map[string]exports.Const{"test": {}}}}}}

	// previous 0.x.x
	// fix with stable
	newVersion, prl, err := CalculateNewVersion(fixChange, "0.5.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.0.0")
	assert.Equal(t, FirstGALabel, prl)

	// fix with beat
	newVersion, prl, err = CalculateNewVersion(fixChange, "0.5.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "0.5.1")
	assert.Equal(t, BetaLabel, prl)

	// breaking with stable
	newVersion, prl, err = CalculateNewVersion(breakingChange, "0.5.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.0.0")
	assert.Equal(t, FirstGABreakingChangeLabel, prl)

	// breaking with beta
	newVersion, prl, err = CalculateNewVersion(breakingChange, "0.5.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "0.6.0")
	assert.Equal(t, BetaBreakingChangeLabel, prl)

	// additive with stable
	newVersion, prl, err = CalculateNewVersion(additiveChange, "0.5.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.0.0")
	assert.Equal(t, FirstGALabel, prl)

	// additive with beta
	newVersion, prl, err = CalculateNewVersion(additiveChange, "0.5.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "0.6.0")
	assert.Equal(t, BetaLabel, prl)

	// previous 1.2.0
	// fix with stable
	newVersion, prl, err = CalculateNewVersion(fixChange, "1.2.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.1")
	assert.Equal(t, StableLabel, prl)

	// fix with beat
	newVersion, prl, err = CalculateNewVersion(fixChange, "1.2.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.1-beta.1")
	assert.Equal(t, BetaLabel, prl)

	// breaking with stable
	newVersion, prl, err = CalculateNewVersion(breakingChange, "1.2.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "2.0.0")
	assert.Equal(t, StableBreakingChangeLabel, prl)

	// breaking with beta
	newVersion, prl, err = CalculateNewVersion(breakingChange, "1.2.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "2.0.0-beta.1")
	assert.Equal(t, BetaBreakingChangeLabel, prl)

	// additive with stable
	newVersion, prl, err = CalculateNewVersion(additiveChange, "1.2.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.3.0")
	assert.Equal(t, StableLabel, prl)

	// additive with beta
	newVersion, prl, err = CalculateNewVersion(additiveChange, "1.2.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.3.0-beta.1")
	assert.Equal(t, BetaLabel, prl)

	// previous 1.2.0-beta.1
	// fix with stable
	newVersion, _, err = CalculateNewVersion(fixChange, "1.2.0-beta.1", false)
	assert.NotEmpty(t, err)

	// fix with beat
	newVersion, prl, err = CalculateNewVersion(fixChange, "1.2.0-beta.1", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.0-beta.2")
	assert.Equal(t, BetaLabel, prl)

	// breaking with stable
	newVersion, _, err = CalculateNewVersion(breakingChange, "1.2.0-beta.1", false)
	assert.NotEmpty(t, err)

	// breaking with beta
	newVersion, prl, err = CalculateNewVersion(breakingChange, "1.2.0-beta.1", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.0-beta.2")
	assert.Equal(t, BetaBreakingChangeLabel, prl)

	// additive with stable
	newVersion, _, err = CalculateNewVersion(additiveChange, "1.2.0-beta.1", false)
	assert.NotEmpty(t, err)

	// additive with beta
	newVersion, prl, err = CalculateNewVersion(additiveChange, "1.2.0-beta.1", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.0-beta.2")
	assert.Equal(t, BetaLabel, prl)
}
