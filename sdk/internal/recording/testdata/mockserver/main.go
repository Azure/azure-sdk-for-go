//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// This is a mock server used to test the sanitizers in sanitizer.go

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Success", time.Now().String())
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "Next-Location")
	err := json.NewEncoder(w).Encode(map[string]string{
		"Tag":  "Value",
		"Tag2": "Value2",
		"Tag3": "https://storageaccount.table.core.windows.net/",
	})
	if err != nil {
		log.Fatalf("error writing the response")
	}
}

func main() {
	http.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
