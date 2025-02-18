package tracing

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestPropagation(t *testing.T) {
	testCases := []struct {
		description  string
		message      *amqp.Message
		isNilMessage bool
	}{
		{
			description:  "nil message",
			message:      nil,
			isNilMessage: true,
		},
		{
			description: "non-nil message",
			message: &amqp.Message{
				Properties: &amqp.MessageProperties{
					MessageID: "message-id",
				},
			},
			isNilMessage: false,
		},
	}

	propagator := tracing.NewPropagator(tracing.PropagatorImpl{
		Inject: func(ctx context.Context, carrier tracing.Carrier) {
			carrier.Set("injected", "true")
		},
		Extract: func(ctx context.Context, carrier tracing.Carrier) context.Context {
			require.Zero(t, carrier.Get("badFlag"))
			return ctx
		},
	})

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			carrier := messageCarrierAdapter(tc.message)
			propagator.Inject(context.TODO(), carrier)
			propagator.Extract(context.TODO(), carrier)

			if !tc.isNilMessage {
				require.EqualValues(t, 1, len(carrier.Keys()))
				require.EqualValues(t, "true", carrier.Get("injected"))
				require.EqualValues(t, 1, len(tc.message.ApplicationProperties))
				require.EqualValues(t, "true", tc.message.ApplicationProperties["injected"])
			}
		})
	}
}
