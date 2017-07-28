package main

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
