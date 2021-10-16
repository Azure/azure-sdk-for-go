package azfile

import (
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"io"
	"strconv"
)

const (
	// CountToEnd indicates a flag for Count parameter. It means the Count of bytes
	// from start Offset to the end of file.
	CountToEnd = 0
)

// HttpRange defines a range of bytes within an HTTP resource, starting at Offset and
// ending at Offset+Count. A zero-value HttpRange indicates the entire resource. An HttpRange
// which has an Offset but na zero value Count indicates from the Offset to the resource's end.
type HttpRange struct {
	Offset int64
	Count  int64
}

func (r HttpRange) format() *string {
	if r.Offset == 0 && r.Count == 0 { // Do common case first for performance
		return nil // No specified range
	}
	endOffset := "" // if Count == CountToEnd (0)
	if r.Count > 0 {
		endOffset = strconv.FormatInt((r.Offset+r.Count)-1, 10)
	}
	return to.StringPtr(fmt.Sprintf("bytes=%v-%s", r.Offset, endOffset))
}

// toRange makes range string adhere to REST API.
// A Count with value CountToEnd means Count of bytes from Offset to the end of file.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/specifying-the-range-header-for-file-service-operations.
func getRangeParam(offset, count *int64) *string {
	if offset == nil && count == nil {
		return nil
	}
	newOffset := int64(0)
	newCount := int64(CountToEnd)

	if offset != nil {
		newOffset = *offset
	}

	if count != nil {
		newCount = *count
	}

	return HttpRange{Offset: newOffset, Count: newCount}.format()
}

func validateSeekableStreamAt0AndGetCount(body io.ReadSeeker) (int64, error) {
	if body == nil { // nil body's are "logically" seekable to 0 and are 0 bytes long
		return 0, nil
	}

	err := validateSeekableStreamAt0(body)
	if err != nil {
		return 0, err
	}

	count, err := body.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, errors.New("body stream must be seekable")
	}

	_, err = body.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// return an error if body is not a valid seekable stream at 0
func validateSeekableStreamAt0(body io.ReadSeeker) error {
	if body == nil { // nil body is "logically" seekable to 0
		return nil
	}
	if pos, err := body.Seek(0, io.SeekCurrent); pos != 0 || err != nil {
		// Help detect programmer error
		if err != nil {
			return errors.New("body stream must be seekable")
		}
		return errors.New("body stream must be set to position 0")
	}
	return nil
}
