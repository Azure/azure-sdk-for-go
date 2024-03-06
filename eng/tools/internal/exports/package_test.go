// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exports_test

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
)

var pkg exports.Package
var exp exports.Content

func init() {
	pkg, _ = exports.LoadPackage("./testdata")
	exp = pkg.GetExports()
}

func Test_Constants(t *testing.T) {
	if l := len(exp.Consts); l != 13 {
		t.Logf("wrong number of constants, got %v", l)
		t.Fail()
	}

	tests := []struct {
		cn string
		exports.Const
	}{
		{"DefaultBaseURI", exports.Const{Type: "string", Value: "https://management.azure.com"}},
		{"Tuesday", exports.Const{Type: "DayOfWeek", Value: "Tuesday"}},
		{"Primary", exports.Const{Type: "KeyType", Value: "Primary"}},
		{"Backup", exports.Const{Type: "KeyType", Value: "Backup"}},
	}

	for _, test := range tests {
		t.Run(test.cn, func(t *testing.T) {
			c := exp.Consts[test.cn]
			if !reflect.DeepEqual(c.Type, test.Type) {
				t.Logf("mismatched types, %s != %s", c.Type, test.Type)
				t.Fail()
			}
			if c.Value != test.Value {
				t.Logf("mismatched values, %s != %s", c.Value, test.Value)
				t.Fail()
			}
		})
	}
}

func Test_TypeAliases(t *testing.T) {
	if l := len(exp.TypeAliases); l != 2 {
		t.Logf("wrong number of type aliases, got %v", l)
		t.Fail()
	}

	tests := []struct {
		cn string
		exports.TypeAlias
	}{
		{"DayOfWeek", exports.TypeAlias{UnderlayingType: "string"}},
		{"KeyType", exports.TypeAlias{UnderlayingType: "string"}},
	}

	for _, test := range tests {
		t.Run(test.cn, func(t *testing.T) {
			c := exp.TypeAliases[test.cn]
			if c.UnderlayingType != test.UnderlayingType {
				t.Logf("mismatched underlaying type, %s != %s", c.UnderlayingType, test.UnderlayingType)
				t.Fail()
			}
		})
	}
}

func Test_Funcs(t *testing.T) {
	if l := len(exp.Funcs); l != 21 {
		t.Logf("wrong number of funcs, got %v", l)
		t.Fail()
	}

	tests := []struct {
		fn string
		exports.Func
	}{
		{"DoNothing", exports.Func{}},
		{"DoNothingWithParam", exports.Func{Params: strPtr("int"), Returns: nil}},
		{"UserAgent", exports.Func{Params: nil, Returns: strPtr("string")}},
		{"Client.Delete", exports.Func{Params: strPtr("context.Context, string, string"), Returns: strPtr("DeleteFuture, error")}},
		{"Client.ListSender", exports.Func{Params: strPtr("*http.Request"), Returns: strPtr("*http.Response, error")}},
	}

	for _, test := range tests {
		t.Run(test.fn, func(t *testing.T) {
			f := exp.Funcs[test.fn]
			if !reflect.DeepEqual(f.Params, test.Params) {
				t.Logf("mismatched params, %s != %s", safeStr(f.Params), safeStr(test.Params))
				t.Fail()
			}
			if !reflect.DeepEqual(f.Returns, test.Returns) {
				t.Logf("mismatched returns, %s != %s", safeStr(f.Returns), safeStr(test.Returns))
				t.Fail()
			}
		})
	}
}

func Test_Interfaces(t *testing.T) {
	if l := len(exp.Interfaces); l != 1 {
		t.Logf("wrong number of interfaces, got %v", l)
		t.Fail()
	}

	tests := []struct {
		in string
		exports.Interface
	}{
		{"SomeInterface", exports.Interface{
			Methods: map[string]exports.Func{
				"One": {
					Params:  nil,
					Returns: nil,
				},
				"Two": {
					Params:  strPtr("bool"),
					Returns: nil,
				},
				"Three": {
					Params:  nil,
					Returns: strPtr("string"),
				},
				"Four": {
					Params:  strPtr("int"),
					Returns: strPtr("error"),
				},
				"Five": {
					Params:  strPtr("int, bool"),
					Returns: strPtr("int, error"),
				},
			},
		}},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			i := exp.Interfaces[test.in]
			if !reflect.DeepEqual(i.Methods, test.Methods) {
				t.Logf("mismatched methods, have %+v, want %+v", i.Methods, test.Methods)
				t.Fail()
			}
		})
	}
}

func Test_Structs(t *testing.T) {
	if l := len(exp.Structs); l != 8 {
		t.Logf("wrong number of structs, got %v", l)
		t.Fail()
	}

	tests := []struct {
		sn string
		exports.Struct
	}{
		{"BaseClient", exports.Struct{
			AnonymousFields: []string{"autorest.Client"},
			Fields: map[string]string{
				"BaseURI":        "string",
				"SubscriptionID": "string",
			}}},
		{"DeleteFuture", exports.Struct{
			AnonymousFields: []string{"azure.Future"},
		}},
		{"ListResultPage", exports.Struct{}},
		{"CreateParameters", exports.Struct{
			AnonymousFields: []string{"*CreateProperties"},
			Fields: map[string]string{
				"Zones":    "*[]string",
				"Location": "*string",
				"Tags":     "map[string]*string",
			}}},
		{"CreateProperties", exports.Struct{
			Fields: map[string]string{
				"SubnetID":           "*string",
				"StaticIP":           "*string",
				"RedisConfiguration": "map[string]*string",
				"EnableNonSslPort":   "*bool",
				"TenantSettings":     "map[string]*string",
				"ShardCount":         "*int32",
			}}},
	}

	for _, test := range tests {
		t.Run(test.sn, func(t *testing.T) {
			s := exp.Structs[test.sn]
			if !reflect.DeepEqual(s.AnonymousFields, test.AnonymousFields) {
				t.Logf("mismatched anonymous fields, have %+v, want %v", s.AnonymousFields, test.AnonymousFields)
				t.Fail()
			}
			if !reflect.DeepEqual(s.Fields, test.Fields) {
				t.Logf("mismatched fields, have %+v want %+v", s.Fields, test.Fields)
				t.Fail()
			}
		})
	}
}

func Test_Name(t *testing.T) {
	if n := pkg.Name(); n != "testdata" {
		t.Logf("incorrect package name, have '%s', want 'testdata'", n)
		t.Fail()
	}
}

func strPtr(s string) *string {
	return &s
}

func safeStr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
