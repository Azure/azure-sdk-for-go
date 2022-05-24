# Hello World Producer / Consumer

This example illustrates a producer sending message round-robbin into an Event Hub instance. The consumer
receives from each Event Hub partition, and outputs the message it receives. Upon entering 'exit' into the
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
	hub, _ := initHub()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		hub.Send(context.Background(), &amqp.Message{Data: []byte(text)})
		if text == "exit\n" {
			break
		}
	}
}
```

## Consumer
```go
func main() {
	hub, partitions := initHub()
	exit := make(chan struct{})

	handler := func(ctx context.Context, msg *amqp.Message) error {
		text := string(msg.Data)
		if text == "exit\n" {
			fmt.Println("Someone told me to exit!")
			exit <- *new(struct{})
		} else {
			fmt.Println(string(msg.Data))
		}
		return nil
	}

	for _, partitionID := range *partitions {
		hub.Receive(partitionID, handler)
	}

	select {
	case <-exit:
		fmt.Println("closing after 2 seconds")
		select {
		case <-time.After(2 * time.Second):
			return
		}
	}
}
```