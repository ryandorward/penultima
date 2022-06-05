package entity

import (
	// "app/pkg/game/event"
	gameModel "app/pkg/game/model"
	serverModel "app/pkg/model"
	// "app/pkg/game/util"
	"app/pkg/model"	
	"github.com/google/uuid"
	"math"
	// "strconv"
// 	"app/pkg/game/util"
		"fmt"
// 	"database/sql"
	// "errors"
	// "os"
	_ "github.com/lib/pq"
	"app/pkg/store"
	"app/pkg/game/event"
	
	"app/pkg/fov"	
	// "app/pkg/game/data"
	"app/pkg/game/util"
	"app/pkg/game/event/network"
	"app/pkg/game/tiles"
)

const (
	warriorTileID  = 21
	playerEnergyThreshold = 3
	slowThresh = 1.0 // need to initialize SlowThresh to 1 so that they can start out moving
)


var ActivePlayers = map[uuid.UUID]*Player{}

// player is an extension of entityData
type Player struct {
	entityData // inheritance, basically. entityData comes from server/pkg/game/entity/entity.go
	startingZone gameModel.Zone
	client       *model.Client
	food float64
	health float64
	gems int
}

func NewPlayer(client *model.Client) *Player {
	// try to get this player from storage?
	p := &Player{
		entityData: entityData{
			UUID: client.Account.UUID,
			Name: client.Account.Username,
			Tile: warriorTileID,
			Type: gameModel.EntityTypePlayer,
			EnergyThreshold: playerEnergyThreshold,
			slowThresh: slowThresh,
		},
		client: client,
		food: 22,
		health: 100,
		gems: 8,
	}
	ActivePlayers[p.UUID] = p
	return p
}

func (p *Player) SetClient(client *model.Client) {
	p.client = client
}

func (p *Player) GetClient() *model.Client {
	return p.client
}

func (p *Player) Act() gameModel.Action {
	a := p.QueuedAction
	p.QueuedAction = nil
	return a
}

func (p *Player) Spawn(z gameModel.Zone) {	
	// this is a bit problematic because if the player's last position really was 0,0 they will get transported to 
	// the spawn location. @todo: fix this (or put a mountain or something there!)
	// actually this only happens when a player is "spawned" from the login screen .. i think
	if p.X == 0 && p.Y == 0 {
		for _, obj := range z.GetAllWorldObjects() {					
			if obj.Type == gameModel.WorldObjectTypePlayerSpawn {		
				p.X = obj.X 
				p.Y = obj.Y
				break
			}
		}
	}
	p.startingZone = z 
	z.AddEntity(p)
}

/*
// Despawn is for log off only, not changing zones (TODO: fix, leave vs. despawn)
func (p *player) Despawn() {
	event.NotifyObservers(event.DespawnEvent{Entity: p})
	p.zone.RemoveEntity(p)
	delete(ActivePlayers, p.UUID)
}

func (p *player) Die() {
	event.NotifyObservers(event.DieEvent{Entity: p})
	p.zone.RemoveEntity(p)
	p.rollStats()           // roll new stats cuz they're dead lol
	p.Spawn(p.startingZone) // send em back to the starting zone
	return
}

func (p *player) GainExp(xp int) {
	originalLevel := p.Stats.Level
	p.entityData.GainExp(xp)
	event.NotifyObservers(event.GainXPEvent{Entity: p, LeveledUp: originalLevel != p.Stats.Level})
}*/

/*
func (p *player) rollStats() {
	p.Stats.Level = 1
	p.Stats.XP = 0
	p.Stats.XPToNextLevel = util.XPForLevel(2)

	// use 3d6 for stats
	r := util.Roll{6, 3, 0} // 3d6 + 0
	p.Stats.Strength = r.Roll()
	p.Stats.Dexterity = r.Roll()
	p.Stats.Constitution = r.Roll()
	p.Stats.Intelligence = r.Roll()
	p.Stats.Wisdom = r.Roll()
	p.Stats.Charisma = r.Roll()

	// hit dice for players is a d8, so HP = 8 + CON (1d8 + CON on level)
	p.Stats.MaxHP = 8 + util.Modifier(p.Stats.Constitution)
	if p.Stats.MaxHP <= 0 {
		p.Stats.MaxHP = 1
	}
	p.Stats.HP = p.Stats.MaxHP

	p.Stats.AC = 10 + util.Modifier(p.Stats.Dexterity)
}

*/

