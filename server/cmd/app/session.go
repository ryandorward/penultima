package main

import (
	"errors"
	"net/http"

	"app/pkg/store"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"os"
	"fmt"
)

var sessionStore = sessions.NewCookieStore([]byte(os.Getenv("API_SECRET")))

func authenticate(w http.ResponseWriter, r *http.Request, newAccount bool) error {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		return errors.New("username or password not set")
	}

	var account *store.Account
	var err error
	if newAccount {
		account, err = store.RegisterAccount(username, password)
		if err != nil {
			return err
		}
	} else {
		account, err = store.LoginAccount(username, password)
		if err != nil {
			return err
		} else {
			fmt.Println("authenticated session value set")
		}
	}

	fmt.Println(r)

	session, _ := sessionStore.Get(r, "penultima")
	session.Values["authenticated"] = true
	session.Values["UUID"] = account.UUID.String()
	fmt.Println(account.UUID.String())
	session.Save(r, w)

	return nil
}

// returns account UUID if authenticated, errors if not
func authenticated(w http.ResponseWriter, r *http.Request) (*store.Account, error) {
	session, _ := sessionStore.Get(r, "penultima")

	fmt.Println(r)

	rawUUID, ok := session.Values["UUID"]
	if !ok || rawUUID == "" {
		return nil, errors.New("not authenticated")
	}

	accountUUID, err := uuid.Parse(rawUUID.(string))
	if err != nil {
		return nil, err
	}

	account, err := store.GetAccount(accountUUID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}

	return account, nil
}
