package game

import (
	"sync"

	"github.com/idalmasso/foxandchicken/server/game/common"
)

//PlayerGameData contains the data for the server for a player in a room
type PlayerGameData struct {
	CharacterType int
	Username      string `json:"username"`
	Position      common.Vector2
	Rotation      float64
	Velocity      common.Vector2
	Acceleration  common.Vector2
	SizeRadius    float64
	timestamp     int64
	mutex         sync.Mutex
}
