package action
import (
	"app/pkg/game/model"
	"app/pkg/game/event"
)

type UpdateAvatarAction struct {
	Mover model.Entity
	Id  int
}

func (a *UpdateAvatarAction) Execute() bool {
	a.Mover.SetTile(int(a.Id))
	eX, eY := a.Mover.GetPosition()
	event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY})
	return true
}