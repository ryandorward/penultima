package action
import (
	"app/pkg/game/model"
	"app/pkg/game/event/network"
)

type TalkAction struct {
	Actor model.Entity	
	X int
	Y int	
}

func (a *TalkAction) Execute() bool {	
	eX, eY := a.Actor.GetPosition()
	lookX := eX + a.X
	lookY := eY + a.Y
	tile := a.Actor.GetZone().GetTile(lookX, lookY)
	a.Actor.GetClient().In <- network.NewServerMessageEvent("(You can't talk to " + tile.Name + ".)")
	a.Actor.GetClient().In <- network.NewServerMessageEvent("You say:")	
	return true
}