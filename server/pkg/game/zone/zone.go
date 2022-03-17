package zone

import (
	"app/pkg/game/model"
	"github.com/google/uuid"
	"fmt"
	"time"
	"app/pkg/game/event/network"
	"math/rand"
)

type Zone struct {
	UUID   uuid.UUID    `json:"uuid"`
	Name   string       `json:"name"`
	Width  int          `json:"width"`
	Height int          `json:"height"`
	Tiles  [][]model.Tile `json:"tiles"`

	Entities     map[string]model.Entity       `json:"entities"`
	WorldObjects map[string]*model.WorldObject `json:"world_objects"`

	Torroidal bool `json:"torroidal"`

	// Time int `json:"time"` // 360 ticks per day

	sunlight int `json:"sunlight"`	// int 1-15

	trammelPhase int `json:"trammel_phase"`
	feluccaPhase int `json:"felucca_phase"`

	wind struct {
		x int
		y int
	}
			
	// Time int `json:"time"` // time of the day, in xoxarian seconds. See: https://www.uoguide.com/Time#:~:text=Each%20day%20is%20divided%20into,from%206%20AM%20until%20midnight
}

func (z *Zone) GetUUID() uuid.UUID {
	return z.UUID
}

var secondsPerDay = 1800 // Xoxarian day is 1800 seconds long!

// @todo implement this properly - taking into account wrap around if current zone is torroidal
func (z *Zone) GetNewLocation(x,y,dx,dy int) (int, int) {
	
	nx := x + dx
	ny := y + dy

	if z.Torroidal {		
		if nx < 0 {
			nx = z.Width - 1
		}
		if ny < 0 {	
			ny = z.Height - 1
		}
		if nx >= z.Width {
			nx = 0
		}
		if ny >= z.Height {	
			ny = 0
		}		
	} else {
		if nx < 0 || ny < 0 || nx >= z.Width || ny >= z.Height {
			return x, y
		}
	}
	return nx, ny
}

func (z *Zone) GetDimensions() (int, int) {
	return z.Width, z.Height
}

func (z *Zone) GetTile(x, y int) *model.Tile {
	if x < 0 || y < 0 || x >= z.Width || y >= z.Height {
		return nil
	}

	// index := (z.Width * y) + x
	return &z.Tiles[x][y]
	// return &z.Tiles[index]
}

func (z *Zone) GetEntities() (entities []model.Entity) {
	for _, e := range z.Entities {
		entities = append(entities, e)
	}
	return
}

func (z *Zone) GetWorldObjects(x, y int) []*model.WorldObject {
	objs := []*model.WorldObject{}
	for _, obj := range z.WorldObjects {
		if obj.X == x && obj.Y == y {
			objs = append(objs, obj)
		}
	}
	return objs
}

func (z *Zone) GetAllWorldObjects() []*model.WorldObject {

	objs := []*model.WorldObject{}
	for _, obj := range z.WorldObjects {
		objs = append(objs, obj)
	}
	return objs
}

func (z *Zone) AddEntity(e model.Entity) {
	z.Entities[e.GetName()] = e
	e.SetZone(z)
}

func (z *Zone) RemoveEntity(e model.Entity) {
	delete(z.Entities, e.GetName())
}

func (z *Zone) Update(dt float64) {

	actions := []model.Action{}
	for _, e := range z.Entities {
		// 	fmt.Println("update entity: ", e.GetName())
		if e.Tick() {
			a := e.Act()
			if a != nil {
				actions = append(actions, a)
			}
		}
	}

	for _, a := range actions {
		a.Execute()
	}
}

func serverDaySeconds() int {
	t := time.Now()
	year, month, day := t.Date()
	t2 := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return int(t.Sub(t2).Seconds())
}

func xoxariaTimeSeconds() int {
	// there are 1800 seconds in a xoxarian day! 
	t := serverDaySeconds()
	xt := t % secondsPerDay
	// fmt.Println("Current xoxarian time of day in seconds: ", xt)
	return xt
}

