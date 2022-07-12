//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

type PageResponse struct {
	Values   []int `json:"values"`
	NextPage bool  `json:"next"`
}

func pageResponseFetcher(ctx context.Context, pl Pipeline, endpoint string) (PageResponse, error) {
	req, err := NewRequest(ctx, http.MethodGet, endpoint)
	if err != nil {
		return PageResponse{}, err
	}
	resp, err := pl.Do(req)
	if err != nil {
		return PageResponse{}, err
	}
	if !HasStatusCode(resp, http.StatusOK) {
		return PageResponse{}, NewResponseError(resp)
	}
	pr := PageResponse{}
	if err := UnmarshalAsJSON(resp, &pr); err != nil {
		return PageResponse{}, err
	}
	return pr, nil
}

func TestPagerSinglePage(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"values": [1, 2, 3, 4, 5]}`)))
	pl := exported.NewPipeline(srv)

	pager := NewPager(PagingHandler[PageResponse]{
		More: func(current PageResponse) bool {
			return current.NextPage
		},
		Fetcher: func(ctx context.Context, current *PageResponse) (PageResponse, error) {
			return pageResponseFetcher(ctx, pl, srv.URL())
		},
	})
	require.True(t, pager.firstPage)

	pageCount := 0
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		require.Equal(t, []int{1, 2, 3, 4, 5}, page.Values)
		require.Empty(t, page.NextPage)
		pageCount++
	}
	require.Equal(t, 1, pageCount)
	page, err := pager.NextPage(context.Background())
	require.Error(t, err)
	require.Empty(t, page)
}

func TestPagerMultiplePages(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"values": [1, 2, 3, 4, 5], "next": true}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"values": [6, 7, 8], "next": true}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"values": [9, 0, 1, 2]}`)))
	pl := exported.NewPipeline(srv)

	pageCount := 0
	pager := NewPager(PagingHandler[PageResponse]{
		More: func(current PageResponse) bool {
			return current.NextPage
		},
		Fetcher: func(ctx context.Context, current *PageResponse) (PageResponse, error) {
			if pageCount == 1 {
				require.Nil(t, current)
			} else {
				require.NotNil(t, current)
			}
			return pageResponseFetcher(ctx, pl, srv.URL())
		},
	})
	require.True(t, pager.firstPage)

	for pager.More() {
		pageCount++
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		switch pageCount {
		case 1:
			require.Equal(t, []int{1, 2, 3, 4, 5}, page.Values)
			require.True(t, page.NextPage)
		case 2:
			require.Equal(t, []int{6, 7, 8}, page.Values)
			require.True(t, page.NextPage)
		case 3:
			require.Equal(t, []int{9, 0, 1, 2}, page.Values)
			require.False(t, page.NextPage)
		}
	}
	require.Equal(t, 3, pageCount)
	page, err := pager.NextPage(context.Background())
	require.Error(t, err)
	require.Empty(t, page)
}

func TestPagerLROMultiplePages(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"values": [6, 7, 8]}`)))
	pl := exported.NewPipeline(srv)

	pager := NewPager(PagingHandler[PageResponse]{
		More: func(current PageResponse) bool {
			return current.NextPage
		},
		Fetcher: func(ctx context.Context, current *PageResponse) (PageResponse, error) {
			return pageResponseFetcher(ctx, pl, srv.URL())
		},
	})
	require.True(t, pager.firstPage)

	require.NoError(t, json.Unmarshal([]byte(`{"values": [1, 2, 3, 4, 5], "next": true}`), pager))

	pageCount := 0
	for pager.More() {
		pageCount++
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		switch pageCount {
		case 1:
			require.Equal(t, []int{1, 2, 3, 4, 5}, page.Values)
			require.True(t, page.NextPage)
		case 2:
			require.Equal(t, []int{6, 7, 8}, page.Values)
			require.False(t, page.NextPage)
		}
	}
	require.Equal(t, 2, pageCount)
	page, err := pager.NextPage(context.Background())
	require.Error(t, err)
	require.Empty(t, page)
}

func TestPagerFetcherError(t *testing.T) {
	pager := NewPager(PagingHandler[PageResponse]{
		More: func(current PageResponse) bool {
			return current.NextPage
		},
		Fetcher: func(ctx context.Context, current *PageResponse) (PageResponse, error) {
			return PageResponse{}, errors.New("fetcher failed")
		},
	})
	require.True(t, pager.firstPage)

	page, err := pager.NextPage(context.Background())
	require.Error(t, err)
	require.Empty(t, page)
}

func TestPagerPipelineError(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetError(errors.New("pipeline failed"))
	pl := exported.NewPipeline(srv)

	pager := NewPager(PagingHandler[PageResponse]{
		More: func(current PageResponse) bool {
			return current.NextPage
		},
		Fetcher: func(ctx context.Context, current *PageResponse) (PageResponse, error) {
			return pageResponseFetcher(ctx, pl, srv.URL())
		},
	})
	require.True(t, pager.firstPage)

	page, err := pager.NextPage(context.Background())
	require.Error(t, err)
	require.Empty(t, page)
}

func TestPagerSecondPageError(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"values": [1, 2, 3, 4, 5], "next": true}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest), mock.WithBody([]byte(`{"message": "didn't work", "code": "PageError"}`)))
	pl := exported.NewPipeline(srv)

	pageCount := 0
	pager := NewPager(PagingHandler[PageResponse]{
		More: func(current PageResponse) bool {
			return current.NextPage
		},
		Fetcher: func(ctx context.Context, current *PageResponse) (PageResponse, error) {
			if pageCount == 1 {
				require.Nil(t, current)
			} else {
				require.NotNil(t, current)
			}
			return pageResponseFetcher(ctx, pl, srv.URL())
		},
	})
	require.True(t, pager.firstPage)

	for pager.More() {
		pageCount++
		page, err := pager.NextPage(context.Background())
		switch pageCount {
		case 1:
			require.NoError(t, err)
			require.Equal(t, []int{1, 2, 3, 4, 5}, page.Values)
			require.True(t, page.NextPage)
		case 2:
			require.Error(t, err)
			var respErr *exported.ResponseError
			require.True(t, errors.As(err, &respErr))
			require.Equal(t, "PageError", respErr.ErrorCode)
			goto ExitLoop
		}
	}
ExitLoop:
	require.Equal(t, 2, pageCount)
}

func TestPagerResponderError(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`incorrect JSON response`)))
	pl := exported.NewPipeline(srv)

	pager := NewPager(PagingHandler[PageResponse]{
		More: func(current PageResponse) bool {
			return current.NextPage
		},
		Fetcher: func(ctx context.Context, current *PageResponse) (PageResponse, error) {
			return pageResponseFetcher(ctx, pl, srv.URL())
		},
	})
	require.True(t, pager.firstPage)

	page, err := pager.NextPage(context.Background())
	require.Error(t, err)
	require.Empty(t, page)
}
