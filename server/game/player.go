package game

import (
	"sync"

	"github.com/idalmasso/foxandchicken/server/game/common"
	"github.com/idalmasso/foxandchicken/server/game/gameobjects"
)

//PlayerGameData contains the data for the server for a player in a room
type PlayerGameData struct {
	CharacterType int
	Username      string `json:"username"`
	gameObject    *gameobjects.GameObject
	playerInput   *gameobjects.PlayerInput
	mutex         sync.RWMutex
}

func (p *PlayerGameData) SetInput(a common.Vector2) {
	p.playerInput.SetInput(a, false)
}

func NewPlayer(username string, characterType int, gameRoom *GameRoom) *PlayerGameData {
	p := PlayerGameData{Username: username, CharacterType: characterType}
	p.gameObject = gameobjects.NewGameObject()
	moving := gameobjects.MovingObject{GameSize: common.Vector2{X: gameRoom.sizeX, Y: gameRoom.sizeY},
		MaxVelocity: gameRoom.MaxVelocity,
		Drag:        gameRoom.Drag}
	p.gameObject.AddBehaviour(&moving)
	p.playerInput = &gameobjects.PlayerInput{}
	p.gameObject.AddBehaviour(p.playerInput)
	p.gameObject.Init()
	return &p
}
