package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	// _ "github.com/mattn/go-sqlite3" // justify it lol
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"github.com/joho/godotenv"
	"math/rand"
	// "app/pkg/game/model"
)

var db *sql.DB

// Account ...
type Account struct {
	UUID           uuid.UUID
	Username       string
	HashedPassword []byte
}

/* 
player (
  uuid TEXT PRIMARY KEY,
  screenname TEXT,
  zone TEXT,
  x smallint,
  y smallint,
  z smallint,
  avatar smallint,
  stats json,
  items json,
  state json
);
*/

type EntityStore struct {
	UUID uuid.UUID
	Name string
	Type string
	X int
	Y int
	Avatar int
	ZoneName string
}


// Init ...
func Init() {
	var err error
	godotenv.Load()
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
		os.Getenv("DB_HOST"), 
		5432, 
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_NAME"))	
	conn, err := sql.Open("postgres",dsn)
	if err != nil {
		log.Fatal(err)
	}	
	err = conn.Ping()
	if err != nil {
		fmt.Println(os.Getenv("DB_HOST"))
		fmt.Println("error:",os.Getenv("DB_HOST"), 
		5432, 
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_NAME"))	
		log.Fatal(err)
	}
	log.Println("Database connection established")
	db = conn

}
 
// LoginAccount ...
func LoginAccount(username, password string) (*Account, error) {
	account, err := GetAccountByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(account.HashedPassword, []byte(password)); err != nil {
		return nil, errors.New("Passwords do not match")
	}

	return account, nil
}

// RegisterAccount ...
func RegisterAccount(username, password string) (*Account, error) {
	account, _ := GetAccountByUsername(username)
	if account != nil {
		return nil, errors.New("account already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	accountUUID := uuid.New()

	

	statement, err := db.Prepare("INSERT INTO accounts (uuid, username, hashed_password) VALUES ($1, $2, $3)")
	if err != nil {
		return nil, err
	}
	_, err = statement.Exec(accountUUID.String(), username, string(hashedPassword))
	if err != nil {
		return nil, err
	}
	statement.Close()

	/* 
player (
  uuid TEXT PRIMARY KEY,
  screenname TEXT,
  zone TEXT,
  x smallint,
  y smallint,
  z smallint,
  avatar smallint,
  stats json,
  items json,
  state json
);
*/

	// generate random avatar tile
	avatar := rand.Intn(4) + 101 

	// initialize row in player. Holds player's data like location, avatar, etc.
	playerStatement, err := db.Prepare("INSERT INTO player (uuid, screenname, zone, x, y, z, avatar, stats, items, state) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")
	if err != nil {
		return nil, err
	}
	_, err = playerStatement.Exec(accountUUID.String(),"", "", 233, 154, 0, avatar, "{}", "{}", "{}")
	if err != nil {
		return nil, err
	}
	playerStatement.Close() 

	return GetAccount(accountUUID)
}

// GetAccount ...
func GetAccount(accountUUID uuid.UUID) (*Account, error) {
	var rawUUID, username, hashedPassword string

	statement, err := db.Prepare("SELECT uuid, username, hashed_password FROM accounts WHERE uuid = $1")
	if err != nil {
		return nil, err
	}
	
	err = statement.QueryRow(accountUUID.String()).Scan(&rawUUID, &username, &hashedPassword)
	if err != nil {
		return nil, err
	}
	statement.Close()

	parsedUUID, err := uuid.Parse(rawUUID)
	if err != nil {
		return nil, err
	}

	a := &Account{
		UUID:           parsedUUID,
		Username:       username,
		HashedPassword: []byte(hashedPassword),
	}

	return a, nil
}

// GetAccountByUsername ...
func GetAccountByUsername(username string) (*Account, error) {
	var rawUUID, resUsername, hashedPassword string

	statement, err := db.Prepare("SELECT uuid, username, hashed_password FROM accounts WHERE username = $1")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = statement.QueryRow(username).Scan(&rawUUID, &resUsername, &hashedPassword)
	if err != nil {
		return nil, err
	}
	statement.Close()

	parsedUUID, err := uuid.Parse(rawUUID)
	if err != nil {
		return nil, err
	}

	a := &Account{
		UUID:           parsedUUID,
		Username:       resUsername,
		HashedPassword: []byte(hashedPassword),
	}

	return a, nil
}


// @todo: make a new struct for entity store. The player.go will use this function, and the entityStore struct
// and then use boilerplate to put that in the entityStore struct. I guess.

// @todo: need to find the right place for this EntityStore model. Currently it creates an import cycle
func GetStoredEntityData(accountUUID uuid.UUID) (*EntityStore, error) {
	var rawUUID string
	var screenname, zone string
	var x, y, z, avatar int	

	statement, err := db.Prepare("SELECT uuid, screenname, zone, x, y, z, avatar FROM player WHERE uuid = $1")
	if err != nil {
		return nil, err
	}
	
	err = statement.QueryRow(accountUUID.String()).Scan(&rawUUID, &screenname, &zone, &x, &y, &z, &avatar)
	if err != nil {
		return nil, err
	}
	statement.Close()

	parsedUUID, err := uuid.Parse(rawUUID)
	if err != nil {
		return nil, err
	}

	fmt.Println("Getting stored entity data. zone is " + zone)

	e := &EntityStore{
		UUID: parsedUUID,
		Name: screenname, 
		Avatar: avatar,
		X: x,
		Y: y,
		ZoneName: zone,
	}

	return e, nil
}	

/*
player (
  uuid TEXT PRIMARY KEY,
  screenname TEXT,
  zone TEXT,
  x smallint,
  y smallint,
  z smallint,
  avatar smallint,
  stats json,
  items json,
  state json
);
*/

/*
	UUID uuid.UUID
	Name string
	Type string
	X int
	Y int
	Avatar int8	
	*/


func SetStoredEntityData(es EntityStore) (error) {	
	statement, err := db.Prepare("UPDATE player SET screenname = $1, zone = $2, x = $3, y = $4, z = $5, avatar = $6, stats = $7, items = $8, state = $9 WHERE uuid = $10")
	if err != nil {
		return err
	}
	_, err = statement.Exec(es.Name, "", es.X, es.Y, 0, es.Avatar, "{}", "{}", "{}", es.UUID.String())
	if err != nil {
		return err
	}
	statement.Close()	
	return nil
}	

func SetStoredAvatar(avatar int, uuid uuid.UUID) (error) {	
	statement, err := db.Prepare("UPDATE player SET avatar = $2 WHERE uuid = $1")
	if err != nil {
		return err
	}
	_, err = statement.Exec(uuid, avatar)
	if err != nil {
		return err
	}
	statement.Close()	
	return nil
}	

func SetStoredName(name string, uuid uuid.UUID) (error) {	
	statement, err := db.Prepare("UPDATE player SET screenname = $2 WHERE uuid = $1")
	if err != nil {
		return err
	}
	_, err = statement.Exec(uuid, name)
	if err != nil {
		return err
	}
	statement.Close()	
	return nil
}	

func SetStoredLocation(x int, y int, uuid uuid.UUID) (error) {	
	statement, err := db.Prepare("UPDATE player SET x = $2, y = $3 WHERE uuid = $1")
	if err != nil {
		return err
	}
	_, err = statement.Exec(uuid, x, y)
	if err != nil {
		return err
	}
	statement.Close()	
	return nil
}	

func SetStoredZone(zone string, uuid uuid.UUID) (error) {	
	fmt.Println("Storing player in zone: " + zone)
	statement, err := db.Prepare("UPDATE player SET zone = $2 WHERE uuid = $1")
	if err != nil {
		return err
	}
	_, err = statement.Exec(uuid, zone)
	if err != nil {
		return err
	}
	statement.Close()	
	return nil
}	
