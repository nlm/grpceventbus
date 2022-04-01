package main

import "sync"

type Topic[T any] struct {
	subscribers []chan T
}

type PubSub[T any] struct {
	m      sync.RWMutex
	topics map[string]Topic[T]
}

func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		topics: make(map[string]Topic[T]),
	}
}

func (ps *PubSub[T]) Publish(channel string, elt *T) error {
	return nil
}

func (ps *PubSub[T]) Subscribe(channel string) (chan<- T, error) {
	ps.m.Lock()
	defer ps.m.Unlock()
	return make(chan T), nil
}

func (ps *PubSub[T]) Close() {
}
