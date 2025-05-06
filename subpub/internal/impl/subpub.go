package impl

import (
	"context"
	"sync"

	"github.com/klimenkokayot/vk-internship/subpub/domain"
)

type SubPub struct {
	topicSubscribes map[string]map[string]*Subscription
	uuidGenerator   domain.UUIDGenerator
	mu              sync.RWMutex
}

func (e *SubPub) removeSubscription(topic, id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.topicSubscribes[topic], id)
}

func (e *SubPub) Subscribe(subject string, callback domain.MessageHandler) (domain.Subscription, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	uid := e.uuidGenerator.NewString()
	subscription := newSubscription(uid, subject, callback, e)
	e.topicSubscribes[subject][uid] = subscription
	return subscription, nil
}

func (e *SubPub) Publish(subject string, msg interface{}) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, subscriber := range e.topicSubscribes[subject] {
		select {
		case subscriber.ch <- msg:
		default:
			// скип
		}
	}
	return nil
}

func (e *SubPub) Close(ctx context.Context) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, topic := range e.topicSubscribes {
		for _, subscriber := range topic {
			subscriber.Unsubscribe()
		}
	}
	return nil
}

func NewSubPub() domain.SubPub {
	return &SubPub{}
}
