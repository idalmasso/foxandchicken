package gameobjects

import (
	"math"

	"github.com/idalmasso/foxandchicken/server/game/common"
)

type MovingObject struct {
	Velocity     common.Vector2
	Acceleration common.Vector2
	MaxVelocity  float64
	Drag         float64
	GameSize     common.Vector2
	gameObject   *GameObject
}

func (o *MovingObject) init(g *GameObject) {
	o.Acceleration = common.NullVector()
	o.Velocity = common.NullVector()
	o.gameObject = g
}

func (o *MovingObject) update(ts float64) {
	o.gameObject.Position = common.VectorSum(o.gameObject.Position, o.Velocity.ScalarProduct(ts))
	o.gameObject.Position = o.gameObject.Position.ClampVector(0, o.GameSize.X, 0, o.GameSize.Y)
	if o.Acceleration.X == 0 && o.Acceleration.Y == 0 {

		magnitude := o.Velocity.SqrtMagnitude()
		if magnitude < 0.15 {
			o.Velocity.X = 0
			o.Velocity.Y = 0
			return
		}
		o.Velocity = common.VectorSum(o.Velocity, o.Velocity.ScalarProduct(-o.Drag*ts))
	} else {
		o.Velocity = common.VectorSum(o.Velocity, o.Acceleration.ScalarProduct(ts))
		magnitude := o.Velocity.SqrtMagnitude()

		if magnitude > o.MaxVelocity {
			o.Velocity = o.Velocity.ScalarProduct(o.MaxVelocity / magnitude)
		}
	}
	if math.Abs(o.Acceleration.X) == 0 && math.Abs(o.Velocity.X) < 0.1 {
		o.Velocity.X = 0
	}
	if math.Abs(o.Acceleration.Y) == 0 && math.Abs(o.Velocity.Y) < 0.1 {
		o.Velocity.Y = 0
	}
}

func (i *MovingObject) getType() GameBehaviourEnum {
	return MovingObjectBehaviour
}
