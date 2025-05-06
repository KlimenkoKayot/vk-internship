package impl

import (
	"sync"

	"github.com/klimenkokayot/vk-internship/subpub/domain"
)

type Subscription struct {
	id    string
	topic string
	ch    chan interface{}

	once sync.Once
	bus  *SubPub
}

func (es *Subscription) Unsubscribe() {
	es.once.Do(func() {
		es.bus.removeSubscription(es.topic, es.id)
		close(es.ch)
	})
	// Отписываемся от топика
	close(es.ch)
}

func newSubscription(id, topic string, callback domain.MessageHandler, bus *SubPub) *Subscription {
	ch := make(chan interface{}, 1)
	go func() {
		for val := range ch {
			callback(val)
		}
	}()
	return &Subscription{
		id:    id,
		topic: topic,
		ch:    ch,

		once: sync.Once{},
		bus:  bus,
	}
}
