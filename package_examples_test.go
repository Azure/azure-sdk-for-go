package servicebus_test

import (
	"context"
	"fmt"
	"os"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func Example_helloWorld() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connStr := mustGetenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println(err)
		return
	}

	const queueName = "helloworld"
	q, err := getQueue(ctx, ns, queueName)
	if err != nil {
		fmt.Printf("failed to build a new queue named %q\n", queueName)
		return
	}

	errs := make(chan error, 2)

	messages := []string{"hello", "world"}

	// Start a receiver that will print messages that it is informed of by Service Bus.
	go func(ctx context.Context, client *servicebus.Queue, quitAfter int) {
		received := make(chan struct{})

		listenHandle, err := client.Receive(ctx, func(ctx context.Context, message *servicebus.Message) servicebus.DispositionAction {
			fmt.Println(string(message.Data))
			received <- struct{}{}
			return message.Complete()
		})
		if err != nil {
			errs <- err
			return
		}
		defer listenHandle.Close(context.Background())
		defer fmt.Println("...no longer listening")
		fmt.Println("listening...")

		for i := 0; i < quitAfter; i++ {
			select {
			case <-received:
				// Intentionally Left Blank
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			case <-listenHandle.Done():
				errs <- listenHandle.Err()
				return
			}
		}
		errs <- nil
	}(ctx, q, len(messages))

	// Publish messages to Service Bus so that the receiver defined above will print them
	go func(ctx context.Context, client *servicebus.Queue, messages ...string) {
		for i := range messages {
			messageSent := make(chan error, 1)

			go func() {
				messageSent <- client.Send(ctx, servicebus.NewMessageFromString(messages[i]))
			}()
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			case err := <-messageSent:
				if err != nil {
					errs <- err
					return
				}
			}
		}
		errs <- nil
	}(ctx, q, messages...)

	for i := 0; i < 2; i++ {
		select {
		case err := <-errs:
			if err != nil {
				fmt.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}

	// Output:
	// listening...
	// hello
	// world
	// ...no longer listening
}

func Example_opentracing() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: "0.0.0.0:6831",
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger

	otHandle, err := cfg.InitGlobalTracer(
		"opentracing_example",
		config.Logger(jLogger),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer otHandle.Close()

	// Create clients for communicating with Service Bus
	connStr := mustGetenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println(err)
		return
	}

	const queueName = "opentracing"
	q, err := getQueue(ctx, ns, queueName)
	if err != nil {
		fmt.Printf("failed to build a new queue named %q\n", queueName)
		return
	}

	errs := make(chan error, 2)

	messages := []string{"diagnose problems", "using tracing", "in development", "or in production"}

	// Start a receiver that will print messages that it is informed of by Service Bus.
	go func(ctx context.Context, client *servicebus.Queue, quitAfter int) {
		received := make(chan struct{})

		listenHandle, err := client.Receive(ctx, func(ctx context.Context, message *servicebus.Message) servicebus.DispositionAction {
			fmt.Println(string(message.Data))
			received <- struct{}{}
			return message.Complete()
		})
		if err != nil {
			errs <- err
			return
		}
		defer listenHandle.Close(context.Background())

		for i := 0; i < quitAfter; i++ {
			select {
			case <-received:
				// Intentionally Left Blank
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			case <-listenHandle.Done():
				errs <- listenHandle.Err()
				return
			}
		}
		errs <- nil
	}(ctx, q, len(messages))

	// Publish messages to Service Bus so that the receiver defined above will print them
	go func(ctx context.Context, client *servicebus.Queue, messages ...string) {
		for i := range messages {
			messageSent := make(chan error, 1)

			go func() {
				messageSent <- client.Send(ctx, servicebus.NewMessageFromString(messages[i]))
			}()
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			case err := <-messageSent:
				if err != nil {
					errs <- err
					return
				}
			}
		}
		errs <- nil
	}(ctx, q, messages...)

	// Wait for both the publisher and receiver to declare themselves done.
	for i := 0; i < 2; i++ {
		select {
		case err := <-errs:
			if err != nil {
				fmt.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}

	// Output:
	// diagnose problems
	// using tracing
	// in development
	// or in production
}

func getQueue(ctx context.Context, ns *servicebus.Namespace, queueName string) (*servicebus.Queue, error) {
	qm := ns.NewQueueManager()
	qe, err := qm.Get(ctx, queueName)
	if err != nil {
		return nil, err
	}

	if qe == nil {
		_, err := qm.Put(ctx, queueName)
		if err != nil {
			return nil, err
		}
	}

	q, err := ns.NewQueue(queueName)
	return q, err
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("Environment variable '" + key + "' required for integration tests.")
	}
	return v
}
