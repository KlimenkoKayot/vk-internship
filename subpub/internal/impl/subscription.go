package impl

import (
	"fmt"
	"sync"
	"time"

	"github.com/klimenkokayot/vk-internship/subpub/domain"
)

const (
	idleDelay      time.Duration = 10 * time.Millisecond
	processingSize uint          = 16
	bufferSize     uint          = 32
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
		// Безопасное переключение флага
		es.mu.Lock()
		es.closed = true
		es.mu.Unlock()

		// Выполнение оставшихся операций
		for msg := range es.processing {
			es.callback(msg)
		}
		for msg := range es.buffer {
			es.callback(msg)
		}

		// Удаление подписки, закрытие каналов
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
		// 1. Обработка
		select {
		case msg, ok := <-es.processing:
			es.callback(msg)
			if !ok {
				return
			}
			continue
		default:
		}
		// 2. Пополнение processing
		select {
		case msg, ok := <-es.buffer:
			es.processing <- msg
			if !ok {
				return
			}
			continue
		case <-time.After(idleDelay):
			// Защита от busy-wait
			continue
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
		es.expandBuffer(msg)
		return nil
	}
}

func (es *Subscription) expandBuffer(msg interface{}) {
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
		processing: make(chan interface{}, processingSize),
		buffer:     make(chan interface{}, bufferSize),
		closed:     false,
		mu:         sync.Mutex{},
		once:       sync.Once{},
		bus:        bus,
	}
	sub.wg.Add(1)
	go sub.process()
	return sub
}
