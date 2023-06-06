//go:build go1.18
// +build go1.18

//go:generate autorest ./autorest.md
//go:generate gofmt -w .

// (not used for now)
//nono:generate autorest ./autorest.md --rawjson-as-bytes

package azeventgrid
