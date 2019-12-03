package azidentity

import (
	"net/url"
	"testing"
)

func Test_IMDSEndpoint_Parse(t *testing.T) {
	_, err := url.Parse(imdsEndpoint)
	if err != nil {
		t.Fatalf("Failed to parse the IMDS endpoint: %v", err)
	}
}
