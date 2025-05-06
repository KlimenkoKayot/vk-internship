package subpub

import (
	"context"
	"sync"

	"github.com/klimenkokayot/vk-internship/domain"
)

// --------------------
// domain.Subscription
// --------------------

type subscription struct {
	id    string
	topic string
	ch    chan interface{}

	once sync.Once
	bus  *eventBus
}

func (es *subscription) Unsubscribe() {
	es.once.Do(func() {
		es.bus.removeSubscription(es.topic, es.id)
		close(es.ch)
	})
	// Отписываемся от топика
	close(es.ch)
}

func newSubscription(id, topic string, callback domain.MessageHandler, bus *eventBus) *subscription {
	ch := make(chan interface{}, 1)
	go func() {
		for val := range ch {
			callback(val)
		}
	}()
	return &subscription{
		id:    id,
		topic: topic,
		ch:    ch,

		once: sync.Once{},
		bus:  bus,
	}
}

// --------------------
// domain.SubPub
// --------------------

type eventBus struct {
	topicSubscribes map[string]map[string]*subscription
	uuidGenerator   domain.IDGenerator
	mu              sync.RWMutex
}

func (e *eventBus) removeSubscription(topic, id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.topicSubscribes[topic], id)
}

func (e *eventBus) Subscribe(subject string, callback domain.MessageHandler) (domain.Subscription, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	uid := e.uuidGenerator.NewString()
	subscription := newSubscription(uid, subject, callback, e)
	e.topicSubscribes[subject][uid] = subscription
	return subscription, nil
}

func (e *eventBus) Publish(subject string, msg interface{}) error {
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

func (e *eventBus) Close(ctx context.Context) error {
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
	return &eventBus{}
}
