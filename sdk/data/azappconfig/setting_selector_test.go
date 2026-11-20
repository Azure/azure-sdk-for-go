// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfig

import (
	"testing"
	"time"
)

func TestSettingSelectorAcceptDateTime(t *testing.T) {
	at := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	want := at.Format(timeFormat)
	sc := SettingSelector{AcceptDateTime: &at}

	kv := sc.toGeneratedGetKeyValues()
	if kv.After != nil {
		t.Errorf("GetKeyValues: AcceptDateTime must not be sent as After, got %q", *kv.After)
	}
	if kv.AcceptDatetime == nil || *kv.AcceptDatetime != want {
		t.Errorf("GetKeyValues: AcceptDatetime = %v, want %q", kv.AcceptDatetime, want)
	}

	rev := sc.toGeneratedGetRevisions()
	if rev.After != nil {
		t.Errorf("GetRevisions: AcceptDateTime must not be sent as After, got %q", *rev.After)
	}
	if rev.AcceptDatetime == nil || *rev.AcceptDatetime != want {
		t.Errorf("GetRevisions: AcceptDatetime = %v, want %q", rev.AcceptDatetime, want)
	}

	chk := sc.toGeneratedCheckKeyValues()
	if chk.AcceptDatetime == nil || *chk.AcceptDatetime != want {
		t.Errorf("CheckKeyValues: AcceptDatetime = %v, want %q", chk.AcceptDatetime, want)
	}
}
