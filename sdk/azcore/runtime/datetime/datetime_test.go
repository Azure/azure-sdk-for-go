package datetime

import (
	"testing"
	"time"
)

func TestParseSuccess(t *testing.T) {
	tests := []struct {
		name  string
		parse func(string) (time.Time, error)
		input string
		want  time.Time
	}{
		{
			name:  "PlainDate",
			parse: func(s string) (time.Time, error) { return Parse[PlainDate](s) },
			input: "2024-01-02",
			want:  time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "PlainTime",
			parse: func(s string) (time.Time, error) { return Parse[PlainTime](s) },
			input: "09:08:07",
			want:  time.Date(0, time.January, 1, 9, 8, 7, 0, time.UTC),
		},
		{
			name:  "RFC3339",
			parse: func(s string) (time.Time, error) { return Parse[RFC3339](s) },
			input: "2024-01-02T03:04:05Z",
			want:  time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC),
		},
		{
			name:  "RFC1123",
			parse: func(s string) (time.Time, error) { return Parse[RFC1123](s) },
			input: "Mon, 02 Jan 2006 15:04:05 GMT",
			want:  time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
		},
		{
			name:  "Unix",
			parse: func(s string) (time.Time, error) { return Parse[Unix](s) },
			input: "1717171717",
			want:  time.Unix(1717171717, 0).UTC(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.parse(tc.input)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
			if !got.Equal(tc.want) {
				t.Fatalf("unexpected result: got %v want %v", got, tc.want)
			}
		})
	}
}

func TestParseFailure(t *testing.T) {
	tests := []struct {
		name  string
		parse func(string) (time.Time, error)
		input string
	}{
		{
			name:  "PlainDateInvalid",
			parse: func(s string) (time.Time, error) { return Parse[PlainDate](s) },
			input: "not-a-date",
		},
		{
			name:  "UnixInvalid",
			parse: func(s string) (time.Time, error) { return Parse[Unix](s) },
			input: "not-a-number",
		},
		{
			name:  "PlainTimeInvalid",
			parse: func(s string) (time.Time, error) { return Parse[PlainTime](s) },
			input: "25:00:00",
		},
		{
			name:  "RFC3339Invalid",
			parse: func(s string) (time.Time, error) { return Parse[RFC3339](s) },
			input: "2024-13-99T25:61:61Z",
		},
		{
			name:  "RFC1123Invalid",
			parse: func(s string) (time.Time, error) { return Parse[RFC1123](s) },
			input: "Bad RFC1123 string",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.parse(tc.input)
			if err == nil {
				t.Fatalf("expected error but got none (result: %v)", got)
			}
			if !got.IsZero() {
				t.Fatalf("expected zero time on error, got %v", got)
			}
		})
	}
}
