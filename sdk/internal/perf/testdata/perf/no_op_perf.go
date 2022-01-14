package main

import "context"

type noOpPerfTest struct{}

func (n *noOpPerfTest) GlobalSetup(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) GlobalTearDown(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) Setup(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) Run(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) TearDown(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) GetMetadata() string {
	return "NoOpTest"
}
