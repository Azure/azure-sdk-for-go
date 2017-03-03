package storage

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// lease constants.
const (
	leaseHeaderPrefix = "x-ms-lease-"
	headerLeaseID     = "x-ms-lease-id"
	leaseAction       = "x-ms-lease-action"
	leaseBreakPeriod  = "x-ms-lease-break-period"
	leaseDuration     = "x-ms-lease-duration"
	leaseProposedID   = "x-ms-proposed-lease-id"
	leaseTime         = "x-ms-lease-time"

	acquireLease = "acquire"
	renewLease   = "renew"
	changeLease  = "change"
	releaseLease = "release"
	breakLease   = "break"
)

// leasePut is common PUT code for the various acquire/release/break etc functions.
func (b *Blob) leaseCommonPut(headers map[string]string, expectedStatus int) (http.Header, error) {
	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), url.Values{"comp": {"lease"}})

	resp, err := b.Container.bsc.client.exec(http.MethodPut, uri, headers, nil, b.Container.bsc.auth)
	if err != nil {
		return nil, err
	}
	defer readAndCloseBody(resp.body)

	if err := checkRespCode(resp.statusCode, []int{expectedStatus}); err != nil {
		return nil, err
	}

	return resp.headers, nil
}

// AcquireLease creates a lease for a blob as per https://msdn.microsoft.com/en-us/library/azure/ee691972.aspx
// returns leaseID acquired
// In API Versions starting on 2012-02-12, the minimum leaseTimeInSeconds is 15, the maximum
// non-infinite leaseTimeInSeconds is 60. To specify an infinite lease, provide the value -1.
func (b *Blob) AcquireLease(leaseTimeInSeconds int, proposedLeaseID string) (returnedLeaseID string, err error) {
	headers := b.Container.bsc.client.getStandardHeaders()
	headers[leaseAction] = acquireLease

	if leaseTimeInSeconds == -1 {
		// Do nothing, but don't trigger the following clauses.
	} else if leaseTimeInSeconds > 60 || b.Container.bsc.client.apiVersion < "2012-02-12" {
		leaseTimeInSeconds = 60
	} else if leaseTimeInSeconds < 15 {
		leaseTimeInSeconds = 15
	}

	headers[leaseDuration] = strconv.Itoa(leaseTimeInSeconds)

	if proposedLeaseID != "" {
		headers[leaseProposedID] = proposedLeaseID
	}

	respHeaders, err := b.leaseCommonPut(headers, http.StatusCreated)
	if err != nil {
		return "", err
	}

	returnedLeaseID = respHeaders.Get(http.CanonicalHeaderKey(headerLeaseID))

	if returnedLeaseID != "" {
		return returnedLeaseID, nil
	}

	return "", errors.New("LeaseID not returned")
}

// BreakLease breaks the lease for a blob as per https://msdn.microsoft.com/en-us/library/azure/ee691972.aspx
// Returns the timeout remaining in the lease in seconds
func (b *Blob) BreakLease() (breakTimeout int, err error) {
	headers := b.Container.bsc.client.getStandardHeaders()
	headers[leaseAction] = breakLease
	return b.breakLeaseCommon(headers)
}

// BreakLeaseWithBreakPeriod breaks the lease for a blob as per https://msdn.microsoft.com/en-us/library/azure/ee691972.aspx
// breakPeriodInSeconds is used to determine how long until new lease can be created.
// Returns the timeout remaining in the lease in seconds
func (b *Blob) BreakLeaseWithBreakPeriod(breakPeriodInSeconds int) (breakTimeout int, err error) {
	headers := b.Container.bsc.client.getStandardHeaders()
	headers[leaseAction] = breakLease
	headers[leaseBreakPeriod] = strconv.Itoa(breakPeriodInSeconds)
	return b.breakLeaseCommon(headers)
}

// breakLeaseCommon is common code for both version of BreakLease (with and without break period)
func (b *Blob) breakLeaseCommon(headers map[string]string) (breakTimeout int, err error) {

	respHeaders, err := b.leaseCommonPut(headers, http.StatusAccepted)
	if err != nil {
		return 0, err
	}

	breakTimeoutStr := respHeaders.Get(http.CanonicalHeaderKey(leaseTime))
	if breakTimeoutStr != "" {
		breakTimeout, err = strconv.Atoi(breakTimeoutStr)
		if err != nil {
			return 0, err
		}
	}

	return breakTimeout, nil
}

// ChangeLease changes a lease ID for a blob as per https://msdn.microsoft.com/en-us/library/azure/ee691972.aspx
// Returns the new LeaseID acquired
func (b *Blob) ChangeLease(currentLeaseID string, proposedLeaseID string) (newLeaseID string, err error) {
	headers := b.Container.bsc.client.getStandardHeaders()
	headers[leaseAction] = changeLease
	headers[headerLeaseID] = currentLeaseID
	headers[leaseProposedID] = proposedLeaseID

	respHeaders, err := b.leaseCommonPut(headers, http.StatusOK)
	if err != nil {
		return "", err
	}

	newLeaseID = respHeaders.Get(http.CanonicalHeaderKey(headerLeaseID))
	if newLeaseID != "" {
		return newLeaseID, nil
	}

	return "", errors.New("LeaseID not returned")
}

// ReleaseLease releases the lease for a blob as per https://msdn.microsoft.com/en-us/library/azure/ee691972.aspx
func (b *Blob) ReleaseLease(currentLeaseID string) error {
	headers := b.Container.bsc.client.getStandardHeaders()
	headers[leaseAction] = releaseLease
	headers[headerLeaseID] = currentLeaseID

	_, err := b.leaseCommonPut(headers, http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

// RenewLease renews the lease for a blob as per https://msdn.microsoft.com/en-us/library/azure/ee691972.aspx
func (b *Blob) RenewLease(currentLeaseID string) error {
	headers := b.Container.bsc.client.getStandardHeaders()
	headers[leaseAction] = renewLease
	headers[headerLeaseID] = currentLeaseID

	_, err := b.leaseCommonPut(headers, http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}
