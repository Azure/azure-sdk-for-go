package servicebus_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

func Example_helloWorld() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	// Create a client to communicate with the queue. (The queue must have already been created, see `QueueManager`)
	q, err := ns.NewQueue("helloworld")
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	err = q.Send(ctx, servicebus.NewMessageFromString("Hello, World!!!"))
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	err = q.ReceiveOne(
		ctx,
		servicebus.HandlerFunc(func(ctx context.Context, message *servicebus.Message) error {
			fmt.Println(string(message.Data))
			return message.Complete(ctx)
		}))
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	// Output: Hello, World!!!
}

type PrioritySubscription struct {
	Name         string
	Expression   string
	MessageCount int
}

type PriorityMessage struct {
	Body     string
	Priority int
}

func Example_auto_forward() {
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

	qm := ns.NewQueueManager()
	target, err := ensureQueue(ctx, qm, "AutoForwardTargetQueue")
	if err != nil {
		fmt.Println(err)
		return
	}

	source, err := ensureQueue(ctx, qm, "AutoForwardSourceQueue", servicebus.QueueEntityWithAutoForward(target))
	if err != nil {
		fmt.Println(err)
		return
	}

	sourceQueue, err := ns.NewQueue(source.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := sourceQueue.Send(ctx, servicebus.NewMessageFromString("forward me to target!")); err != nil {
		fmt.Println(err)
		return
	}

	targetQueue, err := ns.NewQueue(target.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := targetQueue.ReceiveOne(ctx, MessagePrinter{}); err != nil {
		fmt.Println(err)
		return
	}

	// Output: forward me to target!
}

func Example_deadletter_queues() {
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

	qm := ns.NewQueueManager()
	qe, err := ensureQueue(ctx, qm, "DeadletterExample")
	if err != nil {
		fmt.Println(err)
		return
	}

	q, err := ns.NewQueue(qe.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := q.Send(ctx, servicebus.NewMessageFromString("foo")); err != nil {
		fmt.Println(err)
		return
	}

	// Abandon the message 10 times simulating attempting to process the message 10 times. After the 10th time, the
	// message will be placed in the Deadletter Queue.
	for count := 0; count < 10; count++ {
		err = q.ReceiveOne(ctx, servicebus.HandlerFunc(func(ctx context.Context, msg *servicebus.Message) error {
			fmt.Printf("count: %d\n", count+1)
			return msg.Abandon(ctx)
		}))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// receive one from the queue's deadletter queue. It should be the foo message.
	qdl := q.NewDeadLetter()
	if err := qdl.ReceiveOne(ctx, MessagePrinter{}); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// count: 1
	// count: 2
	// count: 3
	// count: 4
	// count: 5
	// count: 6
	// count: 7
	// count: 8
	// count: 9
	// count: 10
	// foo
}

func Example_defer_messages() {
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

	qm := ns.NewQueueManager()
	qe, err := ensureQueue(ctx, qm, "DeferExample")
	if err != nil {
		fmt.Println(err)
		return
	}

	q, err := ns.NewQueue(qe.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	type recipeStep struct {
		Step  int    `json:"step,omitempty"`
		Title string `json:"title,omitempty"`
	}

	steps := []recipeStep{
		{
			Step:  1,
			Title: "Shop",
		},
		{
			Step:  2,
			Title: "Unpack",
		},
		{
			Step:  3,
			Title: "Prepare",
		},
		{
			Step:  4,
			Title: "Cook",
		},
		{
			Step:  5,
			Title: "Eat",
		},
	}

	for _, step := range steps {
		go func(s recipeStep) {
			j, err := json.Marshal(s)
			if err != nil {
				fmt.Println(err)
				return
			}

			msg := &servicebus.Message{
				Data: j,
				ContentType: "application/json",
				Label: "RecipeStep",
			}

			// the way we shuffle the message order is to introduce
			// a tiny random delay before each of the messages is sent
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			if err := q.Send(ctx, msg); err != nil {
				fmt.Println(err)
				return
			}
		}(step)
	}

	sequenceByStepNumber := map[int]int64{}
	// collect and defer messages
	for i := 0; i < len(steps); i++ {
		err = q.ReceiveOne(ctx, servicebus.HandlerFunc(func(ctx context.Context, msg *servicebus.Message) error {
			var step recipeStep
			if err := json.Unmarshal(msg.Data, &step); err != nil {
				return err
			}
			sequenceByStepNumber[step.Step] = *msg.SystemProperties.SequenceNumber
			return msg.Defer(ctx)
		}))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	for i := 0; i < len(steps); i++ {
		err := q.ReceiveDeferred(ctx, servicebus.HandlerFunc(func(ctx context.Context, msg *servicebus.Message) error {
			var step recipeStep
			if err := json.Unmarshal(msg.Data, &step); err != nil {
				return err
			}
			fmt.Printf("step: %d, %s\n", step.Step, step.Title)
			return msg.Complete(ctx)
		}), sequenceByStepNumber[i+1])
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Output:
	// step: 1, Shop
	// step: 2, Unpack
	// step: 3, Prepare
	// step: 4, Cook
	// step: 5, Eat
}

func Example_priority_subscriptions() {
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

type MessagePrinter struct{}

func (mp MessagePrinter) Handle(ctx context.Context, msg *servicebus.Message) error {
	fmt.Println(string(msg.Data))
	return msg.Complete(ctx)
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

func ensureQueue(ctx context.Context, qm *servicebus.QueueManager, name string, opts ...servicebus.QueueManagementOption) (*servicebus.QueueEntity, error) {
	qe, err := qm.Get(ctx, name)
	if err == nil {
		_ = qm.Delete(ctx, name)
	}

	qe, err = qm.Put(ctx, name, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return qe, nil
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
