package game

import "github.com/idalmasso/foxandchicken/server/game/common"

type PlayerGameData struct {
	CharacterType int
	Username      string `json:"username"`
	Position      common.Vector2
	Rotation      float32
	Velocity      common.Vector2
	SizeRadius    float32
}
