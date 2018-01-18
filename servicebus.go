package servicebus

import (
	"context"
	"errors"
	mgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"log"
	"pack.ag/amqp"
	"regexp"
)

var (
	connStrRegex = regexp.MustCompile(`Endpoint=sb:\/\/(?P<Host>.+?);SharedAccessKeyName=(?P<KeyName>.+?);SharedAccessKey=(?P<Key>.+)`)
	receivers    = make(map[string]*amqp.Receiver)
	senders      = make(map[string]*amqp.Sender)
)

// SenderReceiver provides the ability to send and receive messages
type SenderReceiver interface {
	Send(ctx context.Context, entityPath string, msg *amqp.Message) error
	Receive(entityPath string, handler Handler) error
	Close()
}

// EntityManager provides the ability to manage Service Bus entities (Queues, Topics, Subscriptions, etc.)
type EntityManager interface {
	EnsureQueue(ctx context.Context, queueName string) error
	DeleteQueue(ctx context.Context, queueName string) error
}

// SenderReceiverManager provides Service Bus entity management as well as access to send and receive messages
type SenderReceiverManager interface {
	SenderReceiver
	EntityManager
}

// serviceBus provides a simplified facade over the AMQP implementation of Azure Service Bus.
type serviceBus struct {
	client  *amqp.Client
	session *amqp.Session
	token   *adal.ServicePrincipalToken
}

// parsedConn is the structure of a parsed Service Bus connection string.
type parsedConn struct {
	Host    string
	KeyName string
	Key     string
}

// Handler is the function signature for any receiver of AMQP messages
type Handler func(context.Context, *amqp.Message) error

// NewWithConnectionString creates a new connected instance of an Azure Service Bus given a connection string with the
// same format as the Azure portal
// (e.g. Endpoint=sb://XXXXX.servicebus.windows.net/;SharedAccessKeyName=XXXXX;SharedAccessKey=XXXXX). The Service Bus
// instance returned from this function does not have the ability to manage Subscriptions, Topics or Queues. The
// instance is only able to use existing Service Bus entities.
func NewWithConnectionString(connStr string) (SenderReceiver, error) {
	return newWithConnectionString(connStr)
}

func newWithConnectionString(connStr string) (*serviceBus, error) {
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

	return &serviceBus{
		client:  client,
		session: session,
	}, nil
}

// NewWithMSI creates a new connected instance of an Azure Service Bus given a subscription Id, resource group,
// Service Bus namespace, and Service Bus authorization rule name.
func NewWithMSI(subscriptionID, resourceGroup, namespace, ruleName string, environment azure.Environment) (SenderReceiverManager, error) {
	msiEndpoint, err := adal.GetMSIVMEndpoint()
	spToken, err := adal.NewServicePrincipalTokenFromMSI(msiEndpoint, environment.ResourceManagerEndpoint)

	if err != nil {
		return nil, err
	}

	return NewWithSPToken(spToken, subscriptionID, resourceGroup, namespace, ruleName, environment)
}

// NewWithSPToken creates a new connected instance of an Azure Service Bus given a, Azure Active Directory service
// principal token subscription Id, resource group, Service Bus namespace, and Service Bus authorization rule name.
func NewWithSPToken(spToken *adal.ServicePrincipalToken, subscriptionID, resourceGroup, namespace, ruleName string, environment azure.Environment) (SenderReceiverManager, error) {
	authorizer := autorest.NewBearerAuthorizer(spToken)

	nsClient := mgmt.NewNamespacesClientWithBaseURI(environment.ResourceManagerEndpoint, subscriptionID)
	nsClient.Authorizer = authorizer
	nsClient.AddToUserAgent("dataplane-servicebus")

	result, err := nsClient.ListKeys(context.Background(), resourceGroup, namespace, ruleName)
	if err != nil {
		return nil, err
	}

	primary := *result.PrimaryConnectionString
	sb, err := newWithConnectionString(primary)
	if err != nil {
		return nil, err
	}

	sb.token = spToken
	return sb, err
}

// Close closes the Service Bus connection.
func (sb *serviceBus) Close() {
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
func (sb *serviceBus) Receive(entityPath string, handler Handler) error {
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

func (sb *serviceBus) fetchReceiver(entityPath string) (*amqp.Receiver, error) {
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
func (sb *serviceBus) Send(ctx context.Context, entityPath string, msg *amqp.Message) error {
	sender, err := sb.fetchSender(entityPath)
	if err != nil {
		return err
	}

	return sender.Send(ctx, msg)
}

func (sb *serviceBus) fetchSender(entityPath string) (*amqp.Sender, error) {
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

func (sb *serviceBus) EnsureQueue(ctx context.Context, queueName string) error {
	return nil
}

func (sb *serviceBus) DeleteQueue(ctx context.Context, queuename string) error {
	return nil
}
