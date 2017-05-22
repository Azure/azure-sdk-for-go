package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetNamespace(t *testing.T) {
	testCases := []struct {
		given    string
		expected string
	}{
		{
			filepath.Join("plane-package", "2016-02-17", "swagger", "example.json"),
			"plane/package/2016-02-17/example",
		},
		{
			filepath.Join("plane-split-name", "dir3", "2016-02-09-preview", "swagger", "example.json"),
			"plane/split-name/dir3/2016-02-09-preview/example",
		},
		{
			filepath.Join("plane-split-name", "dir3", "dir4", "2016-02-17", "swagger", "example.json"),
			"plane/split-name/dir3/dir4/2016-02-17/example",
		},
		{
			filepath.Join("plane-split-name", "dir3", "2016-02-09-preview", "swagger", "example.json"),
			"plane/split-name/dir3/2016-02-09-preview/example",
		},
		{
			filepath.Join("myService/2015-06-01/swagger/example.json"),
			"services/myService/2015-06-01/example",
		},
		{
			filepath.Join("plane-split-name", "v1.8", "swagger", "example.json"),
			"plane/split-name/v1.8/example",
		},
		{
			filepath.Join("plane-split-name", "1.8", "swagger", "example.json"),
			"plane/split-name/1.8/example",
		},
		{
			filepath.Join("plane-split-name", "subname", "v1.8", "swagger", "example.json"),
			"plane/split-name/subname/v1.8/example",
		},
		{
			filepath.Join("plane-name", "subname", "2015-12-29", "swagger", "Example.json"),
			"plane/name/subname/2015-12-29/example",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.given, func(t *testing.T) {
			result, err := getNamespace(Swagger{
				Path: tc.given,
			})
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
		if err := os.RemoveAll(localAzureRestAPISpecsPath); err != nil {
			fmt.Fprintln(os.Stderr, "Unable to delete folder: ", localAzureRestAPISpecsPath)
		}
	}
	os.Exit(exitStatus)
}
