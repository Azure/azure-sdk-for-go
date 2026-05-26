// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armslis

import (
	"errors"
	"fmt"
	"strings"
)

// conditionInValueSeparator is the literal delimiter the SLI resource provider expects between
// list items for the [ConditionOperatorIn] and [ConditionOperatorNotIn] operators.
const conditionInValueSeparator = "^^"

// Values returns the items encoded into Condition.Value for the [ConditionOperatorIn] /
// [ConditionOperatorNotIn] operators by splitting on the literal "^^" separator. For all other
// operators the returned slice will have a single element equal to Value.
//
// Returns nil when c.Value is nil. An empty Value returns a slice with one empty string, matching
// strings.Split semantics.
func (c *Condition) Values() []string {
	if c == nil || c.Value == nil {
		return nil
	}
	return strings.Split(*c.Value, conditionInValueSeparator)
}

// SetValues populates Condition.Value by joining values with the literal "^^" separator used on
// the wire for the [ConditionOperatorIn] / [ConditionOperatorNotIn] operators. Passing a nil slice
// clears Value; passing an empty slice sets Value to an empty string.
func (c *Condition) SetValues(values []string) {
	if c == nil {
		return
	}
	if values == nil {
		c.Value = nil
		return
	}
	joined := strings.Join(values, conditionInValueSeparator)
	c.Value = &joined
}

// NewConditionWithValues returns a Condition for a list operator ([ConditionOperatorIn] or
// [ConditionOperatorNotIn]) by joining values with the wire "^^" separator. dimensionName,
// scalarFunction, and samplingType are optional and may be nil. An error is returned if op is not
// a list operator, values is empty, or any item contains the "^^" separator.
func NewConditionWithValues(op ConditionOperator, values []string, dimensionName *string, scalarFunction *ScalarFunction, samplingType *SamplingType) (*Condition, error) {
	if op != ConditionOperatorIn && op != ConditionOperatorNotIn {
		return nil, fmt.Errorf("armslis: NewConditionWithValues operator must be ConditionOperatorIn or ConditionOperatorNotIn; got %q", op)
	}
	if len(values) == 0 {
		return nil, errors.New("armslis: NewConditionWithValues requires at least one value")
	}
	for i, v := range values {
		if strings.Contains(v, conditionInValueSeparator) {
			return nil, fmt.Errorf("armslis: value at index %d contains the reserved %q separator", i, conditionInValueSeparator)
		}
	}
	c := &Condition{
		Operator:       &op,
		DimensionName:  dimensionName,
		ScalarFunction: scalarFunction,
		SamplingType:   samplingType,
	}
	c.SetValues(values)
	return c, nil
}
