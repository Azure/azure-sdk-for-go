package azidentity

import (
	"net/url"
	"testing"
)

func TestAuthorityHostParse(t *testing.T) {
	_, err := url.Parse(defaultAuthorityHost)
	if err != nil {
		t.Fatalf("Failed to parse default authority host: %v", err)
	}
}
