package servicebus_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

type RecipeStep struct {
	Step  int    `json:"step,omitempty"`
	Title string `json:"title,omitempty"`
}

func Example_deferMessages() {
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
	defer func() {
		_ = q.Close(ctx)
	}()

	steps := []RecipeStep{
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
		go func(s RecipeStep) {
			j, err := json.Marshal(s)
			if err != nil {
				fmt.Println(err)
				return
			}

			msg := &servicebus.Message{
				Data:        j,
				ContentType: "application/json",
				Label:       "RecipeStep",
			}

			// we shuffle the message order to introduce a random delay before each of the messages is sent to
			// simulate out of order sending
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
			var step RecipeStep
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
			var step RecipeStep
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
