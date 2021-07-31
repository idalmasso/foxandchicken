package game

import (
	"sync"

	"github.com/idalmasso/foxandchicken/server/game/common"
)

type GameObject struct {
	size       float64
	Position   common.Vector2
	rotation   float64
	behaviours map[GameBehaviourEnum]gameBehaviour
	room       *GameRoom
	mutex      sync.RWMutex
}

type gameBehaviour interface {
	init(*GameObject)
	update(ts float64)
	getType() GameBehaviourEnum
}

func NewGameObject(gameRoom *GameRoom) *GameObject {
	g := GameObject{
		size: 1.0,
		Position:   common.Vector2{X: 0, Y: 0},
		rotation:   0,
		behaviours: make(map[GameBehaviourEnum]gameBehaviour),
		room: gameRoom,
	}
	return &g
}

func (g *GameObject) Init() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	priorities := behaviourPriorities()

	for _, p := range priorities {
		if b, ok := g.behaviours[p]; ok {
			b.init(g)
		}
	}
}

func (g *GameObject) Update(ts float64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	priorities := behaviourPriorities()
	for _, p := range priorities {
		if b, ok := g.behaviours[p]; ok {
			b.update(ts)
		}
	}
}

func (g *GameObject) AddBehaviour(behaviour gameBehaviour) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if _, ok := g.behaviours[behaviour.getType()]; !ok {
		g.behaviours[behaviour.getType()] = behaviour
	}
}
