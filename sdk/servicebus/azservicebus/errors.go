package azservicebus

import "errors"

var ErrSenderClosed = errors.New("sender is closed and cannot be used")
var ErrReceiverClosed = errors.New("receiver has been closed and can no longer be used")
var ErrProcessorClosed = errors.New("processor has been closed and can no longer be used")
