// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package fake

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"net/http"
	"sync"
)

type result struct {
	resp *http.Response
	err  error
}

type nonRetriableError struct {
	error
}

func (nonRetriableError) NonRetriable() {
	// marker method
}

func contains[T comparable](s []T, v T) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

func parseWithCast[T any](v string, parse func(v string) (T, error)) (T, error) {
	t, err := parse(v)
	if err != nil {
		return *new(T), err
	}
	return t, err
}

func newTracker[T any]() *tracker[T] {
	return &tracker[T]{
		items: map[string]*T{},
	}
}

type tracker[T any] struct {
	items map[string]*T
	mu    sync.Mutex
}

func (p *tracker[T]) get(req *http.Request) *T {
	p.mu.Lock()
	defer p.mu.Unlock()
	if item, ok := p.items[server.SanitizePagerPollerPath(req.URL.Path)]; ok {
		return item
	}
	return nil
}

func (p *tracker[T]) add(req *http.Request, item *T) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.items[server.SanitizePagerPollerPath(req.URL.Path)] = item
}

func (p *tracker[T]) remove(req *http.Request) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.items, server.SanitizePagerPollerPath(req.URL.Path))
}
