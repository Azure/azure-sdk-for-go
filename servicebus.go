package servicebus

import (
	"errors"
	"regexp"

	"pack.ag/amqp"
)

var (
	connStrRegex = regexp.MustCompile(`Endpoint=sb:\/\/(?P<Host>.+?);SharedAccessKeyName=(?P<KeyName>.+?);SharedAccessKey=(?P<Key>.+)`)
)

// ServiceBus provides a simplified facade over the AMQP implementation of Azure Service Bus
type ServiceBus struct {
	client *amqp.Client
}

type parsedConn struct {
	Host    string
	KeyName string
	Key     string
}

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

	return &ServiceBus{
		client: client,
	}, nil
}

// Close closes the Service Bus connection.
func (sb *ServiceBus) Close() {
	sb.client.Close()
}

func parsedConnectionFromStr(connStr string) (*parsedConn, error) {
	matches := connStrRegex.FindStringSubmatch(connStr)
	parsed, err := newParsedConnection(matches[1], matches[2], matches[3])
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

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
