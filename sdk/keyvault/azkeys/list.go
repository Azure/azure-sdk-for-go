package azkeys

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal"
)

// ListKeyVersionsPager is a Pager for Client.ListKeyVersions results
type ListKeyVersionsPager interface {
	// PageResponse returns the current ListKeyVersionsPage
	PageResponse() ListKeyVersionsPage

	// Err returns true if there is another page of data available, false if not
	Err() error

	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
}

type listKeyVersionsPager struct {
	genPager *internal.KeyVaultClientGetKeyVersionsPager
}

// PageResponse returns the results from the page most recently fetched from the service.
func (l *listKeyVersionsPager) PageResponse() ListKeyVersionsPage {
	return listKeyVersionsPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (l *listKeyVersionsPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *listKeyVersionsPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListKeyVersionsOptions contains the options for the ListKeyVersions operations
type ListKeyVersionsOptions struct {
	// Maximum number of results to return in a page. If not specified the service will return up to 25 results.
	MaxResults *int32
}

// convert the public ListKeyVersionsOptions to the generated version
func (l *ListKeyVersionsOptions) toGenerated() *internal.KeyVaultClientGetKeyVersionsOptions {
	if l == nil {
		return &internal.KeyVaultClientGetKeyVersionsOptions{}
	}
	return &internal.KeyVaultClientGetKeyVersionsOptions{
		Maxresults: l.MaxResults,
	}
}

// The key list result
type ListKeyVersionsPage struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// READ-ONLY; The URL to get the next set of keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of keys in the key vault along with a link to the next page of keys.
	Keys []KeyItem `json:"value,omitempty" azure:"ro"`
}

// create ListKeysPage from generated pager
func listKeyVersionsPageFromGenerated(i internal.KeyVaultClientGetKeyVersionsResponse) ListKeyVersionsPage {
	var keys []KeyItem
	for _, s := range i.Value {
		if s != nil {
			keys = append(keys, *keyItemFromGenerated(s))
		}
	}
	return ListKeyVersionsPage{
		RawResponse: i.RawResponse,
		NextLink:    i.NextLink,
		Keys:        keys,
	}
}

// ListKeyVersions lists all versions of the specified key. The full key identifer and
// attributes are provided in the response. No values are returned for the keys. This operation
// requires the keys/list permission.
func (c *Client) ListKeyVersions(keyName string, options *ListKeyVersionsOptions) ListKeyVersionsPager {
	if options == nil {
		options = &ListKeyVersionsOptions{}
	}

	return &listKeyVersionsPager{
		genPager: c.kvClient.GetKeyVersions(
			c.vaultUrl,
			keyName,
			options.toGenerated(),
		),
	}
}
