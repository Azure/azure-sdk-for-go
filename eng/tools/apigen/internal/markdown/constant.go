// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package markdown

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"strconv"
	"unicode/utf8"
)

func constantExpression(value constant.Value, typ types.Type) (string, error) {
	switch value.Kind() {
	case constant.Int:
		if typ == types.Typ[types.UntypedRune] {
			codePoint, ok := constant.Int64Val(value)
			if ok && codePoint >= 0 && codePoint <= utf8.MaxRune && utf8.ValidRune(rune(codePoint)) {
				return strconv.QuoteRuneToASCII(rune(codePoint)), nil
			}
			// Untyped rune arithmetic can produce values outside the range of rune literals.
			return `'\x00' + (` + value.ExactString() + `)`, nil
		}
		return value.ExactString(), nil
	case constant.Bool, constant.String:
		return value.ExactString(), nil
	case constant.Float:
		expr, err := parser.ParseExpr(value.ExactString())
		if err != nil {
			return "", fmt.Errorf("parse exact floating-point value: %w", err)
		}
		ast.Inspect(expr, func(node ast.Node) bool {
			if literal, ok := node.(*ast.BasicLit); ok && literal.Kind == token.INT {
				literal.Kind = token.FLOAT
				literal.Value += ".0"
			}
			return true
		})
		var output bytes.Buffer
		if err := format.Node(&output, token.NewFileSet(), expr); err != nil {
			return "", fmt.Errorf("format exact floating-point value: %w", err)
		}
		return output.String(), nil
	case constant.Complex:
		real, err := constantExpression(constant.Real(value), nil)
		if err != nil {
			return "", err
		}
		imaginary, err := constantExpression(constant.Imag(value), nil)
		if err != nil {
			return "", err
		}
		return "(" + real + " + (" + imaginary + ")*1i)", nil
	default:
		return "", fmt.Errorf("unsupported constant value kind %v", value.Kind())
	}
}
