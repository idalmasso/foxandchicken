package messaging

import "github.com/idalmasso/foxandchicken/server/game/common"

type CommRoomMessageMovePlayer struct {
	Player             string
	Position, Velocity common.Vector2
	Rotation           float32
}

func (m *CommRoomMessageMovePlayer) ErrorMessage() string {
	return ""
}
func (m *CommRoomMessageMovePlayer) GetMessageType() MessageType {
	return RoomMessageTypeMovePlayer
}

type CommRoomMessageJoinPlayer struct {
	Player string
	Name   string
}

func (m *CommRoomMessageJoinPlayer) GetMessageType() MessageType {
	return RoomMessageTypeJoinPlayer
}
func (m *CommRoomMessageJoinPlayer) ErrorMessage() string {
	return ""
}

type CommRoomMessageLeftPlayer struct {
	Player string
}

func (m *CommRoomMessageLeftPlayer) ErrorMessage() string {
	return ""
}
func (m *CommRoomMessageLeftPlayer) GetMessageType() MessageType {
	return RoomMessageTypeLeftPlayer
}

type CommRoomMessageResponse struct {
	Message string
}

func (m *CommRoomMessageResponse) ErrorMessage() string {
	return m.Message
}
func (m *CommRoomMessageResponse) GetMessageType() MessageType {
	return RoomMessageTypeResponseMessage
}

type RoomMessageValue interface {
	GetMessageType() MessageType
	ErrorMessage() string
}
