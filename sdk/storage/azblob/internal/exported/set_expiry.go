package exported

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"net/http"
	"strconv"
	"time"
)

// ExpiryType defines values for ExpiryType
type ExpiryType interface {
	Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.BlobClientSetExpiryOptions)
}

// ExpiryTypeAbsolute defines the absolute time for the blob expiry
type ExpiryTypeAbsolute time.Time

// ExpiryTypeRelativeToNow defines the duration relative to now for the blob expiry
type ExpiryTypeRelativeToNow time.Duration

// ExpiryTypeRelativeToCreation defines the duration relative to creation for the blob expiry
type ExpiryTypeRelativeToCreation time.Duration

// ExpiryTypeNever defines that the blob will be set to never expire
type ExpiryTypeNever struct {
	// empty struct since NeverExpire expiry type does not require expiry time
}

// ---------------------------------------------------------------------------------------------------------------------

// SetExpiryOptions contains the optional parameters for the Client.SetExpiry method.
type SetExpiryOptions struct {
	// placeholder for future options
}

func (e ExpiryTypeAbsolute) Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.BlobClientSetExpiryOptions) {
	return generated.ExpiryOptionsAbsolute, &generated.BlobClientSetExpiryOptions{
		ExpiresOn: to.Ptr(time.Time(e).UTC().Format(http.TimeFormat)),
	}
}

func (e ExpiryTypeRelativeToNow) Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.BlobClientSetExpiryOptions) {
	return generated.ExpiryOptionsRelativeToNow, &generated.BlobClientSetExpiryOptions{
		ExpiresOn: to.Ptr(strconv.FormatInt(time.Duration(e).Milliseconds(), 10)),
	}
}

func (e ExpiryTypeRelativeToCreation) Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.BlobClientSetExpiryOptions) {
	return generated.ExpiryOptionsRelativeToCreation, &generated.BlobClientSetExpiryOptions{
		ExpiresOn: to.Ptr(strconv.FormatInt(time.Duration(e).Milliseconds(), 10)),
	}
}

func (e ExpiryTypeNever) Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.BlobClientSetExpiryOptions) {
	return generated.ExpiryOptionsNeverExpire, nil
}
