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
	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
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

	// Entry is the Atom wrapper for a management request
	Entry struct {
		XMLName                   xml.Name   `xml:"entry"`
		ID                        string     `xml:"id"`
		Title                     string     `xml:"title"`
		Published                 *date.Time `xml:"published,omitempty"`
		Updated                   *date.Time `xml:"updated,omitempty"`
		Author                    *Author    `xml:"author,omitempty"`
		Link                      *Link      `xml:"link,omitempty"`
		Content                   *Content   `xml:"content"`
		DataServiceSchema         string     `xml:"xmlns:d,attr"`
		DataServiceMetadataSchema string     `xml:"xmlns:m,attr"`
		AtomSchema                string     `xml:"xmlns,attr"`
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

	// BaseEntityDescription provides common fields which are part of Queues, Topics and Subscriptions
	BaseEntityDescription struct {
		InstanceMetadataSchema              string                   `xml:"xmlns:i,attr"`
		ServiceBusSchema                    string                   `xml:"xmlns,attr"`
		DefaultMessageTimeToLive            *string                  `xml:"DefaultMessageTimeToLive,omitempty"`            // DefaultMessageTimeToLive - ISO 8601 default message timespan to live value. This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the default value used when TimeToLive is not set on a message itself.
		MaxSizeInMegabytes                  *int32                   `xml:"MaxSizeInMegabytes,omitempty"`                  // MaxSizeInMegabytes - The maximum size of the queue in megabytes, which is the size of memory allocated for the queue. Default is 1024.
		RequiresDuplicateDetection          *bool                    `xml:"RequiresDuplicateDetection,omitempty"`          // RequiresDuplicateDetection - A value indicating if this queue requires duplicate detection.
		DuplicateDetectionHistoryTimeWindow *string                  `xml:"DuplicateDetectionHistoryTimeWindow,omitempty"` // DuplicateDetectionHistoryTimeWindow - ISO 8601 timeSpan structure that defines the duration of the duplicate detection history. The default value is 10 minutes.
		EnableBatchedOperations             *bool                    `xml:"EnableBatchedOperations,omitempty"`             // EnableBatchedOperations - Value that indicates whether server-side batched operations are enabled.
		SizeInBytes                         *int64                   `xml:"SizeInBytes,omitempty"`                         // SizeInBytes - The size of the queue, in bytes.
		IsAnonymousAccessible               *bool                    `xml:"IsAnonymousAccessible,omitempty"`
		Status                              *servicebus.EntityStatus `xml:"Status,omitempty"`
		CreatedAt                           *date.Time               `xml:"CreatedAt,omitempty"`
		UpdatedAt                           *date.Time               `xml:"UpdatedAt,omitempty"`
		SupportOrdering                     *bool                    `xml:"SupportOrdering,omitempty"`
		AutoDeleteOnIdle                    *string                  `xml:"AutoDeleteOnIdle,omitempty"`
		EnablePartitioning                  *bool                    `xml:"EnablePartitioning,omitempty"`
		EnableExpress                       *bool                    `xml:"EnableExpress,omitempty"`
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
