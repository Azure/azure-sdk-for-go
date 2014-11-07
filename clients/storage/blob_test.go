package storage

import (
	"testing"
)

func TestGetBaseUrl(t *testing.T) {
	cli, err := NewBasicClient("foo", "bar")
	if err != nil {
		t.Error(err)
	}
	blob := cli.GetBlobService()
	output := blob.getBaseUrl()

	if expected := "https://foo.blob.core.windows.net"; output != expected {
		t.Errorf("Wrong base url. Expected: '%s', got: '%s'", expected, output)
	}
}
