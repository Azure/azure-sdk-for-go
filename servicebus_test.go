package servicebus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"os"
)

func TestServiceBus(t *testing.T) {
	connStr := os.Getenv("AZURE_SERVICE_BUS_CONN_STR") // `Endpoint=sb://XXXX.servicebus.windows.net/;SharedAccessKeyName=XXXX;SharedAccessKey=XXXX`
	sb, err := New(connStr)
	defer sb.Close()
	assert.Nil(t, err)
	assert.NotNil(t, sb)
}
