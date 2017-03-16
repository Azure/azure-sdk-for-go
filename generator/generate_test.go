package main

import (
	"path/filepath"
	"testing"
)

func TestGetNamespace_Standard(t *testing.T) {
	basePath := filepath.Join("dir1", "dir2")
	swaggerPath := filepath.Join(basePath, "plane-package", "2016-02-17", "swagger", "example.json")
	result, err := getNamespace(swaggerPath, basePath)
	if err != nil {
		t.Error(err)
	}
	expected := "github.com/Azure/azure-sdk-for-go/plane/package"
	if result != expected {
		t.Logf("got:\n%s\nwant:\n%s", result, expected)
		t.Fail()
	}
}

func TestGetNamespace_Nested(t *testing.T) {
	basePath := filepath.Join("dir1", "dir2")
	swaggerPath := filepath.Join(basePath, "plane-split-name", "dir3", "dir4", "2016-02-09-preview", "swagger", "example.json")
	result, err := getNamespace(swaggerPath, basePath)
	if err != nil {
		t.Error(err)
	}
	expected := "github.com/Azure/azure-sdk-for-go/plane/split-name/dir3/dir4"
	if result != expected {
		t.Logf("got:\n%s\nwant:\n%s", result, expected)
		t.Fail()
	}
}
