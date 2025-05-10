package domain

import (
	"context"

	subpubdomain "github.com/klimenkokayot/vk-internship/subpub/domain"
)

type Subscription interface {
	Unsubscribe()
	Chan() <-chan string
}

type PubSubService interface {
	Subscribe(ctx context.Context, key string) (Subscription, error)
	Publish(ctx context.Context, key, data string) error
}

type Service struct {
	pubsub subpubdomain.SubPub
}

func NewService(pubsub subpubdomain.SubPub) *Service {
	return &Service{
		pubsub: pubsub,
	}
}

func (s *Service) Subscribe(ctx context.Context, key string) (Subscription, error) {
	ch := make(chan string, 32)

	sub, err := s.pubsub.Subscribe(key, func(msg interface{}) {
		if str, ok := msg.(string); ok {
			select {
			case ch <- str:
			case <-ctx.Done():
			}
		}
	})
	if err != nil {
		return nil, err
	}

	return &subscription{
		Subscription: sub,
		ch:           ch,
	}, nil
}

func (s *Service) Publish(ctx context.Context, key, data string) error {
	return s.pubsub.Publish(key, data)
}

type subscription struct {
	subpubdomain.Subscription
	ch chan string
}

func (s *subscription) Chan() <-chan string {
	return s.ch
}

func (s *subscription) Unsubscribe() {
	close(s.ch)
	s.Subscription.Unsubscribe()
}
