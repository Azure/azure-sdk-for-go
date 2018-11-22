// This program simply tests that we can grab a cert from the real keyvault
// and the library works.
package main

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault/beta/easykey"
)

var (
	vaultName   = flag.String("vault", "", "The name of the vault")
	certName    = flag.String("cert", "", "The name of the cert to get")
	certVersion = flag.String("version", easykey.LatestVersion, "The version of the cert, defaults to latest")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	kv, err := easykey.New(ctx, *vaultName, easykey.Config{})
	if err != nil {
		panic(err)
	}

	cer, err := kv.Certificate(ctx, *certName, *certVersion)
	if err != nil {
		panic(err)
	}
	_, err = cer.X509()
	if err != nil {
		panic(err)
	}

	tlsCert, err := kv.TLSCert(ctx, *certName, *certVersion, 1, 0)
	if err != nil {
		panic(err)
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{tlsCert}}
	srv := &http.Server{
		Addr:         ":8080",
		TLSConfig:    cfg,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	go func() {
		if err := srv.ListenAndServeTLS("", ""); err != nil {
			log.Println("server stopped listening, this may just be the server shutting down: ", err)
		}
	}()
	defer srv.Close()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	stop := time.Now().Add(10 * time.Second)
	success := false
	for time.Now().Before(stop) {
		_, err := http.Get("https://127.0.0.1:8080")
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		success = true
		break
	}
	if !success {
		panic("our tls server never came up")
	}
	fmt.Println("Success")
}
