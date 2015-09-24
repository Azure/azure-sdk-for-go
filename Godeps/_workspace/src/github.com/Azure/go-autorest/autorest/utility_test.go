package autorest

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/mocks"
)

const (
	testAuthorizationHeader = "BEARER SECRETTOKEN"
	testBadURL              = ""
	jsonT                   = `
    {
      "name":"Rob Pike",
      "age":42
    }`
)

func TestContainsIntFindsValue(t *testing.T) {
	ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	v := 5
	if !containsInt(ints, v) {
		t.Errorf("autorest: containsInt failed to find %v in %v", v, ints)
	}
}

func TestContainsIntDoesNotFindValue(t *testing.T) {
	ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	v := 42
	if containsInt(ints, v) {
		t.Errorf("autorest: containsInt unexpectedly found %v in %v", v, ints)
	}
}

func TestEscapeStrings(t *testing.T) {
	m := map[string]string{
		"string": "a long string with = odd characters",
		"int":    "42",
		"nil":    "",
	}
	r := map[string]string{
		"string": "a+long+string+with+%3D+odd+characters",
		"int":    "42",
		"nil":    "",
	}
	v := escapeValueStrings(m)
	if !reflect.DeepEqual(v, r) {
		t.Errorf("autorest: ensureValueStrings returned %v\n", v)
	}
}

func TestEnsureStrings(t *testing.T) {
	m := map[string]interface{}{
		"string": "string",
		"int":    42,
		"nil":    nil,
	}
	r := map[string]string{
		"string": "string",
		"int":    "42",
		"nil":    "",
	}
	v := ensureValueStrings(m)
	if !reflect.DeepEqual(v, r) {
		t.Errorf("autorest: ensureValueStrings returned %v\n", v)
	}
}

func doEnsureBodyClosed(t *testing.T) SendDecorator {
	return func(s Sender) Sender {
		return SenderFunc(func(r *http.Request) (*http.Response, error) {
			resp, err := s.Do(r)
			if resp != nil && resp.Body != nil && resp.Body.(*mocks.Body).IsOpen() {
				t.Error("autorest: Expected Body to be closed -- it was left open")
			}
			return resp, err
		})
	}
}

type mockAuthorizer struct{}

func (ma mockAuthorizer) WithAuthorization() PrepareDecorator {
	return WithHeader(headerAuthorization, testAuthorizationHeader)
}

type mockFailingAuthorizer struct{}

func (mfa mockFailingAuthorizer) WithAuthorization() PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			return r, fmt.Errorf("ERROR: mockFailingAuthorizer returned expected error")
		})
	}
}

func withMessage(output *string, msg string) SendDecorator {
	return func(s Sender) Sender {
		return SenderFunc(func(r *http.Request) (*http.Response, error) {
			resp, err := s.Do(r)
			if err == nil {
				*output += msg
			}
			return resp, err
		})
	}
}

type mockInspector struct {
	wasInvoked bool
}

func (mi *mockInspector) WithInspection() PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			mi.wasInvoked = true
			return p.Prepare(r)
		})
	}
}

func (mi *mockInspector) ByInspecting() RespondDecorator {
	return func(r Responder) Responder {
		return ResponderFunc(func(resp *http.Response) error {
			mi.wasInvoked = true
			return r.Respond(resp)
		})
	}
}

func withErrorRespondDecorator(e *error) RespondDecorator {
	return func(r Responder) Responder {
		return ResponderFunc(func(resp *http.Response) error {
			err := r.Respond(resp)
			if err != nil {
				return err
			}
			*e = fmt.Errorf("autorest: Faux Respond Error")
			return *e
		})
	}
}
