package entity

import (
	"app/pkg/game/event"
	"app/pkg/game/model"
	"app/pkg/game/util"
	serverModel "app/pkg/model"
	"github.com/google/uuid"
	"app/pkg/fov"	
	"app/pkg/game/data"
	"app/pkg/game/event/network"
	"app/pkg/store"
	// "fmt" 
)

type entityData struct {
	UUID uuid.UUID        `json:"uuid"`
	Name string           `json:"name"`
	Tile int8              `json:"tile"` // representing tile
	Type model.EntityType `json:"type"`

	Stats model.Stats `json:"stats"`

	X int `json:"x"`
	Y int `json:"y"`

	EnergyThreshold int `json:"-"`
	Energy          int `json:"-"`

	QueuedAction model.Action `json:"-"`
	zone         model.Zone   `json:"-"`

	slowThresh float64	// diminishing threshold if player was slowed last turn
}
 
func (e *entityData) GetUUID() uuid.UUID {
	return e.UUID
}

func (e *entityData) GetName() string {
	return e.Name
}

func (e *entityData) GetZone() model.Zone {
	return e.zone
}

func (e *entityData) SetZone(z model.Zone) {
	e.zone = z
}

func (e *entityData) GetType() model.EntityType {
	return e.Type
}

func (e *entityData) GetPosition() (int, int) {
	return e.X, e.Y
}

func (e *entityData) IsInView(other model.Entity) bool {	
	
	viewWidth := 15 	// @todo: these shouldn't be hardcoded.. 
	viewHeight := 15	
	halfViewWidth := viewWidth / 2 
	halfViewHeight := viewHeight / 2  
	zoneWidth, zoneHeight := e.GetZone().GetDimensions()

	eX, eY := e.GetPosition()
	oX, oY := other.GetPosition()

	dX := util.WrapDiff(eX,oX,zoneWidth)
	dY := util.WrapDiff(eY,oY,zoneHeight)

	return ( (dX < halfViewWidth+2) && (dY < halfViewHeight+2) ) // +2 becasuse: 1 for when something goes from on to off the screen or vice-versa, and the other 1 must be for rounding	
}

func (e *entityData) SetPosition(x, y int) {
	e.X = x
	e.Y = y
	store.SetStoredLocation(x, y, e.GetUUID())
}

func (e *entityData) GetClient() *serverModel.Client {
	return nil
}

func (e *entityData) Act() model.Action {
	return nil // NOP
}

func (e *entityData) Tick() bool {
	e.Energy++
	if e.Energy >= e.EnergyThreshold {
		e.Energy = 0
		return true
	}
	return false
}

func (e *entityData) Die() {
	event.NotifyObservers(event.DieEvent{Entity: e})
	e.zone.RemoveEntity(e)
}

func (e *entityData) GetStats() model.Stats {
	return e.Stats
}

// TakeDamage returns if they would die so XP can be dished out
func (e *entityData) TakeDamage(damage int) bool {
	e.Stats.HP -= damage
	return e.Stats.HP <= 0
}

func (e *entityData) GainExp(xp int) {
	e.Stats.XP += xp
	nextLevelXP := util.XPForLevel(e.Stats.Level)
	for e.Stats.XP >= nextLevelXP {
		e.Stats.Level += 1
		e.Stats.MaxHP += util.Roll{8, 1, util.Modifier(e.Stats.Constitution)}.Roll()
		e.Stats.HP = e.Stats.MaxHP
		nextLevelXP = util.XPForLevel(e.Stats.Level)
	}
	e.Stats.XPToNextLevel = nextLevelXP
}

func (e *entityData) Heal(amount int) {
	e.Stats.HP += amount
	if e.Stats.HP > e.Stats.MaxHP {
		e.Stats.HP = e.Stats.MaxHP
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

func (e *entityData) SetTile(t int8) {
	e.Tile = t
	store.SetStoredAvatar(t, e.GetUUID())
}

func (e *entityData) GetTile()int8 {
	return e.Tile
}

func (e *entityData) SetSlowThresh(t float64) {
	e.slowThresh = t
}

func (e *entityData) GetSlowThresh() float64 {
	return e.slowThresh
}

func (e *entityData) UpdateOwnView(c *serverModel.Client) {
	
	// fmt.Println("game/entity/entity : updating view ", c)
	
	// @todo: should be able to remove newMoveEvent, won't be using it
	// c.In <- newMoveEvent(v.Entity, v.X, v.Y)

	viewWidth := 15
	viewHeight := 15																				
	grid := make([][]*model.Tile, viewHeight) // initialize a slice of viewHeight slices		
	for i:=0; i < viewWidth; i++ {					
		grid[i] = make([]*model.Tile, viewWidth) // initialize a slice of viewWidth in in each of viewHeight slices
	}
																									
	entityX, entityY := e.GetPosition()
	zoneWidth, zoneHeight := e.GetZone().GetDimensions()				
	halfViewWidth := viewWidth / 2
	halfViewHeight := viewHeight / 2
	
	for x := 0; x < viewWidth; x++ {	
		for y := 0; y < viewHeight; y++ {	
			nX := util.WrapMod(entityX+x - halfViewWidth, zoneWidth)
			nY := util.WrapMod(entityY+y - halfViewHeight, zoneHeight)										
			grid[x][y] = e.GetZone().GetTile(nX, nY)										
		} 
	}

	gridmap := data.GridMap{
		Grid: grid,
	}

	// Calculate Field Of View
	fovCalc := fov.New()	
	sunlight := e.GetZone().GetSunlight()
	fovCalc.Compute(gridmap, halfViewWidth, halfViewHeight, sunlight)

	fov := make([][]int8, viewHeight) // initialize a slice of viewHeight slices											
	for i:=0; i < viewWidth; i++ {
		fov[i] = make([]int8, viewWidth) // initialize a slice of viewWidth in in each of viewHeight slices				
	}

	// Creat final client view with visible tiles
	for x :=0; x < viewWidth; x++ {	
		for y := 0; y < viewHeight; y++ {	
			if ! fovCalc.IsVisible(x, y) {
				fov[x][y] = 0
			}	else {
				fov[x][y] = int8(grid[x][y].ID)
			}
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
			fov[nX][nY] = other.GetTile() // @todo: need to look this up from the entity later
		}
	}

	// next calculate FOV

	// c.In <- newMoveEvent(v.Entity, v.X, v.Y) // c.In handled in own goroutine, passes event to client websocket
	c.In <- network.NewUpdateOwnViewEvent(&fov)			
}

func (e *entityData) SetEntityData(eS *store.EntityStore) {
	e.Tile = eS.Avatar
	e.Name = eS.Name
	e.X = eS.X
	e.Y = eS.Y
}