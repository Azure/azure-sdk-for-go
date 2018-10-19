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

func ExampleQueue_sessionsRoundTrip() {
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
	//                                                                                                                //
	// The sessions are deliberately interleaved to demonstrate consumption semantics.                                //
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	const numSessions = 5
	adjectives := []string{"Doltish", "Foolish", "Juvenile"}
	nouns := []string{"Automaton", "Luddite", "Monkey", "Neanderthal"}

	// seed chosen arbitrarily, see https://en.wikipedia.org/wiki/Taxicab_number
	generator := rand.New(rand.NewSource(1729))

	sessionIDs := make([]string, numSessions)

	// Establish a set of sessions
	for i := 0; i < numSessions; i++ {
		if rawSessionID, err := uuid.NewV4(); err == nil {
			sessionIDs[i] = rawSessionID.String()
		} else {
			fmt.Println("FATAL: ", err)
			return
		}
	}

	// Publish an adjective for each session
	for i := 0; i < numSessions; i++ {
		adj := adjectives[generator.Intn(len(adjectives))]
		msg := servicebus.NewMessageFromString(adj)
		msg.GroupID = &sessionIDs[i]
		if err := client.Send(ctx, msg); err != nil {
			fmt.Println("FATAL: ", err)
			return
		}
	}

	// Publish a noun for each session
	for i := 0; i < numSessions; i++ {
		noun := nouns[generator.Intn(len(nouns))]
		msg := servicebus.NewMessageFromString(noun)
		msg.GroupID = &sessionIDs[i]
		if err := client.Send(ctx, msg); err != nil {
			fmt.Println("FATAL: ", err)
			return
		}
	}

	// Publish a numeric suffix for each session
	for i := 0; i < numSessions; i++ {
		suffix := fmt.Sprintf("%02d", generator.Intn(100))
		msg := servicebus.NewMessageFromString(suffix)
		msg.GroupID = &sessionIDs[i]
		if err := client.Send(ctx, msg); err != nil {
			fmt.Println("FATAL: ", err)
			return
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Receive and process the previously published sessions.                                                         //
	//                                                                                                                //
	// The order the sessions are received in is not guaranteed, so the expected output must be "Unordered output".   //
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	handler := &SessionPrinter{}
	for i := 0; i < numSessions; i++ {
		if err := client.ReceiveOneSession(ctx, nil, handler); err != nil {
			fmt.Println("FATAL: ", err)
			return
		}
	}

	// Unordered output:
	// FoolishMonkey63
	// FoolishLuddite05
	// JuvenileMonkey80
	// JuvenileLuddite84
	// FoolishLuddite68
}

type SessionPrinter struct {
	builder          *bytes.Buffer
	messageSession   *servicebus.MessageSession
	messagesReceived uint
}

func (sp *SessionPrinter) Start(ms *servicebus.MessageSession) error {
	if sp.builder == nil {
		sp.builder = &bytes.Buffer{}
	} else {
		sp.builder.Reset()
	}
	sp.messagesReceived = 0
	sp.messageSession = ms
	return nil
}

func (sp *SessionPrinter) Handle(_ context.Context, msg *servicebus.Message) servicebus.DispositionAction {
	sp.builder.Write(msg.Data)

	sp.messagesReceived++

	if sp.messagesReceived >= 3 {
		defer sp.messageSession.Close()
	}

	return msg.Complete()
}

func (sp *SessionPrinter) End() {
	fmt.Println(sp.builder.String())
}
