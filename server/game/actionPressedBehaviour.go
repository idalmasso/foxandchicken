package game

type playerActionPerform func(*playerActionBehaviour)

type playerActionBehaviour struct {
	isPerforming    bool
	pressed         bool
	durationSeconds float64
	actualTime      float64
	playerAction    playerActionPerform
	parentGO        *GameObject
}

func (b *playerActionBehaviour) init(g *GameObject) {
	b.isPerforming = false
	b.pressed = false
	b.parentGO = g
}

func (b *playerActionBehaviour) update(ts float64) {
	if b.pressed && !b.isPerforming {
		b.actualTime = 0
		b.isPerforming = true
		return
	}
	// actually starts the animation attack
	if b.isPerforming {
		b.actualTime += ts
		if b.actualTime >= b.durationSeconds {
			b.isPerforming = false
			b.playerAction(b)
		}
	}
}

func (b *playerActionBehaviour) actionPressed(pressed bool) {
	b.pressed = pressed
}
func (b *playerActionBehaviour) getType() GameBehaviourEnum {
	return PlayerActionBehaviour
}
