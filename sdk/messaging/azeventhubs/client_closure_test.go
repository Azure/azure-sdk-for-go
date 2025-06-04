package azeventhubs

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientClosureBehavior(t *testing.T) {
	t.Run("ConsumerClient operations return ErrorCodeClientClosed after close", func(t *testing.T) {
		// Using fake connection string for unit testing
		connectionString := "Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=test;SharedAccessKey=dGVzdA==;EntityPath=test-eventhub"
		
		consumerClient, err := NewConsumerClientFromConnectionString(connectionString, "", DefaultConsumerGroup, nil)
		require.NoError(t, err)

		// Close the client
		err = consumerClient.Close(context.Background())
		require.NoError(t, err)

		// Test that GetEventHubProperties returns the correct error after close
		_, err = consumerClient.GetEventHubProperties(context.Background(), nil)
		require.Error(t, err, "GetEventHubProperties should return an error after client is closed")

		// Check that the error has the correct ErrorCode
		var ehErr *Error
		if errors.As(err, &ehErr) {
			require.Equal(t, ErrorCodeClientClosed, ehErr.Code)
		}
	})

	t.Run("ProducerClient operations return ErrorCodeClientClosed after close", func(t *testing.T) {
		// Using fake connection string for unit testing
		connectionString := "Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=test;SharedAccessKey=dGVzdA==;EntityPath=test-eventhub"
		
		producerClient, err := NewProducerClientFromConnectionString(connectionString, "", nil)
		require.NoError(t, err)

		// Close the client
		err = producerClient.Close(context.Background())
		require.NoError(t, err)

		// Test that GetEventHubProperties returns the correct error after close
		_, err = producerClient.GetEventHubProperties(context.Background(), nil)
		require.Error(t, err, "GetEventHubProperties should return an error after client is closed")

		// Check that the error has the correct ErrorCode
		var ehErr *Error
		if errors.As(err, &ehErr) {
			require.Equal(t, ErrorCodeClientClosed, ehErr.Code)
		}
	})
}