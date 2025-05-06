package domain

import "context"

// Принимает msg и думает, что с ним делать
type MessageHandler func(msg interface{})

// Управление подпиской
type Subscription interface {
	// Отмена подписки на сообщество
	Unsubscribe()
}

// Участник сообщества
type SubPub interface {
	/* Я хочу подписаться на сообщество subject и управлять подпиской через Subscription
	При получении сообщений буду вызывать MessageHandler от сообщения */
	Subscribe(subject string, cb MessageHandler) (Subscription, error)
	// Я хочу отправить сообщение в сообщество subject
	Publish(subject string, msg interface{}) error
	// Я хочу удалить аккаунт
	Close(ctx context.Context) error
}