func (z *Zone) updateSun() {
	xt := xoxariaTimeSeconds() // range 0-1799
	var sunlight int
	var hourLength = secondsPerDay/16
	if xt < secondsPerDay/2 {
		sunlight = int(2*xt / hourLength) 
	} else {
		sunlight  = int(16 - (2*(xt - secondsPerDay/2) / hourLength)) 
	}
	// fmt.Println("Coeffs:",secondsPerDay, hourLength)
	// fmt.Println("Pre-clipped xoxarian sunlight: ", sunlight)
	sunlight += 2
	if sunlight > 10 {
		sunlight = 10
	} else if sunlight < 2{
		sunlight = 2
	}
	// fmt.Println("Final sunlight: ", sunlight)
	if sunlight != z.sunlight {
		z.sunlight = sunlight
		fmt.Println("Updating sunlight: ", sunlight)
		for _, e := range z.GetEntities() {
			if c := e.GetClient(); c != nil {	
				e.UpdateOwnView(c)
			}		
		}
	}
		
}

// period in seconds
func (z *Zone) getMoonPhase(period int) int {
	st := serverDaySeconds() + 3600*8
	phase := st % period // for an 1800s day, this will give us a number [0,16 200)
	phase = phase * 8 / period // convert to [0,8)
	return phase
}

func (z *Zone) updateMoonPhase() {
	for _, e := range z.GetEntities() {
		if c := e.GetClient(); c != nil {	
			fmt.Println("Trammel, Felucca ", z.trammelPhase,z.feluccaPhase)
			c.In <- network.NewMoonPhaseEvent(z.trammelPhase,z.feluccaPhase)
		}		
	}
}

// Trammel's synodic period is 9 days
func (z *Zone) updateTrammelPhase() {
	// period := 9 * secondsPerDay
	period := 5*secondsPerDay/3
	tramelPhase := z.getMoonPhase(period)
	if tramelPhase != z.trammelPhase {	
		z.trammelPhase = tramelPhase
		z.updateMoonPhase()
	}
}

// Felucca's synodic period is 14
func (z *Zone) updateFeluccaPhase() {
	// period := 14 * secondsPerDay
	period := 5*secondsPerDay/24
	feluccaPhase := z.getMoonPhase(period)
	if feluccaPhase != z.feluccaPhase {	
		z.feluccaPhase = feluccaPhase
		z.updateMoonPhase()
	}
}

func windClip(w int) int {
	if w < -1 {
		return -1
	} else if w > 1 {
		return 1
	}
	return w
}

func (z *Zone) updateClientsWind() {
	for _, e := range z.GetEntities() {
		if c := e.GetClient(); c != nil {			
			c.In <- network.NewWindEvent(z.wind.x,z.wind.y)
		}		
	}
}

func (z *Zone) updateWind() {
	rand.Seed(time.Now().UnixNano())
	windD := rand.Intn(8)
	windx := z.wind.x
	windy := z.wind.y
	
	if windD > 3 { 
		// only change wind half the time
		return
	}
	// should only be equal when no wind
	if windx == windy {
		// if wind is still, change it to one of the 4 cardinals
		switch windD {
			case 0:
				windx = 1
				windy = 0
			case 1:
				windx = -1
				windy = 0
			case 2:
				windx = 0
				windy = 1
			case 3:
				windx = 0
				windy = -1			
		}	
	} else if windD == 3 {
		// if wind is not still, there is 1 1/4 chance it just continues in same direction. Change nothing
		return
	} else {
		windD -= 1 // windD is now -1,0,1 - 0 means move to no wind, -1,1 mean swap x,y, -1 means flip cardinality
		if windx != 0 {			
			windy = windx*windD
			windx = 0
		} else {
			windx = windy*windD
			windy = 0
		}
	}

	if windx != z.wind.x || windy != z.wind.y {
		z.wind.x = windx
		z.wind.y = windy
		fmt.Println("Updating wind: ", z.wind)
		z.updateClientsWind()
	}

}

func (z *Zone) GetSunlight() int {
	return z.sunlight
}
func (z *Zone) GetFelucca() int {
	return z.feluccaPhase
}
func (z *Zone) GetTrammel() int {
	return z.trammelPhase
}
func (z *Zone) GetWind() (int,int) {
	return z.wind.x, z.wind.y
}

// slow updates for zones
func (z *Zone) SlowUpdate() {
	z.updateSun()	
	z.updateTrammelPhase()
	z.updateFeluccaPhase()
	z.updateWind()
}