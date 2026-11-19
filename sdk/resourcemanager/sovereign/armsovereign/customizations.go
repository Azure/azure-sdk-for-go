// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armsovereign

import "encoding/json"

// UnmarshalJSON accepts both the action result and the terminal polling envelope.
// Keeping this on the response preserves the standard poller's status, error,
// final-location, and resume-token behavior.
func (r *ViewOperationsClientRemediateResponse) UnmarshalJSON(data []byte) error {
	var result ViewRemediateResult
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.RemediationID == nil {
		var envelope struct {
			Properties *ViewRemediateResult `json:"properties"`
		}
		if err := json.Unmarshal(data, &envelope); err != nil {
			return err
		}
		if envelope.Properties != nil {
			result.RemediationID = envelope.Properties.RemediationID
			if result.Error == nil {
				result.Error = envelope.Properties.Error
			}
		}
	}
	r.ViewRemediateResult = result
	return nil
}
