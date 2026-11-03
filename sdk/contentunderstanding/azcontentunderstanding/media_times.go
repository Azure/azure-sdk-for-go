// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding

import (
	"errors"
	"math"
	"time"
)

// DurationFromMilliseconds converts service milliseconds without overflowing Go's duration.
func DurationFromMilliseconds(milliseconds int64) (time.Duration, error) {
	if milliseconds > math.MaxInt64/int64(time.Millisecond) || milliseconds < math.MinInt64/int64(time.Millisecond) {
		return 0, errors.New("media timestamp exceeds time.Duration range")
	}
	return time.Duration(milliseconds) * time.Millisecond, nil
}

func mediaTime(value *int64) (time.Duration, error) {
	if value == nil {
		return 0, errors.New("media timestamp is missing")
	}
	return DurationFromMilliseconds(*value)
}

func mediaTimes(values []*int64) ([]time.Duration, error) {
	result := make([]time.Duration, len(values))
	for i, value := range values {
		duration, err := mediaTime(value)
		if err != nil {
			return nil, err
		}
		result[i] = duration
	}
	return result, nil
}

// StartTime returns the media start with millisecond precision.
func (c *AudioVisualContent) StartTime() (time.Duration, error) {
	if c == nil {
		return 0, errors.New("content is required")
	}
	return mediaTime(c.StartTimeMs)
}

// EndTime returns the media end with millisecond precision.
func (c *AudioVisualContent) EndTime() (time.Duration, error) {
	if c == nil {
		return 0, errors.New("content is required")
	}
	return mediaTime(c.EndTimeMs)
}

// KeyFrameTimes returns converted keyframe timestamps.
func (c *AudioVisualContent) KeyFrameTimes() ([]time.Duration, error) {
	if c == nil {
		return nil, errors.New("content is required")
	}
	return mediaTimes(c.KeyFrameTimesMs)
}

// CameraShotTimes returns converted camera-shot timestamps.
func (c *AudioVisualContent) CameraShotTimes() ([]time.Duration, error) {
	if c == nil {
		return nil, errors.New("content is required")
	}
	return mediaTimes(c.CameraShotTimesMs)
}

// StartTime returns a segment's start timestamp.
func (s *AudioVisualContentSegment) StartTime() (time.Duration, error) {
	if s == nil {
		return 0, errors.New("segment is required")
	}
	return mediaTime(s.StartTimeMs)
}

// EndTime returns a segment's end timestamp.
func (s *AudioVisualContentSegment) EndTime() (time.Duration, error) {
	if s == nil {
		return 0, errors.New("segment is required")
	}
	return mediaTime(s.EndTimeMs)
}

// StartTime returns a transcript phrase's start timestamp.
func (p *TranscriptPhrase) StartTime() (time.Duration, error) {
	if p == nil {
		return 0, errors.New("phrase is required")
	}
	return mediaTime(p.StartTimeMs)
}

// EndTime returns a transcript phrase's end timestamp.
func (p *TranscriptPhrase) EndTime() (time.Duration, error) {
	if p == nil {
		return 0, errors.New("phrase is required")
	}
	return mediaTime(p.EndTimeMs)
}

// StartTime returns a transcript word's start timestamp.
func (w *TranscriptWord) StartTime() (time.Duration, error) {
	if w == nil {
		return 0, errors.New("word is required")
	}
	return mediaTime(w.StartTimeMs)
}

// EndTime returns a transcript word's end timestamp.
func (w *TranscriptWord) EndTime() (time.Duration, error) {
	if w == nil {
		return 0, errors.New("word is required")
	}
	return mediaTime(w.EndTimeMs)
}
