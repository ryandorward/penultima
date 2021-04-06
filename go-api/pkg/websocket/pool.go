package websocket

import (
	"fmt"
	"strconv"
	"app/pkg/game"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
	WorldMap *game.WorldMap
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),		
		WorldMap: game.NewWorldMap(),
	}
}

func (pool *Pool) Start() {
	for {
		select {
			case client := <-pool.Register:
				pool.Clients[client] = true			
				fmt.Printf("Client registered: %#v\n", client)
				fmt.Println("Size of Connection Pool: ", len(pool.Clients))
				for client := range pool.Clients {
					fmt.Println("-Client:", &client)
					client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined."})
					client.Conn.WriteJSON(Message{Type: 2, Body: strconv.Itoa(len(pool.Clients))})
				}
			

			case client := <-pool.Unregister:
				delete(pool.Clients, client)
				fmt.Println("Size of Connection Pool: ", len(pool.Clients))
				for client := range pool.Clients {
					client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected."})
					client.Conn.WriteJSON(Message{Type: 2, Body: strconv.Itoa(len(pool.Clients))})
				}
			
			case message := <-pool.Broadcast:
				fmt.Println("Sending message to all clients in Pool")
				for client := range pool.Clients {
					if err := client.Conn.WriteJSON(message); err != nil {
						fmt.Println(err)
						return
					}
				}
								
		}
	}
}