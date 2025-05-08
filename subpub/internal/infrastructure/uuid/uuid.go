package uuid

import (
	"fmt"

	"github.com/klimenkokayot/vk-internship/subpub/internal/infrastructure/uuid/adapter"
	"github.com/klimenkokayot/vk-internship/subpub/internal/infrastructure/uuid/domain"
)

type (
	UUIDGenerator = domain.UUIDGenerator
	GeneratorType int
)

const (
	GoogleUUID GeneratorType = iota
)

var (
	ErrUnknownGenerator = fmt.Errorf("неизвестный тип генератора")
)

func NewUUIDGenerator(generator GeneratorType) (domain.UUIDGenerator, error) {
	switch generator {
	case GoogleUUID:
		return adapter.NewGoogleUUIDGenerator(), nil
	default:
		return nil, ErrUnknownGenerator
	}
}
