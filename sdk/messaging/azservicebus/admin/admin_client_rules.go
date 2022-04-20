// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
)

// Filters allows you to specify a single filter for a particular rule.
// Only one of these properties can be set.
type Filters struct {
	// SQLFilter lets you filter messages using SQL syntax.
	SQLFilter *SQLFilter

	// CorrelationFilter lets you filter based on user or system properties of a message.
	CorrelationFilter *CorrelationFilter

	// TrueFilter evalutes to true for every message.
	TrueFilter *TrueFilter

	// FalseFilter evalutes to false for every message.
	FalseFilter *FalseFilter
}

// Rule specifies a message filter and action for a subscription.
type Rule struct {
	// Filter is the filter that will be used for Rule.
	// NOTE: only one of these filters should be set.
	Filter Filters

	// Action is the action that will be used for Rule.
	// Valid types: *SQLAction
	Action any
}

// SQLAction is an action that updates a message according to its
// expression.
type SQLAction struct {
	// Expression is a SQL Expression
	Expression string

	// Parameters is a map of string to values of type string, number, or boolean.
	Parameters map[string]interface{}
}

// SQLFilter is a filter that evaluates to true for any message that matches
// its expression.
type SQLFilter struct {
	// Expression is a SQL Expression
	Expression string

	// Parameters is a map of string to values of type string, number, or boolean.
	Parameters map[string]interface{}
}

// TrueFilter is a filter that always evaluates to true for any message.
type TrueFilter struct{}

// FalseFilter is a filter that always evaluates to false for any message.
type FalseFilter struct{}

// CorrelationFilter represents a set of conditions that are matched against user
// and system properties of messages for a subscription.
type CorrelationFilter struct {
	// ApplicationProperties will be matched against the application properties for the message.
	ApplicationProperties map[string]interface{}

	// ContentType will be matched against the ContentType property for the message.
	ContentType *string

	// CorrelationID will be matched against the CorrelationID property for the message.
	CorrelationID *string

	// MessageID will be matched against the MessageID property for the message.
	MessageID *string

	// ReplyTo will be matched against the ReplyTo property for the message.
	ReplyTo *string

	// ReplyToSessionID will be matched against the ReplyToSessionID property for the message.
	ReplyToSessionID *string

	// SessionID will be matched against the SessionID property for the message.
	SessionID *string

	// Subject will be matched against the Subject property for the message.
	Subject *string

	// To will be matched against the To property for the message.
	To *string
}

// RuleProperties are the properties for a rule.
type RuleProperties struct {
	// Name is the name of this rule.
	Name string

	// Filter is the filter for this rule.
	Filter *Filters

	// Action is the action for this rule.
	// Valid types: *SQLAction
	Action any
}

// CreateRuleResponse contains the response fields for Client.CreateRule
type CreateRuleResponse struct {
	RuleProperties
}

// CreateRuleOptions contains the optional parameters for Client.CreateRule
type CreateRuleOptions struct {
	// Name is the name of the rule or nil, which will default to $Default
	Name *string

	// Filter is the filter for this rule
	Filter *Filters

	// Action is the action for this rule
	// Valid types: *SQLAction
	Action any
}

// CreateRule creates a rule that can filter and update message for a subscription.
func (ac *Client) CreateRule(ctx context.Context, topicName string, subscriptionName string, options *CreateRuleOptions) (CreateRuleResponse, error) {
	ruleName := ""

	if options != nil && options.Name != nil {
		ruleName = *options.Name
	}

	resp, _, err := ac.createOrUpdateRule(ctx, topicName, subscriptionName, RuleProperties{
		Name:   ruleName,
		Filter: options.Filter,
		Action: options.Action,
	}, true)

	if err != nil {
		return CreateRuleResponse{}, err
	}

	return CreateRuleResponse{RuleProperties: *resp}, nil
}

// GetRuleResponse contains the response fields for Client.GetRule
type GetRuleResponse struct {
	// RuleProperties for the rule.
	RuleProperties
}

// GetRuleOptions contains the optional parameters for Client.GetRule
type GetRuleOptions struct {
	// For future expansion
}

