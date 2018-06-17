package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func main() {
	closer, err := startOpenTracing()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer closer.Close()

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Initialize and create a Service Bus Queue named helloworld if it doesn't exist
	q, err := getQueue(ns, "helloworld")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Send the message "Hello World!" to the Queue named helloworld
	err = q.Send(context.Background(), servicebus.NewMessageFromString("Hello World!"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	done := make(chan interface{}, 1)
	// Receive the message "Hello World!" from queue named helloworld
	listenHandle, err := q.Receive(context.Background(),
		func(ctx context.Context, msg *servicebus.Message) servicebus.DispositionAction {
			fmt.Println(string(msg.Data))
			defer func(){
				done <- nil
			}()
			return msg.Complete()
		})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	// Close the listener handle for the Service Bus Queue
	defer listenHandle.Close(context.Background())

	// Wait for a signal to quit:
	<-done
}

func startOpenTracing() (io.Closer, error) {
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

	return cfg.InitGlobalTracer(
		"opentracing_example",
		config.Logger(jLogger),
	)
}

func getQueue(ns *servicebus.Namespace, queueName string) (*servicebus.Queue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	q, err := ns.NewQueue(ctx, queueName)
	return q, err
}