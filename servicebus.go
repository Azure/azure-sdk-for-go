package servicebus

import (
	"context"
	"errors"
	mgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	log "github.com/sirupsen/logrus"
	"net"
	"pack.ag/amqp"
	"regexp"
	"sync"
	"time"
)

var (
	connStrRegex = regexp.MustCompile(`Endpoint=sb:\/\/(?P<Host>.+?);SharedAccessKeyName=(?P<KeyName>.+?);SharedAccessKey=(?P<Key>.+)`)
)

type receiverSession struct {
	session  *amqp.Session
	receiver *amqp.Receiver
}

type senderSession struct {
	session *amqp.Session
	sender  *amqp.Sender
}

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
	client           *amqp.Client
	token            *adal.ServicePrincipalToken
	environment      azure.Environment
	subscriptionID   string
	resourceGroup    string
	namespace        string
	primaryKey       string
	receiverSessions map[string]*receiverSession
	senderSessions   map[string]*senderSession
	receiverMu       sync.Mutex
	senderMu         sync.Mutex
	Logger           *log.Logger
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

	client, err := amqp.Dial(parsed.Host, amqp.ConnSASLPlain(parsed.KeyName, parsed.Key), amqp.ConnMaxChannels(65535))
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
	sb.senderSessions = make(map[string]*senderSession)
	sb.receiverSessions = make(map[string]*receiverSession)
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
	log.Infof("closing sb %v", sb)
	err := sb.drainReceivers()
	if err != nil {
		return err
	}

	err = sb.drainSenders()
	if err != nil {
		return err
	}

	return sb.client.Close()
}

func (sb *serviceBus) drainReceivers() error {
	log.Infoln("draining receivers")
	sb.receiverMu.Lock()
	defer sb.receiverMu.Unlock()
	for _, item := range sb.receiverSessions {
		err := item.receiver.Close()
		if err != nil {
			return err
		}
		err = item.session.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (sb *serviceBus) drainSenders() error {
	log.Infoln("draining senders")
	sb.senderMu.Lock()
	defer sb.senderMu.Unlock()
	for key, item := range sb.senderSessions {
		//err := item.sender.Close()
		//if err != nil {
		//	return err
		//}
		err := item.session.Close()
		if err != nil {
			return err
		}
		delete(sb.senderSessions, key)
	}
	return nil
}

// Receive subscribes for messages sent to the provided entityPath.
func (sb *serviceBus) Receive(entityPath string, handler Handler) error {
	receiver, err := sb.fetchReceiver(entityPath)
	if err != nil {
		return err
	}

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			// Receive next message
			msg, err := receiver.Receive(ctx)
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue
			} else if err != nil {
				log.Fatalln(err)
			}
			cancel()

			if msg != nil {
				id := interface{}("null")
				if msg.Properties != nil {
					id = msg.Properties.MessageID
				}
				log.Info("Message received: %s, id: %s\n", msg.Data, id)

				err = handler(ctx, msg)
				if err != nil {
					msg.Reject()
					log.Info("Message rejected \n")
				} else {
					// Accept message
					msg.Accept()
					log.Info("Message accepted \n")
				}
			}
		}
	}()
	return nil
}

func (sb *serviceBus) fetchReceiver(entityPath string) (*amqp.Receiver, error) {
	sb.receiverMu.Lock()
	defer sb.receiverMu.Unlock()

	entry, ok := sb.receiverSessions[entityPath]
	if ok {
		log.Infof("found receiver for entity path %s", entityPath)
		return entry.receiver, nil
	} else {
		log.Infof("creating a new receiver for entity path %s", entityPath)
		session, err := sb.client.NewSession()
		if err != nil {
			return nil, err
		}

		receiver, err := session.NewReceiver(
			amqp.LinkAddress(entityPath),
			amqp.LinkCredit(10),
			amqp.LinkBatching(true))
		if err != nil {
			return nil, err
		}

		receiverSession := &receiverSession{receiver: receiver, session: session}
		sb.receiverSessions[entityPath] = receiverSession
		return receiverSession.receiver, nil
	}
}

// Send sends a message to a provided entity path.
func (sb *serviceBus) Send(ctx context.Context, entityPath string, msg *amqp.Message) error {
	senderSession, err := sb.fetchSender(entityPath)
	if err != nil {
		return err
	}

	return senderSession.sender.Send(ctx, msg)
}

func (sb *serviceBus) fetchSender(entityPath string) (*senderSession, error) {
	sb.senderMu.Lock()
	defer sb.senderMu.Unlock()

	entry, ok := sb.senderSessions[entityPath]
	if ok {
		log.Infof("found sender for entity path %s", entityPath)
		return entry, nil
	} else {
		log.Infof("creating a new sender for entity path %s", entityPath)
		session, err := sb.client.NewSession()
		if err != nil {
			return nil, err
		}

		sender, err := session.NewSender(amqp.LinkAddress(entityPath))
		if err != nil {
			return nil, err
		}

		senderSession := &senderSession{session: session, sender: sender}
		sb.senderSessions[entityPath] = senderSession
		return senderSession, nil
	}
}

// EnsureQueue makes sure a queue exists in the given namespace. If the queue doesn't exist, it will create it with
// the specified name and properties. If properties are not specified, it will build a default partitioned queue.
func (sb *serviceBus) EnsureQueue(ctx context.Context, queueName string, properties *mgmt.SBQueueProperties) (*mgmt.SBQueue, error) {
	log.Infof("ensuring exists queue %s", queueName)
	queueClient := sb.getQueueMgmtClient()
	queue, err := queueClient.Get(ctx, sb.resourceGroup, sb.namespace, queueName)

	if properties == nil {
		log.Infof("no properties specified, so using default partitioned queue for %s", queueName)
		properties = &mgmt.SBQueueProperties{
			EnablePartitioning: ptrBool(false),
		}
	}

	if err != nil {
		log.Infof("building a new queue %s", queueName)
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
