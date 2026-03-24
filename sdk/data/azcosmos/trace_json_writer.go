// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const (
	rootTraceTimeFormat   = "2006-01-02T15:04:05.000Z"
	nestedTraceTimeFormat = "2006-01-02T15:04:05.0000000Z07:00"
)

type summaryDiagnostics struct {
	directCalls    map[[2]int]int
	gatewayCalls   map[[2]int]int
	regionsByURI   map[string]struct{}
	regionURICount int
}

func writeTraceJSON(root *trace) string {
	var buffer bytes.Buffer
	writeTrace(&buffer, root, true)
	return buffer.String()
}

func writeTrace(buffer *bytes.Buffer, current *trace, isRoot bool) {
	snapshot := current.snapshot()

	buffer.WriteByte('{')
	firstField := true

	writeField := func(name string, writeValue func()) {
		if !firstField {
			buffer.WriteByte(',')
		}
		firstField = false
		writeJSONString(buffer, name)
		buffer.WriteByte(':')
		writeValue()
	}

	if isRoot {
		writeField("Summary", func() {
			writeSummaryDiagnostics(buffer, current)
		})
	}

	writeField("name", func() {
		writeJSONString(buffer, snapshot.name)
	})

	if isRoot {
		writeField("start datetime", func() {
			writeJSONString(buffer, snapshot.startTime.UTC().Format(rootTraceTimeFormat))
		})
	}

	writeField("duration in milliseconds", func() {
		writeNumber(buffer, durationInMilliseconds(snapshot.startTime, snapshot.endTime))
	})

	if len(snapshot.dataOrder) > 0 {
		writeField("data", func() {
			writeTraceDataObject(buffer, snapshot.dataOrder, snapshot.data)
		})
	}

	if len(snapshot.children) > 0 {
		writeField("children", func() {
			buffer.WriteByte('[')
			for index, child := range snapshot.children {
				if index > 0 {
					buffer.WriteByte(',')
				}
				writeTrace(buffer, child, false)
			}
			buffer.WriteByte(']')
		})
	}

	buffer.WriteByte('}')
}

func writeSummaryDiagnostics(buffer *bytes.Buffer, root *trace) {
	summary := collectSummaryDiagnostics(root)

	buffer.WriteByte('{')
	firstField := true

	writeField := func(name string, writeValue func()) {
		if !firstField {
			buffer.WriteByte(',')
		}
		firstField = false
		writeJSONString(buffer, name)
		buffer.WriteByte(':')
		writeValue()
	}

	if len(summary.directCalls) > 0 {
		writeField("DirectCalls", func() {
			writeCallSummaryObject(buffer, summary.directCalls)
		})
	}

	if summary.regionURICount > 0 {
		writeField("RegionsContacted", func() {
			writeNumber(buffer, float64(summary.regionURICount))
		})
	}

	if len(summary.gatewayCalls) > 0 {
		writeField("GatewayCalls", func() {
			writeCallSummaryObject(buffer, summary.gatewayCalls)
		})
	}

	buffer.WriteByte('}')
}

func collectSummaryDiagnostics(root *trace) summaryDiagnostics {
	summary := summaryDiagnostics{
		directCalls:  map[[2]int]int{},
		gatewayCalls: map[[2]int]int{},
		regionsByURI: map[string]struct{}{},
	}

	var walk func(*trace)
	walk = func(current *trace) {
		snapshot := current.snapshot()
		for _, key := range snapshot.dataOrder {
			value := snapshot.data[key]
			stats, ok := value.(*clientSideRequestStatisticsTraceDatum)
			if !ok {
				continue
			}

			statsSnapshot := stats.snapshot()
			for _, region := range statsSnapshot.regionsContacted {
				if _, exists := summary.regionsByURI[region.uri]; !exists {
					summary.regionsByURI[region.uri] = struct{}{}
					summary.regionURICount++
				}
			}

			for _, gatewayCall := range statsSnapshot.httpResponseStatistics {
				key := [2]int{gatewayCall.statusCode, gatewayCall.subStatusCode}
				summary.gatewayCalls[key]++
			}
		}

		for _, child := range snapshot.children {
			walk(child)
		}
	}

	walk(root)
	return summary
}

