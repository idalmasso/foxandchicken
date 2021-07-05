package messaging

import "github.com/idalmasso/foxandchicken/server/game/common"

type CommRoomMessageMovePlayer struct {
	Player             string `json:"player"`
	Position common.Vector2 `json:"position"`
	Velocity common.Vector2  `json:"velocity"`
	Rotation           float32 `json:"rotation"`
}

func (m *CommRoomMessageMovePlayer) ErrorMessage() string {
	return ""
}
func (m *CommRoomMessageMovePlayer) GetMessageType() MessageType {
	return RoomMessageTypeMovePlayer
}
type CommRoomMessagePlayersMovement []CommRoomMessageMovePlayer

func (m *CommRoomMessagePlayersMovement) ErrorMessage() string {
	return ""
}
func (m *CommRoomMessagePlayersMovement) GetMessageType() MessageType {
	return RoomMessageTypePlayersMovement
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
