package action

import (
	"app/pkg/game/event"
	"app/pkg/game/model"
	//"app/pkg/game/util"
	// "errors" 
	"math/rand"
	"fmt" 
	// "app/pkg/game/event/network"
)

type MoveAction struct {
	Mover model.Entity
	X, Y  int
}

func (a *MoveAction) Execute() bool {
	eX, eY := a.Mover.GetPosition()
	
	// Prohibit moving on diagonal, or by more than one tile, for now anyways
	if ((a.X < -1) || (a.X > 1) || (a.Y < -1) || (a.Y > 1) || ((a.X != 0) && (a.Y != 0))) {
		event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary.
		return false 
	}

	// 	fmt.Println("MoveAction.Execute(): ", eX, eY, a.X, a.Y, a.Mover.GetZone().GetName())

	// Now get the move's translation into a location 
	nX, nY := a.Mover.GetZone().GetNewLocation(eX,eY,a.X,a.Y)

	// If it's -1,-1 that means it's a zone-exit. Is this too kludgey? 
	if (nX == -1) && (nY == -1) {
		if a.Mover.GetType() == model.EntityTypeNPC { // for now, NPCs are restricted to their zone			
			a.Mover.ReceiveResult("", "blocked")
			return false
		} else {
			parentZone := a.Mover.GetZone().GetParentZone() 
			//a.Mover.GetClient().In <- network.NewServerMessageEvent("Exit to " + parentZone.GetName() + ".")
			a.Mover.ReceiveMessage("Exit to " + parentZone.GetName() + ".")
			// find child zone in parent warpObjects
			for _, obj := range parentZone.GetAllWorldObjects() {				
				if (	obj.WarpTarget != nil && obj.WarpTarget.ZoneName != "" && obj.WarpTarget.Zone == a.Mover.GetZone()) {
					fmt.Println("obj.WarpTarget.ZoneName: " + obj.WarpTarget.ZoneName)	
					a.Mover.SetZoneWithTarget(parentZone, obj.X, obj.Y)
					return true						
				}																						
			}
		}
	}

	newTile := a.Mover.GetZone().GetTile(nX, nY)
		
	if newTile == nil || newTile.Solid {	
		// edge of map or or solid, don't move
		// a.Mover.GetClient().In <- network.NewServerMessageEvent("Blocked!")
		a.Mover.ReceiveResult("Blocked!", "blocked")			
		// event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
		return false
	} 

	for _, otherE := range a.Mover.GetZone().GetEntities() {
		otherX, otherY := otherE.GetPosition()
		if a.Mover != otherE && otherX == nX && otherY == nY {		
			a.Mover.ReceiveResult("Blocked!", "blocked")	
			//a.Mover.GetClient().In <- network.NewServerResultEvent("Blocked!", "blocked")
			event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
			return false
		}
	}

	rando := rand.Float64()
	var thresh float64

	lastX, lastY := a.Mover.GetLastMoveTry()
	if (lastX == nX && lastY == nY) {
		thresh = a.Mover.GetSlowThresh() * 1.3
	} else {
		a.Mover.SetLastMoveTry(nX, nY)
		thresh = newTile.Speed
	}

	if rando <= thresh {
		a.Mover.SetSlowThresh(1.0)
		a.Mover.SetPosition(nX, nY)
		var msg string
		if ((a.X == 0) && (a.Y == 0)) {
			msg = "Pass"
		}	else if ((a.X == -1) && (a.Y == 0)) {
			msg = "West"
		}	else if ((a.X == 0) && (a.Y == -1)) {
			msg = "North"
		}	else if ((a.X == 1) && (a.Y == 0)) {
			msg = "East"
		}	else if ((a.X == 0) && (a.Y == 1)) {
			msg = "South"
		}			
		a.Mover.ReceiveResult(msg, "success")
		
	//	a.Mover.GetClient().In <- network.NewServerResultEvent(msg, "success")
	} else {
		a.Mover.SetSlowThresh(thresh)	
		a.Mover.ReceiveResult("Slow progress!", "slow")
	//	a.Mover.GetClient().In <- network.NewServerResultEvent("Slow progress!", "slow")	
	}

	// now is time to calculate mover's view. 
	// And decide do we:
	// 1. Keep track of other observers within FOV and only notify them?
	// 2. or Notify each observer, then let each observer calculate if it can now see mover and update client accordingly?
	// #1 will be more efficient order of n, #2 will be n^2. #2 allows more flexibility, for example
	// if we have invisibility or magical/enhanced seeing abilities.

	// we should NotifyObservers to update their own view
	event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: nX, Y: nY})

	return true // success
}
