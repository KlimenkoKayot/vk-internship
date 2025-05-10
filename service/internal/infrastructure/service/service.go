package service

import (
	"github.com/klimenkokayot/vk-internship/service/internal/domain"
	subpubdomain "github.com/klimenkokayot/vk-internship/subpub/domain"
)

type Service struct {
	pubsub subpubdomain.SubPub
}

func NewService(pubsub subpubdomain.SubPub) domain.PubSubService {
	return domain.NewService(pubsub)
}
