package game

import (
	// "sync"
	"time"
	"fmt"
	"github.com/google/uuid"
	"app/pkg/fov"
)

type Player struct {
	ID   string	
	Location Coord
	TerrainView [][]int8
	AnimateView [][]int8	// view of animating layer, sits "on top" of terrain
	WorldMap *WorldMap
	PlayerCook PlayerCookie
}

const ViewWidth int = 15
const ViewHeight int = 15

var halfViewWidth int = ViewWidth/2 // integer division gives floor, which is what we want
var halfViewHeight int = ViewHeight/2

type PlayerCookie struct {
	Avatar int8    `json:"avatar"`
	Name string `json:"name"`		
}

func (p *Player) GetLocation() (Coord){	
	return p.Location;
}	
func (p *Player) GetID() (string){	
	return p.ID;
}	

func (p *Player) CanSee(other *Player) (bool, int, int) {
	xStart := WrapMod((p.Location.X - halfViewWidth),WorldWidth);
	yStart := WrapMod((p.Location.Y - halfViewHeight),WorldHeight);
	loc := other.GetLocation()
	ox := loc.X
	oy := loc.Y 
	xrel := WrapMod(ox - xStart, WorldWidth) // some head breaking math to figure out if other player is within the view
	yrel := WrapMod(oy - yStart, WorldHeight)
	// if other is within view and not on a shadowed tile
	if ( (xrel < ViewWidth) && (yrel < ViewHeight) && p.TerrainView[xrel][yrel] != 0) {						
		return true, xrel, yrel		
	}			
	return false, 0, 0
}

func NewPlayer(cookieStuff PlayerCookie, worldMap *WorldMap) *Player {

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
		WorldMap: worldMap, 		
	}

	player.SetWorldView() // initialize player's view
	return player
}

func (p *Player) SetWorldView() {	

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

	bb := time.Now()
	fmt.Println("Calculate view time: ", float64(bb.Nanosecond() - aa.Nanosecond()) / 1e9)	  
}	



