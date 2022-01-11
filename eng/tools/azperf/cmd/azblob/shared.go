package cmd

import "io"

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

var size int
var count int

func init() {
	DownloadBlobCmd.Flags().IntVarP(&size, "size", "s", 10240, "Size in bytes of data to be transferred in download tests. Default is 10240.")

	UploadBlobCmd.Flags().IntVarP(&size, "size", "s", 10240, "Size in bytes of data to be transferred in upload tests. Default is 10240.")

	ListBlobCmd.Flags().IntVarP(&count, "count", "c", 100, "Size in bytes of data to be transferred in upload tests. Default is 10240.")
}
