package servicebus_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

type PrioritySubscription struct {
	Name         string
	Expression   string
	MessageCount int
}

type PriorityMessage struct {
	Body     string
	Priority int
}

type PriorityPrinter struct {
	SubName string
}

func (pp PriorityPrinter) Handle(ctx context.Context, msg *servicebus.Message) error {
	i, ok := msg.UserProperties["Priority"].(int64)
	if !ok {
		fmt.Println("Priority is not an int64")
	}

	fmt.Println(strings.Join([]string{pp.SubName, string(msg.Data), strconv.Itoa(int(i))}, "_"))
	return msg.Complete(ctx)
}

func Example_prioritySubscriptions() {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	if connStr == "" {
		fmt.Println("FATAL: expected environment variable SERVICEBUS_CONNECTION_STRING not set")
		return
	}

	// Create a client to communicate with a Service Bus Namespace.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println(err)
		return
	}

	// build the topic for sending priority messages
	tm := ns.NewTopicManager()
	topicEntity, err := ensureTopic(ctx, tm, "PrioritySubscriptionsTopic")
	if err != nil {
		fmt.Println(err)
		return
	}

	sm, err := ns.NewSubscriptionManager(topicEntity.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	// build each priority subscription providing each with a SQL like expression to filter messages from the topic
	prioritySubs := []PrioritySubscription{
		{
			Name:         "Priority1",
			Expression:   "user.Priority=1",
			MessageCount: 1,
		},
		{
			Name:         "Priority2",
			Expression:   "user.Priority=2",
			MessageCount: 1,
		},
		{
			Name:         "PriorityGreaterThan2",
			Expression:   "user.Priority>2",
			MessageCount: 2,
		},
	}
	for _, s := range prioritySubs {
		subEntity, err := ensureSubscription(ctx, sm, s.Name)
		if err != nil {
			fmt.Println(err)
			return
		}

		// remove the default rule, which is the "TrueFilter" that accepts all messages
		err = sm.DeleteRule(ctx, subEntity.Name, "$Default")
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = sm.PutRule(ctx, subEntity.Name, s.Name+"Rule", servicebus.SQLFilter{Expression: s.Expression})
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	priorityMessages := []PriorityMessage{
		{
			Body:     "foo",
			Priority: 1,
		},
		{
			Body:     "bar",
			Priority: 2,
		},
		{
			Body:     "bazz",
			Priority: 3,
		},
		{
			Body:     "buzz",
			Priority: 4,
		},
	}
	topic, err := ns.NewTopic(topicEntity.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = topic.Close(ctx)
	}()

	for _, pMessage := range priorityMessages {
		msg := servicebus.NewMessageFromString(pMessage.Body)
		msg.UserProperties = map[string]interface{}{"Priority": pMessage.Priority}
		if err := topic.Send(ctx, msg); err != nil {
			fmt.Println(err)
			return
		}
	}

	for _, s := range prioritySubs {
		sub, err := topic.NewSubscription(s.Name)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i := 0; i < s.MessageCount; i++ {
			err := sub.ReceiveOne(ctx, PriorityPrinter{SubName: sub.Name})
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		err = sub.Close(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Output:
	// Priority1_foo_1
	// Priority2_bar_2
	// PriorityGreaterThan2_bazz_3
	// PriorityGreaterThan2_buzz_4
}

func ensureTopic(ctx context.Context, tm *servicebus.TopicManager, name string, opts ...servicebus.TopicManagementOption) (*servicebus.TopicEntity, error) {
	te, err := tm.Get(ctx, name)
	if err == nil {
		_ = tm.Delete(ctx, name)
	}

	te, err = tm.Put(ctx, name, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return te, nil
}

func ensureSubscription(ctx context.Context, sm *servicebus.SubscriptionManager, name string, opts ...servicebus.SubscriptionManagementOption) (*servicebus.SubscriptionEntity, error) {
	subEntity, err := sm.Get(ctx, name)
	if err == nil {
		_ = sm.Delete(ctx, name)
	}

	subEntity, err = sm.Put(ctx, name, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return subEntity, nil
}
