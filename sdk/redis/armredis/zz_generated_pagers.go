// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armredis

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"reflect"
)

type FirewallRulesListPager interface {
	azcore.Pager
	// PageResponse returns the current FirewallRulesListResponse.
	PageResponse() FirewallRulesListResponse
}

type firewallRulesListPager struct {
	client    *FirewallRulesClient
	current   FirewallRulesListResponse
	err       error
	requester func(context.Context) (*azcore.Request, error)
	advancer  func(context.Context, FirewallRulesListResponse) (*azcore.Request, error)
}

func (p *firewallRulesListPager) Err() error {
	return p.err
}

func (p *firewallRulesListPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.RedisFirewallRuleListResult.NextLink == nil || len(*p.current.RedisFirewallRuleListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(http.StatusOK) {
		p.err = p.client.listHandleError(resp)
		return false
	}
	result, err := p.client.listHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *firewallRulesListPager) PageResponse() FirewallRulesListResponse {
	return p.current
}

type LinkedServerListPager interface {
	azcore.Pager
	// PageResponse returns the current LinkedServerListResponse.
	PageResponse() LinkedServerListResponse
}

type linkedServerListPager struct {
	client    *LinkedServerClient
	current   LinkedServerListResponse
	err       error
	requester func(context.Context) (*azcore.Request, error)
	advancer  func(context.Context, LinkedServerListResponse) (*azcore.Request, error)
}

func (p *linkedServerListPager) Err() error {
	return p.err
}

func (p *linkedServerListPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.RedisLinkedServerWithPropertiesList.NextLink == nil || len(*p.current.RedisLinkedServerWithPropertiesList.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(http.StatusOK) {
		p.err = p.client.listHandleError(resp)
		return false
	}
	result, err := p.client.listHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *linkedServerListPager) PageResponse() LinkedServerListResponse {
	return p.current
}

type OperationsListPager interface {
	azcore.Pager
	// PageResponse returns the current OperationsListResponse.
	PageResponse() OperationsListResponse
}

type operationsListPager struct {
	client    *OperationsClient
	current   OperationsListResponse
	err       error
	requester func(context.Context) (*azcore.Request, error)
	advancer  func(context.Context, OperationsListResponse) (*azcore.Request, error)
}

func (p *operationsListPager) Err() error {
	return p.err
}

func (p *operationsListPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.OperationListResult.NextLink == nil || len(*p.current.OperationListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(http.StatusOK) {
		p.err = p.client.listHandleError(resp)
		return false
	}
	result, err := p.client.listHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *operationsListPager) PageResponse() OperationsListResponse {
	return p.current
}

type PatchSchedulesListByRedisResourcePager interface {
	azcore.Pager
	// PageResponse returns the current PatchSchedulesListByRedisResourceResponse.
	PageResponse() PatchSchedulesListByRedisResourceResponse
}

type patchSchedulesListByRedisResourcePager struct {
	client    *PatchSchedulesClient
	current   PatchSchedulesListByRedisResourceResponse
	err       error
	requester func(context.Context) (*azcore.Request, error)
	advancer  func(context.Context, PatchSchedulesListByRedisResourceResponse) (*azcore.Request, error)
}

func (p *patchSchedulesListByRedisResourcePager) Err() error {
	return p.err
}

func (p *patchSchedulesListByRedisResourcePager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.RedisPatchScheduleListResult.NextLink == nil || len(*p.current.RedisPatchScheduleListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(http.StatusOK) {
		p.err = p.client.listByRedisResourceHandleError(resp)
		return false
	}
	result, err := p.client.listByRedisResourceHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *patchSchedulesListByRedisResourcePager) PageResponse() PatchSchedulesListByRedisResourceResponse {
	return p.current
}

type RedisListByResourceGroupPager interface {
	azcore.Pager
	// PageResponse returns the current RedisListByResourceGroupResponse.
	PageResponse() RedisListByResourceGroupResponse
}

type redisListByResourceGroupPager struct {
	client    *RedisClient
	current   RedisListByResourceGroupResponse
	err       error
	requester func(context.Context) (*azcore.Request, error)
	advancer  func(context.Context, RedisListByResourceGroupResponse) (*azcore.Request, error)
}

func (p *redisListByResourceGroupPager) Err() error {
	return p.err
}

func (p *redisListByResourceGroupPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.RedisListResult.NextLink == nil || len(*p.current.RedisListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(http.StatusOK) {
		p.err = p.client.listByResourceGroupHandleError(resp)
		return false
	}
	result, err := p.client.listByResourceGroupHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *redisListByResourceGroupPager) PageResponse() RedisListByResourceGroupResponse {
	return p.current
}

type RedisListBySubscriptionPager interface {
	azcore.Pager
	// PageResponse returns the current RedisListBySubscriptionResponse.
	PageResponse() RedisListBySubscriptionResponse
}

type redisListBySubscriptionPager struct {
	client    *RedisClient
	current   RedisListBySubscriptionResponse
	err       error
	requester func(context.Context) (*azcore.Request, error)
	advancer  func(context.Context, RedisListBySubscriptionResponse) (*azcore.Request, error)
}

func (p *redisListBySubscriptionPager) Err() error {
	return p.err
}

func (p *redisListBySubscriptionPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.RedisListResult.NextLink == nil || len(*p.current.RedisListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(http.StatusOK) {
		p.err = p.client.listBySubscriptionHandleError(resp)
		return false
	}
	result, err := p.client.listBySubscriptionHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *redisListBySubscriptionPager) PageResponse() RedisListBySubscriptionResponse {
	return p.current
}

type RedisListUpgradeNotificationsPager interface {
	azcore.Pager
	// PageResponse returns the current RedisListUpgradeNotificationsResponse.
	PageResponse() RedisListUpgradeNotificationsResponse
}

type redisListUpgradeNotificationsPager struct {
	client    *RedisClient
	current   RedisListUpgradeNotificationsResponse
	err       error
	requester func(context.Context) (*azcore.Request, error)
	advancer  func(context.Context, RedisListUpgradeNotificationsResponse) (*azcore.Request, error)
}

func (p *redisListUpgradeNotificationsPager) Err() error {
	return p.err
}

func (p *redisListUpgradeNotificationsPager) NextPage(ctx context.Context) bool {
	var req *azcore.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.NotificationListResponse.NextLink == nil || len(*p.current.NotificationListResponse.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(http.StatusOK) {
		p.err = p.client.listUpgradeNotificationsHandleError(resp)
		return false
	}
	result, err := p.client.listUpgradeNotificationsHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *redisListUpgradeNotificationsPager) PageResponse() RedisListUpgradeNotificationsResponse {
	return p.current
}
