package servicebus

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"

	mgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"pack.ag/amqp"
)

const (
	banner = `
   _____                 _               ____            
  / ___/___  ______   __(_)________     / __ )__  _______
  \__ \/ _ \/ ___/ | / / // ___/ _ \   / __  / / / / ___/
 ___/ /  __/ /   | |/ / // /__/  __/  / /_/ / /_/ (__  ) 
/____/\___/_/    |___/_/ \___/\___/  /_____/\__,_/____/

`
)

var (
	connStrRegex = regexp.MustCompile(`Endpoint=sb:\/\/(?P<Host>.+?);SharedAccessKeyName=(?P<KeyName>.+?);SharedAccessKey=(?P<Key>.+)`)
)

// SenderReceiver provides the ability to send and receive messages
type (
	SenderReceiver interface {
		Send(ctx context.Context, entityPath string, msg *amqp.Message, opts ...SendOption) error
		Receive(entityPath string, handler Handler) error
		Close() error
	}

	// EntityManager provides the ability to manage Service Bus entities (Queues, Topics, Subscriptions, etc.)
	EntityManager interface {
		EnsureQueue(ctx context.Context, name string, opts ...QueueOption) (*mgmt.SBQueue, error)
		DeleteQueue(ctx context.Context, name string) error
		EnsureTopic(ctx context.Context, name string, opts ...TopicOption) (*mgmt.SBTopic, error)
		DeleteTopic(ctx context.Context, name string) error
		EnsureSubscription(ctx context.Context, topicName, name string, opts ...SubscriptionOption) (*mgmt.SBSubscription, error)
		DeleteSubscription(ctx context.Context, topicName, name string) error
	}

	// SenderReceiverManager provides Service Bus entity management as well as access to send and receive messages
	SenderReceiverManager interface {
		SenderReceiver
		EntityManager
	}

	// ServicePrincipalCredentials contains the details needed to authenticate to Azure Active Directory with a Service
	// Principal. For more info on Service Principals see: https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-create-service-principal-portal
	ServicePrincipalCredentials struct {
		TenantID      string
		ApplicationID string
		Secret        string
	}

	// serviceBus provides a simplified facade over the AMQP implementation of Azure Service Bus.
	serviceBus struct {
		name             uuid.UUID
		client           *amqp.Client
		clientMu         sync.Mutex
		armToken         *adal.ServicePrincipalToken
		sbToken          *adal.ServicePrincipalToken
		connectionString string
		environment      azure.Environment
		subscriptionID   string
		resourceGroup    string
		namespace        string
		primaryKey       string
		receivers        []*receiver
		senders          map[string]*sender
		receiverMu       sync.Mutex
		senderMu         sync.Mutex
		Logger           *log.Logger
		cbsMu            sync.Mutex
		cbsLink          *cbsLink
	}

	// parsedConn is the structure of a parsed Service Bus connection string.
	parsedConn struct {
		Host    string
		KeyName string
		Key     string
	}

	// Handler is the function signature for any receiver of AMQP messages
	Handler func(context.Context, *amqp.Message) error
)

// NewWithConnectionString creates a new connected instance of an Azure Service Bus given a connection string with the
// same format as the Azure portal
// (e.g. Endpoint=sb://XXXXX.servicebus.windows.net/;SharedAccessKeyName=XXXXX;SharedAccessKey=XXXXX). The Service Bus
// instance returned from this function does not have the ability to manage Subscriptions, Topics or Queues. The
// instance is only able to use existing Service Bus entities.
func NewWithConnectionString(connStr string) (SenderReceiver, error) {
	return newWithConnectionString(connStr)
}

// NewWithServicePrincipal builds a Service Bus SenderReceiverManager which authenticates with Azure Active Directory
// using Claims-based Security
func NewWithServicePrincipal(subscriptionID, namespace string, credentials ServicePrincipalCredentials, env azure.Environment) (SenderReceiverManager, error) {
	armToken, err := getArmTokenProvider(credentials, env)
	if err != nil {
		return nil, err
	}

	sbToken, err := getServiceBusTokenProvider(credentials, env)
	if err != nil {
		return nil, err
	}

	return NewWithTokenProviders(subscriptionID, namespace, armToken, sbToken, env)
}

// NewWithTokenProviders builds a Service Bus SenderReceiverManager which authenticates with Azure Active Directory
// using Claims-based Security using Azure Active Directory token providers
func NewWithTokenProviders(subscriptionID, namespace string, armToken, serviceBusToken *adal.ServicePrincipalToken, env azure.Environment) (SenderReceiverManager, error) {
	sb := &serviceBus{
		name:           uuid.NewV4(),
		sbToken:        serviceBusToken,
		armToken:       armToken,
		namespace:      namespace,
		subscriptionID: subscriptionID,
		environment:    env,
		Logger:         log.New(),
		senders:        make(map[string]*sender),
	}
	sb.Logger.SetLevel(log.WarnLevel)

	ns, err := sb.GetNamespace(context.Background())
	if err != nil {
		return nil, err
	}

	parsedID, err := parseAzureResourceID(*ns.ID)
	if err != nil {
		return nil, err
	}
	sb.resourceGroup = parsedID.ResourceGroup

	return sb, nil
}

