package gameobjects

import "github.com/idalmasso/foxandchicken/server/game/common"

type PlayerInput struct {
	directionInput       common.Vector2
	actionPressed        bool
	movingObjectBehavour *MovingObject
}

func (i *PlayerInput) init(g *GameObject) {
	if mo, ok := g.behaviours[MovingObjectBehaviour]; ok {
		i.directionInput = common.NullVector()
		i.actionPressed = false
		i.movingObjectBehavour = mo.(*MovingObject)
	} else {
		panic("Error in code, no behaviour correct")
	}
}

func (i *PlayerInput) update(ts float64) {
	i.movingObjectBehavour.Acceleration.X = i.directionInput.X
	i.movingObjectBehavour.Acceleration.Y = i.directionInput.Y
	//TODO: dispatch here the "action" to some other behaviour
}
func (i *PlayerInput) getType() GameBehaviourEnum {
	return PlayerInputBehaviour
}

func (i *PlayerInput) SetInput(input common.Vector2, actionPressed bool) {
	i.movingObjectBehavour.gameObject.mutex.Lock()
	defer i.movingObjectBehavour.gameObject.mutex.Unlock()
	i.directionInput.X = input.X
	i.directionInput.Y = input.Y
	i.actionPressed = actionPressed
}
