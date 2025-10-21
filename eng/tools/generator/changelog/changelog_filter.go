// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package changelog

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/delta"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/report"
)

func FilterChangelog(changelog *Changelog, opts ...func(changelog *Changelog)) {
	if changelog.Modified != nil {
		for _, opt := range opts {
			opt(changelog)
		}
	}
}

func MarshalUnmarshalFilter(changelog *Changelog) {
	if changelog.Modified != nil {
		if changelog.Modified.AdditiveChanges != nil {
			removeMarshalUnmarshalFunc(changelog.Modified.AdditiveChanges.Added.Funcs)
		}
		if changelog.Modified.BreakingChanges != nil && changelog.Modified.BreakingChanges.Removed != nil {
			removeMarshalUnmarshalFunc(changelog.Modified.BreakingChanges.Removed.Funcs)
		}
	}
}

func removeMarshalUnmarshalFunc(funcs map[string]exports.Func) {
	for k := range funcs {
		if strings.HasSuffix(k, ".MarshalJSON") || strings.HasSuffix(k, ".UnmarshalJSON") {
			delete(funcs, k)
		}
	}
}

func EnumFilter(changelog *Changelog) {
	if changelog.Modified.HasAdditiveChanges() {
		if changelog.Modified.AdditiveChanges != nil && changelog.Modified.AdditiveChanges.Added.TypeAliases != nil {
			for typeAliases := range changelog.Modified.AdditiveChanges.Added.TypeAliases {
				funcKeys, funcExist := searchKey(changelog.Modified.AdditiveChanges.Added.Funcs, typeAliases, "Possible")
				if funcExist && len(funcKeys) == 1 {
					for _, f := range funcKeys {
						delete(changelog.Modified.AdditiveChanges.Added.Funcs, f)
					}
				}
			}
		}
	}

	if changelog.Modified.HasBreakingChanges() {
		enumOperation(changelog.Modified.BreakingChanges.Removed)
	}
}

func enumOperation(content *delta.Content) {
	if content != nil && content.TypeAliases != nil {
		for typeAliases := range content.TypeAliases {
			constKeys, constExist := searchKey(content.Consts, typeAliases, "")
			funcKeys, funcExist := searchKey(content.Funcs, typeAliases, "Possible")

			if constExist && funcExist && len(funcKeys) == 1 {
				for _, c := range constKeys {
					delete(content.Consts, c)
				}
				for _, f := range funcKeys {
					delete(content.Funcs, f)
				}
			}
		}
	}
}

func searchKey[T exports.Const | exports.Func | exports.Struct](m map[string]T, key1, prefix string) ([]string, bool) {
	keys := make([]string, 0)
	for k := range m {
		if regexp.MustCompile(fmt.Sprintf("^%s%s\\w*", prefix, key1)).MatchString(k) {
			keys = append(keys, k)
		}
	}
	if len(keys) != 0 {
		return keys, true
	}
	return nil, false
}

func FuncFilter(changelog *Changelog) {
	if changelog.Modified.HasAdditiveChanges() {
		funcOperation(changelog.Modified.AdditiveChanges.Added)
	}

	if changelog.Modified.HasBreakingChanges() {
		funcOperation(changelog.Modified.BreakingChanges.Removed)

		// function operation parameters from interface{} to any is not a breaking change
		for f, v := range changelog.Modified.BreakingChanges.Funcs {
			if v.Params == nil {
				continue
			}
			from := strings.Split(v.Params.From, ",")
			to := strings.Split(v.Params.To, ",")
			if len(from) != len(to) {
				continue
			}

			flag := false
			for i := range from {
				if strings.TrimSpace(from[i]) != strings.TrimSpace(to[i]) {
					if strings.TrimSpace(from[i]) == "interface{}" && strings.TrimSpace(to[i]) == "any" {
						flag = true
					} else {
						flag = false
						break
					}
				}
			}

			if flag {
				delete(changelog.Modified.BreakingChanges.Funcs, f)
			}
		}
	}
}

func funcOperation(content *delta.Content) {
	if content != nil && content.Funcs != nil {
		for funcName, funcValue := range content.Funcs {
			clientFunc := strings.Split(funcName, ".")
			if len(clientFunc) == 2 {
				// the last parameter
				if len(funcValue.Params) > 0 {
					lastParam := funcValue.Params[len(funcValue.Params)-1]
					clientFuncOptions := strings.TrimLeft(lastParam.Type, "*")
					if clientFuncOptions != "" && content.CompleteStructs != nil {
						delete(content.Structs, clientFuncOptions)
						for i, v := range content.CompleteStructs {
							if v == clientFuncOptions {
								content.CompleteStructs = append(content.CompleteStructs[:i],
									content.CompleteStructs[i+1:]...)
								break
							}
						}
					}
				}

				// the first return value
				if funcValue.Returns != nil {
					rs := strings.Split(*funcValue.Returns, ",")
					clientFuncResponse := rs[0]
					if strings.Contains(clientFuncResponse, "runtime") {
						re := regexp.MustCompile(`\[(?P<response>.*)\]`)
						clientFuncResponse = re.FindString(clientFuncResponse)
						clientFuncResponse = re.ReplaceAllString(clientFuncResponse, "${response}")
					} else {
						clientFuncResponse = strings.TrimLeft(clientFuncResponse, "*")
					}
					if clientFuncResponse != "" && content.CompleteStructs != nil {
						delete(content.Structs, clientFuncResponse)
						for i, v := range content.CompleteStructs {
							if v == clientFuncResponse {
								content.CompleteStructs = append(content.CompleteStructs[:i],
									content.CompleteStructs[i+1:]...)
								break
							}
						}
					}
				}
			}
		}
	}
}

