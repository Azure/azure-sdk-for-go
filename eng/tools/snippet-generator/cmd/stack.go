//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cmd

import (
	"container/list"
	"fmt"
)

type Stack struct {
	l *list.List
}

func (s *Stack) Push(value string) {
	s.l.PushBack(value)
}

func NewStack() *Stack {
	return &Stack{
		l: list.New(),
	}
}

func (s *Stack) Pop() (string, error) {
	e := s.l.Back()
	if e == nil {
		return "", fmt.Errorf("the stack is empty") // we get nothing in the stack
	}
	return s.l.Remove(e).(string), nil
}

func (s *Stack) Peek() (string, error) {
	e := s.l.Back()
	if e == nil {
		return "", fmt.Errorf("the stack is empty")
	}
	return e.Value.(string), nil
}

func (s *Stack) Len() int {
	return s.l.Len()
}
