// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package markdown

import (
	"go/constant"
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstantExpression(t *testing.T) {
	real := constant.BinaryOp(constant.MakeFloat64(1), token.QUO, constant.MakeFloat64(3))
	imaginary := constant.BinaryOp(constant.MakeFloat64(-2), token.QUO, constant.MakeFloat64(7))
	values := []struct {
		name  string
		value constant.Value
	}{
		{name: "boolean", value: constant.MakeBool(true)},
		{name: "string", value: constant.MakeString("quoted \"text\"\n")},
		{name: "large integer", value: constant.Shift(constant.MakeInt64(1), token.SHL, 200)},
		{name: "negative integer", value: constant.MakeInt64(-42)},
		{name: "integral float", value: constant.MakeFloat64(1)},
		{name: "rational float", value: constant.BinaryOp(constant.MakeFloat64(1), token.QUO, constant.MakeFloat64(3))},
		{name: "large float", value: constant.MakeFromLiteral("1e1000", token.FLOAT, 0)},
		{name: "rational complex", value: constant.BinaryOp(real, token.ADD, constant.MakeImag(imaginary))},
	}
	for _, tt := range values {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := constantExpression(tt.value)
			require.NoError(t, err)
			parsed, err := types.Eval(token.NewFileSet(), nil, token.NoPos, expr)
			require.NoError(t, err, expr)
			assert.Equal(t, tt.value.Kind(), parsed.Value.Kind(), "constant kind changed: %s", expr)
			assert.True(t, constant.Compare(tt.value, token.EQL, parsed.Value), "constant value changed: %s", expr)
		})
	}
	_, err := constantExpression(constant.MakeUnknown())
	require.ErrorContains(t, err, "unsupported constant value kind")
}
