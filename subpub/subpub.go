package subpub

import (
	"log"

	"github.com/klimenkokayot/vk-internship/subpub/domain"
	"github.com/klimenkokayot/vk-internship/subpub/internal/impl"
	"github.com/klimenkokayot/vk-internship/subpub/internal/infrastructure/uuid"
)

type (
	MessageHandler = domain.MessageHandler
	SubPub         = domain.SubPub
	Subscription   = domain.Subscription
)

func NewSubPub() SubPub {
	gen, err := uuid.NewUUIDGenerator(uuid.GoogleUUID)
	if err != nil {
		log.Fatal(err)
	}
	return impl.NewSubPub(gen)
}
