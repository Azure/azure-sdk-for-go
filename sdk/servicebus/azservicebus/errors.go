package azservicebus

import "errors"

var ErrSenderClosed = errors.New("sender is closed and cannot be used")
