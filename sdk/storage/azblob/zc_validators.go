// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (pr *PageRange) Raw() (start, end int64) {
	if pr.Start != nil {
		start = *pr.Start
	}
	if pr.End != nil {
		end = *pr.End
	}

	return
}

// HttpRange defines a range of bytes within an HTTP resource, starting at Start and
// ending at End. A zero-value HttpRange indicates the entire resource. An HttpRange
// which has an Start but na zero value End indicates from the Start to the resource's end.
type HttpRange struct {
	Start int64
	End   int64
}

func (r HttpRange) pointers() *string {
	if r.Start == 0 && r.End == 0 { // Do common case first for performance
		return nil // No specified range
	}
	end := "" // if End == CountToEnd (0)
	if r.End > 0 {
		end = strconv.FormatInt(r.End, 10)
	}
	dataRange := fmt.Sprintf("bytes=%v-%s", r.Start, end)
	return &dataRange
}

func getSourceRange(start, end *int64) *string {
	if start == nil && end == nil {
		return nil
	}
	newStart := int64(0)
	newEnd := int64(CountToEnd)

	if start != nil {
		newStart = *start
	}

	if end != nil {
		newEnd = *end
	}

	return HttpRange{Start: newStart, End: newEnd}.pointers()
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

	body.Seek(0, io.SeekStart)
	return count, nil
}

// return an error if body is not a valid seekable stream at 0
func validateSeekableStreamAt0(body io.ReadSeeker) error {
	if body == nil { // nil body's are "logically" seekable to 0
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
