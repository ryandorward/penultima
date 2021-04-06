package websocket

import (
	"app/pkg/game"
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
	player *game.Player
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

	fmt.Println("Avatar: ", avatari8)

	player := game.NewPlayer(
		game.PlayerCookie{
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
 
func (c *Client) GetLocation() (game.Coord){	
	return c.player.Location
}	 
func (c *Client) GetID() (string){	
	return c.player.ID
}	 
func (c *Client) GetAvatar() (int8) {	
	fmt.Println("Avatar is", c.player.PlayerCook.Avatar)
	return c.player.PlayerCook.Avatar
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
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Conn.WriteJSON(WorldViewUpdate{Type: 4, TerrainView: c.player.TerrainView, AnimalView: c.player.AnimateView})
}

// combine with showOtherClientsPlayers function
func (c *Client) getVisibleClients()(map[*Client]bool) {
	visibleClients := make(map[*Client]bool)
	for other := range c.Pool.Clients {			
		if (other.GetID() == c.GetID()) { continue }
		canSee, _ , _ := c.player.CanSee(other.player)
		if canSee {			
			visibleClients[other] = true
		}
	}
	return visibleClients
}	

// @todo refactor this func 
func (c *Client) showOtherClientsPlayers()(map[*Client]bool) {
	// Add the other players onto the worldview, keep track of which ones are visible, 
// they will need to be updated too
	
	updatableClients := make(map[*Client]bool)

	for other := range c.Pool.Clients {			
		if (other.GetID() == c.GetID()) { continue }
		canSee, xrel, yrel := c.player.CanSee(other.player)
		if canSee {
			fmt.Println("Can see")
			c.player.AnimateView[xrel][yrel] = other.GetAvatar()						
			updatableClients[other] = true
		}		 
	}

	return updatableClients
}

// Retrieve a slice of all oter Client
func (c *Client) GetOthers() ([]*Client) {
	others := []*Client{}
	for other := range c.Pool.Clients { 		
		if (other.player.ID != c.player.ID) {
			others = append(others,other)
		}
	}
	return others
}

func (c *Client) makeMove(newLocation game.Coord) {
	// get list of clients with players in the view before the move 
	updatableClientsBeforeMove := c.getVisibleClients()

	// make move
	c.player.Location = newLocation;

	// Update player's view
	c.player.SetWorldView() 

	/* 
	
	@todo: this whole section is a good use case for channels.
	
	each player listens to a channel
	
	each player broadcasts to that chanel if they move
	
	if a move is picked up on the channel, update own view including new positions of other players in view

	*/

	// Get new list of visible clients, and update them to this clients view. @todo refactor 
	updatableClients := c.showOtherClientsPlayers()	

	// publish to the client
	c.updateClientViews()  

	// now merge updatableClients with those from before move:
	for client, val := range updatableClientsBeforeMove {
		updatableClients[client] = val
	}

	// Update the other clients - only updating others that are within current view
	// Trigger other players to update theirown worldview, since this player has updated		
	for other := range updatableClients { 
		if other.GetID() != c.GetID() {
			go func(other *Client) {
				fmt.Println("Getting other player to update its view")		
				other.player.SetWorldView()		
				other.showOtherClientsPlayers()								
				other.updateClientViews() 
			}(other)	
		}		      		  
	}	
}

func (c *Client) readMove( move int)  {

	if move== 0 { 
		return 
}
	
	newLocation, error := game.GetNewPosition( move, c.player.Location )	
	if error != nil {
		fmt.Println("No new position granted, bad move request ", c.player.ID)
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
	if (! game.IsLocationValid(newLocation, *worldMap, le ) ){ 
		return
	}

c.makeMove(newLocation)
 
}

func (c *Client) Read() {	 

	defer func() {
		fmt.Printf("unregistering") 
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {		
		var m =ClientMessage{}	
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