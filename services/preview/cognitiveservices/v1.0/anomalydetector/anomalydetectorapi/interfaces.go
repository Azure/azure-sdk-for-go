package anomalydetectorapi

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/preview/cognitiveservices/v1.0/anomalydetector"
	"github.com/Azure/go-autorest/autorest"
)

// BaseClientAPI contains the set of methods on the BaseClient type.
type BaseClientAPI interface {
	ChangePointDetect(ctx context.Context, body anomalydetector.ChangePointDetectRequest) (result anomalydetector.ChangePointDetectResponse, err error)
	EntireDetect(ctx context.Context, body anomalydetector.Request) (result anomalydetector.EntireDetectResponse, err error)
	LastDetect(ctx context.Context, body anomalydetector.Request) (result anomalydetector.LastDetectResponse, err error)
}

var _ BaseClientAPI = (*anomalydetector.BaseClient)(nil)

// TimeSeriesClientAPI contains the set of methods on the TimeSeriesClient type.
type TimeSeriesClientAPI interface {
	ChangePointDetectInTimeRange(ctx context.Context, timeSeriesID string, body anomalydetector.ChangePointDetectInTimeRangeRequest) (result anomalydetector.ChangePointDetectInTimeRangeResponse, err error)
	Create(ctx context.Context, timeSeriesID string, body anomalydetector.TimeSeriesCreateRequest) (result autorest.Response, err error)
	Delete(ctx context.Context, timeSeriesID string) (result autorest.Response, err error)
	EntireDetectInTimeRange(ctx context.Context, timeSeriesID string, body anomalydetector.AnomalyDetectInTimeRangeRequest) (result anomalydetector.AnomalyDetectInTimeRangeResponse, err error)
	Get(ctx context.Context, timeSeriesID string) (result anomalydetector.TimeSeries, err error)
	Label(ctx context.Context, timeSeriesID string, body anomalydetector.LabelRequest) (result autorest.Response, err error)
	LastDetectInTimeRange(ctx context.Context, timeSeriesID string, body anomalydetector.AnomalyDetectInTimeRangeRequest) (result anomalydetector.AnomalyDetectInTimeRangeResponse, err error)
	List(ctx context.Context, next string) (result anomalydetector.TimeSeriesList, err error)
	ListGroups(ctx context.Context, timeSeriesID string, next string) (result anomalydetector.TimeSeriesGroupList, err error)
	Query(ctx context.Context, timeSeriesID string, body anomalydetector.TimeSeriesQueryRequest) (result anomalydetector.TimeSeriesQueryResponse, err error)
	Write(ctx context.Context, timeSeriesID string, body []anomalydetector.Point) (result autorest.Response, err error)
}

var _ TimeSeriesClientAPI = (*anomalydetector.TimeSeriesClient)(nil)

// TimeSeriesGroupClientAPI contains the set of methods on the TimeSeriesGroupClient type.
type TimeSeriesGroupClientAPI interface {
	AddTimeSeries(ctx context.Context, timeSeriesGroupID string, timeSeriesID string) (result autorest.Response, err error)
	Create(ctx context.Context, timeSeriesGroupID string, body anomalydetector.TimeSeriesGroupCreateRequest) (result autorest.Response, err error)
	Delete(ctx context.Context, timeSeriesGroupID string) (result autorest.Response, err error)
	DeleteTimeSeries(ctx context.Context, timeSeriesGroupID string, timeSeriesID string) (result autorest.Response, err error)
	Get(ctx context.Context, timeSeriesGroupID string) (result anomalydetector.TimeSeriesGroup, err error)
	InconsistencyDetect(ctx context.Context, timeSeriesGroupID string, body anomalydetector.InconsistencyDetectRequest) (result anomalydetector.Inconsistency, err error)
	InconsistencyQuery(ctx context.Context, timeSeriesGroupID string, body anomalydetector.InconsistencyQueryRequest) (result anomalydetector.ListInconsistency, err error)
	List(ctx context.Context, next string) (result anomalydetector.TimeSeriesGroupList, err error)
	ListSeries(ctx context.Context, timeSeriesGroupID string, next string) (result anomalydetector.TimeSeriesList, err error)
}

var _ TimeSeriesGroupClientAPI = (*anomalydetector.TimeSeriesGroupClient)(nil)
