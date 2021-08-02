package messaging

import "github.com/idalmasso/foxandchicken/server/game/common"

//CommRoomMessageMovePlayer movement of a player message
type CommRoomMessageMovePlayer struct {
	Player        string         `json:"player"`
	Position      common.Vector2 `json:"position"`
	Velocity      common.Vector2 `json:"velocity"`
	Acceleration  common.Vector2 `json:"acceleration"`
	ActionPressed bool           `json:"action"`
	Rotation      float64        `json:"rotation"`
	Timestamp     int64          `json:"ts"`
}

func (m *CommRoomMessageMovePlayer) ErrorMessage() string {
	return ""
}
func (m *CommRoomMessageMovePlayer) GetMessageType() MessageType {
	return RoomMessageTypeMovePlayer
}

//CommRoomMessagePlayerStatus movement and other for players
type CommRoomMessagePlayerStatus struct {
	Player           string         `json:"player"`
	Position         common.Vector2 `json:"position"`
	Velocity         common.Vector2 `json:"velocity"`
	ActionPressed    bool           `json:"action"`
	PerformingAction bool           `json:"performingaction"`
	Rotation         float64        `json:"rotation"`
	HitPoints        int            `json:"hitpoints"`
	Timestamp        int64          `json:"ts"`
}

func (m *CommRoomMessagePlayerStatus) ErrorMessage() string {
	return ""
}
func (m *CommRoomMessagePlayerStatus) GetMessageType() MessageType {
	return RoomMessageTypeMovePlayer
}

type CommRoomMessagePlayersStatuses []CommRoomMessagePlayerStatus

func (m *CommRoomMessagePlayersStatuses) ErrorMessage() string {
	return ""
}
func (m *CommRoomMessagePlayersStatuses) GetMessageType() MessageType {
	return RoomMessageTypePlayersStatuses
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
