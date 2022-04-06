package event

import "app/pkg/game/model"

type Event interface {
}

type JoinEvent struct {
	Entity model.Entity
}

type LeaveEvent struct {
	Entity model.Entity
}

type SpawnEvent struct {
	Entity model.Entity
}

type DespawnEvent struct {
	Entity model.Entity
}

type DieEvent struct {
	Entity model.Entity
}

type MoveEvent struct {
	Entity model.Entity
	X, Y   int
}

type UpdateOwnViewEvent struct {
	Entity model.Entity
}

type UpdateAvatarEvent struct { 
	Entity model.Entity
	id		 int8
}

type ServerMessageEvent struct {
	Message string 
}

type PeerGemEvent struct {
	Entity model.Entity 
}

/*
type PeerGemEvent struct {
	Entity model.Entity	
}
*/

/*
type ChatEvent struct {
	Entity  model.Entity
	Message string
}

type AttackEvent struct {
	Attacker, Target model.Entity
	Hit              bool
	Damage           int
	TargetHP         int
}

type HealEvent struct {
	Entity model.Entity
	Amount int
	Full   bool
}

type GainXPEvent struct {
	Entity    model.Entity
	LeveledUp bool
}
*/