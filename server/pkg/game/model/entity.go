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

	GetType() string // EntityType

	GetZone() Zone
	GetZoneName() string
	SetZone(Zone)
	SetZoneWithTarget(Zone, int, int)

	GetPosition() (int, int)
	SetPosition(int, int)	

	SetQuantity(int)	

	GetStats() Stats

	Tick() bool // increments energy, returns if can act
	Act() Action

	RollToHit(int) bool
	RollDamage() int

	TakeDamage(float64) bool
	Die()
	GainExp(int)
	Heal(float64)

	GetClient() *model.Client

	IsInView(Entity) bool
	 
	SetTile(int)
	GetTile() int

	SetSlowThresh(float64)
	GetSlowThresh()float64

	GetLastMoveTry() (int, int)
	SetLastMoveTry(int, int)

	UpdateOwnView() 
	UpdateOwnStats()
	UpdateClientStat(string, int)

	SetEntityData(*store.EntityStore)

	ReceiveMessage(string) string

	ReceiveResult(string, string)

	AddFood(int) int
	AddGems(int) int
	AddSilver(int) int
	GetGemCount() int
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
	HP            float64 `json:"hp"`
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
	Type EntityType `json:"type"`
	Tile string `json:"tileName"`	
	Health int `json:"health"`
	IsMortal bool `json:"isMortal"`
	Movement struct {
		SpeedMod int `json:"speedMod"` // speed = 1/SpeedMod
		Jitter int `json:"jitter"` // better be less than SpeedMod or weirdness will ensue
		Algorithm string `json:"algorithm"`
		DirectionChangeProbability int `json:"directionChangeProbability"` // 0-100	
	} `json:"movement"`
	Ordinality int `json:"ordinality"`
}  