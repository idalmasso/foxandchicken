package messaging

type MessageType int

const (
	MessageTypeCreateRoom MessageType = iota
	MessageTypeDeleteRoom
	MessageResponse
	MessageResponseCreateRoom
	MessageResponseJoinRoom
	RoomMessageTypeMovePlayer
	RoomMessageTypePlayerStatus
	RoomMessageTypePlayersStatuses
	RoomMessageTypeJoinPlayer
	RoomMessageTypeLeftPlayer
	RoomMessageTypeResponseMessage
	RoomMessageTypeActionPlayer
)
