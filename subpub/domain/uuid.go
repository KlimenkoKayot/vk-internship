package domain

type IDGenerator interface {
	NewString() string
}
