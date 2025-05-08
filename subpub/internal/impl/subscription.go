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
		es.logger.Info("Начало процедуры отписки")

		// Безопасное переключение флага
		es.mu.Lock()
		es.closed = true
		es.mu.Unlock()
		es.logger.Debug("Флаг closed установлен")

		// Закрытие каналов
		close(es.processing)
		close(es.buffer)
		es.logger.Debug("Каналы processing и buffer закрыты")

		// Обработка оставшихся сообщений
		es.logger.Debug("Обработка оставшихся сообщений")
		processedCount := 0
		for msg := range es.processing {
			es.callback(msg)
			processedCount++
		}
		for msg := range es.buffer {
			es.callback(msg)
			processedCount++
		}
		es.logger.Info("Обработаны оставшиеся сообщения",
			logger.Field{Key: "processed_messages", Value: processedCount},
		)

		// Удаление подписки
		es.bus.removeSubscription(es.topic, es.id)
		es.logger.Debug("Подписка удалена из bus")

		// Ожидание завершения обработки
		es.wg.Wait()
		es.logger.OK("Отписка завершена успешно")
	})
}

func (es *Subscription) process() {
	defer func() {
		if r := recover(); r != nil {
			es.logger.Error("Паника в обработчике сообщений",
				logger.Field{Key: "error", Value: r},
			)
		}
		es.wg.Done()
		es.logger.Debug("Горутина обработчика завершена")
	}()

	es.logger.Debug("Запуск обработчика сообщений")
	for {
		// Проверяем флаг закрытия
		es.mu.Lock()
		if es.closed {
			es.mu.Unlock()
			es.logger.Debug("Подписка закрыта, завершение работы")
			return
		}
		es.mu.Unlock()

		// 1. Обработка сообщений из processing
		select {
		case msg, ok := <-es.processing:
			if !ok {
				es.logger.Debug("Канал processing закрыт, завершение работы")
				return
			}
			es.logger.Info("Начало обработки сообщения из processing")
			startTime := time.Now()

			es.callback(msg)

			es.logger.Info("Сообщение обработано",
				logger.Field{Key: "processing_time", Value: time.Since(startTime)},
			)
			continue
		default:
		}

		// 2. Пополнение processing из buffer
		es.mu.Lock()
		if es.closed {
			es.mu.Unlock()
			es.logger.Debug("Подписка закрыта, завершение работы")
			return
		}
		select {
		case msg, ok := <-es.buffer:
			if !ok {
				es.logger.Info("Канал buffer закрыт, завершение работы")
				es.mu.Unlock()
				return
			}

			es.logger.Info("Перемещение сообщения из buffer в processing")
			select {
			case es.processing <- msg:
				es.logger.Info("Сообщение успешно перемещено в processing")
			default:
				es.logger.Warn("Не удалось переместить сообщение в processing (канал полон)")
			}
		case <-time.After(idleDelay):
		}
		es.mu.Unlock()
	}
}

func (es *Subscription) send(msg interface{}) error {
	es.logger.Debug("Попытка отправки сообщения")
	startTime := time.Now()

	es.mu.Lock()
	if es.closed {
		es.mu.Unlock()
		es.logger.Warn("Отказ в отправке: подписка закрыта")
		return fmt.Errorf("подписка закрыта")
	}
	defer es.mu.Unlock()

	select {
	case es.buffer <- msg:
		es.logger.Debug("Сообщение успешно помещено в buffer",
			logger.Field{Key: "buffer_size", Value: len(es.buffer)},
			logger.Field{Key: "processing_time", Value: time.Since(startTime)},
		)
		return nil
	default:
		es.logger.Warn("Буфер переполнен, попытка расширения",
			logger.Field{Key: "current_buffer_size", Value: cap(es.buffer)},
		)
		es.expandBuffer(msg)
		es.logger.Debug("Сообщение обработано после расширения буфера",
			logger.Field{Key: "new_buffer_size", Value: cap(es.buffer)},
			logger.Field{Key: "total_processing_time", Value: time.Since(startTime)},
		)
		return nil
	}
}

func (es *Subscription) expandBuffer(msg interface{}) {
	es.mu.Lock()
	defer es.mu.Unlock()

	oldSize := cap(es.buffer)
	newSize := oldSize * 2
	if newSize == 0 {
		newSize = int(bufferSize)
	}

	es.logger.Info("Расширение буфера",
		logger.Field{Key: "old_capacity", Value: oldSize},
		logger.Field{Key: "new_capacity", Value: newSize},
	)

	newBuffer := make(chan interface{}, newSize)
	messageCount := 0

	// Перенос существующих сообщений
	for {
		select {
		case m := <-es.buffer:
			newBuffer <- m
			messageCount++
		default:
			// Добавление нового сообщения
			newBuffer <- msg
			messageCount++

			// Замена буфера
			close(es.buffer)
			es.buffer = newBuffer

			es.logger.Info("Буфер успешно расширен",
				logger.Field{Key: "messages_transferred", Value: messageCount},
				logger.Field{Key: "new_capacity", Value: cap(es.buffer)},
			)
			return
		}
	}
}

func newSubscription(id, topic string, callback domain.MessageHandler, bus *SubPub, log logger.Logger) *Subscription {
	log = log.WithFields(
		logger.Field{Key: "subscription_id", Value: id},
		logger.Field{Key: "topic", Value: topic},
	)

	log.Info("Создание новой подписки",
		logger.Field{Key: "processing_size", Value: processingSize},
		logger.Field{Key: "initial_buffer_size", Value: bufferSize},
	)

	sub := &Subscription{
		id:         id,
		topic:      topic,
		callback:   callback,
		processing: make(chan interface{}, processingSize),
		buffer:     make(chan interface{}, bufferSize),
		logger:     log,
		closed:     false,
		mu:         sync.Mutex{},
		once:       sync.Once{},
		bus:        bus,
	}

	sub.wg.Add(1)
	go sub.process()

	log.OK("Подписка успешно создана",
		logger.Field{Key: "processing_capacity", Value: cap(sub.processing)},
		logger.Field{Key: "buffer_capacity", Value: cap(sub.buffer)},
	)
	return sub
}
