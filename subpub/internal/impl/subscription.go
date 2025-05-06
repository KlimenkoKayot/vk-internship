package impl

import (
	"fmt"
	"sync"

	"github.com/klimenkokayot/vk-internship/subpub/domain"
)

type Subscription struct {
	id       string
	topic    string
	callback domain.MessageHandler

	processing chan interface{}
	buffer     chan interface{}

	closed bool
	wg     sync.WaitGroup
	mu     sync.Mutex
	once   sync.Once
	bus    *SubPub
}

func (es *Subscription) Unsubscribe() {
	es.once.Do(func() {
		es.bus.removeSubscription(es.topic, es.id)
		close(es.processing)
		close(es.buffer)
		es.wg.Wait()
	})
}

func (es *Subscription) process() {
	defer es.wg.Done()
	// Сначала обрабатываем первостепенные задачи из processing
	// Если их нет, то стараемся перелить задачи из buffer в processing
	for {
		select {
		case msg, ok := <-es.processing:
			if !ok {
				return
			}
			es.callback(msg)
		case msg, ok := <-es.buffer:
			if !ok {
				return
			}
			es.processing <- msg
		}
	}
}

func (es *Subscription) send(msg interface{}) error {
	es.mu.Lock()
	if es.closed {
		es.mu.Unlock()
		return fmt.Errorf("подписка закрыта")
	}
	es.mu.Unlock()
	select {
	case es.buffer <- msg:
		return nil
	default:
		// переполнение буффера
		es.channelTransfusion(msg)
		return nil
	}
}

func (es *Subscription) channelTransfusion(msg interface{}) {
	es.mu.Lock()
	defer es.mu.Unlock()
	newBuffer := make(chan interface{}, len(es.buffer)*2)
	for {
		select {
		case msg := <-es.buffer:
			newBuffer <- msg
		default:
			close(es.buffer)
			newBuffer <- msg
			es.buffer = newBuffer
			return
		}
	}
}

func newSubscription(id, topic string, callback domain.MessageHandler, bus *SubPub) *Subscription {
	sub := &Subscription{
		id:         id,
		topic:      topic,
		callback:   callback,
		processing: make(chan interface{}, 16),
		buffer:     make(chan interface{}, 1024),
		closed:     false,
		mu:         sync.Mutex{},
		once:       sync.Once{},
		bus:        bus,
	}
	sub.wg.Add(1)
	go sub.process()
	return sub
}
