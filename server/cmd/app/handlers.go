package main

import (
	"net/http"
	"fmt"
	"os"
)

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://" + os.Getenv("CLIENT_HOST") + ":3002") 
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../public/index.html")
} 

func handleGame(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	_, err := authenticated(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.ServeFile(w, r, "../public/game.html")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {	
	setHeaders(w)
	fmt.Println("Login Endpoint Hit")
	err := authenticate(w, r, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Println("Authenticated?")
	_, err2 := authenticated(w, r)
	if err2 != nil {
		fmt.Println("Auth Error in Login: ")
		fmt.Println(err2)	
		return
	}
	fmt.Println("Authenticated.")
	w.WriteHeader(200)
	 
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	err := authenticate(w, r, true)
	if err != nil { 
		fmt.Println("Registration error:")	
		http.Error(w, err.Error(), http.StatusBadRequest)
	}	
	w.WriteHeader(200)
} 

func handleLogout(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	fmt.Println("Logout Endpoint Hit")
	session, _ := sessionStore.Get(r, "penultima")
	session.Values["UUID"] = ""
	session.Save(r, w)
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}
