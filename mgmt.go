package servicebus

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-amqp-common-go/auth"
	"github.com/Azure/go-autorest/autorest/date"
)

const (
	instanceMetadataSchema    = "http://www.w3.org/2001/XMLSchema-instance"
	serviceBusSchema          = "http://schemas.microsoft.com/netservices/2010/10/servicebus/connect"
	dataServiceSchema         = "http://schemas.microsoft.com/ado/2007/08/dataservices"
	dataServiceMetadataSchema = "http://schemas.microsoft.com/ado/2007/08/dataservices/metadata"
	atomSchema                = "http://www.w3.org/2005/Atom"
	applicationXML            = "application/xml"
)

type (
	// EntityManager provides CRUD functionality for Service Bus entities (Queues, Topics, Subscriptions...)
	EntityManager struct {
		TokenProvider auth.TokenProvider
		Host          string
	}

	// Feed is an Atom feed which contains entries
	Feed struct {
		XMLName xml.Name   `xml:"feed"`
		ID      string     `xml:"id"`
		Title   string     `xml:"title"`
		Updated *date.Time `xml:"updated,omitempty"`
		Entries []Entry    `xml:"entry"`
	}

	// Author is an Atom author used in an Entry
	Author struct {
		XMLName xml.Name `xml:"author"`
		Name    *string  `xml:"name,omitempty"`
	}

	// Link is an Atom link used in an Entry
	Link struct {
		XMLName xml.Name `xml:"link"`
		Rel     string   `xml:"rel,attr"`
		HREF    string   `xml:"href,attr"`
	}

	// Content is a generic body for an Atom Entry
	Content struct {
		XMLName xml.Name `xml:"content"`
		Type    string   `xml:"type,attr"`
		Body    string   `xml:",innerxml"`
	}
)

// NewEntityManager creates a new instance of an EntityManager given a token provider and host
func NewEntityManager(host string, tokenProvider auth.TokenProvider) *EntityManager {
	return &EntityManager{
		Host:          host,
		TokenProvider: tokenProvider,
	}
}

// Get performs an HTTP Get for a given entity path
func (em *EntityManager) Get(ctx context.Context, entityPath string) (*http.Response, error) {
	return em.Execute(ctx, http.MethodGet, entityPath, http.NoBody)
}

// Put performs an HTTP PUT for a given entity path and body
func (em *EntityManager) Put(ctx context.Context, entityPath string, body []byte) (*http.Response, error) {
	return em.Execute(ctx, http.MethodPut, entityPath, bytes.NewReader(body))
}

// Delete performs an HTTP DELETE for a given entity path
func (em *EntityManager) Delete(ctx context.Context, entityPath string) (*http.Response, error) {
	return em.Execute(ctx, http.MethodDelete, entityPath, http.NoBody)
}

// Post performs an HTTP POST for a given entity path and body
func (em *EntityManager) Post(ctx context.Context, entityPath string, body []byte) (*http.Response, error) {
	return em.Execute(ctx, http.MethodPost, entityPath, bytes.NewReader(body))
}

// Execute performs an HTTP request given a http method, path and body
func (em *EntityManager) Execute(ctx context.Context, method string, entityPath string, body io.Reader) (*http.Response, error) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	req, err := http.NewRequest(method, em.Host+strings.TrimPrefix(entityPath, "/"), body)
	if err != nil {
		return nil, err
	}
	req = addAtomXMLContentType(req)
	req = addAPIVersion201704(req)
	req, err = em.addAuthorization(req)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return client.Do(req)
}

func (em *EntityManager) addAuthorization(req *http.Request) (*http.Request, error) {
	signature, err := em.TokenProvider.GetToken(req.URL.String())
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", signature.Token)
	return req, nil
}

func addAtomXMLContentType(req *http.Request) *http.Request {
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		req.Header.Add("Content-Type", "application/atom+xml;type=entry;charset=utf-8")
	}
	return req
}

func addAPIVersion201704(req *http.Request) *http.Request {
	q := req.URL.Query()
	q.Add("api-version", "2017-04")
	req.URL.RawQuery = q.Encode()
	return req
}

func xmlDoc(content []byte) []byte {
	return []byte(xml.Header + string(content))
}

// ptrBool takes a boolean and returns a pointer to that bool. For use in literal pointers, ptrBool(true) -> *bool
func ptrBool(toPtr bool) *bool {
	return &toPtr
}

// ptrString takes a string and returns a pointer to that string. For use in literal pointers,
// ptrString(fmt.Sprintf("..", foo)) -> *string
func ptrString(toPtr string) *string {
	return &toPtr
}

// durationTo8601Seconds takes a duration and returns a string period of whole seconds (int cast of float)
func durationTo8601Seconds(duration *time.Duration) *string {
	return ptrString(fmt.Sprintf("PT%dS", int(duration.Seconds())))
}