func (p *Player) SetZoneWithTarget(z gameModel.Zone, tx, ty int) {
 	x, y := p.GetPosition()
	p.GetZone().RemoveEntity(p)
	event.NotifyObservers(event.MoveEvent{Entity: p, X: x, Y: y}) // @todo this must need a zone? Or does this only happen in the current zone?
	z.AddEntity(p)					
	p.SetZone(z)
	p.SetPosition(tx, ty)
	event.NotifyObservers(event.MoveEvent{Entity: p, X: tx, Y: ty})	 	
}


func (e *Player) UpdateOwnView(c *serverModel.Client) {
	
	viewWidth := 15
	viewHeight := 15																				
	grid := make([][]*gameModel.Tile, viewHeight) // initialize a slice of viewHeight slices		
	for i:=0; i < viewWidth; i++ {					
		grid[i] = make([]*gameModel.Tile, viewWidth) // initialize a slice of viewWidth in in each of viewHeight slices
	}
																									
	entityX, entityY := e.GetPosition()
	zone := e.GetZone()
	zoneWidth, zoneHeight := zone.GetDimensions()				
	halfViewWidth := viewWidth / 2
	halfViewHeight := viewHeight / 2
	
	for x := 0; x < viewWidth; x++ {	
		for y := 0; y < viewHeight; y++ {	
			nX := entityX+x - halfViewWidth
			nY := entityY+y - halfViewHeight							
			if zone.GetTorroidal() {
				nX = util.WrapMod(nX, zoneWidth)
				nY = util.WrapMod(nY, zoneHeight)	
			}			
			grid[x][y] = e.GetZone().GetTile(nX, nY)												
		} 
	}

	gridmap := fov.MyGridMap{
		Grid: grid,
	}

	// Calculate Field Of View
	fovCalc := fov.New()	
	sunlightRange := e.GetZone().GetSunlight()
	fovCalc.Compute(gridmap, halfViewWidth, halfViewHeight, sunlightRange)

	fov := make([][]int, viewHeight) // initialize a slice of viewHeight slices											
	for i:=0; i < viewWidth; i++ {
		fov[i] = make([]int, viewWidth) // initialize a slice of viewWidth in in each of viewHeight slices				
	}

	// new fov algo:
	// 1. calc fov without sunlight, so full range of viewport 
	// 2. find any gems/lightsources withing the fov
	// 3. for each lightsource, calc another fov, but this time with the lightsource's range, and centered around lightsource
	// 4. calc another fov with the sunlight range
	// 5. OR all the FOVs, except the first one. 
	// there is a lot of waste here, should be a way to optimize. 
	// definitely a better way...  

	// how about this:
	// 1. calc fov without sunlight, so full range of viewport
	// 2. find any gems/lightsources withing the fov. 

	// Create final client view with visible tiles
	for x :=0; x < viewWidth; x++ {	 
		for y := 0; y < viewHeight; y++ {	
			if ! fovCalc.IsVisible(x, y) {
				fov[x][y] = 0
			}	else {
				fov[x][y] = int(grid[x][y].ID)
			}
		} 
	}

	// add in world objects!
	for _, obj := range e.GetZone().GetAllWorldObjects() {							
		nX := obj.X - entityX + halfViewWidth
		nY := obj.Y - entityY + halfViewHeight			
		nX = util.WrapMod(nX, zoneWidth)
		nY = util.WrapMod(nY, zoneHeight)

		// crude way to make sure in range. Later we won't need to do this 
		// because we'll actually be calculating if this entity is in view
		// (currently just returning true for that check)
		if nX < viewWidth && nY < viewHeight && nX >= 0 && nY >= 0 && fov[nX][nY] != 0 {																					
			fov[nX][nY] = tiles.Tiles[obj.Tile].ID
		}				
	}

	// now add the avatars of the entities in view. Other _actually_ includes self
	//for _, other := range inView {
	for _, other := range e.GetZone().GetEntities() {
		otherX, otherY := other.GetPosition()
		nX := otherX - entityX + halfViewWidth
		nY := otherY - entityY + halfViewHeight			
		nX = util.WrapMod(nX, zoneWidth)
		nY = util.WrapMod(nY, zoneHeight)

		// crude way to make sure in range. Later we won't need to do this 
		// because we'll actually be calculating if this entity is in view
		// (currently just returning true for that check)
		if nX < viewWidth && nY < viewHeight && nX >= 0 && nY >= 0 && fov[nX][nY] != 0 {
			fov[nX][nY] = other.GetTile()
		}
	}

	c.In <- network.NewUpdateOwnViewEvent(&fov)			
}

