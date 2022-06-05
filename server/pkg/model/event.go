package model

import (
	"github.com/google/uuid"	
)

// ClientEvent ...
type ClientEvent struct {
	Join   *ClientJoinEvent   `json:"join,omitempty"`
	Leave  *ClientLeaveEvent  `json:"leave,omitempty"`
	Chat   *ClientChatEvent   `json:"chat,omitempty"`
	Move   *ClientMoveEvent   `json:"move,omitempty"`
	Attack *ClientAttackEvent `json:"attack,omitempty"`
	UpdateAvatar *ClientUpdateAvatarEvent `json:"avatar,omitempty"`
	UpdateName *ClientUpdateNameEvent `json:"updateName,omitempty"`
	PeerGem *ClientPeerGemEvent `json:"peerGem,omitempty"`
	CastSpell *ClientCastSpellEvent `json:"castSpell,omitempty"`
	Look *ClientLookEvent `json:"look,omitempty"`
	Talk *ClientTalkEvent `json:"talk,omitempty"`
	SimpleAction *ClientSimpleActionEvent `json:"simpleAction,omitempty"`
	Sender *Client `json:"-"`
}


// ClientJoinEvent ...
type ClientJoinEvent struct {
	Ok bool
}

// ClientLeaveEvent ...
type ClientLeaveEvent struct {
	Ok bool
}

// ClientChatEvent ...
type ClientChatEvent struct {
	Message string `json:"message"`
}

type ClientMoveEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ClientUpdateAvatarEvent struct {
	ID int `json:"id"`
}

type ClientUpdateNameEvent struct {
	Name string `json:"name"`
}

type ClientViewUpdateEvent struct {
	grid [][]int8
}

// ClientAttackEvent ...
type ClientAttackEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ServerEvent ...
type ServerEvent struct {
	Connect *ServerConnectEvent `json:"connect,omitempty"`
	Error   *ServerErrorEvent   `json:"error,omitempty"`
}

// ServerConnectEvent
type ServerConnectEvent struct {
	UUID uuid.UUID `json:"uuid"`
}

// ServerErrorEvent ...
type ServerErrorEvent struct {
	Message string `json:"message"`
}

type ServerMessageEvent struct {
	Message string `json:"message"`
}

type ClientPeerGemEvent struct {}

type ClientCastSpellEvent struct {
	Spell string `json:"spell"`
}

type ClientLookEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ClientTalkEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
	Message string `json:"message"`
}

type ClientSimpleActionEvent struct {
	Action string `json:"action"`
}