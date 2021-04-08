package main

import (
    "fmt"
    "net/http"    
	"app/pkg/websocket"          
    "app/pkg/game" 
)

func serveWs(pool *game.Pool, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
    fmt.Println("WebSocket Endpoint Hit")

    playerCookie, err := r.Cookie("player")
    if err != nil {
        fmt.Println("Cookie error:", err.Error())
    } else {
        fmt.Println("Player Cookie",playerCookie.Value)
    }
        
    conn, err := websocket.Upgrade(w, r)
    if err != nil {
        fmt.Fprintf(w, "%+v\n", err)
    }

    client := game.NewClient(conn, pool, playerCookie)

    /* New plan: ? 

    make Client in it's own package
    make new player here 
    player is dependent on client package
    player has a field called client that uses the Client struct 
    
    make new Client here in this serveWS func
    make client the player's client, as in player.client := client
    pool will no longer register the clients, just the players

    look closely at Client.Read() method. Will need to be able to do that even when client is an 
    independent property of player. 

    Client.Read will need to maybe receive a pointer to the player, and satisfy the clients interface, doing:
    - unregister <- will trigger unregister of player from pool
    - readText(string text)
    - readMove(move int)

    Client will expect a connection to the player via a playerInterface, which has those above methods
    so pseudo code will look like

    player := NewPlayer()
    client := NewClient()
    client.player = player // <- player satisfies clients playerInterce, which expects the above 3 methods
    player.client = client // <- player is dependent on client's package so this should be fine
    pool.Register <- player
    player.client.Read();

    Now how does pool fit in in the package heirarchy?
    same package as client

    */

    pool.Register <- client    
    client.Read()
}

func setupRoutes() {

	pool := game.NewPool()

	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
    fmt.Println("Ultima go-api v0.0.1")
    setupRoutes()
    http.ListenAndServe(":8080", nil)
}