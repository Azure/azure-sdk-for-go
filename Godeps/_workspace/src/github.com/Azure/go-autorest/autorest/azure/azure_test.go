package azure

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	. "github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/mocks"
)

// Use a Client Inspector to set the request identifier.
func ExampleWithClientID() {
	uuid := "71FDB9F4-5E49-4C12-B266-DE7B4FD999A6"
	req, _ := Prepare(&http.Request{},
		AsGet(),
		WithBaseURL("https://microsoft.com/a/b/c/"))

	c := Client{Sender: mocks.NewSender()}
	c.RequestInspector = WithReturningClientID(uuid)

	c.Send(req)
	fmt.Printf("Inspector added the %s header with the value %s\n",
		HeaderClientID, req.Header.Get(HeaderClientID))
	fmt.Printf("Inspector added the %s header with the value %s\n",
		HeaderReturnClientID, req.Header.Get(HeaderReturnClientID))
	// Output:
	// Inspector added the x-ms-client-request-id header with the value 71FDB9F4-5E49-4C12-B266-DE7B4FD999A6
	// Inspector added the x-ms-return-client-request-id header with the value true
}

func TestWithReturningClientIDReturnsError(t *testing.T) {
	var errIn error
	uuid := "71FDB9F4-5E49-4C12-B266-DE7B4FD999A6"
	_, errOut := Prepare(&http.Request{},
		withErrorPrepareDecorator(&errIn),
		WithReturningClientID(uuid))

	if errOut == nil || errIn != errOut {
		t.Errorf("azure: WithReturningClientID failed to exit early when receiving an error -- expected (%v), received (%v)",
			errIn, errOut)
	}
}

func TestWithClientID(t *testing.T) {
	uuid := "71FDB9F4-5E49-4C12-B266-DE7B4FD999A6"
	req, _ := Prepare(&http.Request{},
		WithClientID(uuid))

	if req.Header.Get(HeaderClientID) != uuid {
		t.Errorf("azure: WithClientID failed to set %s -- expected %s, received %s",
			HeaderClientID, uuid, req.Header.Get(HeaderClientID))
	}
}

func TestWithReturnClientID(t *testing.T) {
	b := false
	req, _ := Prepare(&http.Request{},
		WithReturnClientID(b))

	if req.Header.Get(HeaderReturnClientID) != strconv.FormatBool(b) {
		t.Errorf("azure: WithReturnClientID failed to set %s -- expected %s, received %s",
			HeaderClientID, strconv.FormatBool(b), req.Header.Get(HeaderClientID))
	}
}

func TestExtractClientID(t *testing.T) {
	uuid := "71FDB9F4-5E49-4C12-B266-DE7B4FD999A6"
	resp := mocks.NewResponse()
	mocks.SetResponseHeader(resp, HeaderClientID, uuid)

	if ExtractClientID(resp) != uuid {
		t.Errorf("azure: ExtractClientID failed to extract the %s -- expected %s, received %s",
			HeaderClientID, uuid, ExtractClientID(resp))
	}
}

func TestExtractRequestID(t *testing.T) {
	uuid := "71FDB9F4-5E49-4C12-B266-DE7B4FD999A6"
	resp := mocks.NewResponse()
	mocks.SetResponseHeader(resp, HeaderRequestID, uuid)

	if ExtractRequestID(resp) != uuid {
		t.Errorf("azure: ExtractRequestID failed to extract the %s -- expected %s, received %s",
			HeaderRequestID, uuid, ExtractRequestID(resp))
	}
}

func withErrorPrepareDecorator(e *error) PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			*e = fmt.Errorf("autorest: Faux Prepare Error")
			return r, *e
		})
	}
}
