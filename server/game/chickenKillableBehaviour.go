package game

type chickenKillableObject struct {
	hitPoints int
}

func (k *chickenKillableObject) init(g *GameObject) {
	k.hitPoints =0
}

func (o *chickenKillableObject) update(ts float64) {
}

func (i *chickenKillableObject) getType() GameBehaviourEnum {
	return ChickenkillableObjectBehaviour
}

func (i *chickenKillableObject) hit(damage int) {
	i.hitPoints -= damage
}
