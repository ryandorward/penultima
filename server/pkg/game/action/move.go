package action

import (
	"app/pkg/game/event"
	"app/pkg/game/model"
	//"app/pkg/game/util"
	// "errors" 
	"math/rand"
	"fmt"
	"app/pkg/game/event/network"

)

/*
const (
	maxValidMoveDist = 1
)
*/

type MoveAction struct {
	Mover model.Entity
	X, Y  int
}

func (a *MoveAction) Execute() bool {
	eX, eY := a.Mover.GetPosition()

	// Prohibit moving on diagonal, or by more than one tile, for now anyways
	if ((a.X < -1) || (a.X > 1) || (a.Y < -1) || (a.Y > 1) || ((a.X != 0) && (a.Y != 0))) {
		event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
		return false
	}

	// Now get the move's translation into a location
	nX, nY := a.Mover.GetZone().GetNewLocation(eX,eY,a.X,a.Y)
	currentTile := a.Mover.GetZone().GetTile(eX, eY)
	t := a.Mover.GetZone().GetTile(nX, nY)
		
	if t == nil || t.Solid {	
		// edge of map or or solid, don't move
		a.Mover.GetClient().In <- network.NewServerMessageEvent("Blocked!")
		event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
		return false
	} 

	/*
	if a.Mover.GetType() == model.EntityTypePlayer {
		objs := a.Mover.GetZone().GetWorldObjects(a.X, a.Y)
		for _, obj := range objs {
			if obj.WarpTarget != nil {
				event.NotifyObservers(event.DespawnEvent{Entity: a.Mover})
				a.Mover.GetZone().RemoveEntity(a.Mover)
				obj.WarpTarget.Zone.AddEntity(a.Mover)
				a.Mover.SetPosition(obj.WarpTarget.X, obj.WarpTarget.Y)
				event.NotifyObservers(event.SpawnEvent{Entity: a.Mover})
				return true
			}		
			if obj.HealZone != nil {
				if obj.HealZone.Full {
					a.Mover.Heal(a.Mover.GetStats().MaxHP)
					event.NotifyObservers(event.HealEvent{Entity: a.Mover, Amount: a.Mover.GetStats().MaxHP, Full: true})
					event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
					return true
				}
			}		
		}
	}
	*/

	for _, otherE := range a.Mover.GetZone().GetEntities() {
		otherX, otherY := otherE.GetPosition()
		if a.Mover != otherE && otherX == nX && otherY == nY {
			// someone is there, block the way		
			a.Mover.GetClient().In <- network.NewServerMessageEvent("Blocked!")
			event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
			return false
		}
	}

	rando := rand.Float64()
	var thresh float64
	if (a.Mover.GetSlowThresh() < 1.0) {	
		thresh = a.Mover.GetSlowThresh() * 1.3
	} else {
		thresh = currentTile.Speed
	}

	if rando <= thresh {
		a.Mover.SetSlowThresh(1.0) 
		a.Mover.SetPosition(nX, nY)

		var msg string
		if ((a.X == 0) && (a.Y == 0)) {
			msg = "Stationary"
		}	else if ((a.X == -1) && (a.Y == 0)) {
			msg = "West"
		}	else if ((a.X == 0) && (a.Y == -1)) {
			msg = "North"
		}	else if ((a.X == 1) && (a.Y == 0)) {
			msg = "East"
		}	else if ((a.X == 0) && (a.Y == 1)) {
			msg = "South"
		}	
		a.Mover.GetClient().In <- network.NewServerMessageEvent(msg)
	} else {
		a.Mover.SetSlowThresh(thresh)
		a.Mover.GetClient().In <- network.NewServerMessageEvent("Slow progress!")	
	}

	// now is time to calculate mover's view. 
	// And decide do we:
	// 1. Keep track of other observers within FOV and only notify them?
	// 2. or Notify each observer, then let each observer calculate if it can now see mover and update client accordingly?
	// #1 will be more efficient order of n, #2 will be n^2. #2 allows more flexibility, for example
	// if we have invisibility or magical/enhanced seeing abilities.

	// we should NotifyObservers to update their own view
	event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: nX, Y: nY})
	fmt.Printf("MoveAction: Execute: X: %d, Y: %d\n", nX, nY) 

	return true // success
}

/*
// implementation of zone interface for main world will have a function like this 
// to return new location with wrapping
func GetNewLocation(y, x, dx, dy) (nX, nY int) {	
	nX := util.WrapMod(X+dX, zoneWidth)
	nY := util.WrapMod(Y+Y, zoneHeight)
	return (nX,nY)
}	

*/