// LROFilter LROFilter after OperationFilter
func LROFilter(changelog *Changelog) {
	if changelog.Modified.HasBreakingChanges() && changelog.Modified.HasAdditiveChanges() && changelog.Modified.BreakingChanges.Removed != nil && changelog.Modified.BreakingChanges.Removed.Funcs != nil {
		removedContent := changelog.Modified.BreakingChanges.Removed
		for bFunc, v := range removedContent.Funcs {
			var beginFunc string
			clientFunc := strings.Split(bFunc, ".")
			if len(clientFunc) == 2 {
				if strings.Contains(clientFunc[1], "Begin") {
					clientFunc[1] = strings.TrimPrefix(clientFunc[1], "Begin")
					beginFunc = fmt.Sprintf("%s.%s", clientFunc[0], clientFunc[1])
				} else {
					beginFunc = fmt.Sprintf("%s.Begin%s", clientFunc[0], clientFunc[1])
				}
				if _, ok := changelog.Modified.AdditiveChanges.Added.Funcs[beginFunc]; ok {
					delete(changelog.Modified.AdditiveChanges.Added.Funcs, beginFunc)
					v.ReplacedBy = &beginFunc
					removedContent.Funcs[bFunc] = v
				}
			}
		}
	}
}

// PageableFilter PageableFilter after OperationFilter
func PageableFilter(changelog *Changelog) {
	if changelog.Modified.HasBreakingChanges() && changelog.Modified.HasAdditiveChanges() && changelog.Modified.BreakingChanges.Removed != nil && changelog.Modified.BreakingChanges.Removed.Funcs != nil {
		removedContent := changelog.Modified.BreakingChanges.Removed
		for bFunc, v := range removedContent.Funcs {
			var pagination string
			clientFunc := strings.Split(bFunc, ".")
			if len(clientFunc) == 2 {
				if strings.Contains(clientFunc[1], "New") && strings.Contains(clientFunc[1], "Pager") {
					clientFunc[1] = strings.TrimPrefix(strings.TrimSuffix(clientFunc[1], "Pager"), "New")
					pagination = fmt.Sprintf("%s.%s", clientFunc[0], clientFunc[1])
				} else {
					pagination = fmt.Sprintf("%s.New%sPager", clientFunc[0], clientFunc[1])
				}
				if _, ok := changelog.Modified.AdditiveChanges.Added.Funcs[pagination]; ok {
					delete(changelog.Modified.AdditiveChanges.Added.Funcs, pagination)
					v.ReplacedBy = &pagination
					removedContent.Funcs[bFunc] = v
				}
			}
		}
	}
}

func InterfaceToAnyFilter(changelog *Changelog) {
	if changelog.HasBreakingChanges() {
		for structName, s := range changelog.Modified.BreakingChanges.Structs {
			for k, v := range s.Fields {
				if strings.Contains(v.From, "interface{}") && strings.Contains(v.To, "any") {
					delete(s.Fields, k)
				}
			}

			if len(s.Fields) == 0 {
				delete(changelog.Modified.BreakingChanges.Structs, structName)
			}
		}
	}
}

func NonExportedFilter(changelog *Changelog) {
	if !changelog.Modified.IsEmpty() {
		if changelog.Modified.HasAdditiveChanges() {
			nonExportOperation(changelog.Modified.AdditiveChanges.Added)
		}

		if changelog.Modified.HasBreakingChanges() {
			breakingChanges := changelog.Modified.BreakingChanges
			for fName := range breakingChanges.Funcs {
				before, after, _ := strings.Cut(fName, ".")
				if !ast.IsExported(strings.TrimLeft(before, "*")) || (after != "" && !ast.IsExported(after)) {
					delete(changelog.Modified.BreakingChanges.Funcs, fName)
				}
			}

			for sName := range breakingChanges.Structs {
				if !ast.IsExported(sName) {
					delete(changelog.Modified.BreakingChanges.Structs, sName)
				}
			}

			if breakingChanges.Removed != nil && !breakingChanges.Removed.IsEmpty() {
				nonExportOperation(breakingChanges.Removed)
			}

		}
	}
}

func nonExportOperation(content *delta.Content) {
	if content.IsEmpty() {
		return
	}

	for fName := range content.Funcs {
		before, after, _ := strings.Cut(fName, ".")
		if !ast.IsExported(strings.TrimLeft(before, "*")) || (after != "" && !ast.IsExported(after)) {
			delete(content.Funcs, fName)
		}
	}

	for sName := range content.Structs {
		if !ast.IsExported(sName) {
			delete(content.Structs, sName)
		}
	}
}

func TypeToAnyFilter(changelog *Changelog) {
	if changelog.Modified.HasBreakingChanges() {
		for structName, s := range changelog.Modified.BreakingChanges.Changes.Structs {
			for k, v := range s.Fields {
				if v.To == "any" {
					delete(s.Fields, k)
					if changelog.Modified.AdditiveChanges == nil {
						changelog.Modified.AdditiveChanges = &report.AdditiveChanges{}
					}
					if changelog.Modified.AdditiveChanges.Changes.Structs == nil {
						changelog.Modified.AdditiveChanges.Changes.Structs = map[string]delta.StructDef{}
					}
					if _, ok := changelog.Modified.AdditiveChanges.Changes.Structs[structName]; !ok {
						changelog.Modified.AdditiveChanges.Changes.Structs[structName] = delta.StructDef{
							Fields: make(map[string]delta.Signature),
						}
					}
					changelog.Modified.AdditiveChanges.Changes.Structs[structName].Fields[k] = v
				}
			}
			if len(s.Fields) == 0 {
				delete(changelog.Modified.BreakingChanges.Structs, structName)
			}
		}
	}
}
