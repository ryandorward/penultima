package action
import (
	"app/pkg/game/model"
	"app/pkg/game/event/network"
	// "app/pkg/game/data"
	"app/pkg/game/tiles"
	"strconv"
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
			var name, quantity string
			if (obj.Name != "" ) {
				name = obj.Name
			} else {
				name = tiles.Tiles[obj.Tile].Name					
			}
			if (obj.Quantity != 0 ) {
				quantity = strconv.Itoa(obj.Quantity) + " "
			}

			a.Looker.GetClient().In <- network.NewServerResultEvent("You see " + quantity + name + ".", "success")		
			return true 
		}
	}


	tile := a.Looker.GetZone().GetTile(lookX, lookY)
	a.Looker.GetClient().In <- network.NewServerResultEvent("You see " + tile.Name + ".", "success")
	return true
} 