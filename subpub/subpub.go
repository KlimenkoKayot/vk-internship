package subpub

import (
	"context"

	"github.com/klimenkokayot/vk-internship/domain"
)

type EventBus struct {
}

func (e *EventBus) Close(ctx context.Context) error {
	panic("unimplemented")
}

func (e *EventBus) Publish(subject string, msg interface{}) error {
	panic("unimplemented")
}

func (e *EventBus) Subscribe(subject string, cb domain.MessageHandler) (domain.Subscription, error) {
	panic("unimplemented")
}

func NewSubPub() domain.SubPub {
	return &EventBus{}
}
