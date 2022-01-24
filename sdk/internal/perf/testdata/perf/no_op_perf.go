package main

import "context"

type noOpPerfTest struct{}

func (n *noOpPerfTest) GlobalSetup(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) GlobalCleanup(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) Setup(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) Run(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) Cleanup(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) GetMetadata() string {
	return "NoOpTest"
}
