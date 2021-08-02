package game

import (
	"math"

	"github.com/golang/glog"
	"github.com/idalmasso/foxandchicken/server/game/common"
)

type MovingObject struct {
	Velocity     common.Vector2
	Acceleration common.Vector2
	MaxVelocity  float64
	Drag         float64
	gameObject   *GameObject
}

func (o *MovingObject) init(g *GameObject) {
	o.Acceleration = common.NullVector()
	o.Velocity = common.NullVector()
	o.gameObject = g
}

func (o *MovingObject) update(ts float64) {
	var selectedPosition common.Vector2
	selectedPosition = common.VectorSum(o.gameObject.Position, o.Velocity.ScalarProduct(ts))
	selectedPosition = selectedPosition.ClampVector(0, o.gameObject.room.sizeX, 0, o.gameObject.room.sizeY)
	objects := o.gameObject.room.gameObjectsInPoint(selectedPosition)
	found := false
	for _, obj := range objects {
		if obj != o.gameObject {
			found = true
			break
		}
	}
	if !found {
		o.gameObject.room.objectMove(o.gameObject, o.gameObject.Position, selectedPosition)
		o.gameObject.Position = selectedPosition
	}
	magnitude := o.Velocity.SqrtMagnitude()
	if o.Acceleration.X == 0 && o.Acceleration.Y == 0 {
		if magnitude < 0.15 {
			o.Velocity.X = 0
			o.Velocity.Y = 0
			return
		}
		o.Velocity = common.VectorSum(o.Velocity, o.Velocity.ScalarProduct(-o.Drag*ts))
	} else {
		o.Velocity = common.VectorSum(o.Velocity, o.Acceleration.ScalarProduct(ts))

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
	if o.Velocity.X == 0 {
		if o.Velocity.Y > 0 {
			o.gameObject.rotation = math.Pi
		} else {
			o.gameObject.rotation = -math.Pi
		}
	} else {
		o.gameObject.rotation = math.Atan(o.Velocity.Y / o.Velocity.X)
		if o.Velocity.X < 0 {
			o.gameObject.rotation += math.Pi
		}
		if glog.V(4) {
			glog.Infoln("DEBUG - change rotation: velocity", o.Velocity, "rotation", o.gameObject.rotation*180/math.Pi)
		}
	}

}

func (i *MovingObject) getType() GameBehaviourEnum {
	return MovingObjectBehaviour
}