func (p *Player) SetTile(t int) {
	p.entityData.SetTile(t)
	store.SetStoredAvatar(t, p.GetUUID())
}

func (p *Player) UpdateOwnStats() {
	p.GetClient().In <- network.NewUpdateStatsEvent(p.GetName())
	p.GetClient().In <- network.NewStatEvent("gems",p.gems)
	p.GetClient().In <- network.NewStatEvent("food",int(math.Floor(p.food)))
	p.GetClient().In <- network.NewStatEvent("health",int(math.Floor(p.health)))
}	

func (p *Player) UpdateClientStat(name string, value int) {
	p.GetClient().In <- network.NewStatEvent(name, value)	
}	

func (p *Player) SetName(n string) { 	
	p.entityData.SetName(n)
	store.SetStoredName(n, p.GetUUID())  
}

func (p *Player) SetZone(z gameModel.Zone) {
	p.entityData.SetZone(z) 
	store.SetStoredZone(z.GetName(), p.GetUUID() )
}

func (p *Player) SetPosition(x, y int) {
	p.entityData.SetPosition(x, y)
	store.SetStoredLocation(x, y, p.GetUUID()) 
}

func (p *Player) ReceiveMessage(m string) string {
	p.GetClient().In <- network.NewServerMessageEvent(m)
	return ""
}

func (p *Player) ReceiveResult(msg string, code string) {
	p.GetClient().In <- network.NewServerResultEvent(msg, code)
}

func (p *Player) Tick() bool {	
	p.updateFood()
	return p.entityData.Tick()
}

func (p *Player) updateFood() { 
	oldFood := p.food
	p.food = p.food - 0.001	
	if p.food <= 0 {
		p.food = 0
		oldHealth := p.health
		p.health = p.health - 0.001
		if math.Floor(oldHealth) - math.Floor(p.health) > 0 {	
			p.ReceiveResult("You are starving. You need to get some food" , "status")	
			p.GetClient().In <- network.NewStatEvent("health",int(math.Floor(p.health)))
		}
	}
	if math.Floor(oldFood) - math.Floor(p.food) > 0 {	
		p.updateClientFood()
	}
}

// update client with current food
func (p *Player) updateClientFood() { 
	foodInt := int(math.Floor(p.food))
	foodStr := fmt.Sprintf("%d", foodInt)
	p.ReceiveResult("You have " + foodStr + " food." , "food")	
	p.GetClient().In <- network.NewStatEvent("food",foodInt)
}

func (p *Player) AddFood(food float64) { 
	p.food += food
	p.updateClientFood()
	// p.entityData.AddFood(food)
}

func (p *Player) AddGems(gems int) { 
	p.gems += gems
	p.GetClient().In <- network.NewStatEvent("gems",p.gems)
}

func (p *Player) GetGems() int {
	return p.gems
}

