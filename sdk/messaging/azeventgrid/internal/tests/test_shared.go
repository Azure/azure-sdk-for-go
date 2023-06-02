package tests

import "os"

type TestVars struct {
	Key          string
	Endpoint     string
	Topic        string
	Subscription string
}

func LoadEnv() TestVars {
	key := os.Getenv("EVENTGRID_KEY")
	ep := os.Getenv("EVENTGRID_ENDPOINT")
	topic := os.Getenv("EVENTGRID_TOPIC")
	sub := os.Getenv("EVENTGRID_SUBSCRIPTION")

	return TestVars{
		Key:          key,
		Endpoint:     ep,
		Topic:        topic,
		Subscription: sub,
	}
}
