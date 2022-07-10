package action
import (
	"app/pkg/game/model"
	"app/pkg/game/event/network"	
	// "strconv"
	// "fmt"
)

type DirectionalAction struct {
	Actor model.Entity	
	Action string	
	X, Y  int
}

func (a *DirectionalAction) Execute() bool {	
	
	// can only act N, S, E, or W
	if ((a.X != 0 && a.Y != 0) || a.X > 1 || a.X < -1 || a.Y > 1 || a.Y < -1) {
		return false
	} 

	x, y := a.Actor.GetPosition()
	zone := a.Actor.GetZone() 

	actionX := x + a.X
	actionY := y + a.Y

	if (a.Action == "attack") {		
		for _, otherE := range zone.GetEntities() {
			otherX, otherY := otherE.GetPosition()
			if a.Actor != otherE && otherX == actionX && otherY == actionY {					
				a.Actor.GetClient().In <- network.NewServerResultEvent("Attack " + otherE.GetName() + " with hands.", "success")		
				otherE.TakeDamage(1.0)			
				return true 
			}
		}
		a.Actor.GetClient().In <- network.NewServerResultEvent("Attack with hands. Nothing to attack!", "fail")			
	} 

	return false

}