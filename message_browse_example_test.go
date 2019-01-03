package servicebus_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

type (
	Scientist struct {
		Surname   string `json:"surname,omitempty"`
		FirstName string `json:"firstname,omitempty"`
	}
)

func Example_messageBrowse() {
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
	qEntity, err := ensureQueue(ctx, qm, "MessageBrowseExample")
	if err != nil {
		fmt.Println(err)
		return
	}

	q, err := ns.NewQueue(qEntity.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	txRxCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	go sendMessages(txRxCtx, q)
	time.Sleep(1 * time.Second) // wait a second to ensure a message has landed in the queue
	go peekMessages(txRxCtx, q)

	<-txRxCtx.Done() // wait for the context to finish

	// Output:
	// Firstname: Albert, Surname: Einstein
	// Firstname: Werner, Surname: Heisenberg
	// Firstname: Marie, Surname: Curie
	// Firstname: Steven, Surname: Hawking
	// Firstname: Isaac, Surname: Newton
	// Firstname: Niels, Surname: Bohr
	// Firstname: Michael, Surname: Faraday
	// Firstname: Galileo, Surname: Galilei
	// Firstname: Johannes, Surname: Kepler
	// Firstname: Nikolaus, Surname: Kopernikus
}

func sendMessages(ctx context.Context, q *servicebus.Queue) {

	scientists := []Scientist{
		{
			Surname:   "Einstein",
			FirstName: "Albert",
		},
		{
			Surname:   "Heisenberg",
			FirstName: "Werner",
		},
		{
			Surname:   "Curie",
			FirstName: "Marie",
		},
		{
			Surname:   "Hawking",
			FirstName: "Steven",
		},
		{
			Surname:   "Newton",
			FirstName: "Isaac",
		},
		{
			Surname:   "Bohr",
			FirstName: "Niels",
		},
		{
			Surname:   "Faraday",
			FirstName: "Michael",
		},
		{
			Surname:   "Galilei",
			FirstName: "Galileo",
		},
		{
			Surname:   "Kepler",
			FirstName: "Johannes",
		},
		{
			Surname:   "Kopernikus",
			FirstName: "Nikolaus",
		},
	}

	for _, scientist := range scientists {
		bits, err := json.Marshal(scientist)
		if err != nil {
			fmt.Println(err)
			return
		}

		ttl := 2 * time.Minute
		msg := servicebus.NewMessage(bits)
		msg.ContentType = "application/json"
		msg.TTL = &ttl
		if err := q.Send(ctx, msg); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func peekMessages(ctx context.Context, q *servicebus.Queue) {
	var opts []servicebus.PeekOption
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := q.PeekOne(ctx, opts...)
			if err != nil {
				switch err.(type) {
				case servicebus.ErrNoMessages:
					// all done
					return
				default:
					fmt.Println(err)
					return
				}
			}

			var scientist Scientist
			if err := json.Unmarshal(msg.Data, &scientist); err != nil {
				fmt.Println(err)
				return
			}

			opts = []servicebus.PeekOption{servicebus.PeekFromSequenceNumber(*msg.SystemProperties.SequenceNumber)}
			fmt.Printf("Firstname: %s, Surname: %s\n", scientist.FirstName, scientist.Surname)
		}
	}
}
