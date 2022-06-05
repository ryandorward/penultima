package action
import (
	"app/pkg/game/model"	
)

type UpdateNameAction struct {
	Actor model.Entity
	Name string
}

func (a *UpdateNameAction) Execute() bool {
	a.Actor.SetName(a.Name)
	return true
}