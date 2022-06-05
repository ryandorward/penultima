package model

import (
	"app/pkg/model"
	"github.com/google/uuid"
	"app/pkg/store"
)

type Entity interface {
	GetUUID() uuid.UUID
	GetName() string
	SetName(string)

	GetType() EntityType

	GetZone() Zone
	GetZoneName() string
	SetZone(Zone)
	SetZoneWithTarget(Zone, int, int)

	GetPosition() (int, int)
	SetPosition(int, int)

	GetStats() Stats

	Tick() bool // increments energy, returns if can act
	Act() Action

	RollToHit(int) bool
	RollDamage() int

	TakeDamage(int) bool
	Die()
	GainExp(int)
	Heal(int)

	GetClient() *model.Client

	IsInView(Entity) bool
	 
	SetTile(int)
	GetTile() int

	SetSlowThresh(float64)
	GetSlowThresh()float64

	GetLastMoveTry() (int, int)
	SetLastMoveTry(int, int)

	UpdateOwnView(c *model.Client)
	UpdateOwnStats()
	UpdateClientStat(string, int)

	SetEntityData(*store.EntityStore)

	ReceiveMessage(string) string

	ReceiveResult(string, string)

	AddFood(float64)
	AddGems(int)
	GetGems() int

}

type EntityType string

const (
	EntityTypePlayer  = "player"
	EntityTypeMonster = "monster"
	EntityTypeNPC = "npc"
)

type Stats struct {
	Level         int `json:"level"`
	MaxHP         int `json:"max_hp"`
	HP            int `json:"hp"`
	XP            int `json:"xp"`
	AC            int `json:"ac"`
	XPToNextLevel int `json:"xp_to_next_level"`

	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

type NPCProperties struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Type string `json:"type"`
	Tile string `json:"tileName"`	
	Movement struct {
		SpeedMod int `json:"speedMod"` // speed = 1/SpeedMod
		Jitter int `json:"jitter"` // better be less than SpeedMod or weirdness will ensue
		Algorithm string `json:"algorithm"`
		DirectionChangeProbability int `json:"directionChangeProbability"` // 0-100	
	} `json:"movement"`
}  