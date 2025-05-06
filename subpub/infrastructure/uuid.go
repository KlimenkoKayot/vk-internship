package infrastructure

import (
	"github.com/google/uuid"
	"github.com/klimenkokayot/vk-internship/domain"
)

type IDGenerator struct{}

func (g *IDGenerator) NewString() string {
	return uuid.NewString()
}

func NewIDGenerator() domain.IDGenerator {
	return &IDGenerator{}
}
