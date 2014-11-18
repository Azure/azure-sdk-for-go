package storage

import (
	"net/url"
	"testing"
)

func TestGetBaseUrl_Basic_Https(t *testing.T) {
	cli, err := NewBasicClient("foo", "bar")

	if cli.apiVersion != DefaultApiVersion {
		t.Errorf("Wrong api version. Expected: '%s', got: '%s'", DefaultApiVersion, cli.apiVersion)
	}

	if err != nil {
		t.Error(err)
	}
	output := cli.getBaseUrl("table")

	if expected := "https://foo.table.core.windows.net"; output != expected {
		t.Errorf("Wrong base url. Expected: '%s', got: '%s'", expected, output)
	}
}

func TestGetBaseUrl_Custom_NoHttps(t *testing.T) {
	apiVersion := DefaultApiVersion
	cli, err := NewClient("foo", "bar", "core.chinacloudapi.cn", apiVersion, false)
	if err != nil {
		t.Error(err)
	}

	if cli.apiVersion != apiVersion {
		t.Errorf("Wrong api version. Expected: '%s', got: '%s'", apiVersion, cli.apiVersion)
	}

	output := cli.getBaseUrl("table")

	if expected := "http://foo.table.core.chinacloudapi.cn"; output != expected {
		t.Errorf("Wrong base url. Expected: '%s', got: '%s'", expected, output)
	}
}

func TestGetEndpoint_None(t *testing.T) {
	cli, err := NewBasicClient("foo", "bar")
	if err != nil {
		t.Error(err)
	}
	output := cli.getEndpoint(blobServiceName, "", url.Values{})

	if expected := "https://foo.blob.core.windows.net/"; output != expected {
		t.Errorf("Wrong endpoint url. Expected: '%s', got: '%s'", expected, output)
	}
}

func TestGetEndpoint_PathOnly(t *testing.T) {
	cli, err := NewBasicClient("foo", "bar")
	if err != nil {
		t.Error(err)
	}
	output := cli.getEndpoint(blobServiceName, "path", url.Values{})

	if expected := "https://foo.blob.core.windows.net/path"; output != expected {
		t.Errorf("Wrong endpoint url. Expected: '%s', got: '%s'", expected, output)
	}
}

func TestGetEndpoint_ParamsOnly(t *testing.T) {
	cli, err := NewBasicClient("foo", "bar")
	if err != nil {
		t.Error(err)
	}
	params := url.Values{}
	params.Set("a", "b")
	params.Set("c", "d")
	output := cli.getEndpoint(blobServiceName, "", params)

	if expected := "https://foo.blob.core.windows.net/?a=b&c=d"; output != expected {
		t.Errorf("Wrong endpoint url. Expected: '%s', got: '%s'", expected, output)
	}
}

func TestGetEndpoint_Mixed(t *testing.T) {
	cli, err := NewBasicClient("foo", "bar")
	if err != nil {
		t.Error(err)
	}
	params := url.Values{}
	params.Set("a", "b")
	params.Set("c", "d")
	output := cli.getEndpoint(blobServiceName, "path", params)

	if expected := "https://foo.blob.core.windows.net/path?a=b&c=d"; output != expected {
		t.Errorf("Wrong endpoint url. Expected: '%s', got: '%s'", expected, output)
	}
}
