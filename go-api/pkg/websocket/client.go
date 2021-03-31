package websocket

import (
    "fmt"
    "log"   	
		"github.com/gorilla/websocket"
		"app/pkg/game" 
		"app/pkg/fov" 
		"time"
		"sync"
		"net/http"
		"net/url"
		"encoding/json"
		"strconv" 
		"github.com/google/uuid" 
		"errors"  
)

const ViewWidth int = 15
const ViewHeight int = 15

var halfViewWidth int = ViewWidth/2 // gives floor, which is what we want
var halfViewHeight int = ViewHeight/2

type Player struct {
	Avatar int8    `json:"avatar"`
	Name string `json:"name"`		
}

type Coord struct {
	X int
	Y int
}

type gridMap struct {
	grid [][]int8	
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	Position Coord
	TerrainView [][]int8
	AnimalView [][]int8	// view of animated layer
	mutex sync.Mutex
	Player Player
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`		
}

type Move struct {
	Type int    `json:"type"`
	ID string    `json:"id"`
	Direction int `json:"direction"`		
}

type WorldViewUpdate struct {
	Type int    `json:"type"`	
	TerrainView [][]int8 `json:"terrainView"`		
	AnimalView [][]int8 `json:"animalView"`		
}

type ClientMessage struct {
	Text string `json:"message"`
	Move int `json:"move"`	
}

func wrapMod (x, mod int) int{
	return (x%mod + mod)%mod;
}

func (g gridMap) InBounds(x, y int) bool {
	if (x < 0 || y < 0 || x > ViewWidth-1 || y > ViewHeight-1) {
		return false
	}
	return true;
}
func (g gridMap) IsOpaque(x, y int) bool {	
	if g.grid[x][y] == 8 { // high mountain
		return true
	}
	if g.grid[x][y] == 10 { // heavy forest
		return true
	}
	return false
}


func NewClient(c *websocket.Conn, p *Pool, playerCookie *http.Cookie) *Client {

	unescaped, _ := url.QueryUnescape(playerCookie.Value)
        
	var playerMiddle struct {
			Avatar string    `json:"avatar"`
			Name string `json:"name"`		
	}

	er := json.Unmarshal([]byte(unescaped), &playerMiddle)
	if er != nil { 
			fmt.Println("json Unmarshal error: ", er)
	}

	// fmt.Println("Unmarshalled Player Cookie",playerMiddle.Name, playerMiddle.Avatar)

	avatari, _ := strconv.Atoi(playerMiddle.Avatar) // read the avatar as a string from the JSON from the cookie, need it as int    
	avatari8 := int8(avatari)
	// fmt.Println(playerMiddle.Avatar,avatari, avatari8, playerCookie)

	// Initialize client TerrainView slices
	view := make([][]int8, ViewHeight)       // initialize a slice of viewHeight slices
	for i:=0; i < ViewHeight; i++ {
			view[i] = make([]int8, ViewWidth)  // initialize a slice of viewWidth in in each of viewHeight slices
	}

	// Initialize client animalView slices
	animal := make([][]int8, ViewHeight)       // initialize a slice of viewHeight slices
	for i:=0; i < ViewHeight; i++ {
			animal[i] = make([]int8, ViewWidth)  // initialize a slice of viewWidth in in each of viewHeight slices
	}

	client := &Client{ 
		Conn: c,
		Pool: p,
		ID: uuid.New().String(),
		Position: Coord{
				
			 // X: 240, // yew/lagoon
			 // Y: 120, // yew/lagoon
				X: 244, // britain
				Y: 148, // britain
			 // X: 10, // origin
			 // Y: 10, // origin
		},
		TerrainView: view, 
		AnimalView: animal,
		Player: Player{
				Avatar: avatari8,
				Name: playerMiddle.Name,
		},     
	}

	client.SetWorldView()

	return client;

}

func (c *Client) SetWorldView() {	

	aa := time.Now()
	
	xStart := wrapMod((c.Position.X - halfViewWidth),game.Width);
	yStart := wrapMod((c.Position.Y - halfViewHeight),game.Height);
	
	board := c.Pool.WorldMap.Grid
	 
	for x := 0; x < ViewWidth; x++ {	
		xMap := wrapMod((xStart + x), game.Width)
		for y := 0; y < ViewHeight; y++ {			
			yMap := wrapMod((yStart + y), game.Height)			
			c.TerrainView[x][y] = board[xMap][yMap]
			c.AnimalView[x][y] = 0 // while we're at it, clear the AnimalView
		}
	}	

	grid := make([][]int8, ViewHeight) 

	// Initialize Gridmap - this might be funky, flipping the axes?
	for x := 0; x < ViewWidth; x++ {		
		grid[x] = make([]int8, ViewHeight)  
		for y := 0; y < ViewHeight; y++ {			
			// fmt.Println("Initialize gridmap",x,y)
			grid[x][y]=c.TerrainView[x][y]
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
				c.TerrainView[x][y] = 0
			}			
		} 
	}
 
	// Add the other players onto the worldview:
	for player, _ := range c.Pool.Clients {
		if (player.ID == c.ID) { continue }
		px := player.Position.X
		py := player.Position.Y   		
		xrel := wrapMod(px - xStart, game.Width) // some head breaking math to figure out if other player is within the view
		yrel := wrapMod(py - yStart, game.Height)
		if ( (xrel < ViewWidth) && (yrel < ViewHeight) && c.TerrainView[xrel][yrel] != 0) {			
			// c.TerrainView[xrel][yrel] = -1	
			c.AnimalView[xrel][yrel] = player.Player.Avatar
			fmt.Println("Adding an animal" , player.Player.Avatar)
		}				 
	}	

	bb := time.Now()
	fmt.Println("Calculate view time: ", float64(bb.Nanosecond() - aa.Nanosecond()) / 1e9)
	  
}	

func (c *Client) readText(text string ) {
	if text != "" {
		message := Message{Type: 1, Body: string(text)}
		c.Pool.Broadcast <- message
		fmt.Printf("Got message: %#v\n",text)
	}
}

func (c *Client) getNewPosition(move int) (Coord, error) {
	newPosition := c.Position;
	switch move {
		case 38: // up  
			newPosition.Y = wrapMod((c.Position.Y - 1), game.Height); 		
		case 40: // down
			newPosition.Y = wrapMod((c.Position.Y + 1), game.Height); 		
		case 37: // left
			newPosition.X = wrapMod((c.Position.X - 1), game.Width);		
		case 39: // right
			newPosition.X = wrapMod((c.Position.X + 1), game.Width); 
		case 13: // enter/ping, no change				
			return newPosition, nil
		default: 		
			return Coord{X: -1, Y: -1}, errors.New("Requested a non-move")
	}
	return newPosition, nil;
}	

// Check if new position is valid
func (c *Client) isPositionValid(position Coord) bool {
	
	// Check if the tile is impassible
	newTerrain := c.Pool.WorldMap.Grid[position.X][position.Y]
	
	if (newTerrain) <= 3 { // water
		return false
	}
	if (newTerrain) == 8 { // high mountain
		return false
	}

	// Check if a player's in the way:
	for player, _ := range c.Pool.Clients { 
		if (player.ID == c.ID) { continue }
		if (player.Position.X == position.X) && (player.Position.Y == position.Y) {								
			return false
		}								
	}		

	return true;
}

// Updates World View of self, and all other clients that are now within World View
func (c *Client) updateAllClientsAfterMove() {

	// Always update this client first
	c.SetWorldView() // Now that client has moved, update their worldview

	c.mutex.Lock()
 	c.Conn.WriteJSON(WorldViewUpdate{Type: 4, TerrainView: c.TerrainView, AnimalView: c.AnimalView})
	c.mutex.Unlock()

	// Then update the other clients
	for player, _ := range c.Pool.Clients { // Trigger other players to update their worldview, since this player has updated
		// @todo: IMPORTANT FIX BEFORE THIS GROWS
		// should just update clients that are in view range to save work. 
		// This would also prevent other players being able to "spy" moves by watching socket traffic
		if player.ID != c.ID {
			go func(player *Client) {
				player.SetWorldView()			
				player.mutex.Lock()
				defer player.mutex.Unlock()
				player.Conn.WriteJSON(WorldViewUpdate{Type: 4, TerrainView: player.TerrainView, AnimalView: player.AnimalView })   
			}(player)			      
		}
	}	

}

func (c *Client) readMove( move int)  {

	fmt.Println("Client.readMove", move) 
	
	if move == 0 { return }
	
	newPosition, error := c.getNewPosition( move )	
	if error != nil {
		fmt.Println("No new position granted, bad move request ", c.ID)
		return 
	}
	
	if (!c.isPositionValid(newPosition)) { return }

	c.Position = newPosition;
	
	c.updateAllClientsAfterMove()
			
}


func (c *Client) Read() {	

	defer func() {
		fmt.Printf("unregistering") 
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {		
		var m = ClientMessage{}	
		err :=  c.Conn.ReadJSON(&m)	
		if err != nil {
			fmt.Printf("Read JSON Error")
			log.Println(err)
			return
		}	

		c.readText(m.Text)		
		c.readMove(m.Move)			
	}
}