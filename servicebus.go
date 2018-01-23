package servicebus

import (
	"context"
	"errors"
	mgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	log "github.com/sirupsen/logrus"
	"pack.ag/amqp"
	"regexp"
	"sync"
)

var (
	connStrRegex = regexp.MustCompile(`Endpoint=sb:\/\/(?P<Host>.+?);SharedAccessKeyName=(?P<KeyName>.+?);SharedAccessKey=(?P<Key>.+)`)
)

// SenderReceiver provides the ability to send and receive messages
type SenderReceiver interface {
	Send(ctx context.Context, entityPath string, msg *amqp.Message) error
	Receive(entityPath string, handler Handler) error
	Close() error
}

// EntityManager provides the ability to manage Service Bus entities (Queues, Topics, Subscriptions, etc.)
type EntityManager interface {
	EnsureQueue(ctx context.Context, queueName string, properties *mgmt.SBQueueProperties) (*mgmt.SBQueue, error)
	DeleteQueue(ctx context.Context, queueName string) error
}

// SenderReceiverManager provides Service Bus entity management as well as access to send and receive messages
type SenderReceiverManager interface {
	SenderReceiver
	EntityManager
}

// serviceBus provides a simplified facade over the AMQP implementation of Azure Service Bus.
type serviceBus struct {
	client         *amqp.Client
	token          *adal.ServicePrincipalToken
	environment    azure.Environment
	subscriptionID string
	resourceGroup  string
	namespace      string
	primaryKey     string
	receivers      []*Receiver
	senders        map[string]*Sender
	receiverMu     sync.Mutex
	senderMu       sync.Mutex
	Logger         *log.Logger
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

func newClient(connStr string) (*amqp.Client, error) {
	if connStr == "" {
		return nil, errors.New("connection string can not be null")
	}
	parsed, err := parsedConnectionFromStr(connStr)
	if err != nil {
		return nil, errors.New("connection string was not in expected format (Endpoint=sb://XXXXX.servicebus.windows.net/;SharedAccessKeyName=XXXXX;SharedAccessKey=XXXXX)")
	}

	client, err := amqp.Dial(parsed.Host, amqp.ConnSASLPlain(parsed.KeyName, parsed.Key), amqp.ConnMaxSessions(65535))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func newWithConnectionString(connStr string) (*serviceBus, error) {
	client, err := newClient(connStr)
	if err != nil {
		return nil, err
	}

	sb := &serviceBus{
		Logger: log.New(),
		client: client,
	}
	sb.Logger.SetLevel(log.WarnLevel)
	sb.senders = make(map[string]*Sender)
	return sb, nil
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
	sb.environment = environment
	sb.subscriptionID = subscriptionID
	sb.resourceGroup = resourceGroup
	sb.namespace = namespace
	sb.primaryKey = primary
	return sb, err
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

// Close drains and closes all of the existing senders, receivers and connections
func (sb *serviceBus) Close() error {
	// TODO: add some better error handling for cleaning up on Close
	sb.drainReceivers()
	sb.drainSenders()
	log.Debugf("closing sb %v", sb)
	sb.client.Close()
	return nil
}

func (sb *serviceBus) drainReceivers() error {
	log.Debugln("draining receivers")
	sb.receiverMu.Lock()
	defer sb.receiverMu.Unlock()

	for _, receiver := range sb.receivers {
		receiver.Close()
	}
	sb.receivers = []*Receiver{}
	return nil
}

func (sb *serviceBus) drainSenders() error {
	log.Debugln("draining senders")
	sb.senderMu.Lock()
	defer sb.senderMu.Unlock()

	for key, sender := range sb.senders {
		sender.Close()
		delete(sb.senders, key)
	}
	return nil
}

// Listen subscribes for messages sent to the provided entityPath.
func (sb *serviceBus) Receive(entityPath string, handler Handler) error {
	sb.receiverMu.Lock()
	defer sb.receiverMu.Unlock()

	receiver, err := NewReceiver(sb.client, entityPath)
	if err != nil {
		return err
	}

	sb.receivers = append(sb.receivers, receiver)
	receiver.Listen(handler)
	return nil
}

// Send sends a message to a provided entity path.
func (sb *serviceBus) Send(ctx context.Context, entityPath string, msg *amqp.Message) error {
	sender, err := sb.fetchSender(entityPath)
	if err != nil {
		return err
	}

	return sender.Send(ctx, msg)
}

func (sb *serviceBus) fetchSender(entityPath string) (*Sender, error) {
	sb.senderMu.Lock()
	defer sb.senderMu.Unlock()

	entry, ok := sb.senders[entityPath]
	if ok {
		log.Debugf("found sender for entity path %s", entityPath)
		return entry, nil
	}

	log.Debugf("creating a new sender for entity path %s", entityPath)
	sender, err := NewSender(sb.client, entityPath)
	if err != nil {
		return nil, err
	}
	sb.senders[entityPath] = sender
	return sender, nil
}

// EnsureQueue makes sure a queue exists in the given namespace. If the queue doesn't exist, it will create it with
// the specified name and properties. If properties are not specified, it will build a default partitioned queue.
func (sb *serviceBus) EnsureQueue(ctx context.Context, queueName string, properties *mgmt.SBQueueProperties) (*mgmt.SBQueue, error) {
	log.Debugf("ensuring exists queue %s", queueName)
	queueClient := sb.getQueueMgmtClient()
	queue, err := queueClient.Get(ctx, sb.resourceGroup, sb.namespace, queueName)
	// TODO: check if the queue properties are the same as the requested. If not, throw error or build new queue??

	if properties == nil {
		log.Debugf("no properties specified, so using default partitioned queue for %s", queueName)
		properties = &mgmt.SBQueueProperties{
			EnablePartitioning: ptrBool(true),
		}
	}

	if err != nil {
		log.Debugf("building a new queue %s", queueName)
		newQueue := mgmt.SBQueue{
			Name:              &queueName,
			SBQueueProperties: properties,
		}
		queue, err = queueClient.CreateOrUpdate(ctx, sb.resourceGroup, sb.namespace, queueName, newQueue)
		if err != nil {
			return nil, err
		}
	}
	return &queue, nil
}

// DeleteQueue deletes an existing queue
func (sb *serviceBus) DeleteQueue(ctx context.Context, queueName string) error {
	queueClient := sb.getQueueMgmtClient()
	_, err := queueClient.Delete(ctx, sb.resourceGroup, sb.namespace, queueName)
	return err
}

func (sb *serviceBus) getQueueMgmtClient() mgmt.QueuesClient {
	client := mgmt.NewQueuesClientWithBaseURI(sb.environment.ResourceManagerEndpoint, sb.subscriptionID)
	client.Authorizer = autorest.NewBearerAuthorizer(sb.token)
	return client
}
