package action
import (
	"app/pkg/game/model"
// 	"app/pkg/game/event"
	//"app/pkg/game/data"
	"app/pkg/game/event/network"	
	"strconv"
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
		for _, obj := range zone.GetAllWorldObjects() {
			if obj.X == x && obj.Y == y {						
				added := -1								
				if obj.Type == "food" {						
					a.Actor.GetClient().In <- network.NewServerMessageEvent("Get " +  obj.Name + "!")				
					added = a.Actor.AddFood(obj.Quantity)																									
				}	
				if obj.Type == "gem" {
					a.Actor.GetClient().In <- network.NewServerMessageEvent("Get " +  obj.Name + "!")									
					added = a.Actor.AddGems(obj.Quantity)															
				}	
				if obj.Type == "silver" {	
					a.Actor.GetClient().In <- network.NewServerMessageEvent("Get " +  obj.Name + "!")								
					added = a.Actor.AddSilver(obj.Quantity)						
				}					
				if added == 0 {
					a.Actor.GetClient().In <- network.NewServerMessageEvent("You can't carry more " + obj.Name + "!")
					return true
				} else if added > 0 {
					fmt.Println("Added ", added, " ", obj.Type)
					message := "You got " + strconv.Itoa(added) + " " + obj.Name + "."
					a.Actor.GetClient().In <- network.NewServerMessageEvent(message)				
					if added == obj.Quantity {
						zone.RemoveWorldObjectByUUID(obj.UUID) 
					} else {
						obj.Quantity = obj.Quantity - added
					}
					return true
				}	
			}
		}
		a.Actor.GetClient().In <- network.NewServerMessageEvent("Get. Nothing to get!")
		return true
	}


	return false

}