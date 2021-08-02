package game

import (
	"math"

	"github.com/golang/glog"
	"github.com/idalmasso/foxandchicken/server/game/common"
)

func foxAction(b *playerActionObject) {

	sumVect := common.Vector2{X: math.Cos(b.parentGO.rotation) * 1.5, Y: math.Sin(b.parentGO.rotation) * 1.5}
	sumVect = common.VectorSum(sumVect, b.parentGO.Position)
	if glog.V(3) {
		glog.Infoln("Fox attacking, actual position:", b.parentGO.Position, "rotation:", 180/math.Pi*b.parentGO.rotation, "attacking position:", sumVect)
	}
	objects := b.parentGO.room.gameObjectsInPoint(sumVect)
	for _, target := range objects {
		if target.gameObjectType == GameObjectTypeChicken {
			if b, ok := target.behaviours[KillableObjectBehaviour]; ok {
				targetBehaviour := b.(*killableObject)
				targetBehaviour.hit(10)
			}
		}
	}
}
