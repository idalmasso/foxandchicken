package game

import (
	"math"

	"github.com/idalmasso/foxandchicken/server/game/common"
)


func foxAction(b *playerActionObject) {
	sumVect := common.Vector2{X: math.Cos(b.parentGO.rotation), Y: math.Sin(b.parentGO.rotation)}
	sumVect = common.VectorSum(sumVect, b.parentGO.Position)
	objects:= b.parentGO.room.gameObjectsInPoint(sumVect)
	for _, target:=range(objects){
		if b, ok := target.behaviours[ChickenkillableObjectBehaviour]; ok {
			targetBehaviour := b.(*chickenKillableObject)
			targetBehaviour.hit(10)
		}
	}
}
