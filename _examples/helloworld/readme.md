# Hello World Producer / Consumer

This example illustrates a producer sending messages into a Service Bus FIFO Queue. The consumer
receives from each message in FIFO order from the queue, and outputs the message it receives. Upon entering 'exit' into the
producer, the producer will send, then exit, and the receiver will receive the message and close.

## To Run
- from this directory execute `make`
- open two terminal windows
  - in the first terminal, execute `./bin/consumer`
  - in the second terminal, execute `./bin/producer`
  - in the second terminal, type some works and press enter
- see the words you typed in the second terminal in the first
- type 'exit' in the second terminal when you'd like to end your session

## Producer
```go
func main() {
	// Connect
    connStr := mustGetenv("SERVICEBUS_CONNECTION_STRING")
    ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    queueName := "helloworld"
    // Create the queue if it doesn't exist
    err = ensureQueue(ns, queueName)
    q := ns.NewQueue(queueName)
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("Enter text: ")
        text, _ := reader.ReadString('\n')
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        q.Send(ctx, servicebus.NewEventFromString(text))
        if text == "exit\n" {
            break
        }
        cancel()
    }
}
```

## Consumer
```go
func main() {
	// Connect
    connStr := mustGetenv("SERVICEBUS_CONNECTION_STRING")
    ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    queueName := "helloworld"
    // Create the queue if it doesn't exist
    err = ensureQueue(ns, queueName)
    q := ns.NewQueue(queueName)

    // Start listening to events on the queue
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    exit := make(chan struct{})
    listenHandle, err := q.Receive(ctx, func(ctx context.Context, event *servicebus.Event) error {
        text := string(event.Data)
        if text == "exit\n" {
            fmt.Println("Oh snap!! Someone told me to exit!")
            exit <- *new(struct{})
        } else {
            fmt.Println(string(event.Data))
        }
        return nil
    })
    defer listenHandle.Close(context.Background())
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println("I am listening...")

    select {
    case <-exit:
        fmt.Println("closing after 2 seconds")
        select {
        case <-time.After(2 * time.Second):
            listenHandle.Close(context.Background())
            return
        }
    }
}
```