func (sb *serviceBus) Start() error {
	log.Println(banner)
	return nil
}

// Close drains and closes all of the existing senders, receivers and connections
func (sb *serviceBus) Close() error {
	// TODO: add some better error handling for cleaning up on Close
	sb.drainReceivers()
	sb.drainSenders()
	if sb.client != nil {
		log.Debugf("closing sb amqp connection %v", sb)
		sb.client.Close()
	}
	return nil
}

// Listen subscribes for messages sent to the provided entityPath.
func (sb *serviceBus) Receive(entityPath string, handler Handler) error {
	sb.receiverMu.Lock()
	defer sb.receiverMu.Unlock()

	receiver, err := sb.newReceiver(entityPath)
	if err != nil {
		return err
	}

	sb.receivers = append(sb.receivers, receiver)
	receiver.Listen(handler)
	return nil
}

// Send sends a message to a provided entity path with options
func (sb *serviceBus) Send(ctx context.Context, entityPath string, msg *amqp.Message, opts ...SendOption) error {
	sender, err := sb.fetchSender(entityPath)
	if err != nil {
		return err
	}

	return sender.Send(ctx, msg, opts...)
}

func (sb *serviceBus) connection() (*amqp.Client, error) {
	sb.clientMu.Lock()
	defer sb.clientMu.Unlock()

	if sb.client == nil && sb.claimsBasedSecurityEnabled() {
		host := getHostName(sb.namespace)
		client, err := amqp.Dial(host, amqp.ConnSASLAnonymous(), amqp.ConnMaxSessions(65535))
		if err != nil {
			return nil, err
		}
		sb.client = client
	}
	return sb.client, nil
}

func (sb *serviceBus) newSession() (*amqp.Session, error) {
	conn, err := sb.connection()
	if err != nil {
		return nil, err
	}
	return conn.NewSession()
}

func (sb *serviceBus) fetchSender(entityPath string) (*sender, error) {
	sb.senderMu.Lock()
	defer sb.senderMu.Unlock()

	entry, ok := sb.senders[entityPath]
	if ok {
		return entry, nil
	}

	sender, err := sb.newSender(entityPath)
	if err != nil {
		return nil, err
	}
	sb.senders[entityPath] = sender
	return sender, nil
}

func newClientWithConnectionString(connStr string) (*amqp.Client, error) {
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
	client, err := newClientWithConnectionString(connStr)
	if err != nil {
		return nil, err
	}

	sb := &serviceBus{
		name:             uuid.NewV4(),
		Logger:           log.New(),
		client:           client,
		connectionString: connStr,
	}
	sb.Logger.SetLevel(log.WarnLevel)
	sb.senders = make(map[string]*sender)
	return sb, nil
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

func getHostName(namespace string) string {
	return fmt.Sprintf("amqps://%s.%s", namespace, "servicebus.windows.net")
}

// claimsBasedSecurityEnabled indicates that the connection will use AAD JWT RBAC to authenticate in connections
func (sb *serviceBus) claimsBasedSecurityEnabled() bool {
	return sb.sbToken != nil
}

func getArmTokenProvider(credential ServicePrincipalCredentials, env azure.Environment) (*adal.ServicePrincipalToken, error) {
	return getTokenProvider(azure.PublicCloud.ResourceManagerEndpoint, credential, env)
}

func getServiceBusTokenProvider(credential ServicePrincipalCredentials, env azure.Environment) (*adal.ServicePrincipalToken, error) {
	return getTokenProvider("https://servicebus.azure.net/", credential, env)
}

func getTokenProvider(resourceURI string, cred ServicePrincipalCredentials, env azure.Environment) (*adal.ServicePrincipalToken, error) {
	oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, cred.TenantID)
	if err != nil {
		log.Fatalln(err)
	}

	tokenProvider, err := adal.NewServicePrincipalToken(*oauthConfig, cred.ApplicationID, cred.Secret, resourceURI)
	if err != nil {
		return nil, err
	}

	err = tokenProvider.Refresh()
	if err != nil {
		return nil, err
	}

	return tokenProvider, nil
}

func (sb *serviceBus) drainReceivers() error {
	log.Debugln("draining receivers")
	sb.receiverMu.Lock()
	defer sb.receiverMu.Unlock()

	for _, receiver := range sb.receivers {
		receiver.Close()
	}
	sb.receivers = []*receiver{}
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
