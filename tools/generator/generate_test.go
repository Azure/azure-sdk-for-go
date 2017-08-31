package main

// Copyright 2017 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetNamespace(t *testing.T) {
	testCases := []struct {
		given    string
		expected string
	}{
		{
			"testdata/azure-rest-api-specs/specification/cdn/resource-manager/Microsoft.Cdn/2016-10-02/cdn.json",
			"services/cdn/management/2016-10-02/cdn",
		},
		{
			`testdata\azure-rest-api-specs\specification\keyvault\data-plane\Microsoft.KeyVault\2015-06-01\keyvault.json`,
			"services/keyvault/2015-06-01/keyvault",
		},
		{
			`testdata\azure-rest-api-specs\specification\keyvault\resource-manager\Microsoft.KeyVault\2015-06-01\keyvault.json`,
			"services/keyvault/management/2015-06-01/keyvault",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.given, func(t *testing.T) {
			var subject SwaggerFile
			contents, err := ioutil.ReadFile(tc.given)
			if err != nil {
				t.Error(err)
				return
			}

			err = json.Unmarshal(contents, &subject)
			if err != nil {
				t.Error(err)
				return
			}
			subject.Path = tc.given

			result, err := getNamespace(subject)
			if err != nil {
				t.Error(err)
			}

			if result != tc.expected {
				t.Logf("got:\n%s\nwant:\n%s", result, tc.expected)
				t.Fail()
			}
		})
	}
}

func TestMain(m *testing.M) {
	exitStatus := m.Run()
	if noClone == false {
		if err := os.RemoveAll(localAzureRESTAPISpecsPath); err != nil {
			fmt.Fprintln(os.Stderr, "Unable to delete folder: ", localAzureRESTAPISpecsPath)
		}
	}
	os.Exit(exitStatus)
}
