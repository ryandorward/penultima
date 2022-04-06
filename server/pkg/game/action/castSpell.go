package action
import (
	"app/pkg/game/model"
	// "fmt"
	"app/pkg/game/event/network"
	// "app/pkg/game/util"
)

type CastSpellAction struct {
	Caster model.Entity
	Spell  string
}


func (a *CastSpellAction) Execute() bool {

	if (a.Spell == "rel por bori") {
		a.Caster.GetClient().In <- network.NewServerMessageEvent("Blink north!")	
	} else if (a.Spell == "rel por ori") {
		a.Caster.GetClient().In <- network.NewServerMessageEvent("Blink east!")	
	} else if (a.Spell == "rel por ozi") {
		a.Caster.GetClient().In <- network.NewServerMessageEvent("Blink south!")	
	} else if (a.Spell == "rel por ocsi") {
		a.Caster.GetClient().In <- network.NewServerMessageEvent("Blink west!")	
	} else {
		a.Caster.GetClient().In <- network.NewServerMessageEvent("Cast " + a.Spell +"! Wvooshabadabada! ... it fizzzles.")	
	}
	return true
}