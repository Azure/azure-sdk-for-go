package servicebus

import (
	"errors"
	"regexp"

	"context"
	"log"
	"pack.ag/amqp"
)

var (
	connStrRegex = regexp.MustCompile(`Endpoint=sb:\/\/(?P<Host>.+?);SharedAccessKeyName=(?P<KeyName>.+?);SharedAccessKey=(?P<Key>.+)`)
	receivers    = make(map[string]*amqp.Receiver)
	senders      = make(map[string]*amqp.Sender)
)

// ServiceBus provides a simplified facade over the AMQP implementation of Azure Service Bus.
type ServiceBus struct {
	client  *amqp.Client
	session *amqp.Session
}

// parsedConn is the structure of a parsed Service Bus connection string.
type parsedConn struct {
	Host    string
	KeyName string
	Key     string
}

// Handler is the function signature for any receiver of AMQP messages
type Handler func(context.Context, *amqp.Message) error

// New creates a new connected instance of an Azure Service Bus given a connection string with the same format as the
// Azure portal (e.g. Endpoint=sb://XXXXX.servicebus.windows.net/;SharedAccessKeyName=XXXXX;SharedAccessKey=XXXXX).
func New(connStr string) (*ServiceBus, error) {
	if connStr == "" {
		return nil, errors.New("connection string can not be null")
	}
	parsed, err := parsedConnectionFromStr(connStr)
	if err != nil {
		return nil, errors.New("connection string was not in expected format (Endpoint=sb://XXXXX.servicebus.windows.net/;SharedAccessKeyName=XXXXX;SharedAccessKey=XXXXX)")
	}

	client, err := amqp.Dial(parsed.Host, amqp.ConnSASLPlain(parsed.KeyName, parsed.Key))
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	return &ServiceBus{
		client:  client,
		session: session,
	}, nil
}

// Close closes the Service Bus connection.
func (sb *ServiceBus) Close() {
	sb.client.Close()
}

// parsedConnectionFromStr takes a string connection string from the Azure portal and returns the parsed representation.
func parsedConnectionFromStr(connStr string) (*parsedConn, error) {
	matches := connStrRegex.FindStringSubmatch(connStr)
	parsed, err := newParsedConnection(matches[1], matches[2], matches[3])
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

// newParsedConnection is a constructor for a parsedConn and verifies each of the inputs is non-null.
func newParsedConnection(host string, keyName string, key string) (*parsedConn, error) {
	if host == "" || keyName == "" || key == "" {
		return nil, errors.New("connection string contains an empty entry")
	}
	return &parsedConn{
		Host:    "amqps://" + host,
		KeyName: keyName,
		Key:     key,
	}, nil
}

// Receive subscribes for messages sent to the provided entityPath.
func (sb *ServiceBus) Receive(entityPath string, handler Handler) error {
	receiver, err := sb.fetchReceiver(entityPath)
	if err != nil {
		return err
	}

	ctx := context.Background()

	go func() {
		for {
			// Receive next message
			msg, err := receiver.Receive(ctx)
			if err != nil {
				log.Fatal("Reading message from AMQP:", err)
			}

			if msg != nil {
				id := interface{}("null")
				if msg.Properties != nil {
					id = msg.Properties.MessageID
				}
				log.Printf("Message received: %s, id: %s\n", msg.Data, id)

				err = handler(ctx, msg)
				if err != nil {
					msg.Reject()
					log.Printf("Message rejected \n")
				} else {
					// Accept message
					msg.Accept()
					log.Printf("Message accepted \n")
				}
			}
		}
	}()
	return nil
}

func (sb *ServiceBus) fetchReceiver(entityPath string) (*amqp.Receiver, error) {
	receiver, ok := receivers[entityPath]
	if ok {
		return receiver, nil
	}

	receiver, err := sb.session.NewReceiver(amqp.LinkAddress(entityPath), amqp.LinkCredit(10))
	if err != nil {
		return nil, err
	}
	receivers[entityPath] = receiver
	return receiver, nil
}

// Send sends a message to a provided entity path.
func (sb *ServiceBus) Send(ctx context.Context, entityPath string, msg *amqp.Message) error {
	sender, err := sb.fetchSender(entityPath)
	if err != nil {
		return err
	}

	return sender.Send(ctx, msg)
}

func (sb *ServiceBus) fetchSender(entityPath string) (*amqp.Sender, error) {
	sender, ok := senders[entityPath]
	if ok {
		return sender, nil
	}

	sender, err := sb.session.NewSender(amqp.LinkAddress(entityPath))
	if err != nil {
		return nil, err
	}
	senders[entityPath] = sender
	return sender, nil
}
