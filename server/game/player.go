package game

import (
	"sync"

	"github.com/idalmasso/foxandchicken/server/game/common"
	"github.com/idalmasso/foxandchicken/server/game/messaging"
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

func (p *PlayerGameData) SetInput(a common.Vector2, actionPressed bool) {
	p.playerInput.SetInput(a, actionPressed)
}

func NewPlayer(username string, characterType CharacterType, gameRoom *GameRoom) *PlayerGameData {
	p := PlayerGameData{Username: username, playerCharacterType: characterType}
	action := playerActionObject{durationSeconds: 0.5}
	switch characterType {
	case CharacterTypeFox:
		action.playerAction = foxAction
		p.gameObject = NewGameObject(gameRoom, GameObjectTypeFox)
		killable := killableObject{hitPoints: 200}
		p.gameObject.AddBehaviour(&killable)
	case CharacterTypeChicken:
		action.playerAction = chickenAction
		p.gameObject = NewGameObject(gameRoom, GameObjectTypeChicken)
		killable := killableObject{hitPoints: 100}
		p.gameObject.AddBehaviour(&killable)
	}

	moving := MovingObject{MaxVelocity: gameRoom.MaxVelocity,
		Drag: gameRoom.Drag}
	p.gameObject.AddBehaviour(&moving)
	p.playerInput = &PlayerInput{}
	p.gameObject.AddBehaviour(p.playerInput)

	p.gameObject.AddBehaviour(&action)

	p.gameObject.Init()
	return &p
}

func (p *PlayerGameData) getStatusMessage() messaging.CommRoomMessagePlayerStatus {
	killable := p.gameObject.behaviours[KillableObjectBehaviour].(*killableObject)
	movement := p.gameObject.behaviours[MovingObjectBehaviour].(*MovingObject)
	m := messaging.CommRoomMessagePlayerStatus{Position: p.gameObject.Position,
		Player:           p.Username,
		ActionPressed:    p.playerInput.actionPressed,
		HitPoints:        killable.hitPoints,
		PerformingAction: p.playerInput.actionBehaviour.isPerforming,
		Velocity:         movement.Velocity,
	}
	return m
}
