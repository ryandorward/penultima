package main

import (
	"fmt"
	"log"
	"net/http"
	"app/pkg/game"
	"app/pkg/model"
  "app/pkg/store"
	//"github.com/floralbit/dungeon/store"
	"github.com/gorilla/websocket"
	// "os"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}, 
}

func setResponseHeaders (w http.ResponseWriter) {


}
 
func main() { 	
	// pages
	// http.HandleFunc("/", handleIndex)
	// http.HandleFunc("/game", handleGame)
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../public/static"))))

	// endpoints
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/logout", handleLogout)

	// websocket
	http.HandleFunc("/ws", handleWs)  

	store.Init() // connect to db

	go game.Run() // kick off gameloop

	log.Println("http server started on :8085")
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// spawned in a goroutine by http
func handleWs(w http.ResponseWriter, r *http.Request) {

	setHeaders(w)
	fmt.Println("handling WS")
 
	// upgrade to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// auth
	account, err := authenticated(w, r)
	if err != nil {
		fmt.Println("Auth Error: ")
		fmt.Println(err)
		ws.Close()
		return
	}

	fmt.Printf("%s connected and authenticated\n", account.UUID.String())

	c := model.NewClient(ws, game.In, account)
	go c.HandleOutbound()
	c.HandleInbound()
}