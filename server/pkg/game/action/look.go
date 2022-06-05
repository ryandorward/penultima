package action
import (
	"app/pkg/game/model"
	"app/pkg/game/event/network"
	// "app/pkg/game/data"
	"app/pkg/game/tiles"
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

	zone := a.Looker.GetZone()

	for _, otherE := range zone.GetEntities() {
		otherX, otherY := otherE.GetPosition()
		if a.Looker != otherE && otherX == lookX && otherY == lookY {
			// someone is there, talk to them			
			a.Looker.GetClient().In <- network.NewServerResultEvent("You see " + otherE.GetName() + ".", "success")		
			return true 
		}
	}

	// now try world objects
	for _, obj := range zone.GetAllWorldObjects() {
		if obj.X == lookX && obj.Y == lookY {					
			var Name string
			if (obj.Name != "" ) {
				Name = obj.Name
			} else {
				Name = tiles.Tiles[obj.Tile].Name					
			}

			a.Looker.GetClient().In <- network.NewServerResultEvent("You see " + Name + ".", "success")		
			return true 
		}
	}


	tile := a.Looker.GetZone().GetTile(lookX, lookY)
	a.Looker.GetClient().In <- network.NewServerResultEvent("You see " + tile.Name + ".", "success")
	return true
} 