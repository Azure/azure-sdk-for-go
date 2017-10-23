// +build !go1.8

// Copyright 2017 Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package storage

import (
	"io"
	"net/http"
	"os"
)

func setContentLengthFromLimitedReader(req *http.Request, lr *io.LimitedReader) error {
	req.ContentLength = lr.N
	return nil
}

func setContentLengthFromFile(req *http.Request, f *os.File) error {
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	req.ContentLength = fi.Size()
	return nil
}
