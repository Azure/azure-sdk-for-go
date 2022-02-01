package main

import (
	"io"

	"github.com/spf13/pflag"
)

type nopCloser struct {
	io.ReadSeeker
}

func (n nopCloser) Close() error {
	return nil
}

// NopCloser returns a ReadSeekCloser with a no-op close method wrapping the provided io.ReadSeeker.
func NopCloser(rs io.ReadSeeker) io.ReadSeekCloser {
	return nopCloser{rs}
}

var size *int64
var count *int32

func RegisterArguments() {
	count = pflag.Int32("num-blobs", 100, "Number of blobs to list. Defaults to 100.")
	size = pflag.Int64("size", 10240, "Size in bytes of data to be transferred in upload or download tests. Default is 10240.")
}
