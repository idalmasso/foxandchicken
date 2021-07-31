package game

import (
	"sync"

	"github.com/idalmasso/foxandchicken/server/game/common"
)

type CharacterType int

const (
	CharacterTypeFox CharacterType = iota
	CharacterTypeChicken
)

//PlayerGameData contains the data for the server for a player in a room
type PlayerGameData struct {
	playerCharacterType CharacterType
	Username            string `json:"username"`
	gameObject          *GameObject
	playerInput         *PlayerInput
	mutex               sync.RWMutex
}

func (p *PlayerGameData) SetInput(a common.Vector2) {
	p.playerInput.SetInput(a, false)
}

func NewPlayer(username string, characterType CharacterType, gameRoom *GameRoom) *PlayerGameData {
	p := PlayerGameData{Username: username, playerCharacterType: characterType}
	p.gameObject = NewGameObject(gameRoom)
	moving := MovingObject{MaxVelocity: gameRoom.MaxVelocity,
		Drag: gameRoom.Drag}
	p.gameObject.AddBehaviour(&moving)
	p.playerInput = &PlayerInput{}
	p.gameObject.AddBehaviour(p.playerInput)
	action := playerActionBehaviour{durationSeconds: 0.5}
	switch characterType {
	case CharacterTypeFox:
		action.playerAction = foxAction
	case CharacterTypeChicken:
		action.playerAction = chickenAction
	}
	p.gameObject.AddBehaviour(&action)
	p.gameObject.Init()
	return &p
}
