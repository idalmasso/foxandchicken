package messaging

type MessageType int

const (
	MessageTypeCreateRoom MessageType = iota
	MessageTypeDeleteRoom
	MessageResponse
	MessageResponseCreateRoom
	RoomMessageTypeMovePlayer
	RoomMessageTypeJoinPlayer
	RoomMessageTypeLeftPlayer
	RoomMessageTypeResponseMessage
)

type CommMessageCreateRoom struct {
	Player string
	Name   string
}

func (m *CommMessageCreateRoom) GetMessageType() MessageType {
	return MessageTypeCreateRoom
}
func (m *CommMessageCreateRoom) ErrorMessage() string {
	return ""
}

type CommMessageDeleteRoom struct {
	Player string
	Name   string
}

func (m *CommMessageDeleteRoom) GetMessageType() MessageType {
	return MessageTypeDeleteRoom
}
func (m *CommMessageDeleteRoom) ErrorMessage() string {
	return ""
}

type CommMessageResponse struct {
	Message string
}

func (m *CommMessageResponse) GetMessageType() MessageType {
	return MessageResponse
}
func (m *CommMessageResponse) ErrorMessage() string {
	return m.Message
}

type CommMessageResponseCreateRoom struct {
	Message     string
	RoomChannel chan<- RoomMessageValue
}

func (m *CommMessageResponseCreateRoom) GetMessageType() MessageType {
	return MessageResponseCreateRoom
}
func (m *CommMessageResponseCreateRoom) ErrorMessage() string {
	return m.Message
}

type MessageValue interface {
	GetMessageType() MessageType
	ErrorMessage() string
}
