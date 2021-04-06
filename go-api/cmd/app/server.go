package main

import (
    "fmt"
    "net/http"    
	"app/pkg/websocket"          
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {

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

    client := websocket.NewClient(conn, pool, playerCookie)

    pool.Register <- client    
    client.Read()
}

func setupRoutes() {

	pool := websocket.NewPool()

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