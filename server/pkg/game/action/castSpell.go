package action
import (
	"app/pkg/game/model"
	"app/pkg/game/event"
	"fmt"
	"app/pkg/game/event/network"
	"regexp"
	"strconv"
//	"strings"
	// "app/pkg/game/util"
	
)

type CastSpellAction struct {
	Caster model.Entity
	Spell  string
} 

func (a *CastSpellAction) Execute() bool {

	matchTransPort, _ := regexp.MatchString(`trans port:?\([0-9]{1,3},[0-9]{1,3}\)`, a.Spell)		
	matchWheraAmI, _ := regexp.MatchString(`wer ami`, a.Spell)

	if (a.Spell == "rel por bori") {
		a.Caster.GetClient().In <- network.NewServerMessageEvent("Blink north!")	
	} else if (a.Spell == "rel por ori") {
		a.Caster.GetClient().In <- network.NewServerMessageEvent("Blink east!")	
	} else if (a.Spell == "rel por ozi") {
		a.Caster.GetClient().In <- network.NewServerMessageEvent("Blink south!")	
	} else if (a.Spell == "rel por ocsi") { 
		a.Caster.GetClient().In <- network.NewServerMessageEvent("Blink west!")		
	} else if (matchWheraAmI) {
		cX, cY := a.Caster.GetPosition()
		a.Caster.GetClient().In <- network.NewServerMessageEvent("(" + strconv.Itoa(cX) +", " + strconv.Itoa(cY) +")")	
	} else if (matchTransPort) {
	
		re := regexp.MustCompile(`[0-9]{1,3}`)
		coords := re.FindAllString(a.Spell, -1) 
		fmt.Println(coords)

		a.Caster.GetClient().In <- network.NewServerMessageEvent("Trans Port!")	
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		fmt.Println(x,y)
		a.Caster.SetPosition(x, y)
		a.Caster.SetSlowThresh(1.0) // mover may have been stuck on a non-traversible tile before this
		event.NotifyObservers(event.MoveEvent{Entity: a.Caster, X: x, Y: y})
	}	else {
		a.Caster.GetClient().In <- network.NewServerMessageEvent("Cast " + a.Spell +"! Wvooshabadabada! ... it fizzzles.")	
	}
	return true
}