// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"
)

// FaultInjectionOperation selects a supported item or metadata operation; zero matches any.
type FaultInjectionOperation int32

const (
	// FaultInjectionReadItem matches point reads.
	FaultInjectionReadItem FaultInjectionOperation = 1
	// FaultInjectionCreateItem matches item creation.
	FaultInjectionCreateItem FaultInjectionOperation = 3
	// FaultInjectionUpsertItem matches item creation or replacement.
	FaultInjectionUpsertItem FaultInjectionOperation = 4
	// FaultInjectionReplaceItem matches replacements.
	FaultInjectionReplaceItem FaultInjectionOperation = 5
	// FaultInjectionDeleteItem matches deletes.
	FaultInjectionDeleteItem FaultInjectionOperation = 6
	// FaultInjectionPatchItem matches patches.
	FaultInjectionPatchItem FaultInjectionOperation = 7
	// FaultInjectionReadContainer matches container metadata reads.
	FaultInjectionReadContainer FaultInjectionOperation = 10
	// FaultInjectionReadAccount matches account metadata reads.
	FaultInjectionReadAccount FaultInjectionOperation = 11
	// FaultInjectionReadPartitionRanges matches partition topology reads.
	FaultInjectionReadPartitionRanges FaultInjectionOperation = 13
)

// FaultInjectionError selects a native fault. Zero requires a custom HTTP status.
type FaultInjectionError int32

const (
	// FaultInjectionInternalServerError injects an internal server error.
	FaultInjectionInternalServerError FaultInjectionError = 1
	// FaultInjectionTooManyRequests injects throttling.
	FaultInjectionTooManyRequests FaultInjectionError = 2
	// FaultInjectionRetryWith injects a retry-with response.
	FaultInjectionRetryWith FaultInjectionError = 3
	// FaultInjectionReadSessionNotAvailable injects session unavailability.
	FaultInjectionReadSessionNotAvailable FaultInjectionError = 4
	// FaultInjectionTimeout injects a timeout.
	FaultInjectionTimeout FaultInjectionError = 5
	// FaultInjectionServiceUnavailable injects service unavailability.
	FaultInjectionServiceUnavailable FaultInjectionError = 6
	// FaultInjectionPartitionIsGone injects a gone partition.
	FaultInjectionPartitionIsGone FaultInjectionError = 7
	// FaultInjectionWriteForbidden injects a write-forbidden response.
	FaultInjectionWriteForbidden FaultInjectionError = 8
	// FaultInjectionAccountNotFound injects account-not-found.
	FaultInjectionAccountNotFound FaultInjectionError = 9
	// FaultInjectionConnectionError injects a connection error.
	FaultInjectionConnectionError FaultInjectionError = 10
	// FaultInjectionResponseTimeout injects a response timeout before service execution.
	FaultInjectionResponseTimeout FaultInjectionError = 11
	// FaultInjectionResponseTimeoutAfterService injects a timeout after service execution.
	FaultInjectionResponseTimeoutAfterService FaultInjectionError = 12
)

// FaultInjectionTransport selects a transport; zero matches any.
type FaultInjectionTransport int32

const (
	// FaultInjectionGateway matches Gateway V1.
	FaultInjectionGateway FaultInjectionTransport = 1
	// FaultInjectionGatewayV2 matches Gateway V2.
	FaultInjectionGatewayV2 FaultInjectionTransport = 2
)

// FaultInjectionCondition selects requests affected by a test rule. Zero fields match any.
type FaultInjectionCondition struct {
	// Operation selects an item or metadata operation.
	Operation FaultInjectionOperation
	// Region restricts matching to a region.
	Region Region
	// ContainerID restricts matching to a container.
	ContainerID string
	// Transport restricts matching to a gateway transport.
	Transport FaultInjectionTransport
}

// FaultInjectionResult configures an injected native transport result for testing.
// Error and CustomStatusCode are mutually exclusive. Native code validates fault-specific constraints.
type FaultInjectionResult struct {
	// Error selects a predefined native fault.
	Error FaultInjectionError
	// Delay delays the fault. Nil uses the native default.
	Delay *time.Duration
	// Probability is between zero and one. Nil defaults to one.
	Probability *float32
	// CustomStatusCode selects a custom HTTP response; zero means unset.
	CustomStatusCode int
	// CustomSubStatus sets the custom Cosmos substatus. Nil means unset.
	CustomSubStatus *uint16
	// RetryAfter sets the custom retry delay. Nil means unset.
	RetryAfter *time.Duration
	// Headers supplies custom response headers.
	Headers map[string]string
	// Body is copied verbatim into a custom response.
	Body []byte
}

