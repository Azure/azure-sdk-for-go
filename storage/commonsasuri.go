package storage

import (
	"net/url"
	"time"
)

// SASOptions includes options used by SAS URIs for different
// services and resources.
type SASOptions struct {
	APIVersion string
	Start      time.Time
	Expiry     time.Time
	IP         string
	UseHTTPS   bool
	Identifier string
}

func addQueryParameter(query url.Values, key, value string) url.Values {
	if value != "" {
		query.Add(key, value)
	}
	return query
}
