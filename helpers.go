package servicebus

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type (
	// resourceID represents a parsed long-form Azure Resource Manager ID
	// with the Subscription ID, Resource Group and the Provider as top-
	// level fields, and other key-value pairs available via a map in the
	// Path field.
	resourceID struct {
		SubscriptionID string
		ResourceGroup  string
		Provider       string
		Path           map[string]string
	}
)

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

// parseAzureResourceID converts a long-form Azure Resource Manager ID
// into a ResourceID. We make assumptions about the structure of URLs,
// which is obviously not good, but the best thing available given the
// SDK.
func parseAzureResourceID(id string) (*resourceID, error) {
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("cannot parse Azure Id: %s", err)
	}

	path := idURL.Path

	path = strings.TrimSpace(path)
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	components := strings.Split(path, "/")

	// We should have an even number of key-value pairs.
	if len(components)%2 != 0 {
		return nil, fmt.Errorf("the number of path segments is not divisible by 2 in %q", path)
	}

	var subscriptionID string

	// Put the constituent key-value pairs into a map
	componentMap := make(map[string]string, len(components)/2)
	for current := 0; current < len(components); current += 2 {
		key := components[current]
		value := components[current+1]

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return nil, fmt.Errorf("key/value cannot be empty strings. Key: '%s', Value: '%s'", key, value)
		}

		// Catch the subscriptionID before it can be overwritten by another "subscriptions"
		// value in the ID which is the case for the Service Bus subscription resource
		if key == "subscriptions" && subscriptionID == "" {
			subscriptionID = value
		} else {
			componentMap[key] = value
		}
	}

	// Build up a ResourceID from the map
	idObj := &resourceID{}
	idObj.Path = componentMap

	if subscriptionID != "" {
		idObj.SubscriptionID = subscriptionID
	} else {
		return nil, fmt.Errorf("no subscription ID found in: %q", path)
	}

	if resourceGroup, ok := componentMap["resourceGroups"]; ok {
		idObj.ResourceGroup = resourceGroup
		delete(componentMap, "resourceGroups")
	} else {
		// Some Azure APIs are weird and provide things in lower case...
		// However it's not clear whether the casing of other elements in the URI
		// matter, so we explicitly look for that case here.
		if resourceGroup, ok := componentMap["resourcegroups"]; ok {
			idObj.ResourceGroup = resourceGroup
			delete(componentMap, "resourcegroups")
		} else {
			return nil, fmt.Errorf("no resource group name found in: %q", path)
		}
	}

	// It is OK not to have a provider in the case of a resource group
	if provider, ok := componentMap["providers"]; ok {
		idObj.Provider = provider
		delete(componentMap, "providers")
	}

	return idObj, nil
}

type retryable struct {
	message string
}

func (r *retryable) Error() string {
	return r.message
}

func retry(times int, delay time.Duration, action func() (interface{}, error)) (interface{}, error) {
	var lastErr error
	for i := 0; i < times; i++ {
		item, err := action()
		if err != nil {
			if err, ok := err.(*retryable); ok {
				lastErr = err
				time.Sleep(delay)
				continue
			} else {
				return nil, err
			}
		}
		return item, nil
	}
	return nil, lastErr
}
