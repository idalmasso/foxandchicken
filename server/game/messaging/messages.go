package messaging

type MessageType int

const (
	MessageTypeCreateRoom MessageType = iota
	MessageTypeDeleteRoom
	MessageResponse
	MessageResponseCreateRoom
	MessageResponseJoinRoom
	RoomMessageTypeMovePlayer
	RoomMessageTypePlayersMovment
	RoomMessageTypeJoinPlayer
	RoomMessageTypeLeftPlayer
	RoomMessageTypeResponseMessage
)
