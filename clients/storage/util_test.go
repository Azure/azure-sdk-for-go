package storage

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func Test_timeRfc1123Formatted(t *testing.T) {
	now := time.Now().UTC()

	expectedLayout := "Mon, 02 Jan 2006 15:04:05 GMT"
	expected := now.Format(expectedLayout)

	if output := timeRfc1123Formatted(now); output != expected {
		t.Errorf("Expected: %s, got: %s", expected, output)
	}
}

func Test_mergeParams(t *testing.T) {
	v1 := url.Values{
		"k1": {"v1"},
		"k2": {"v2"}}
	v2 := url.Values{
		"k1": {"v11"},
		"k3": {"v3"}}

	out := mergeParams(v1, v2)
	if v := out.Get("k1"); v != "v1" {
		t.Errorf("Wrong value for k1: %s", v)
	}

	if v := out.Get("k2"); v != "v2" {
		t.Errorf("Wrong value for k2: %s", v)
	}

	if v := out.Get("k3"); v != "v3" {
		t.Errorf("Wrong value for k3: %s", v)
	}

	if v := out["k1"]; !reflect.DeepEqual(v, []string{"v1", "v11"}) {
		t.Errorf("Wrong multi-value for k1: %s", v)
	}
}
