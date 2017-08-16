package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/marstr/collection"
)

func ExampleSwaggerFinder_Enumerate() {
	subject, err := NewSwaggerFinder("testdata/azure-rest-api-specs", ioutil.Discard)
	if err != nil {
		return
	}

	var seen interface{}
	seen, err = collection.First(subject)
	if err != nil {
		return
	}
	fmt.Println(filepath.Base(seen.(SwaggerFile).Path))
	// Output: advisor.json
}
