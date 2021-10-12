package azkeys

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal"
)

// ListDeletedSecrets is the interface for the Client.ListDeletedKeys operation
type ListDeletedKeysPager interface {
	// PageResponse returns the current ListDeletedSecretsPage
	PageResponse() ListDeletedKeysPage

	// Err returns true if there is another page of data available, false if not
	Err() error

	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
}

// listDeletedSecretsPager is the pager returned by Client.ListDeletedSecrets
type listDeletedSecretsPager struct {
	genPager *internal.KeyVaultClientGetDeletedKeysPager
}

// PageResponse returns the current page of results
func (l *listDeletedSecretsPager) PageResponse() ListDeletedKeysPage {
	resp := l.genPager.PageResponse()

	var values []*internal.DeletedKeyItem
	for _, d := range resp.Value {
		values = append(values, d)
	}

	return ListDeletedKeysPage{
		RawResponse: resp.RawResponse,
		NextLink:    resp.NextLink,
		DeletedKeys: values,
	}
}

// Err returns an error if the last operation resulted in an error.
func (l *listDeletedSecretsPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next page of results.
func (l *listDeletedSecretsPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListDeletedKeysPage holds the data for a single page.
type ListDeletedKeysPage struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// READ-ONLY; The URL to get the next set of deleted keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of the deleted keys in the vault along with a link to the next page of deleted keys
	DeletedKeys []*internal.DeletedKeyItem `json:"value,omitempty" azure:"ro"`
}

// ListDeletedKeysOptions contains the optional parameters for the Client.ListDeletedSecrets operation.
type ListDeletedKeysOptions struct {
	// Maximum number of results to return in a page. If not specified the service will return up to 25 results.
	MaxResults *int32
}

// Convert publicly exposed options to the generated version.a
func (l *ListDeletedKeysOptions) toGenerated() *internal.KeyVaultClientGetDeletedKeysOptions {
	return &internal.KeyVaultClientGetDeletedKeysOptions{
		Maxresults: l.MaxResults,
	}
}

// ListDeletedSecrets lists all versions of the specified key. The full key identifier and attributes are provided
// in the response. No values are returned for the keys. This operation requires the keys/list permission.
func (c *Client) ListDeletedSecrets(options *ListDeletedKeysOptions) ListDeletedKeysPager {
	if options == nil {
		options = &ListDeletedKeysOptions{}
	}

	return &listDeletedSecretsPager{
		genPager: c.kvClient.GetDeletedKeys(c.vaultUrl, options.toGenerated()),
	}

}
