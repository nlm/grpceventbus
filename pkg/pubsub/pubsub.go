package pubsub

import (
	"log"
	"sync"
)

type Topic[T any] struct {
	name string
	mu   sync.RWMutex
	subs map[chan T]struct{}
}

type Subscription[T any] struct {
	topic *Topic[T]
	c     chan T
}

func (s Subscription[T]) C() <-chan T {
	return s.c
}

type PubSub[T any] struct {
	mu          sync.RWMutex
	topics      map[string]*Topic[T]
	logger      *log.Logger
	chanBufSize int
}

func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		logger: log.New(
			log.Default().Writer(),
			"[PubSub]",
			log.Default().Flags(),
		),
		chanBufSize: 1,
		topics:      make(map[string]*Topic[T]),
	}
}

func (ps *PubSub[T]) Publish(topic string, elt T) {
	ps.logger.Println("publish:", topic, elt)
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	var queued uint
	if t, ok := ps.topics[topic]; ok {
		func() {
			t.mu.RLock()
			defer t.mu.RUnlock()
			for sub := range t.subs {
				select {
				case sub <- elt:
					//ps.logger.Println("publish: message queued")
					queued++
				default:
					ps.logger.Println("publish: message skipped")
				}
			}
		}()
		ps.logger.Println("publish: messages queued:", queued)
	} else {
		ps.logger.Println("publish: no subscribers")
	}
}

func (ps *PubSub[T]) Subscribe(topic string) *Subscription[T] {
	ps.logger.Println("subscribe:", topic)
	newchan := make(chan T, ps.chanBufSize)
	ps.mu.Lock()
	defer ps.mu.Unlock()
	t, ok := ps.topics[topic]
	if ok {
		t.mu.Lock()
		defer t.mu.Unlock()
		t.subs[newchan] = struct{}{}
	} else {
		t = &Topic[T]{
			name: topic,
			subs: map[chan T]struct{}{newchan: {}},
		}
		ps.logger.Println("new topic:", t.name)
		ps.topics[topic] = t
	}
	return &Subscription[T]{
		topic: t,
		c:     newchan,
	}
}

func (ps *PubSub[T]) Unsubscribe(s *Subscription[T]) {
	ps.logger.Println("unsubscribe:", s.topic.name)
	func() {
		s.topic.mu.Lock()
		defer s.topic.mu.Unlock()
		delete(s.topic.subs, s.c)
		close(s.c)
	}()
	if len(s.topic.subs) == 0 {
		ps.logger.Println("delete topic:", s.topic.name)
		ps.mu.Lock()
		defer ps.mu.Unlock()
		delete(ps.topics, s.topic.name)
	}
}

func (ps *PubSub[T]) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	for _, t := range ps.topics {
		t.mu.Lock()
		for sub := range t.subs {
			close(sub)
		}
	}
	ps.topics = make(map[string]*Topic[T])
}