// GetRule gets a rule for a subscription.
func (ac *Client) GetRule(ctx context.Context, topicName string, subscriptionName string, ruleName string, options *GetRuleOptions) (*GetRuleResponse, error) {
	var ruleEnv *atom.RuleEnvelope

	_, err := ac.em.Get(ctx, fmt.Sprintf("/%s/Subscriptions/%s/Rules/%s", topicName, subscriptionName, ruleName), &ruleEnv)

	if err != nil {
		return mapATOMError[GetRuleResponse](err)
	}

	props, err := newRuleProperties(ruleEnv)

	if err != nil {
		return nil, err
	}

	return &GetRuleResponse{
		RuleProperties: *props,
	}, nil
}

// ListRulesResponse contains the response fields for the pager returned from Client.ListRules.
type ListRulesResponse struct {
	// Rules are all the rules for the page.
	Rules []RuleProperties
}

// ListRulesOptions contains the optional parameters for Client.ListRules
type ListRulesOptions struct {
	// MaxPageSize is the maximum size of each page of results.
	MaxPageSize int32
}

// NewListRulesPager creates a pager that can list rules for a subscription.
func (ac *Client) NewListRulesPager(topicName string, subscriptionName string, options *ListRulesOptions) *runtime.Pager[ListRulesResponse] {
	var pageSize int32

	if options != nil {
		pageSize = options.MaxPageSize
	}

	ep := &entityPager[atom.RuleFeed, atom.RuleEnvelope, RuleProperties]{
		convertFn:    newRuleProperties,
		baseFragment: fmt.Sprintf("/%s/Subscriptions/%s/Rules/", topicName, subscriptionName),
		maxPageSize:  pageSize,
		em:           ac.em,
	}

	return runtime.NewPager(runtime.PageProcessor[ListRulesResponse]{
		More: func(ltr ListRulesResponse) bool {
			return ep.More()
		},
		Fetcher: func(ctx context.Context, t *ListRulesResponse) (ListRulesResponse, error) {
			items, err := ep.Fetcher(ctx)

			if err != nil {
				return ListRulesResponse{}, err
			}

			return ListRulesResponse{
				Rules: items,
			}, nil
		},
	})
}

// UpdateRuleResponse contains the response fields for Client.UpdateRule
type UpdateRuleResponse struct {
	// RuleProperties for the updated rule.
	RuleProperties
}

// UpdateRuleOptions can be used to configure the UpdateRule method.
type UpdateRuleOptions struct {
	// For future expansion
}

// UpdateRule updates a rule for a subscription.
func (ac *Client) UpdateRule(ctx context.Context, topicName string, subscriptionName string, properties RuleProperties) (UpdateRuleResponse, error) {
	resp, _, err := ac.createOrUpdateRule(ctx, topicName, subscriptionName, properties, false)

	if err != nil {
		return UpdateRuleResponse{}, err
	}

	return UpdateRuleResponse{RuleProperties: *resp}, nil
}

// DeleteRuleResponse contains the response fields for Client.DeleteRule
type DeleteRuleResponse struct {
	// For future expansion
}

// DeleteRuleOptions can be used to configure the Client.DeleteRule method.
type DeleteRuleOptions struct {
	// For future expansion
}

// DeleteRule deletes a rule for a subscription.
func (ac *Client) DeleteRule(ctx context.Context, topicName string, subscriptionName string, ruleName string, options *DeleteRuleOptions) (DeleteRuleResponse, error) {
	_, err := ac.em.Delete(ctx, fmt.Sprintf("/%s/Subscriptions/%s/Rules/%s", topicName, subscriptionName, ruleName))

	return DeleteRuleResponse{}, err
}

