package game

import (	
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"	
) 
 
type Client struct {	
	Conn   *websocket.Conn 
	Pool * Pool
	mutex sync.Mutex 
	player *Player
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`		
}

type Move struct {
	Type int `json:"type"`
	ID string    `json:"id"`
	Direction int `json:"direction"`		
}

type WorldViewUpdate struct {
	Type int `json:"type"`	
	TerrainView [][]int8 `json:"terrainView"`
	AnimalView [][]int8 `json:"animalView"`		
}

type ClientMessage struct {
	Text string `json:"mesage"`
	Move int `json:"move"`	
}

func NewClient(c *websocket.Conn, p *Pool, playerCookie *http.Cookie) *Client {

	unescaped, _ := url.QueryUnescape(playerCookie.Value)
        
	var playerMiddle struct {
		Avatar string `json:"avatar"`
		Name string `json:"name"`		
	}

	er := json.Unmarshal([]byte(unescaped), &playerMiddle)
	if er != nil { 
		fmt.Println("json Unmarshal error: ", er)
	}

	avatari, _ := strconv.Atoi(playerMiddle.Avatar) // read the avatar as a string from the JSON from the cookie, need it as int    
	avatari8 := int8(avatari)	

	player := NewPlayer(
		PlayerCookie{
			Avatar  : avatari8,
			Name: playerMiddle.Name, 
		}, p.WorldMap) 

	client := &Client{ 
		Conn: c,
		Pool: p,		
		player: player,
	}

	client.player.SetWorldView() 
	return client;
}
 
func (c *Client) GetLocation() (Coord){	
	return c.player.Location
}	 

func (c *Client) GetID() (string){	
	return c.player.ID
}	 

func (c *Client) GetAvatar() (int8) {		
	return c.player.PlayerCook.Avatar
}	

// Retrieve slice of all other Clients as LocatableEntities
func (c *Client) GetOthers() (map[LocatableEntity]bool) { 
	others := c.getOtherClients()
	// Create slice of LocatableEntity from others
	les := make(map[LocatableEntity]bool)
	for other := range others {
		les[other] = true
	}	
	return les
}

func (c *Client) readText(text string ) {
	if text != "" {
		message := Message{Type: 1, Body: string(text)}
		c.Pool.Broadcast <- message 
		fmt.Printf("Got message: %#v\n",text)   
	}
}

// Flag updateView switches whether or not to update clients view with the visible others
func (c *Client) getVisibleClients(updateView bool)(map[*Client]bool) {
	others := c.getOtherClients()
	visibleClients := make(map[*Client]bool)
	for other := range others {			
		if (other.GetID() == c.GetID()) { continue }
		canSee, xrel, yrel := c.player.CanSee(other.player)
		if canSee {			
			visibleClients[other] = true
			if updateView {
				c.player.AnimateView[xrel][yrel] = other.GetAvatar()
			}
		}
	}
	return visibleClients
}	

// Retrieve slice of all other Client
func (c *Client) getOtherClients() (map[*Client]bool) {
	others := make(map[*Client]bool)
	for other := range c.Pool.Clients { 		
		if (other.GetID() != c.GetID()) {
			others[other] = true
		}
	}
	return others
}

// Send the view JSON to client to update their view of the world
func (c *Client) updateClientViews() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Conn.WriteJSON(WorldViewUpdate{Type: 4, TerrainView: c.player.TerrainView, AnimalView: c.player.AnimateView})
}

func (c *Client) makeMove(newLocation Coord) {
	// get list of clients with players in the view before the move 
	updatableClientsBeforeMove := c.getVisibleClients(false)	
	c.player.Location = newLocation // make move	
	c.player.SetWorldView() // Update player's view

	// Get new list of visible clients, and update them to this clients view, 
	// then merge with those visible before
	updatableClients := c.getVisibleClients(true)		
	for client, val := range updatableClientsBeforeMove { 
		updatableClients[client] = val
	}

	// publish to the client
	c.updateClientViews()  

	// Update the other clients - only updating others that are within current view
	// Trigger other players to update their own worldview, since this player has updated		
	for other := range updatableClients { 	
			go func(other *Client) {				
				other.player.SetWorldView()					
				other.getVisibleClients(true)							
				other.updateClientViews() 
			}(other)				      		  
	}	

}

func (c *Client)locationCheck(location Coord) bool {		
	les := c.GetOthers()
	worldMap := c.player.WorldMap
	return IsLocationValid(location, *worldMap, les )	
}

func (c *Client) readMove( move int)  {

	if move==0 { 
		return 
	}
	
	newLocation, error := GetNewLocation( move, c.player.Location )	
	if error != nil {
		fmt.Println("No new position granted, bad move request ", c.player.ID)
		return 
	}
	
	if c.locationCheck(newLocation) {
		c.makeMove(newLocation)
	}
 
}

func (c *Client) Read() {	 

	defer func() {
		fmt.Printf("unregistering") 
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {		
		var m =ClientMessage{}	
		err := c.Conn.ReadJSON(&m)	
		if err != nil {
			fmt.Printf("Read JSON Error")
			log.Println(err)
			return
		}	

		c.readText(m.Text)
		c.readMove(m.Move)			
	}

}