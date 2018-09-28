package servicebus_test

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Azure/azure-amqp-common-go/uuid"
	"github.com/Azure/azure-service-bus-go"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func ExampleQueue_getOrBuildQueue() {
	const queueName = "myqueue"

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	if connStr == "" {
		fmt.Println("FATAL: expected environment variable SERVICEBUS_CONNECTION_STRING not set")
		return
	}

	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	qm := ns.NewQueueManager()
	qe, err := qm.Get(ctx, queueName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if qe == nil {
		_, err := qm.Put(ctx, queueName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	q, err := ns.NewQueue(queueName)

	fmt.Println(q.Name)
	// Output: myqueue
}

func ExampleQueue_ReceiveSessions() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Setup the required clients for communicating with Service Bus.                                                 //
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	if connStr == "" {
		fmt.Println("FATAL: expected environment variable SERVICEBUS_CONNECTION_STRING not set")
		return
	}

	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	client, err := ns.NewQueue("receivesession")
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Publish five session's worth of data.                                                                          //
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	const numSessions = 5
	adjectives := []string{"Doltish", "Foolish", "Juvenile"}
	nouns := []string{"Automaton", "Luddite", "Monkey", "Neanderthal"}

	// seed chosen arbitrarily, see https://en.wikipedia.org/wiki/Taxicab_number
	generator := rand.New(rand.NewSource(1729))

	for i := 0; i < numSessions; i++ {
		var sessionID string
		if id, err := uuid.NewV4(); err == nil {
			sessionID = id.String()
		} else {
			fmt.Println("FATAL: ", err)
			return
		}

		if err != nil {
			fmt.Println("FATAL: ", err)
			return
		}

		prepareMessage := func(payload string) (retval *servicebus.Message) {
			retval = servicebus.NewMessageFromString(payload)
			retval.GroupID = &sessionID
			return retval
		}

		adj := adjectives[generator.Intn(len(adjectives))]
		err = client.Send(ctx, prepareMessage(adj))
		if err != nil {
			fmt.Println("FATAL: ", err)
			return
		}

		noun := nouns[generator.Intn(len(nouns))]
		err = client.Send(ctx, prepareMessage(noun))
		if err != nil {
			fmt.Println("FATAL: ", err)
			return
		}

		num := fmt.Sprintf("%02d", generator.Intn(100))
		client.Send(ctx, prepareMessage(num))
		if err != nil {
			fmt.Println("FATAL: ", err)
			return
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Receive and process the previously published sessions.                                                         //
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	inner, innerCancel := context.WithCancel(ctx)

	builder := &bytes.Buffer{}
	messagesReceived := 0

	handler := servicebus.NewSessionHandler(
		servicebus.HandlerFunc(func(ctx context.Context, msg *servicebus.Message) servicebus.DispositionAction {
			builder.Write(msg.Data)

			// The following clause is needed to quit receiving after 5 messages are received.
			messagesReceived++
			if messagesReceived >= numSessions {
				innerCancel()
			}

			return msg.Complete()
		}),
		func(_ *servicebus.MessageSession) error {
			builder.Reset()
			return nil
		},
		func() {
			fmt.Println(builder.String())
		})

	err = client.ReceiveSessions(inner, handler)
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	// Output:
	// FoolishMonkey10
	// JuvenileNeanderthal50
	// FoolishMonkey37
	// JuvenileNeanderthal05
	// JuvenileAutomaton68
}