func (ac *Client) createOrUpdateRule(ctx context.Context, topicName string, subscriptionName string, putProps RuleProperties, creating bool) (*RuleProperties, *http.Response, error) {
	ruleDesc := atom.RuleDescription{}

	ourFilter := &atom.FilterDescription{}
	ruleDesc.Filter = ourFilter
	theirFilter := putProps.Filter

	if theirFilter != nil {
		if theirFilter.FalseFilter != nil {
			ourFilter.Type = "FalseFilter"
			ourFilter.SQLExpression = to.Ptr("1=0")
		}

		if theirFilter.TrueFilter != nil {
			ourFilter.Type = "TrueFilter"
			ourFilter.SQLExpression = to.Ptr("1=1")
		}

		if theirFilter.SQLFilter != nil {
			params, err := publicSQLParametersToInternal(theirFilter.SQLFilter.Parameters)

			if err != nil {
				return nil, nil, err
			}

			ourFilter.Type = "SqlFilter"
			ourFilter.SQLExpression = &theirFilter.SQLFilter.Expression
			ourFilter.Parameters = params
		}

		if theirFilter.CorrelationFilter != nil {
			ourFilter.Type = "CorrelationFilter"

			ourFilter.CorrelationFilter.ContentType = theirFilter.CorrelationFilter.ContentType
			ourFilter.CorrelationFilter.CorrelationID = theirFilter.CorrelationFilter.CorrelationID
			ourFilter.CorrelationFilter.MessageID = theirFilter.CorrelationFilter.MessageID
			ourFilter.CorrelationFilter.ReplyTo = theirFilter.CorrelationFilter.ReplyTo
			ourFilter.CorrelationFilter.ReplyToSessionID = theirFilter.CorrelationFilter.ReplyToSessionID
			ourFilter.CorrelationFilter.SessionID = theirFilter.CorrelationFilter.SessionID
			ourFilter.CorrelationFilter.Label = theirFilter.CorrelationFilter.Subject
			ourFilter.CorrelationFilter.To = theirFilter.CorrelationFilter.To

			appProps, err := publicSQLParametersToInternal(theirFilter.CorrelationFilter.ApplicationProperties)

			if err != nil {
				return nil, nil, err
			}

			ourFilter.CorrelationFilter.Properties = appProps
		}
	} else {
		ourFilter.Type = "TrueFilter"
		ourFilter.SQLExpression = to.Ptr("1=1")
	}

	theirAction := putProps.Action

	if theirAction != nil {
		switch actualAction := theirAction.(type) {
		case *SQLAction:
			ourAction := &atom.ActionDescription{
				Type: "SqlRuleAction",
			}
			ruleDesc.Action = ourAction

			params, err := publicSQLParametersToInternal(actualAction.Parameters)

			if err != nil {
				return nil, nil, err
			}

			ourAction.SQLExpression = actualAction.Expression
			ourAction.Parameters = params
		default:
			return nil, nil, fmt.Errorf("unknown action type %T", theirAction)
		}
	}

	ruleDesc.Name = "$Default"

	if putProps.Name != "" {
		ruleDesc.Name = putProps.Name
	}

	var mw []atom.MiddlewareFunc

	if !creating {
		// an update requires the entity to already exist.
		mw = append(mw, func(next atom.RestHandler) atom.RestHandler {
			return func(ctx context.Context, req *http.Request) (*http.Response, error) {
				req.Header.Set("If-Match", "*")
				return next(ctx, req)
			}
		})
	}

	putEnv := atom.WrapWithRuleEnvelope(&ruleDesc)

	var respEnv *atom.RuleEnvelope

	httpResp, err := ac.em.Put(ctx, fmt.Sprintf("/%s/Subscriptions/%s/Rules/%s", topicName, subscriptionName, putProps.Name), putEnv, &respEnv, mw...)

	if err != nil {
		return nil, nil, err
	}

	respProps, err := newRuleProperties(respEnv)

	return respProps, httpResp, err
}

