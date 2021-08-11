package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	ctx := context.Background()

	client, err := azservicebus.NewServiceBusClient(azservicebus.ServiceBusWithConnectionString(cs))
	if err != nil {
		log.Fatalf("Failed to create service bus client: %s", err.Error())
	}

	defer func() {
		err := client.Close(ctx)

		if err != nil {
			log.Printf("Error when closing client: %s", err.Error())
		}
	}()

	var process = true

	if process {
		processor, err := client.NewProcessor(azservicebus.ProcessorWithQueue("samples"))

		if err != nil {
			log.Fatalf("Failed when creating receiver: %s", err.Error())
		}

		ctr := int64(0)

		go func() {
			err = processor.Start(func(msg *azservicebus.ReceivedMessage) error {
				currCtr := atomic.AddInt64(&ctr, 1)
				log.Printf("[%d] Got message: ID: %s, Body: %s", currCtr, msg.ID, string(msg.Body))

				// If you don't disposition you don't get a new message.
				// if err := receiver.CompleteMessage(subscribeCtx, msg); err != nil {
				// 	log.Printf("Failed to complete message with ID %s: %s", msg.ID, err.Error())
				// }

				// auto complete is on by default.
				return nil
			}, func(err error) {
				log.Printf("Demo: got error: %s", err.Error())
			})

			if err != nil {
				log.Fatalf("Processor failed to run: %s", err.Error())
			}

			<-processor.Done()

			log.Printf("Processor has closed")
		}()
	} else {

		go func() {
			sender, err := client.NewSender("samples")

			if err != nil {
				log.Fatalf("Failed to create the sender: %s", err.Error())
			}

			ticker := time.NewTicker(time.Second)

			for t := range ticker.C {
				err := sender.SendMessage(context.Background(), &azservicebus.Message{
					Body: []byte(fmt.Sprintf("hello world: %s", t.String())),
				})

				if err != nil {
					log.Printf("Failed to send message: %s", err.Error())
				}
			}
		}()
	}

	ch := make(chan struct{})
	<-ch
}
