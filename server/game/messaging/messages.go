package messaging

type MessageType int

const (
	MessageTypeCreateRoom MessageType = iota
	MessageTypeDeleteRoom
	MessageOkOrError
)

type CommMessageCreateRoom struct {
	Player string
	Name   string
}

func (m *CommMessageCreateRoom) GetMessageType() MessageType {
	return MessageTypeCreateRoom
}

type CommMessageDeleteRoom struct {
	Player string
	Name   string
}

func (m *CommMessageDeleteRoom) GetMessageType() MessageType {
	return MessageTypeDeleteRoom
}

type CommMessageOkOrError struct {
	Message string
}

func (m *CommMessageOkOrError) GetMessageType() MessageType {
	return MessageOkOrError
}

type MessageValue interface {
	GetMessageType() MessageType
}
