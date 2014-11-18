package storage

import (
	"encoding/base64"
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

func Test_getStandardHeaders(t *testing.T) {
	cli, err := NewBasicClient("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	headers := cli.getStandardHeaders()
	if len(headers) != 2 {
		t.Error("Wrong standard header count")
	}
	if v, ok := headers["x-ms-version"]; !ok || v != cli.apiVersion {
		t.Error("Wrong version header")
	}
	if _, ok := headers["x-ms-date"]; !ok {
		t.Error("Missing date header")
	}
}

func Test_buildCanonicalizedResource(t *testing.T) {
	cli, err := NewBasicClient("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	type test struct{ url, expected string }
	tests := []test{
		{"https://foo.blob.core.windows.net/path?a=b&c=d", "/foo/path\na:b\nc:d"},
		{"https://foo.blob.core.windows.net/?comp=list", "/foo/\ncomp:list"},
		{"https://foo.blob.core.windows.net/cnt/blob", "/foo/cnt/blob"},
	}

	for _, i := range tests {
		if out, err := cli.buildCanonicalizedResource(i.url); err != nil {
			t.Error(err)
		} else if out != i.expected {
			t.Errorf("Wrong canonicalized resource. Expected:\n'%s', Got:\n'%s'", i.expected, out)
		}
	}
}

func Test_buildCanonicalizedHeader(t *testing.T) {
	cli, err := NewBasicClient("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	type test struct {
		headers  map[string]string
		expected string
	}
	tests := []test{
		{map[string]string{}, ""},
		{map[string]string{"x-ms-foo": "bar"}, "x-ms-foo:bar"},
		{map[string]string{"foo:": "bar"}, ""},
		{map[string]string{"foo:": "bar", "x-ms-foo": "bar"}, "x-ms-foo:bar"},
		{map[string]string{
			"x-ms-version":   "9999-99-99",
			"x-ms-blob-type": "BlockBlob"}, "x-ms-blob-type:BlockBlob\nx-ms-version:9999-99-99"}}

	for _, i := range tests {
		if out := cli.buildCanonicalizedHeader(i.headers); out != i.expected {
			t.Errorf("Wrong canonicalized resource. Expected:\n'%s', Got:\n'%s'", i.expected, out)
		}
	}
}

func Test_createAuthorizationHeader(t *testing.T) {
	key := base64.StdEncoding.EncodeToString([]byte("bar"))
	cli, err := NewBasicClient("foo", key)
	if err != nil {
		t.Error(err)
	}

	canonicalizedString := `foobarzoo`
	expected := `SharedKey foo:h5U0ATVX6SpbFX1H6GNuxIMeXXCILLoIvhflPtuQZ30=`

	if out, err := cli.createAuthorizationHeader(canonicalizedString); err != nil {
		t.Error(err)
	} else if out != expected {
		t.Errorf("Wrong authoriztion header. Expected: '%s', Got:'%s'", expected, out)
	}
}
