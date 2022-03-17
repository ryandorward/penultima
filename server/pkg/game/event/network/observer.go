package network

import	(
	"app/pkg/game/event"
	// "fmt"
	
	 "app/pkg/game/model"
	// "app/pkg/game/util"
	//"app/pkg/game/data"
//	"app/pkg/fov"
)

const MOTD = "Welcome to Penultima! Stay safe."

type networkObserver struct {
}

func NewObserver() event.Observer {
	return &networkObserver{}
}

func (o *networkObserver) Notify(e event.Event) {
	switch v := e.(type) {

	case event.PeerGemEvent:
		if c := v.Entity.GetClient(); c != nil {
			c.In <- NewServerMessageEvent(MOTD)
		}

	case event.JoinEvent:
		if c := v.Entity.GetClient(); c != nil {
			c.In <- NewServerMessageEvent(MOTD)
		}

	/*
	case event.ServerMessageEvent:
		if c := v.Entity.GetClient(); c != nil {
			c.In <- newServerMessageEvent(e.Message)
			fmt.Println(e)
		}		
	*/		

	case event.LeaveEvent:
		break

	// Don't think we'll need this. Purpose is to let clients know new entity in zone
	case event.SpawnEvent:
		for _, ent := range v.Entity.GetZone().GetEntities() {
			if c := ent.GetClient(); c != nil {
				c.In <- newSpawnEvent(v.Entity)
			}
		}
		if c := v.Entity.GetClient(); c != nil {
			c.In <- newZoneLoadEvent(v.Entity.GetZone())
		}

	// let clients know entity left zone
	case event.DespawnEvent:
		for _, ent := range v.Entity.GetZone().GetEntities() {
			if c := ent.GetClient(); c != nil {
				c.In <- newDespawnEvent(v.Entity, false)
			}
		}
 	
	/*
	case event.DieEvent:
		for _, ent := range v.Entity.GetZone().GetEntities() {
			if c := ent.GetClient(); c != nil {
				c.In <- newDespawnEvent(v.Entity, true)
			}
		}
		if c := v.Entity.GetClient(); c != nil {
			c.In <- newServerMessageEvent("You died.")
			c.In <- newServerMessageEvent("Your soul enters a new body. You are reborn.")
		}	
	*/

	case event.MoveEvent:

		// fmt.Printf("game/event/observer MoveEvent: X: %d, Y: %d\n", v.X, v.Y)
		// tell all clients about this entity's move. Won't be doing this, we'll
		// simply get each concerned client to recalc it's view when an entity moves.				
		
		inView := make([]model.Entity, 0) 
		for _, ent := range v.Entity.GetZone().GetEntities() {
			// This is the point where we can calculate FOV and send it to each client? 
			// Well it won't be this, it will be something like event.UpdateViewEvent 
			// Original code tells client about how an entity moved, and leaves it up 
			// to client to update view. But we are going to handle view on thie end												
			if v.Entity.IsInView(ent) {	
				inView = append(inView, ent)																	
			}								
		}

		for _, ent := range inView  {
			// ent.UpdateOwnView()
			if c := ent.GetClient(); c != nil {	
				ent.UpdateOwnView(c)
			}
		}			

	/*
	case event.ChatEvent:
		for _, ent := range v.Entity.GetZone().GetEntities() {
			if c := ent.GetClient(); c != nil {
				c.In <- newChatEvent(v.Entity, v.Message)
			}
		}

	case event.AttackEvent:
		for _, ent := range v.Attacker.GetZone().GetEntities() {
			if c := ent.GetClient(); c != nil {
				c.In <- newAttackEvent(v.Attacker, v.Target.GetUUID(), v.Hit, v.Damage, v.TargetHP)
			}
		}

	case event.HealEvent:
		if c := v.Entity.GetClient(); c != nil {
			c.In <- newUpdateEvent(v.Entity)
			c.In <- newServerMessageEvent("You pray to your gods and are fully healed in their light.")
		}

	case event.GainXPEvent:
		if c := v.Entity.GetClient(); c != nil {
			c.In <- newUpdateEvent(v.Entity)
			if v.LeveledUp {
				c.In <- newServerMessageEvent("You leveled up! You have a newfound strength coursing through your veins.")
			} 
		}
	*/

	}
}
