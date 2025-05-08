package domain

type UUIDGenerator interface {
	NewString() string
}