func newRuleProperties(env *atom.RuleEnvelope) (*RuleProperties, error) {
	desc := env.Content.RuleDescription

	props := RuleProperties{
		Name:   env.Title,
		Filter: &Filters{},
	}

	switch desc.Filter.Type {
	case "TrueFilter":
		props.Filter.TrueFilter = &TrueFilter{}
	case "FalseFilter":
		props.Filter.FalseFilter = &FalseFilter{}
	case "CorrelationFilter":
		cf := desc.Filter.CorrelationFilter

		appProps, err := internalSQLParametersToPublic(cf.Properties)

		if err != nil {
			return nil, err
		}

		props.Filter.CorrelationFilter = &CorrelationFilter{
			ContentType:           cf.ContentType,
			CorrelationID:         cf.CorrelationID,
			MessageID:             cf.MessageID,
			ReplyTo:               cf.ReplyTo,
			ReplyToSessionID:      cf.ReplyToSessionID,
			SessionID:             cf.SessionID,
			Subject:               cf.Label,
			To:                    cf.To,
			ApplicationProperties: appProps,
		}
	case "SqlFilter":
		params, err := internalSQLParametersToPublic(desc.Filter.Parameters)

		if err != nil {
			return nil, err
		}

		props.Filter.SQLFilter = &SQLFilter{
			Expression: *desc.Filter.SQLExpression,
			Parameters: params,
		}
	default:
		return nil, fmt.Errorf("filter for rule %s, with type %s, is not handled", env.Title, desc.Filter.Type)
	}

	switch desc.Action.Type {
	case "EmptyRuleAction":
	case "SqlRuleAction":
		params, err := internalSQLParametersToPublic(desc.Action.Parameters)

		if err != nil {
			return nil, err
		}

		props.Action = &SQLAction{
			Expression: desc.Action.SQLExpression,
			Parameters: params,
		}
	default:
		return nil, fmt.Errorf("action for rule %s, with type %s, is not handled", env.Title, desc.Action.Type)
	}

	return &props, nil
}

func publicSQLParametersToInternal(publicParams map[string]interface{}) ([]atom.KeyValueOfstringanyType, error) {
	var params []atom.KeyValueOfstringanyType

	for k, v := range publicParams {
		switch asType := v.(type) {
		case string:
			params = append(params, atom.KeyValueOfstringanyType{
				Key: k,
				Value: atom.Value{
					Type:  "l28:string",
					L28NS: "http://www.w3.org/2001/XMLSchema",
					Text:  asType,
				},
			})
		case bool:
			params = append(params, atom.KeyValueOfstringanyType{
				Key: k,
				Value: atom.Value{
					Type:  "l28:boolean",
					L28NS: "http://www.w3.org/2001/XMLSchema",
					Text:  fmt.Sprintf("%t", v),
				},
			})
		case int, int64, int32:
			params = append(params, atom.KeyValueOfstringanyType{
				Key: k,
				Value: atom.Value{
					Type:  "l28:int",
					L28NS: "http://www.w3.org/2001/XMLSchema",
					Text:  fmt.Sprintf("%d", v),
				},
			})
		case float32, float64:
			params = append(params, atom.KeyValueOfstringanyType{
				Key: k,
				Value: atom.Value{
					Type:  "l28:double",
					L28NS: "http://www.w3.org/2001/XMLSchema",
					Text:  fmt.Sprintf("%f", v),
				},
			})
		case time.Time:
			params = append(params, atom.KeyValueOfstringanyType{
				Key: k,
				Value: atom.Value{
					Type:  "l28:dateTime",
					L28NS: "http://www.w3.org/2001/XMLSchema",
					Text:  asType.UTC().Format(time.RFC3339Nano),
				},
			})
		default:
			// TODO: 'duration'
			return nil, fmt.Errorf("type %T of parameter %s is not a handled type for SQL parameters", v, k)
		}
	}

	return params, nil
}

func internalSQLParametersToPublic(internalParams []atom.KeyValueOfstringanyType) (map[string]interface{}, error) {
	params := map[string]interface{}{}

	for _, p := range internalParams {
		switch p.Value.Type {
		case "d6p1:string":
			params[p.Key] = p.Value.Text
		case "d6p1:boolean":
			val, err := strconv.ParseBool(p.Value.Text)

			if err != nil {
				return nil, err
			}

			params[p.Key] = val
		case "d6p1:int":
			val, err := strconv.ParseInt(p.Value.Text, 10, 64)

			if err != nil {
				return nil, err
			}

			params[p.Key] = val
		case "d6p1:double":
			val, err := strconv.ParseFloat(p.Value.Text, 64)

			if err != nil {
				return nil, err
			}

			params[p.Key] = val
		case "d6p1:dateTime":
			val, err := time.Parse(time.RFC3339Nano, p.Value.Text)

			if err != nil {
				return nil, err
			}

			params[p.Key] = val.UTC()
		default:
			// TODO: timespan
			return nil, fmt.Errorf("type %s of parameter %s is not a handled type for SQL parameters", p.Value.Type, p.Key)
		}
	}

	if len(params) == 0 {
		return nil, nil
	}

	return params, nil
}
