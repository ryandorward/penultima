package model

import "github.com/google/uuid"

type WorldObject struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
	// Tile int       `json:"tile"` // representing tile
	Tile string       `json:"tile"` // representing tile
	X    int       `json:"x"`
	Y    int       `json:"y"`
	Quantity int `json:"quantity"`

	// special features
	LightRadius int `json:"lightRadius"`
	Type       WorldObjectType `json:"type"`
	WarpTarget *WarpTarget     `json:"warp_target,omitemtpy"`
	HealZone   *HealZone       `json:"heal_zone,omitempty"`
}

type WorldObjectType string

const (
	WorldObjectTypePlayerSpawn = "playerSpawn"
	WorldObjectTypePortal      = "portal"
	WorldObjectTypeHealing     = "healing"
)

type WarpTarget struct {
	ZoneUUID uuid.UUID `json:"zone_uuid"`
	ZoneName string `json:"zone_name"`
	Zone     Zone      `json:"-"`
	X        int       `json:"x"`
	Y        int       `json:"y"`
}

type HealZone struct {
	Full bool
}

func (w *WorldObject) SetPosition(x, y int) {
	w.X = x
	w.Y = y
}

func (w *WorldObject) SetQuantity(q int) {
	w.Quantity = q
}

func (w *WorldObject) GetType() string {
	return string(w.Type)
}