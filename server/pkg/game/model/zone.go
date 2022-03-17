package model

import "github.com/google/uuid"

type Tile struct {
	ID    int    `json:"id"`
	Solid bool   `json:"solid"`
	Opaque bool  `json:"opaque"`
	Name  string `json:"name"`
	Speed float64 `json:"speed"`
}

type Zone interface {
	GetUUID() uuid.UUID
	GetDimensions() (int, int)

	GetTile(x, y int) *Tile
	GetEntities() []Entity
	AddEntity(Entity)
	RemoveEntity(Entity)

	GetAllWorldObjects() []*WorldObject
	GetWorldObjects(x, y int) []*WorldObject

	GetNewLocation(x,y,dx,dy int) (int, int)

	Update(float64)
	SlowUpdate()

	GetSunlight() int

	GetTrammel() int
	GetFelucca() int
	GetWind() (int,int)

}
