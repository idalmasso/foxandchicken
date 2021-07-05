package game

import "github.com/idalmasso/foxandchicken/server/game/common"

//PlayerGameData contains the data for the server for a player in a room
type PlayerGameData struct {
	CharacterType int
	Username      string `json:"username"`
	Position      common.Vector2
	Rotation      float32
	Velocity      common.Vector2
	SizeRadius    float32
}
