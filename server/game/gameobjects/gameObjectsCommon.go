package gameobjects

type GameBehaviourEnum int

const (
	MovingObjectBehaviour GameBehaviourEnum = iota
	PlayerInputBehaviour
)

func behaviourPriorities() []GameBehaviourEnum {
	return []GameBehaviourEnum{PlayerInputBehaviour, MovingObjectBehaviour}
}
