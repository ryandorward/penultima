package entity

import (
	// "app/pkg/game/event"
	"app/pkg/game/model"
	"app/pkg/game/util"
	serverModel "app/pkg/model"
	"github.com/google/uuid"
	// "app/pkg/fov"	
	// "app/pkg/game/data"
	// "app/pkg/game/event/network"
	"app/pkg/store"	
	"fmt"
	"math/rand"
)

type entityData struct {
	UUID uuid.UUID        `json:"uuid"`
	Name string           `json:"name"`
	Tile int              `json:"tile"` // representing tile
	Type model.EntityType `json:"type"`

	Stats model.Stats `json:"stats"`

	X int `json:"x"`
	Y int `json:"y"` 

	lastMoveTry struct {
		X int `json:"X"`
		Y int `json:"Y"`
	} `json:"lastMoveTry"`

	EnergyThreshold int `json:"-"`
	Energy          int `json:"-"` 

	QueuedAction model.Action `json:"-"`
	zone         model.Zone   `json:"-"`
	zoneName 	 string       `json:"-"`

	slowThresh float64	// diminishing threshold if player was slowed last turn

}

const viewWidth = 15
const viewHeight = 15
 
func (e *entityData) GetUUID() uuid.UUID {
	return e.UUID
}

func (e *entityData) GetName() string {
	return e.Name
}

func (e *entityData) SetName(n string) { 	
	e.Name = n
}

func (e *entityData) GetZone() model.Zone {
	return e.zone
}
func (e *entityData) GetZoneName() string {
	fmt.Println("GetZoneName(): ", e.zoneName)
	return e.zoneName
}
 
func (e *entityData) SetZone(z model.Zone) {
	e.zone = z
}

func (e *entityData) SetZoneWithTarget(z model.Zone, tx, ty int) {
	// override this in player.go	
}

func (e *entityData) GetType() string { // model.EntityType {
	return string (e.Type)
}

func (e *entityData) SetPosition(x, y int) {
	e.X = x
	e.Y = y	
}

func (e *entityData) GetPosition() (int, int) {
	return e.X, e.Y
}

// meaningless fn just to make it comply to Positionable interface
func (e *entityData) SetQuantity(q int) {}

func (e *entityData) SetLastMoveTry(x, y int) {
	e.lastMoveTry.X = x
	e.lastMoveTry.Y = y	
}

func (e *entityData) GetLastMoveTry() (int, int) {
	return e.lastMoveTry.X, e.lastMoveTry.Y
}

func (e *entityData) IsInView(other model.Entity) bool {	
	
	//viewWidth := 15 	// @todo: these shouldn't be hardcoded.. 
	// viewHeight := 15	
	halfViewWidth := viewWidth / 2 
	halfViewHeight := viewHeight / 2  
	zoneWidth, zoneHeight := e.GetZone().GetDimensions()

	eX, eY := e.GetPosition()
	oX, oY := other.GetPosition()

	dX := util.WrapDiff(eX,oX,zoneWidth)
	dY := util.WrapDiff(eY,oY,zoneHeight)

	return ( (dX < halfViewWidth+2) && (dY < halfViewHeight+2) ) // +2 becasuse: 1 for when something goes from on to off the screen or vice-versa, and the other 1 must be for rounding	
}

func (e *entityData) GetClient() *serverModel.Client {
	return nil
}

func (e *entityData) Act() model.Action {
	return nil // NOP
}

func (e *entityData) Tick() bool {
	
	e.checkDie()

	e.Energy++	
	if e.Energy >= e.EnergyThreshold {
		e.Energy = 0
		return true
	}
	return false
}

func (e *entityData) Die() {

	zone := e.GetZone()
	
	// now add a dead body object in that spot
	body := &model.WorldObject{		
		UUID: uuid.New(),
		Name: "Dead Body",		
		Tile: "deadBody",		
		Type: "deadBody",
		X: e.X,
		Y:	e.Y,
	}
	if e.Type == "chicken" {
		body.Name = "Food"
		body.Type = "food"
		body.Tile = "cookedChicken"		
		body.Quantity = rand.Intn(14) + 8
	} else {
		fmt.Println(e.Type)
	}

	zone.AddWorldObject(body)
	zone.RemoveEntity(e)
	e.UpdateOwnView()
}

func (e *entityData) GetStats() model.Stats {
	return e.Stats
}

// TakeDamage returns if they would die so XP can be dished out
func (e *entityData) TakeDamage(damage float64) bool {
	e.Stats.HP -= damage
	return e.Stats.HP <= 0
}

func (e *entityData) GainExp(xp int) {
	e.Stats.XP += xp
	nextLevelXP := util.XPForLevel(e.Stats.Level)
	for e.Stats.XP >= nextLevelXP {
		e.Stats.Level += 1
		e.Stats.MaxHP += util.Roll{8, 1, util.Modifier(e.Stats.Constitution)}.Roll()
		e.Stats.HP = float64(e.Stats.MaxHP)
		nextLevelXP = util.XPForLevel(e.Stats.Level)
	}
	e.Stats.XPToNextLevel = nextLevelXP
}

func (e *entityData) Heal(amount float64) {
	e.Stats.HP += amount
	if e.Stats.HP > float64(e.Stats.MaxHP) {
		e.Stats.HP = float64(e.Stats.MaxHP)
	}
}

func (e *entityData) RollToHit(targetAC int) bool {
	toHit := util.Roll{Sides: 20, N: 1, Plus: util.Modifier(e.Stats.Strength)}.Roll() // TODO: swap modifier based on weapon
	return toHit >= targetAC
}

func (e *entityData) RollDamage() int {
	damage := util.Roll{Sides: 3, N: 1, Plus: util.Modifier(e.Stats.Strength)}.Roll()
	if damage <= 0 {
		damage = 1 // minimum 1 dmg
	}
	return damage 
}

func (e *entityData) SetTile(t int) {
	e.Tile = t	
}

func (e *entityData) GetTile()int {
	return e.Tile
}

func (e *entityData) SetSlowThresh(t float64) {
	e.slowThresh = t
}

func (e *entityData) GetSlowThresh() float64 {
	return e.slowThresh
}

func (e *entityData) UpdateOwnStats() {}	

func (e *entityData) UpdateClientStat(name string, value int) {}

func (e *entityData) UpdateOwnView() {}

func (e *entityData) SetEntityData(eS *store.EntityStore) {
	e.Tile = eS.Avatar
	e.Name = eS.Name
	e.X = eS.X
	e.Y = eS.Y
	e.zoneName = eS.ZoneName	
	fmt.Println("game/entity/entity : setting entity data. Zone is: ", e.zoneName)
}

func (e *entityData) ReceiveMessage(m string) string{
	return ""
}

func (e *entityData) checkDie() {
	if e.Stats.HP <= 0 {		
		e.Die()
	}	
}

func (e *entityData) ReceiveResult(msg string, code string) {}

func (e *entityData) AddFood(food int) int { return 0 }
func (e *entityData) AddSilver(silver int) int { return 0}
func (e *entityData) AddGems(gems int) int {return 0 }
func (e *entityData) GetGemCount() int { return 0 }