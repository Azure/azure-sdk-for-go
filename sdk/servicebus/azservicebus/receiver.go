package azservicebus

type ReceiveMode string

const (
	ReceiveModePeekLock         = "peekLock"
	ReceiveModeReceiveAndDelete = "receiveAndDelete"
)

type SubQueue string

const (
	SubQueueNone       = ""
	SubQueueDeadLetter = "deadLetter"
	SubQueueTransfer   = "transferDeadLetter"
)

type ServiceBusReceiver struct{}

// used for batch APIs

// TODO: needs manual credit management.
// func (r *ServiceBusReceiver) ReceiveMessages(ctx context.Context, numMessages int) ([]*ServiceBusReceivedMessage, error) {
// 	return nil, nil
// }
