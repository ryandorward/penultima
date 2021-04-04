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
	//	"errors"  
)

const ViewWidth int = 15
const ViewHeight int = 15

var halfViewWidth int = ViewWidth/2 // gives floor, which is what we want
var halfViewHeight int = ViewHeight/2


type Client struct {
	ID   string
	Conn *websocket.Conn 
	Pool *Pool
	Position game.Coord
	TerrainView [][]int8
	AnimalView [][]int8	// view of animated layer
	mutex sync.Mutex
	PlayerCook player
	Player game.Player
}

type player struct {
	Avatar int8    `json:"avatar"`
	Name string `json:"name"`		
}

type gridMap struct {
	grid [][]int8	
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
		Position: game.Coord{				
			 // X: 240, // yew/lagoon
			 // Y: 120, // yew/lagoon
				X: 244, // britain
				Y: 148, // britain
			 // X: 10, // origin
			 // Y: 10, // origin
		},
		TerrainView: view, 
		AnimalView: animal,
		PlayerCook: player{
				Avatar: avatari8,
				Name: playerMiddle.Name,
		},    
		Player: game.Player{},
	}

	client.SetWorldView()

	return client;

}

func (c *Client) GetLocation() (game.Coord){	
	return c.Position;
}	
func (c *Client) GetID() (string){	
	return c.ID;
}	

func (c *Client) SetWorldView() ([]*Client){	

	aa := time.Now()
	 
	xStart := game.WrapMod((c.Position.X - halfViewWidth),game.WorldWidth);
	yStart := game.WrapMod((c.Position.Y - halfViewHeight),game.WorldHeight);
	
	board := c.Pool.WorldMap.Grid
	 
	for x := 0; x < ViewWidth; x++ {	
		xMap := game.WrapMod((xStart + x), game.WorldWidth)
		for y := 0; y < ViewHeight; y++ {			
			yMap := game.WrapMod((yStart + y), game.WorldHeight)			
			c.TerrainView[x][y] = board[xMap][yMap]
			c.AnimalView[x][y] = 0 // while we're at it, clear the AnimalView
		}
	}	

	grid := make([][]int8, ViewHeight) 

	// Initialize Gridmap - this might be funky, flipping the axes?
	for x := 0; x < ViewWidth; x++ {		
		grid[x] = make([]int8, ViewHeight)  
		for y := 0; y < ViewHeight; y++ {						
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

	updatableOthers := []*Client{}
 
	// Add the other players onto the worldview, keep track of which ones are visible, they will need to be updated too
	for other, _ := range c.Pool.Clients {
		if (other.ID == c.ID) { continue }
		px := other.Position.X
		py := other.Position.Y   		
		xrel := game.WrapMod(px - xStart, game.WorldWidth) // some head breaking math to figure out if other player is within the view
		yrel := game.WrapMod(py - yStart, game.WorldHeight)
		
		// if other is within view and not on a shadowed tile
		if ( (xrel < ViewWidth) && (yrel < ViewHeight) && c.TerrainView[xrel][yrel] != 0) {	
			// c.TerrainView[xrel][yrel] = -1	
			c.AnimalView[xrel][yrel] = other.PlayerCook.Avatar
			fmt.Println("Adding an animal" , other.PlayerCook.Avatar)
			
			updatableOthers = append(updatableOthers,other)
			
		}				 
	}	

	bb := time.Now()
	fmt.Println("Calculate view time: ", float64(bb.Nanosecond() - aa.Nanosecond()) / 1e9)

	return updatableOthers;
	  
}	

func (c *Client) readText(text string ) {
	if text != "" {
		message := Message{Type: 1, Body: string(text)}
		c.Pool.Broadcast <- message
		fmt.Printf("Got message: %#v\n",text)
	}
}

// Send the view JSON to client to update their view of the world
func (c *Client) updateClientViews() {
	c.Conn.WriteJSON(WorldViewUpdate{Type: 4, TerrainView: c.TerrainView, AnimalView: c.AnimalView})
}

// Updates World View of self, and all other clients that are now within World View
func (c *Client) updateAllClientsAfterMove() {

	// Always update this client first
	others := c.SetWorldView() // Now that client has moved, update their worldview

	c.mutex.Lock()
 	c.updateClientViews() 
	c.mutex.Unlock()

	// Update the other clients - only updating others that are within current view
	// Trigger other players to update their worldview, since this player has updated		
	for _, other := range others { 
		if other.ID != c.ID {
			go func(other *Client) {
				other.SetWorldView()			
				other.mutex.Lock()
				defer other.mutex.Unlock()				
				other.updateClientViews() 
			}(other)			      
		} 
	}	
}

func (c *Client) CanSee(other *Client) bool {
	return true
}

// Retrieve a slice of all other Clients
func (c *Client) GetOthers() []*Client{
	others := []*Client{}
	for other, _ := range c.Pool.Clients { 		
		if (other.ID != c.ID) {
			others = append(others,other)
		}
	}
	return others
}

func (c *Client) readMove( move int)  {

	fmt.Println("Client.readMove", move) 
	
	if move == 0 { 
		return 
	}
	
	newPosition, error := game.GetNewPosition( move, c.Position )	
	if error != nil {
		fmt.Println("No new position granted, bad move request ", c.ID)
		return 
	}
	
	// get "others": the slice of LocatableEntity (which is satisfied by Client) that are not self		
	others := c.GetOthers()
	// Create slice of LocatableEntity from others
	le := make([]game.LocatableEntity, len(others))
	for i, other := range others {
		le[i] = other
	}
	
	worldMap := c.Pool.WorldMap
	if (! game.IsLocationValid(newPosition, *worldMap, le ) ){ 
		return
	}

	// get list of clients in the view:

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