func writeCallSummaryObject(buffer *bytes.Buffer, values map[[2]int]int) {
	keys := make([][2]int, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i int, j int) bool {
		if keys[i][0] == keys[j][0] {
			return keys[i][1] < keys[j][1]
		}
		return keys[i][0] < keys[j][0]
	})

	buffer.WriteByte('{')
	firstField := true
	for _, key := range keys {
		if !firstField {
			buffer.WriteByte(',')
		}
		firstField = false
		writeJSONString(buffer, fmt.Sprintf("(%d, %d)", key[0], key[1]))
		buffer.WriteByte(':')
		writeNumber(buffer, float64(values[key]))
	}
	buffer.WriteByte('}')
}

func writeTraceDataObject(buffer *bytes.Buffer, order []string, values map[string]any) {
	buffer.WriteByte('{')
	for index, key := range order {
		if index > 0 {
			buffer.WriteByte(',')
		}
		writeJSONString(buffer, key)
		buffer.WriteByte(':')
		writeTraceDatum(buffer, values[key])
	}
	buffer.WriteByte('}')
}

func writeTraceDatum(buffer *bytes.Buffer, value any) {
	switch typed := value.(type) {
	case *clientSideRequestStatisticsTraceDatum:
		writeClientSideRequestStatistics(buffer, typed.snapshot())
	case pointOperationStatisticsTraceDatum:
		writePointOperationStatistics(buffer, typed)
	case *pointOperationStatisticsTraceDatum:
		writePointOperationStatistics(buffer, *typed)
	case json.RawMessage:
		buffer.Write(typed)
	case string:
		writeJSONString(buffer, typed)
	case float64:
		writeNumber(buffer, typed)
	case float32:
		writeNumber(buffer, float64(typed))
	case int:
		writeNumber(buffer, float64(typed))
	case int32:
		writeNumber(buffer, float64(typed))
	case int64:
		writeNumber(buffer, float64(typed))
	case bool:
		if typed {
			buffer.WriteString("true")
		} else {
			buffer.WriteString("false")
		}
	case []string:
		buffer.WriteByte('[')
		for index, item := range typed {
			if index > 0 {
				buffer.WriteByte(',')
			}
			writeJSONString(buffer, item)
		}
		buffer.WriteByte(']')
	case []any:
		buffer.WriteByte('[')
		for index, item := range typed {
			if index > 0 {
				buffer.WriteByte(',')
			}
			writeTraceDatum(buffer, item)
		}
		buffer.WriteByte(']')
	case map[string]any:
		buffer.WriteByte('{')
		firstField := true
		for key, item := range typed {
			if !firstField {
				buffer.WriteByte(',')
			}
			firstField = false
			writeJSONString(buffer, key)
			buffer.WriteByte(':')
			writeTraceDatum(buffer, item)
		}
		buffer.WriteByte('}')
	case nil:
		buffer.WriteString("null")
	default:
		writeJSONString(buffer, fmt.Sprint(value))
	}
}

