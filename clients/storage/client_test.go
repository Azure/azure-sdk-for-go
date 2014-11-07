package storage

import (
	"testing"
)

func TestGetBaseUrl_Basic_Https(t *testing.T) {
	cli, err := NewBasicClient("foo", "bar")
	if err != nil {
		t.Error(err)
	}
	output := cli.getBaseUrl("table")

	if expected := "https://foo.table.core.windows.net"; output != expected {
		t.Errorf("Wrong base url. Expected: '%s', got: '%s'", expected, output)
	}
}

func TestGetBaseUrl_Custom_NoHttps(t *testing.T) {
	cli, err := NewClient("foo", "bar", "core.chinacloudapi.cn", false)
	if err != nil {
		t.Error(err)
	}

	output := cli.getBaseUrl("table")

	if expected := "http://foo.table.core.chinacloudapi.cn"; output != expected {
		t.Errorf("Wrong base url. Expected: '%s', got: '%s'", expected, output)
	}
}
