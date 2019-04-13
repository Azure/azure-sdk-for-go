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
	if err := godotenv.Load(); err != nil {
		fmt.Println("FATAL: ", err)
	}
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
	if err != nil && !servicebus.IsErrNotFound(err) {
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
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(q.Name)
	// Output: myqueue
}

func ExampleQueue_Send() {
	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString("<your connection string here>"))
	if err != nil {
		return
	}

	client, err := ns.NewQueue("myqueue")
	if err != nil {
		return
	}

	// Create a context to limit how long we will try to send, then push the message over the wire.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Send(ctx, servicebus.NewMessageFromString("Hello World!!!")); err != nil {
		fmt.Println("FATAL: ", err)
	}
}

func ExampleQueue_Receive() {
	// Define a function that should be executed when a message is received.
	var printMessage servicebus.HandlerFunc = func(ctx context.Context, msg *servicebus.Message) error {
		fmt.Println(string(msg.Data))
		return msg.Complete(ctx)
	}

	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString("<your connection string here>"))
	if err != nil {
		return
	}

	client, err := ns.NewQueue("myqueue")
	if err != nil {
		return
	}

	// Define a context to limit how long we will block to receive messages, then start serving our function.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := client.Receive(ctx, printMessage); err != nil {
		fmt.Println("FATAL: ", err)
	}
}

func ExampleQueue_Receive_second() {
	// Set concurrent number
	const concurrentNum = 5
	// Define msg chan
	msgChan := make(chan *servicebus.Message, concurrentNum)
	// Define a function that should be executed when a message is received.
	var concurrentHandler servicebus.HandlerFunc = func(ctx context.Context, msg *servicebus.Message) error {
		msgChan <- msg
		return nil
	}

	// Define msg workers
	for i := 0; i < concurrentNum; i++ {
		go func() {
			for msg := range msgChan {

				ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Millisecond)
				defer cancel()

				fmt.Println(string(msg.Data))
				msg.Complete(ctx)
			}
		} ()
	}

	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString("<your connection string here>"))
	if err != nil {
		close(msgChan)
		return
	}

	// Init queue client with prefetch count
	client, err := ns.NewQueue("myqueue", servicebus.QueueWithPrefetchCount(concurrentNum))
	if err != nil {
		close(msgChan)
		return
	}

	// Define a context to limit how long we will block to receive messages, then start serving our function.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := client.Receive(ctx, concurrentHandler); err != nil {
		fmt.Println("FATAL: ", err)
	}
	// Close the message chan
	close(msgChan)
}

func ExampleQueue_scheduleAndCancelMessages() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute+40*time.Second)
	defer cancel()

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	if connStr == "" {
		fmt.Println("FATAL: expected environment variable SERVICEBUS_CONNECTION_STRING not set")
		return
	}

	// Create a client to communicate with a Service Bus Namespace.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	client, err := ns.NewQueue("schedulewithqueue")
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	// The delay that we should schedule a message for.
	const waitTime = 1 * time.Minute
	expectedTime := time.Now().Add(waitTime)
	msg := servicebus.NewMessageFromString("to the future!!")

	scheduled, err := client.ScheduleAt(ctx, expectedTime, msg)
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	err = client.CancelScheduled(ctx, scheduled...)
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	fmt.Println("All Messages Scheduled and Cancelled")

	// Output: All Messages Scheduled and Cancelled
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
		msg.SessionID = &sessionIDs[i]
		if err := client.Send(ctx, msg); err != nil {
			fmt.Println("FATAL: ", err)
			return
		}
	}

	// Publish a noun for each session
	for i := 0; i < numSessions; i++ {
		noun := nouns[generator.Intn(len(nouns))]
		msg := servicebus.NewMessageFromString(noun)
		msg.SessionID = &sessionIDs[i]
		if err := client.Send(ctx, msg); err != nil {
			fmt.Println("FATAL: ", err)
			return
		}
	}

	// Publish a numeric suffix for each session
	for i := 0; i < numSessions; i++ {
		suffix := fmt.Sprintf("%02d", generator.Intn(100))
		msg := servicebus.NewMessageFromString(suffix)
		msg.SessionID = &sessionIDs[i]
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
	for i := 0; i < numSessions; i++ {
		handler := &SessionPrinter{}
		qs := client.NewSession(nil)
		if err := qs.ReceiveOne(ctx, handler); err != nil {
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

func (sp *SessionPrinter) Handle(ctx context.Context, msg *servicebus.Message) error {
	sp.builder.Write(msg.Data)

	sp.messagesReceived++

	if sp.messagesReceived >= 3 {
		defer sp.messageSession.Close()
	}

	return msg.Complete(ctx)
}

func (sp *SessionPrinter) End() {
	fmt.Println(sp.builder.String())
}
