package game

import "github.com/golang/glog"

type killableObject struct {
	hitPoints int
}

func (k *killableObject) init(g *GameObject) {
	k.hitPoints = 0
}

func (o *killableObject) update(ts float64) {
}

func (i *killableObject) getType() GameBehaviourEnum {
	return KillableObjectBehaviour
}

func (i *killableObject) hit(damage int) {
	if glog.V(3) {
		glog.Infoln("DEBUG - killableObject hit for damage", damage)
	}
	if i.hitPoints > 0 {
		i.hitPoints -= damage
		if i.hitPoints <= 0 {
			i.die()
		}
	}
}
func (i *killableObject) die() {
	if glog.V(2) {
		glog.Infoln("killableObject dead")
	}
	//TODO: Todo
}
