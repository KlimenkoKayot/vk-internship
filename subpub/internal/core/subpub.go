package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/klimenkokayot/vk-internship/libs/logger"
	"github.com/klimenkokayot/vk-internship/subpub/domain"
	"github.com/klimenkokayot/vk-internship/subpub/internal/infrastructure/uuid"
)

var (
	ErrSubPubClosed  = fmt.Errorf("SubPub closed")
	ErrPanicRecover  = fmt.Errorf("panic recover")
	ErrTopicNotExist = fmt.Errorf("topic does not exist")
)

type SubPub struct {
	topicSubscribes map[string]map[string]*Subscription

	uuidGenerator uuid.UUIDGenerator
	logger        logger.Logger

	closed bool
	once   sync.Once
	mu     sync.RWMutex
}

func (e *SubPub) removeSubscription(topic, id string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.logger.Debug(fmt.Sprintf("Удаление подписки %s из темы %s", id, topic))

	if subs, ok := e.topicSubscribes[topic]; ok {
		delete(subs, id)
		if len(subs) == 0 {
			delete(e.topicSubscribes, topic)
			e.logger.Debug(fmt.Sprintf("Тема %s удалена, так как больше не содержит подписок", topic))
		}
	} else {
		e.logger.Warn(fmt.Sprintf("Попытка удаления подписки из несуществующей темы %s", topic))
	}
}

func (e *SubPub) Subscribe(subject string, callback domain.MessageHandler) (subscription domain.Subscription, err error) {
	e.mu.Lock()
	defer func() {
		if r := recover(); r != nil {
			e.logger.Error(fmt.Sprintf("PANIC в Subscribe: %v", r))
			err = fmt.Errorf("%w: %w", ErrPanicRecover, r.(error))
		}
		defer e.mu.Unlock()
	}()

	e.logger.Info(fmt.Sprintf("Попытка подписки на тему %s", subject))

	if e.closed {
		e.logger.Warn("Отказ в подписке: сервис закрыт")
		return nil, ErrSubPubClosed
	}

	uid := e.uuidGenerator.NewString()
	loggerSub := e.logger.WithLayer(fmt.Sprintf("SUB-%s", uid))

	sub := newSubscription(uid, subject, callback, e, loggerSub)

	if _, ok := e.topicSubscribes[subject]; !ok {
		e.topicSubscribes[subject] = make(map[string]*Subscription, 1024)
		e.logger.Debug(fmt.Sprintf("Создана новая тема %s", subject))
	}

	e.topicSubscribes[subject][uid] = sub
	e.logger.OK(fmt.Sprintf("Успешная подписка %s на тему %s", uid, subject))

	return sub, nil
}

func (e *SubPub) Publish(subject string, msg interface{}) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	e.logger.Debug(fmt.Sprintf("Публикация сообщения в тему %s", subject))

	if e.closed {
		e.logger.Warn("Отказ в публикации: сервис закрыт")
		return ErrSubPubClosed
	}

	subs, ok := e.topicSubscribes[subject]
	if !ok {
		e.logger.Warn(fmt.Sprintf("Попытка публикации в несуществующую тему %s", subject))
		return ErrTopicNotExist
	}

	for id, subscriber := range subs {
		e.logger.Debug(fmt.Sprintf("Отправка сообщения подписчику %s темы %s", id, subject))
		if err := subscriber.send(msg); err != nil {
			e.logger.Warn(fmt.Sprintf("Ошибка отправки сообщения подписчику %s: %v", id, err))
		}
	}

	e.logger.OK(fmt.Sprintf("Сообщение успешно опубликовано в тему %s (%d подписчиков)", subject, len(subs)))
	return nil
}

func (e *SubPub) Close(ctx context.Context) error {
	e.logger.Info("Запуск процедуры закрытия SubPub")

	e.once.Do(func() {
		e.mu.Lock()
		e.closed = true

		subscribers := make([]*Subscription, 0)
		for topic, subs := range e.topicSubscribes {
			e.logger.Debug(fmt.Sprintf("Обработка темы %s (%d подписчиков)", topic, len(subs)))
			for id, sub := range subs {
				subscribers = append(subscribers, sub)
				e.logger.Debug(fmt.Sprintf("Добавлена подписка %s для отписки", id))
			}
		}
		e.mu.Unlock()

		e.logger.Info(fmt.Sprintf("Начата отписка %d подписчиков", len(subscribers)))
		for _, sub := range subscribers {
			sub.Unsubscribe()
			e.logger.Debug(fmt.Sprintf("Подписчик %s отписан", sub.id))
		}
		e.logger.OK("SubPub успешно закрыт")
		return
	})

	e.logger.Warn("Попытка повторного закрытия.")
	return nil
}

func NewSubPub(uuidGenerator uuid.UUIDGenerator, logger logger.Logger) domain.SubPub {
	logger.Info("Инициализация нового SubPub")
	topicSubscribes := make(map[string]map[string]*Subscription, 0)
	subPub := &SubPub{
		topicSubscribes: topicSubscribes,
		uuidGenerator:   uuidGenerator,
		logger:          logger,
		closed:          false,
		once:            sync.Once{},
		mu:              sync.RWMutex{},
	}

	logger.OK("SubPub успешно инициализирован")
	return subPub
}
