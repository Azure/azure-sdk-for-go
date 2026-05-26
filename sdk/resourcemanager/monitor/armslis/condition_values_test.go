// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armslis

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestCondition_Values_RoundTrip(t *testing.T) {
	c := &Condition{Value: to.Ptr("east^^west^^north")}
	require.Equal(t, []string{"east", "west", "north"}, c.Values())
}

func TestCondition_SetValues_Joins(t *testing.T) {
	c := &Condition{}
	c.SetValues([]string{"east", "west", "north"})
	require.NotNil(t, c.Value)
	require.Equal(t, "east^^west^^north", *c.Value)
}

func TestCondition_SetValues_NilClearsValue(t *testing.T) {
	c := &Condition{Value: to.Ptr("east^^west")}
	c.SetValues(nil)
	require.Nil(t, c.Value)
}

func TestCondition_Values_NilWhenValueNil(t *testing.T) {
	c := &Condition{}
	require.Nil(t, c.Values())
}

func TestNewConditionWithValues_In(t *testing.T) {
	c, err := NewConditionWithValues(ConditionOperatorIn, []string{"east", "west"}, to.Ptr("region"), nil, nil)
	require.NoError(t, err)
	require.NotNil(t, c)
	require.Equal(t, ConditionOperatorIn, *c.Operator)
	require.Equal(t, "east^^west", *c.Value)
	require.Equal(t, "region", *c.DimensionName)
}

func TestNewConditionWithValues_NotIn(t *testing.T) {
	c, err := NewConditionWithValues(ConditionOperatorNotIn, []string{"only"}, nil, nil, nil)
	require.NoError(t, err)
	require.Equal(t, ConditionOperatorNotIn, *c.Operator)
	require.Equal(t, "only", *c.Value)
}

func TestNewConditionWithValues_RejectsWrongOperator(t *testing.T) {
	_, err := NewConditionWithValues(ConditionOperatorEqual, []string{"x"}, nil, nil, nil)
	require.Error(t, err)
}

func TestNewConditionWithValues_RejectsEmpty(t *testing.T) {
	_, err := NewConditionWithValues(ConditionOperatorIn, nil, nil, nil, nil)
	require.Error(t, err)
}

func TestNewConditionWithValues_RejectsSeparatorInItem(t *testing.T) {
	_, err := NewConditionWithValues(ConditionOperatorIn, []string{"ok", "bad^^value"}, nil, nil, nil)
	require.Error(t, err)
}
