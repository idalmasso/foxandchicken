package game

type GameBehaviourEnum int

const (
	MovingObjectBehaviour GameBehaviourEnum = iota
	PlayerInputBehaviour
	PlayerActionBehaviour
)

func behaviourPriorities() []GameBehaviourEnum {
	return []GameBehaviourEnum{PlayerInputBehaviour, PlayerActionBehaviour, MovingObjectBehaviour}
}
