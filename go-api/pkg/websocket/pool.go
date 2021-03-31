package websocket

import (
	"fmt"
	"strconv"
	// "github.com/google/uuid"
	"app/pkg/game"
	// "app/pkg/fov" 
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
	MovesChan chan Move
	// FOVCalc *fov.View
	WorldMap *game.WorldMap
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		MovesChan:  make(chan Move),		
		WorldMap: game.NewWorldMap(),
		// FOVCalc: fov.New(),
	}
}

func (pool *Pool) Start() {
	for {
		select {
			case client := <-pool.Register:
				pool.Clients[client] = true
			
				fmt.Printf("Client registered: %#v\n", client)
				fmt.Println("Size of Connection Pool: ", len(pool.Clients))
				for client, _ := range pool.Clients {
					fmt.Println("-Client:", &client)
					client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined."})
					client.Conn.WriteJSON(Message{Type: 2, Body: strconv.Itoa(len(pool.Clients))})
				}
				break

			case client := <-pool.Unregister:
				delete(pool.Clients, client)
				fmt.Println("Size of Connection Pool: ", len(pool.Clients))
				for client, _ := range pool.Clients {
					client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected."})
					client.Conn.WriteJSON(Message{Type: 2, Body: strconv.Itoa(len(pool.Clients))})
				}
				break

			case message := <-pool.Broadcast:
				fmt.Println("Sending message to all clients in Pool")
				for client, _ := range pool.Clients {
					if err := client.Conn.WriteJSON(message); err != nil {
						fmt.Println(err)
						return
					}
				}
				break
			
			case move := <-pool.MovesChan:
				fmt.Println("Sending moves to all clients in Pool")
				for client, _ := range pool.Clients {
					if err := client.Conn.WriteJSON(move); err != nil {
						fmt.Println(err)
						return
					}
				}
		}
	}
}