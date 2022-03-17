package game

import (
	"time"	
	/*
	"github.com/floralbit/dungeon/game/data"
	"github.com/floralbit/dungeon/game/dungeon"
	"github.com/floralbit/dungeon/game/entity"
	"github.com/floralbit/dungeon/game/event"
	"github.com/floralbit/dungeon/game/event/network"
	"github.com/floralbit/dungeon/game/zone"
	*/

	"app/pkg/game/data"
	"app/pkg/game/event"
	"app/pkg/game/event/network"
	"app/pkg/game/entity"
	"app/pkg/game/action"
	"app/pkg/game/zone"
	"app/pkg/model"
	gamemodel "app/pkg/game/model" 
	"app/pkg/game/util"
	// "github.com/google/uuid"
	"fmt"	
	"app/pkg/store"
)

const tickLength = 60 // 100 // in ms 
const eventBufferSize = 256

// var startingZoneUUID = uuid.MustParse("10f8b073-cbd7-46b7-a6e3-9cbdf68a933f")

var startingZoneUUID = "xoxaria"

// In ...
var In = make(chan model.ClientEvent, eventBufferSize)

var zones = map[string]*zone.Zone{}

// Run ...
func Run() {  

	util.PrettyPrint("Run!")
	event.Observers = append(event.Observers, network.NewObserver())

	zones = data.LoadZones()

	fmt.Println("The zones are:")
	// print all the keys of zones
	for k, z := range zones {
		z.SlowUpdate() // trigger initialization/update of states like sunlight, wind, tides
		fmt.Println(k)	
	}

	// util.PrettyPrint(zones["xoxaria"])
	// util.PrettyPrint(zones)

	// build the dungeon!
	/*
	dungeonFloors := dungeon.BuildDungeon(zones[startingZoneUUID])
	for _, floor := range dungeonFloors {
		zones[floor.UUID] = floor
	}
	*/

	ticker := time.NewTicker(tickLength * time.Millisecond)
	lastTime := time.Now()

	// this is the interval for slow updates like weather, wind, moons, sun
	// maybe monster generation, etc.
	slowTicker := time.NewTicker(5 * time.Second)

	for {
		select {
		case now := <-ticker.C:
			dt := now.Sub(lastTime).Seconds()
			lastTime = now		
			update(dt) // update all zones
		case e := <-In:
			// fmt.Println("game run loop ", e )
			processEvent(e) 
		case <-slowTicker.C:
			// fmt.Println("slow tick")
			slowUpdate() // update all zones
		}
	}
}

func update(dt float64) { // dt in seconds	
	for _, z := range zones {
		// util.PrettyPrint(z.Entities)
	  // util.PrettyPrint(dt)
		z.Update(dt)
	}	
}

func slowUpdate() {
	for _, z := range zones {
		z.SlowUpdate()
	}	
}

func processEvent(e model.ClientEvent) {

	if (e.Join != nil) {
		fmt.Println("game/processEvent: join")
		handleJoinEvent(e)
	}

	p, ok := entity.ActivePlayers[e.Sender.Account.UUID]
	if !ok {
		util.PrettyPrint("not an active player")
		return // ignore inactive players
	}

	switch {
		case e.Move != nil:			
			handleMoveEvent(e,p)	
		case e.UpdateAvatar != nil:
			fmt.Println("game/processEvent: update avatar: ", e.UpdateAvatar.ID)
			handleUpdateAvatar(e,p)
		case e.PeerGem != nil:
			fmt.Println("game/processEvent: peer gem ")		
			p.QueuedAction = &action.PeerGemAction{
				Peerer: p,
			}	
		default:
			fmt.Println("game/processEvent: default")
			fmt.Println(e)
	}	
	/* @todo eventually we will handle events other than e.Move! 
	switch {
		case e.Move != nil:
			handleMoveEvent(e)		
		case e.Join != nil:
			handleJoinEvent(e)
		case e.Leave != nil:
			handleLeaveEvent(e)
		case e.Chat != nil:
			handleChatEvent(e)		
		case e.Attack != nil:
			handleAttackEvent(e)		
	}
	*/	
}

func handleUpdateAvatar(e model.ClientEvent, p *entity.Player) {
	p.QueuedAction = &action.UpdateAvatarAction{
		Mover: p,
		Id: e.UpdateAvatar.ID,
	}	
}

func handleMoveEvent(e model.ClientEvent,  p *entity.Player) {
	p.QueuedAction = &action.MoveAction{
		Mover: p,
		X:     e.Move.X,
		Y:     e.Move.Y,
	}
}

func handleJoinEvent(e model.ClientEvent) {
	activePlayer, ok := entity.ActivePlayers[e.Sender.Account.UUID]	
	if ok {
		util.PrettyPrint("player already logged in!")
		activePlayer.SetClient(e.Sender) // update player with new connection!
		initializeJoin(e, activePlayer)
		return	
	}
	util.PrettyPrint("join") 
	p := entity.NewPlayer(e.Sender) // TODO: pull from storage	
	// event.NotifyObservers(event.JoinEvent{Entity: p}) // don't need to do we?
	p.Spawn(zones[startingZoneUUID]) 	
	initializeJoin(e, p)	
}


func initializeJoin(e model.ClientEvent, p gamemodel.Entity) {	
	zone := p.GetZone()
	e.Sender.In <- network.NewMoonPhaseEvent(zone.GetTrammel(),zone.GetFelucca()) // initialize moons
	e.Sender.In <- network.NewWindEvent(zone.GetWind()) // initialize wind
	eS, err := store.GetStoredEntityData(p.GetUUID()) // initialize stored entity data
	if err != nil {
		util.PrettyPrint("error getting stored entity data")
	} else {
		p.SetEntityData(eS)
	}	
	p.UpdateOwnView(e.Sender) // initialize player's view
}

/*

func handleLeaveEvent(e model.ClientEvent) {
	p, ok := entity.ActivePlayers[e.Sender.Account.UUID]
	if !ok {
		return
	}
	p.Despawn()
}

func handleChatEvent(e model.ClientEvent) {
	p, ok := entity.ActivePlayers[e.Sender.Account.UUID]
	if !ok {
		return // ignore inactive players, TODO: send an error?
	}
	event.NotifyObservers(event.ChatEvent{Entity: p, Message: e.Chat.Message})
}


func handleAttackEvent(e model.ClientEvent) {
	p, ok := entity.ActivePlayers[e.Sender.Account.UUID]
	if !ok {
		return // ignore inactive players
	}
	p.QueuedAction = &action.LightAttackAction{
		Attacker: p,
		X:        e.Attack.X,
		Y:        e.Attack.Y,
	}
}
*/