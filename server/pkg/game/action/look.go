package action
import (
	"app/pkg/game/model"
	"app/pkg/game/event/network"
)

type LookAction struct {
	Looker model.Entity	
	X int
	Y int	
}

func (a *LookAction) Execute() bool {	
	eX, eY := a.Looker.GetPosition()
	lookX := eX + a.X
	lookY := eY + a.Y
	tile := a.Looker.GetZone().GetTile(lookX, lookY)
	a.Looker.GetClient().In <- network.NewServerMessageEvent("You see " + tile.Name + ".")
	return true
} 