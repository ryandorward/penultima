package game

import (
	"sync"
	"time"
	"fmt"
	"github.com/google/uuid"
	"app/pkg/fov"
)

type Player struct {
	ID   string	
	others []*Player
	Location Coord
	TerrainView [][]int8
	AnimateView [][]int8	// view of animating layer, sits "on top" of terrain
	mutex sync.Mutex
	WorldMap *WorldMap
	PlayerCook playerCookie
}

const ViewWidth int = 15
const ViewHeight int = 15

var halfViewWidth int = ViewWidth/2 // integer division gives floor, which is what we want
var halfViewHeight int = ViewHeight/2

type playerCookie struct {
	Avatar int8    `json:"avatar"`
	Name string `json:"name"`		
}

func (p *Player) CanSee(other *Player) bool {
	return true
}

func NewPlayer(cookieStuff playerCookie) *Player {

	// Initialize client TerrainView slices
	view := make([][]int8, ViewHeight)       // initialize a slice of viewHeight slices
	for i:=0; i < ViewHeight; i++ {
			view[i] = make([]int8, ViewWidth)  // initialize a slice of viewWidth in in each of viewHeight slices
	}

	// Initialize client animalView slices
	animal := make([][]int8, ViewHeight) // initialize a slice of viewHeight slices
	for i:=0; i < ViewHeight; i++ {
			animal[i] = make([]int8, ViewWidth)  // initialize a slice of viewWidth in in each of viewHeight slices
	}

	player := &Player{ 	
		ID: uuid.New().String(),
		Location: Coord{							 
			X: 244, // britain
			Y: 148, // britain			 
		},
		TerrainView: view, 
		AnimateView: animal,
		PlayerCook: cookieStuff,    
	}

	player.SetWorldView() // initialize player's view
	return player
}


func (p *Player) SetWorldView() ([]*Player){	

	aa := time.Now()
	
	xStart := WrapMod((p.Location.X - halfViewWidth),WorldWidth);
	yStart := WrapMod((p.Location.Y - halfViewHeight),WorldHeight);
	
	board := p.WorldMap.Grid
	 
	for x := 0; x < ViewWidth; x++ {	
		xMap := WrapMod((xStart + x), WorldWidth)
		for y := 0; y < ViewHeight; y++ {			
			yMap := WrapMod((yStart + y), WorldHeight)			
			p.TerrainView[x][y] = board[xMap][yMap]
			p.AnimateView[x][y] = 0 // while we're at it, clear the AnimalView
		}
	}	

	grid := make([][]int8, ViewHeight) 

	// Initialize Gridmap - this might be funky, flipping the axes?
	for x := 0; x < ViewWidth; x++ {		
		grid[x] = make([]int8, ViewHeight)  
		for y := 0; y < ViewHeight; y++ {						
			grid[x][y]=p.TerrainView[x][y]
		}
	}

	gridmap := gridMap{
		grid: grid,
	}
 
	// Calculate Field Of View
	fovCalc := fov.New()	
	fovCalc.Compute(gridmap, halfViewWidth, halfViewHeight, 10)
	// Update View with visible tiles
	for x :=0; x < ViewWidth; x++ {	
		for y := 0; y < ViewHeight; y++ {	
			if ! fovCalc.IsVisible(x, y) {
				p.TerrainView[x][y] = 0
			}			
		} 
	}

	updatableOthers := []*Player{}
 
	// Add the other players onto the worldview, keep track of which ones are visible, they will need to be updated too
	for _, other := range p.others {	
		if (other.ID == p.ID) { continue }
		px := other.Location.X
		py := other.Location.Y   		
		xrel := WrapMod(px - xStart, WorldWidth) // some head breaking math to figure out if other player is within the view
		yrel := WrapMod(py - yStart, WorldHeight)
		
		// if other is within view and not on a shadowed tile
		if ( (xrel < ViewWidth) && (yrel < ViewHeight) && p.TerrainView[xrel][yrel] != 0) {	
			// c.TerrainView[xrel][yrel] = -1	
			p.AnimateView[xrel][yrel] = other.PlayerCook.Avatar
			fmt.Println("Adding to animated layer" , other.PlayerCook.Avatar)
			
			updatableOthers = append(updatableOthers,other)
			
		}				 
	}	

	bb := time.Now()
	fmt.Println("Calculate view time: ", float64(bb.Nanosecond() - aa.Nanosecond()) / 1e9)

	return updatableOthers;
	  
}	



