package model

import (
	"log"
	"app/pkg/store"
	"fmt"
	"github.com/gorilla/websocket"
)

const eventBufferSize = 256

// Client ...
type Client struct { 
	Conn    *websocket.Conn
	Account *store.Account

	Out chan<- ClientEvent // to gameloop
	In  chan interface{}   // from gameLoop
}

// NewClient ...  
func NewClient(conn *websocket.Conn, outChan chan<- ClientEvent, account *store.Account) *Client {
	c := &Client{
		Conn:    conn,
		Out:     outChan,
		In:      make(chan interface{}, eventBufferSize),
		Account: account,
	}
   
	fmt.Printf("model/client/NewClient New client: %v \n", c)

	// notify join	
	outChan <- ClientEvent{
		Sender: c,
		Join: &ClientJoinEvent{
			Ok: true,
		},
	}

	return c
}

// Close ...
func (c *Client) Close() {
	c.Conn.Close()

	// let game know they've bailed
	c.Out <- ClientEvent{
		Sender: c,
		Leave: &ClientLeaveEvent{
			Ok: true,
		},
	}
}

// HandleInbound runs in websocket handler's goroutine (per conn)
func (c *Client) HandleInbound() {
	for {
		var e ClientEvent
		err := c.Conn.ReadJSON(&e)
		
		if err != nil {
			log.Printf("model/client/HandleInbound error: %v", err)
			fmt.Println("model/client/HandleInbound closing client ", c) 
			c.Close()
			return
		}

		e.Sender = c // label sender
		c.Out <- e   // send event to game loop
	}
}

// HandleOutbound runs in its own goroutine too
func (c *Client) HandleOutbound() {
	// tell client they've joined successfully
	connectEvent := ServerEvent{
		Connect: &ServerConnectEvent{
			UUID: c.Account.UUID,
		},
	}
	err := c.Conn.WriteJSON(connectEvent)
	if err != nil {
		log.Printf("model/client/HandleOutbound error1: %v", err)
		c.Close()
		return
	}

	for e := range c.In {
		err := c.Conn.WriteJSON(e)
		if err != nil {
			log.Printf("model/client/HandleOutbound error2: %v", err)
			c.Close()
			return
		}
	}
}

// SendError ...
func (c *Client) SendError(err error) {
	c.In <- ServerEvent{
		Error: &ServerErrorEvent{
			Message: err.Error(),
		},
	}
}
