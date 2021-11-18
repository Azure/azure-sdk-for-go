package sas

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	sas = "SharedAccessSignature"
)

type (
	sig struct {
		sr  string
		se  string
		skn string
		sig string
	}
)

func TestNewSigner(t *testing.T) {
	keyName, key := "foo", "superSecret"
	signer := NewSigner(keyName, key)
	before := time.Now().UTC().Add(-2 * time.Second)
	sigStr, expiry := signer.SignWithDuration("http://microsoft.com", 1*time.Hour)
	nixExpiry, err := strconv.ParseInt(expiry, 10, 64)
	assert.WithinDuration(t, before.Add(1*time.Hour), time.Unix(nixExpiry, 0), 10*time.Second, "signing expiry")

	sig, err := parseSig(sigStr)
	assert.Nil(t, err)
	assert.Equal(t, "http%3a%2f%2fmicrosoft.com", sig.sr)
	assert.Equal(t, keyName, sig.skn)
	assert.Equal(t, expiry, sig.se)
	assert.NotNil(t, sig.sig)
}

func parseSig(sigStr string) (*sig, error) {
	if !strings.HasPrefix(sigStr, sas+" ") {
		return nil, errors.New("should start with " + sas)
	}
	sigStr = strings.TrimPrefix(sigStr, sas+" ")
	parts := strings.Split(sigStr, "&")
	parsed := new(sig)
	for _, part := range parts {
		keyValue := strings.Split(part, "=")
		if len(keyValue) != 2 {
			return nil, errors.New("key value is malformed")
		}
		switch keyValue[0] {
		case "sr":
			parsed.sr = keyValue[1]
		case "se":
			parsed.se = keyValue[1]
		case "sig":
			parsed.sig = keyValue[1]
		case "skn":
			parsed.skn = keyValue[1]
		default:
			return nil, fmt.Errorf(fmt.Sprintf("unknown key / value: %q", keyValue))
		}
	}
	return parsed, nil
}
