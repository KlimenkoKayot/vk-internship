//go:generate mockgen -source=../domain/subpub.go -destination=mocks/subpub_mock.go -package=mocks
//go:generate mockgen -source=../internal/infrastructure/uuid/domain/interfaces.go -destination=mocks/uuid_mock.go -package=mocks
//go:generate mockgen -destination=mocks/logger_mock.go -package=mocks github.com/klimenkokayot/vk-internship/libs/logger Logger

package testutils
