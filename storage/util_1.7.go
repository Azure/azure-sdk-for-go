// +build !go1.8

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
