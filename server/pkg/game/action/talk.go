package action
import (
	"app/pkg/game/model"
	// "app/pkg/game/data"
	"app/pkg/game/event/network"	
	"fmt"
	"strings"
	"regexp"
	"app/pkg/game/tiles"
)

type TalkAction struct {
	Actor model.Entity	
	X int
	Y int
	Message string
}

func (a *TalkAction) Execute() bool {	
	eX, eY := a.Actor.GetPosition()
	tX := eX + a.X
	tY := eY + a.Y
	tile := a.Actor.GetZone().GetTile(tX, tY)
	message := strings.TrimSpace(a.Message)

	for _, otherE := range a.Actor.GetZone().GetEntities() {
		otherX, otherY := otherE.GetPosition()
		if a.Actor != otherE && otherX == tX && otherY == tY {
			// someone is there, talk to them	
			if (message == "") {
				a.Actor.GetClient().In <- network.NewServerMessageEvent("You are talking to " + otherE.GetName() + ":" + message)
						
				otherE.ReceiveMessage(a.Actor.GetName() + " is talking to you.")

			} else {
				response := otherE.ReceiveMessage(a.Actor.GetName() + " says: " + message)
				a.Actor.GetClient().In <- network.NewServerMessageEvent(otherE.GetName() + " says: " + response)
			}			
			return true
		}
	} 

	if (*tile == tiles.Tiles["shallow_water"]) {
		fmt.Println("shallow water dialogue",a.Message)			
		match, _ := regexp.MatchString("water spirit", message)
		if (match) {
			a.Actor.GetClient().In <- network.NewServerMessageEvent("we flow ~ sometimes we are still ~ sometimes we wave")
		} else if (a.Message != "") {
			a.Actor.GetClient().In <- network.NewServerMessageEvent("we don't know about " + a.Message)
		}	else {
			a.Actor.GetClient().In <- network.NewServerMessageEvent("hello ~ we are water spirit")
		}
	} else {
		a.Actor.GetClient().In <- network.NewServerMessageEvent("(You can't talk to " + tile.Name + ".)")
	}
	a.Actor.GetClient().In <- network.NewServerMessageEvent("You say:")	
	return true
}