package model

type UserCommandState int

const (
	UserCommandStateNothing UserCommandState = iota
	UserCommandStateNewMessage
	UserCommandStateNewGIFCaption
)

type UserState struct {
	CommandState UserCommandState
	Payload      map[string]interface{}
}
