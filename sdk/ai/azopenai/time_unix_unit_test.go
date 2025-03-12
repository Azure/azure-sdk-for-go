package azopenai

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPopulateTimeUnix(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		objectMap := make(map[string]any)

		populateTimeUnix(objectMap, "time", nil)
		_, exists := objectMap["time"]
		assert.False(t, exists)

		time_rfc, _ := time.Parse(time.RFC3339, "2025-03-03T15:04:05.000Z")
		populateTimeUnix(objectMap, "time", &time_rfc)
		assert.Equal(t, objectMap["time"], (*timeUnix)(&time_rfc))
	})
}

func TestUnpopulateTimeUnix(t *testing.T) {
	type ResponseData struct {
		Timestamp *time.Time
	}

	testCases := []struct {
		name     string
		data     json.RawMessage
		expected any
		isError  bool
	}{
		{
			name:     "valid time",
			data:     json.RawMessage("1741791845"),
			expected: time.Unix(1741791845, 0).Local(),
			isError:  false,
		},
		{
			name:     "null time",
			data:     json.RawMessage("null"),
			expected: nil,
			isError:  false,
		},
		{
			name:     "nil time",
			data:     nil,
			expected: nil,
			isError:  false,
		},
		{
			name:     "invalid time",
			data:     json.RawMessage("invalid"),
			expected: nil,
			isError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response_data := ResponseData{Timestamp: nil}
			err := unpopulateTimeUnix(tc.data, "Timestamp", &response_data.Timestamp)
			if tc.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tc.expected == nil {
				assert.Nil(t, response_data.Timestamp)
			} else {
				expected := tc.expected.(time.Time)
				assert.Equal(t, &expected, response_data.Timestamp)
			}
		})
	}
}

func TestTimeUnix_MarshalJSON(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		time_ip := time.Unix(1741791845, 0)
		time_unix := timeUnix(time_ip)
		result, err := time_unix.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, []byte("1741791845"), result)
	})
}

func TestTimeUnix_UnmarshalJSON(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		time_ip := time.Unix(1741791845, 0)
		time_unix := timeUnix(time_ip)
		err := time_unix.UnmarshalJSON([]byte("1741791846"))
		assert.NoError(t, err)
		assert.Equal(t, timeUnix(time.Unix(1741791846, 0)), time_unix)
	})
}
