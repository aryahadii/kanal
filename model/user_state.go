package model

type UserCommandState int

const (
	UserCommandStateNothing UserCommandState = iota
	UserCommandStateNewMessage
)

type UserState struct {
	CommandState UserCommandState
}
