package adapter

import (
	"github.com/google/uuid"
	"github.com/klimenkokayot/vk-internship/subpub/internal/infrastructure/uuid/domain"
)

type GoogleUUIDGenerator struct{}

func (g *GoogleUUIDGenerator) NewString() string {
	return uuid.NewString()
}

func NewGoogleUUIDGenerator() domain.UUIDGenerator {
	return &GoogleUUIDGenerator{}
}
