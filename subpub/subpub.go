package subpub

import (
	"github.com/klimenkokayot/vk-internship/libs/logger"
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
	log, _ := logger.NewAdapter(&logger.Config{
		Adapter: logger.AdapterZap,
		Level:   logger.LevelDebug,
	})
	gen, _ := uuid.NewUUIDGenerator(uuid.GoogleUUID)
	return impl.NewSubPub(gen, log)
}
