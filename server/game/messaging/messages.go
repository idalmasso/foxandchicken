package messaging

type MessageType int

const (
	MessageTypeCreateRoom MessageType = iota
	MessageTypeDeleteRoom
	MessageResponse
	MessageResponseCreateRoom
	MessageResponseJoinRoom
	RoomMessageTypeMovePlayer
	RoomMessageTypePlayersMovement
	RoomMessageTypeJoinPlayer
	RoomMessageTypeLeftPlayer
	RoomMessageTypeResponseMessage
)
