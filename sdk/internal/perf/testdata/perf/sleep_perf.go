package main

import (
	"context"
	"math"
	"time"
)

type sleepPerfTest struct {
	InstanceCount int
	secondsPerOp  int
}

func (s *sleepPerfTest) GlobalSetup(ctx context.Context) error {
	s.InstanceCount += 1
	s.secondsPerOp = int(math.Pow(2.0, float64(s.InstanceCount)))
	return nil
}

func (s *sleepPerfTest) GlobalCleanup(ctx context.Context) error {
	return nil
}

func (s *sleepPerfTest) Setup(ctx context.Context) error {
	return nil
}

func (s *sleepPerfTest) Run(ctx context.Context) error {
	time.Sleep(time.Duration(s.secondsPerOp) * time.Second)
	return nil
}

func (s *sleepPerfTest) Cleanup(ctx context.Context) error {
	return nil
}

func (s *sleepPerfTest) GetMetadata() string {
	return "SleepTest"
}
