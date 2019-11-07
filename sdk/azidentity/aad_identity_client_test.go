package azidentity

import (
	"net/url"
	"testing"
)

func TestCreateClientAssertionJWT(t *testing.T) {
	_, err := url.Parse(defaultAuthorityHost)
	if err != nil {
		t.Fatalf("Failed to parse default authority host: %v", err)
	}
}
