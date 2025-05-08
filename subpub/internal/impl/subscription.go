package impl

import (
	"fmt"
	"sync"
	"time"

	"github.com/klimenkokayot/vk-internship/libs/logger"
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

	logger logger.Logger

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
			if !ok {
				es.logger.Info("Канал processing закрыт.")
				return
			}
			es.logger.Info("Получена публикация")
			es.callback(msg)
			es.logger.OK("Получена публикация")
			continue
		default:
		}
		// 2. Пополнение processing
		select {
		case msg, ok := <-es.buffer:
			if !ok {
				es.logger.Info("Канал buffer закрыт.")
				return
			}
			es.processing <- msg
			continue
		case <-time.After(idleDelay):
			// Защита от busy-wait
			continue
		}
	}
}

func (es *Subscription) send(msg interface{}) error {
	es.logger.Info("Новая публикация.")
	es.mu.Lock()
	if es.closed {
		es.logger.Warn("Subscription closed.")
		es.mu.Unlock()
		return fmt.Errorf("подписка закрыта")
	}
	es.mu.Unlock()
	select {
	case es.buffer <- msg:
	default:
		// переполнение буффера
		es.expandBuffer(msg)
	}
	es.logger.OK("Публикация отправлена в buffer.")
	return nil
}

func (es *Subscription) expandBuffer(msg interface{}) {
	es.mu.Lock()
	defer es.mu.Unlock()
	es.logger.Info("Увеличение буффера.",
		logger.Field{Key: "old_buffer_size", Value: len(es.buffer)},
	)
	newBuffer := make(chan interface{}, len(es.buffer)*2)
	for {
		select {
		case msg := <-es.buffer:
			newBuffer <- msg
		default:
			close(es.buffer)
			newBuffer <- msg
			es.buffer = newBuffer
			es.logger.OK("Буффер успешно учеличен.",
				logger.Field{Key: "new_buffer_size", Value: len(es.buffer)},
			)
			return
		}
	}
}

func newSubscription(id, topic string, callback domain.MessageHandler, bus *SubPub, logger logger.Logger) *Subscription {
	logger.Info("Инициализация экземпляра Subscription.")
	sub := &Subscription{
		id:       id,
		topic:    topic,
		callback: callback,

		processing: make(chan interface{}, processingSize),
		buffer:     make(chan interface{}, bufferSize),

		logger: logger,

		closed: false,
		mu:     sync.Mutex{},
		once:   sync.Once{},
		bus:    bus,
	}
	sub.wg.Add(1)
	go sub.process()
	logger.OK("Инициализация Subscription успешно выполнена.")
	return sub
}