// FaultInjectionRule configures native fault injection on one client. Use only for testing.
// Rules and nested data are copied at client construction; they do not mutate a shared runtime.
type FaultInjectionRule struct {
	// ID uniquely identifies the rule within the client.
	ID string
	// Condition selects requests to match.
	Condition FaultInjectionCondition
	// Result defines the injected outcome.
	Result FaultInjectionResult
	// HitLimit bounds matching injections. Nil is unlimited; zero disables the rule.
	HitLimit *uint32
	// StartDelay delays activation. Nil starts immediately.
	StartDelay *time.Duration
	// ExpireAfter bounds the active duration after StartDelay. Nil has no expiration.
	ExpireAfter *time.Duration
}

func cloneFaultInjectionRules(rules []FaultInjectionRule) []FaultInjectionRule {
	copy := slices.Clone(rules)
	for i := range copy {
		r := &copy[i]
		r.HitLimit = clonePointer(r.HitLimit)
		r.StartDelay = clonePointer(r.StartDelay)
		r.ExpireAfter = clonePointer(r.ExpireAfter)
		r.Result.Delay = clonePointer(r.Result.Delay)
		r.Result.Probability = clonePointer(r.Result.Probability)
		r.Result.CustomSubStatus = clonePointer(r.Result.CustomSubStatus)
		r.Result.RetryAfter = clonePointer(r.Result.RetryAfter)
		r.Result.Headers = (OperationOptions{CustomHeaders: r.Result.Headers}).clone().CustomHeaders
		r.Result.Body = slices.Clone(r.Result.Body)
	}
	return copy
}

func validateFaultInjectionRules(rules []FaultInjectionRule) error {
	seen := make(map[string]bool, len(rules))
	for _, r := range rules {
		if r.ID == "" || strings.ContainsRune(r.ID, 0) || seen[r.ID] {
			return fmt.Errorf("azcosmos: fault injection rule IDs must be nonempty, unique, and contain no NUL")
		}
		seen[r.ID] = true
		switch r.Condition.Operation {
		case 0, FaultInjectionReadItem, FaultInjectionCreateItem, FaultInjectionUpsertItem,
			FaultInjectionReplaceItem, FaultInjectionDeleteItem, FaultInjectionPatchItem,
			FaultInjectionReadContainer, FaultInjectionReadAccount, FaultInjectionReadPartitionRanges:
		default:
			return fmt.Errorf("azcosmos: invalid fault injection operation")
		}
		if r.Condition.Transport < 0 || r.Condition.Transport > 2 ||
			r.Result.Error < 0 || r.Result.Error > 12 {
			return fmt.Errorf("azcosmos: invalid fault injection transport or error")
		}
		if (r.Result.Error == 0) == (r.Result.CustomStatusCode == 0) ||
			r.Result.CustomStatusCode != 0 && (r.Result.CustomStatusCode < 100 || r.Result.CustomStatusCode > 599) {
			return fmt.Errorf("azcosmos: select a fault injection error or a custom HTTP status")
		}
		for _, value := range []*time.Duration{r.StartDelay, r.ExpireAfter, r.Result.Delay, r.Result.RetryAfter} {
			if value != nil && *value < 0 {
				return fmt.Errorf("azcosmos: fault injection durations must not be negative")
			}
		}
		if p := r.Result.Probability; p != nil && (math.IsNaN(float64(*p)) || *p < 0 || *p > 1) {
			return fmt.Errorf("azcosmos: fault injection probability must be between zero and one")
		}
		if strings.ContainsRune(string(r.Condition.Region), 0) || strings.ContainsRune(r.Condition.ContainerID, 0) {
			return fmt.Errorf("azcosmos: fault injection condition must not contain NUL")
		}
		if err := (OperationOptions{CustomHeaders: r.Result.Headers}).validate(); err != nil {
			return err
		}
	}
	return nil
}
