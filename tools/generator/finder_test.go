package main

// Copyright 2017 Microsoft Corporation
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

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
