package subpub

import (
	"github.com/klimenkokayot/vk-internship/subpub/domain"
	"github.com/klimenkokayot/vk-internship/subpub/internal/impl"
)

type (
	MessageHandler = domain.MessageHandler
	SubPub         = domain.SubPub
	Subscription   = domain.Subscription
)

func NewSubPub() SubPub {
	return impl.NewSubPub()
}
