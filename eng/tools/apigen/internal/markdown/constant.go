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
)

func constantExpression(value constant.Value) (string, error) {
	switch value.Kind() {
	case constant.Bool, constant.String, constant.Int:
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
		real, err := constantExpression(constant.Real(value))
		if err != nil {
			return "", err
		}
		imaginary, err := constantExpression(constant.Imag(value))
		if err != nil {
			return "", err
		}
		return "(" + real + " + (" + imaginary + ")*1i)", nil
	default:
		return "", fmt.Errorf("unsupported constant value kind %v", value.Kind())
	}
}
