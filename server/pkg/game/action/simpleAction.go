package action
import (
	"app/pkg/game/model"
// 	"app/pkg/game/event"
	//"app/pkg/game/data"
	"app/pkg/game/event/network"	
	"fmt"
//	"app/pkg/game/util"
//	"strings"
//	"regexp"
)

type SimpleAction struct {
	Actor model.Entity	
	Action string	
}

func (a *SimpleAction) Execute() bool {	

	x, y := a.Actor.GetPosition()
	zone := a.Actor.GetZone() 

	if (a.Action == "enter") {
		// check if on enterable world object
		for _, obj := range zone.GetAllWorldObjects() {
			if obj.X == x && obj.Y == y {						
				if obj.WarpTarget != nil {
					a.Actor.GetClient().In <- network.NewServerMessageEvent("Entering " + obj.Name)											
					a.Actor.SetZoneWithTarget(obj.WarpTarget.Zone, obj.WarpTarget.X, obj.WarpTarget.Y)																
					return true
				}					
			}
		}
		a.Actor.GetClient().In <- network.NewServerMessageEvent("Enter.")
		return true
	} else if (a.Action == "get") {
		// check if any world object is food. Eventually we might have some "gettable" flag
		for key, obj := range zone.GetAllWorldObjects() {
			if obj.X == x && obj.Y == y {						
				if obj.Type == "food" {
					a.Actor.GetClient().In <- network.NewServerMessageEvent("Get Food!")	
					a.Actor.AddFood(20.0)			
					fmt.Println(obj, key)						
					zone.RemoveWorldObjectByUUID(obj.UUID) 
					return true
				}	
				if obj.Type == "gem" {
					a.Actor.GetClient().In <- network.NewServerMessageEvent("Get Gems!")	
					a.Actor.AddGems(1)										
					zone.RemoveWorldObjectByUUID(obj.UUID) 
					return true
				}					
			}
		}
		a.Actor.GetClient().In <- network.NewServerMessageEvent("Get. Nothing to get!")
		return true
	}


	return false

}