func writeClientSideRequestStatistics(buffer *bytes.Buffer, snapshot clientSideRequestStatisticsSnapshot) {
	buffer.WriteByte('{')
	firstField := true

	writeField := func(name string, writeValue func()) {
		if !firstField {
			buffer.WriteByte(',')
		}
		firstField = false
		writeJSONString(buffer, name)
		buffer.WriteByte(':')
		writeValue()
	}

	writeField("Id", func() {
		writeJSONString(buffer, "AggregatedClientSideRequestStatistics")
	})

	writeField("ContactedReplicas", func() {
		buffer.WriteByte('[')
		buffer.WriteByte(']')
	})

	writeField("RegionsContacted", func() {
		buffer.WriteByte('[')
		for index, region := range snapshot.regionsContacted {
			if index > 0 {
				buffer.WriteByte(',')
			}
			writeJSONString(buffer, region.uri)
		}
		buffer.WriteByte(']')
	})

	writeField("FailedReplicas", func() {
		buffer.WriteByte('[')
		buffer.WriteByte(']')
	})

	if len(snapshot.forceAddressRefreshes) > 0 {
		writeField("ForceAddressRefresh", func() {
			buffer.WriteByte('[')
			for index, refresh := range snapshot.forceAddressRefreshes {
				if index > 0 {
					buffer.WriteByte(',')
				}
				buffer.WriteByte('{')
				firstRefreshField := true
				writeRefreshField := func(name string, values []string) {
					if !firstRefreshField {
						buffer.WriteByte(',')
					}
					firstRefreshField = false
					writeJSONString(buffer, name)
					buffer.WriteByte(':')
					writeTraceDatum(buffer, values)
				}
				if len(refresh.noChangeToCache) > 0 {
					writeRefreshField("No change to cache", refresh.noChangeToCache)
				} else {
					writeRefreshField("Original", refresh.original)
					writeRefreshField("New", refresh.newValues)
				}
				buffer.WriteByte('}')
			}
			buffer.WriteByte(']')
		})
	}

	writeField("AddressResolutionStatistics", func() {
		buffer.WriteByte('[')
		for index, stat := range snapshot.addressResolutionStatistics {
			if index > 0 {
				buffer.WriteByte(',')
			}
			buffer.WriteByte('{')
			writeJSONString(buffer, "StartTimeUTC")
			buffer.WriteByte(':')
			writeJSONString(buffer, stat.startTimeUTC.UTC().Format(nestedTraceTimeFormat))
			buffer.WriteByte(',')
			writeJSONString(buffer, "EndTimeUTC")
			buffer.WriteByte(':')
			if stat.endTimeUTC != nil {
				writeJSONString(buffer, stat.endTimeUTC.UTC().Format(nestedTraceTimeFormat))
			} else {
				writeJSONString(buffer, "EndTime Never Set.")
			}
			buffer.WriteByte(',')
			writeJSONString(buffer, "TargetEndpoint")
			buffer.WriteByte(':')
			if stat.targetEndpoint == "" {
				buffer.WriteString("null")
			} else {
				writeJSONString(buffer, stat.targetEndpoint)
			}
			buffer.WriteByte('}')
		}
		buffer.WriteByte(']')
	})

	writeField("StoreResponseStatistics", func() {
		buffer.WriteByte('[')
		buffer.WriteByte(']')
	})

	if len(snapshot.httpResponseStatistics) > 0 {
		writeField("HttpResponseStats", func() {
			buffer.WriteByte('[')
			for index, stat := range snapshot.httpResponseStatistics {
				if index > 0 {
					buffer.WriteByte(',')
				}
				writeHTTPResponseStatistic(buffer, stat)
			}
			buffer.WriteByte(']')
		})
	}

	buffer.WriteByte('}')
}

