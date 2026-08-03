// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exports

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddConst_SelectorExprValue(t *testing.T) {
	dir := t.TempDir()
	src := `package foo

type BlobType string

const (
	BlobTypeBlockBlob BlobType = generated.BlobTypeBlockBlob
)
`
	if err := os.WriteFile(filepath.Join(dir, "foo.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	c, err := Get(dir)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := c.Consts["BlobTypeBlockBlob"]
	if !ok {
		t.Fatalf("expected const BlobTypeBlockBlob, got %#v", c.Consts)
	}
	if got.Type != "BlobType" {
		t.Errorf("expected type BlobType, got %q", got.Type)
	}
	if got.Value != "generated.BlobTypeBlockBlob" {
		t.Errorf("expected value generated.BlobTypeBlockBlob, got %q", got.Value)
	}
}

func TestAddConst_UntypedSelectorExprValue(t *testing.T) {
	dir := t.TempDir()
	src := `package foo

const (
	EventUpload = exported.EventUpload
)
`
	if err := os.WriteFile(filepath.Join(dir, "foo.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	c, err := Get(dir)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := c.Consts["EventUpload"]
	if !ok {
		t.Fatalf("expected const EventUpload, got %#v", c.Consts)
	}
	if got.Type != "*ast.SelectorExpr" {
		t.Errorf("expected sentinel type *ast.SelectorExpr, got %q", got.Type)
	}
	if got.Value != "exported.EventUpload" {
		t.Errorf("expected value exported.EventUpload, got %q", got.Value)
	}
}
