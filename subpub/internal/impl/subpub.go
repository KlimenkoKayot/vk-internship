package impl

import (
	"context"
	"log"
	"sync"

	"github.com/klimenkokayot/vk-internship/subpub/domain"
)

type SubPub struct {
	topicSubscribes map[string]map[string]*Subscription

	uuidGenerator domain.UUIDGenerator

	once sync.Once
	mu   sync.RWMutex
}

// OK
func (e *SubPub) removeSubscription(topic, id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if subs, ok := e.topicSubscribes[topic]; ok {
		delete(subs, id)
	}
}

func (e *SubPub) Subscribe(subject string, callback domain.MessageHandler) (domain.Subscription, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	uid := e.uuidGenerator.NewString()
	subscription := newSubscription(uid, subject, callback, e)
	if e.topicSubscribes[subject] == nil {
		e.topicSubscribes[subject] = make(map[string]*Subscription, 1024)
	}
	e.topicSubscribes[subject][uid] = subscription
	return subscription, nil
}

func (e *SubPub) Publish(subject string, msg interface{}) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, subscriber := range e.topicSubscribes[subject] {
		subscriber.send(msg) // резона ошибку фиксить нет, все равно скоро удалится из map в removeSubscriber()
	}
	return nil
}

func (e *SubPub) Close(ctx context.Context) error {
	e.once.Do(func() {
		e.mu.RLock()
		defer e.mu.RUnlock()
		for _, topic := range e.topicSubscribes {
			for _, subscriber := range topic {
				subscriber.Unsubscribe()
			}
		}
	})
	return nil
}

func NewSubPub(uuidGenerator domain.UUIDGenerator) domain.SubPub {
	subPub := &SubPub{
		topicSubscribes: make(map[string]map[string]*Subscription, 1024),
		uuidGenerator:   uuidGenerator,
		mu:              sync.RWMutex{},
	}
	return subPub
}