func writeHTTPResponseStatistic(buffer *bytes.Buffer, stat httpResponseStatistics) {
	buffer.WriteByte('{')
	firstField := true

	writeField := func(name string, writeValue func()) {
		if !firstField {
			buffer.WriteByte(',')
		}
		firstField = false
		writeJSONString(buffer, name)
		buffer.WriteByte(':')
		writeValue()
	}

	writeField("StartTimeUTC", func() {
		writeJSONString(buffer, stat.startTimeUTC.UTC().Format(nestedTraceTimeFormat))
	})
	writeField("DurationInMs", func() {
		writeNumber(buffer, float64(stat.duration)/float64(time.Millisecond))
	})
	writeField("RequestUri", func() {
		writeJSONString(buffer, stat.requestURI)
	})
	writeField("ResourceType", func() {
		writeJSONString(buffer, resourceTypeName(stat.resourceType))
	})
	writeField("HttpMethod", func() {
		writeJSONString(buffer, stat.httpMethod)
	})
	writeField("ActivityId", func() {
		if stat.activityID == "" {
			buffer.WriteString("null")
		} else {
			writeJSONString(buffer, stat.activityID)
		}
	})

	if stat.exceptionType != "" {
		writeField("ExceptionType", func() {
			writeJSONString(buffer, stat.exceptionType)
		})
		writeField("ExceptionMessage", func() {
			writeJSONString(buffer, stat.exceptionMessage)
		})
	}

	if stat.statusCodeText != "" {
		writeField("StatusCode", func() {
			writeJSONString(buffer, stat.statusCodeText)
		})
		if stat.statusCode >= http.StatusBadRequest && stat.reasonPhrase != "" {
			writeField("ReasonPhrase", func() {
				writeJSONString(buffer, stat.reasonPhrase)
			})
		}
	}

	buffer.WriteByte('}')
}

func writePointOperationStatistics(buffer *bytes.Buffer, stat pointOperationStatisticsTraceDatum) {
	buffer.WriteByte('{')
	firstField := true

	writeField := func(name string, writeValue func()) {
		if !firstField {
			buffer.WriteByte(',')
		}
		firstField = false
		writeJSONString(buffer, name)
		buffer.WriteByte(':')
		writeValue()
	}

	writeField("Id", func() {
		writeJSONString(buffer, "PointOperationStatistics")
	})
	writeField("ActivityId", func() {
		if stat.ActivityID == "" {
			buffer.WriteString("null")
		} else {
			writeJSONString(buffer, stat.ActivityID)
		}
	})
	writeField("ResponseTimeUtc", func() {
		writeJSONString(buffer, stat.ResponseTimeUTC.UTC().Format(nestedTraceTimeFormat))
	})
	writeField("StatusCode", func() {
		writeNumber(buffer, float64(stat.StatusCode))
	})
	writeField("SubStatusCode", func() {
		writeNumber(buffer, float64(stat.SubStatusCode))
	})
	writeField("RequestCharge", func() {
		writeNumber(buffer, float64(stat.RequestCharge))
	})
	writeField("RequestUri", func() {
		if stat.RequestURI == "" {
			buffer.WriteString("null")
		} else {
			writeJSONString(buffer, stat.RequestURI)
		}
	})
	writeField("ErrorMessage", func() {
		if stat.ErrorMessage == "" {
			buffer.WriteString("null")
		} else {
			writeJSONString(buffer, stat.ErrorMessage)
		}
	})
	writeField("RequestSessionToken", func() {
		if stat.RequestSessionToken == "" {
			buffer.WriteString("null")
		} else {
			writeJSONString(buffer, stat.RequestSessionToken)
		}
	})
	writeField("ResponseSessionToken", func() {
		if stat.ResponseSessionToken == "" {
			buffer.WriteString("null")
		} else {
			writeJSONString(buffer, stat.ResponseSessionToken)
		}
	})
	writeField("BELatencyInMs", func() {
		if stat.BELatencyInMs == "" {
			buffer.WriteString("null")
		} else {
			writeJSONString(buffer, stat.BELatencyInMs)
		}
	})

	buffer.WriteByte('}')
}

func writeJSONString(buffer *bytes.Buffer, value string) {
	encoded, _ := json.Marshal(value)
	buffer.Write(encoded)
}

func writeNumber(buffer *bytes.Buffer, value float64) {
	buffer.WriteString(strconv.FormatFloat(value, 'f', -1, 64))
}

func durationInMilliseconds(start time.Time, end *time.Time) float64 {
	if end != nil {
		return float64(end.Sub(start)) / float64(time.Millisecond)
	}

	return float64(time.Since(start)) / float64(time.Millisecond)
}
