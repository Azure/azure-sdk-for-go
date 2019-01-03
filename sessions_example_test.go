package servicebus_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

type StepSessionHandler struct {
	messageSession *servicebus.MessageSession
}

// Start is called when a new session is started
func (ssh *StepSessionHandler) Start(ms *servicebus.MessageSession) error {
	ssh.messageSession = ms
	fmt.Println("Begin session: ", *ssh.messageSession.SessionID())
	return nil
}

// Handle is called when a new session message is received
func (ssh *StepSessionHandler) Handle(ctx context.Context, msg *servicebus.Message) error {
	var step RecipeStep
	if err := json.Unmarshal(msg.Data, &step); err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("  Step: %d, %s\n", step.Step, step.Title)

	if step.Step == 5 {
		ssh.messageSession.Close()
	}
	return msg.Complete(ctx)
}

// End is called when the message session is closed. Service Bus will not automatically end your message session. Be
// sure to know when to terminate your own session.
func (ssh *StepSessionHandler) End() {
	fmt.Println("End session: ", *ssh.messageSession.SessionID())
	fmt.Println("")
}

func Example_messageSessions() {
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

	// Create a Service Bus Queue with required sessions enabled. This will ensure that all messages sent and received
	// are bound to a session.
	qm := ns.NewQueueManager()
	qEntity, err := ensureQueue(ctx, qm, "MessageSessionsExample", servicebus.QueueEntityWithRequiredSessions())
	if err != nil {
		fmt.Println(err)
		return
	}

	q, err := ns.NewQueue(qEntity.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	sessions := []string{"foo", "bar", "bazz", "buzz"}
	for _, session := range sessions {
		// send recipe steps
		// note that order is preserved within a given session
		sendSessionRecipeSteps(ctx, session, q)
	}

	// receive messages for each session
	// you can also call q.NewSession(nil) to receive from any available session
	for _, session := range sessions {
		queueSession := q.NewSession(&session)
		err := queueSession.ReceiveOne(ctx, new(StepSessionHandler))
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := queueSession.Close(ctx); err != nil {
			fmt.Println(err)
			return
		}
	}

	// Output:
	// Begin session:  foo
	//   Step: 1, Shop
	//   Step: 2, Unpack
	//   Step: 3, Prepare
	//   Step: 4, Cook
	//   Step: 5, Eat
	// End session:  foo
	//
	// Begin session:  bar
	//   Step: 1, Shop
	//   Step: 2, Unpack
	//   Step: 3, Prepare
	//   Step: 4, Cook
	//   Step: 5, Eat
	// End session:  bar
	//
	// Begin session:  bazz
	//   Step: 1, Shop
	//   Step: 2, Unpack
	//   Step: 3, Prepare
	//   Step: 4, Cook
	//   Step: 5, Eat
	// End session:  bazz
	//
	// Begin session:  buzz
	//   Step: 1, Shop
	//   Step: 2, Unpack
	//   Step: 3, Prepare
	//   Step: 4, Cook
	//   Step: 5, Eat
	// End session:  buzz
}

func sendSessionRecipeSteps(ctx context.Context, sessionID string, q *servicebus.Queue) {
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
		bits, err := json.Marshal(step)
		if err != nil {
			fmt.Println(err)
			return
		}

		msg := servicebus.NewMessage(bits)
		msg.ContentType = "application/json"
		msg.SessionID = &sessionID
		if err := q.Send(ctx, msg); err != nil {
			fmt.Println(err)
			return
		}
	}
}
