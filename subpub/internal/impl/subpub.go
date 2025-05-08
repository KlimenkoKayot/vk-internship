package impl

import (
	"context"
	"fmt"
	"sync"

	"github.com/klimenkokayot/vk-internship/libs/logger"
	"github.com/klimenkokayot/vk-internship/subpub/domain"
)

var (
	ErrSubPubClosed  = fmt.Errorf("SubPub closed")
	ErrPanicRecover  = fmt.Errorf("panic recover")
	ErrTopicNotExist = fmt.Errorf("topic does not exist")
)

type SubPub struct {
	topicSubscribes map[string]map[string]*Subscription

	uuidGenerator domain.UUIDGenerator
	logger        logger.Logger

	closed bool
	once   sync.Once
	mu     sync.RWMutex
}

func (e *SubPub) removeSubscription(topic, id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if subs, ok := e.topicSubscribes[topic]; ok {
		delete(subs, id)
		if len(topic) == 0 {
			delete(e.topicSubscribes, topic)
		}
	}
}

func (e *SubPub) Subscribe(subject string, callback domain.MessageHandler) (subscription domain.Subscription, err error) {
	e.mu.Lock()
	if e.closed {
		return nil, ErrSubPubClosed
	}
	defer func() {
		if err := recover(); err != nil {
			err = ErrPanicRecover
		}
		e.mu.Unlock()
	}()
	uid := e.uuidGenerator.NewString()
	sub := newSubscription(uid, subject, callback, e)
	if e.topicSubscribes[subject] == nil {
		e.topicSubscribes[subject] = make(map[string]*Subscription, 1024)
	}
	e.topicSubscribes[subject][uid] = sub
	return sub, nil
}

func (e *SubPub) Publish(subject string, msg interface{}) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.closed {
		return ErrSubPubClosed
	}
	if _, ok := e.topicSubscribes[subject]; ok {
		return ErrTopicNotExist
	}
	for _, subscriber := range e.topicSubscribes[subject] {
		subscriber.send(msg) // резона ошибку фиксить нет, все равно скоро удалится из map в removeSubscriber()
	}
	return nil
}

func (e *SubPub) Close(ctx context.Context) error {
	e.once.Do(func() {
		e.mu.Lock()
		// Меняем флаг
		e.closed = true
		// Копируем все подписки
		subscribers := make([]*Subscription, 0)
		for _, topic := range e.topicSubscribes {
			for _, sub := range topic {
				subscribers = append(subscribers, sub)
			}
		}
		e.mu.Unlock()
		for _, sub := range subscribers {
			sub.Unsubscribe()
		}
	})
	return nil
}

func NewSubPub(uuidGenerator domain.UUIDGenerator, logger logger.Logger) domain.SubPub {
	subPub := &SubPub{
		topicSubscribes: make(map[string]map[string]*Subscription, 1024),

		uuidGenerator: uuidGenerator,
		logger:        logger,

		closed: false,
		once:   sync.Once{},
		mu:     sync.RWMutex{},
	}
	return subPub
}
