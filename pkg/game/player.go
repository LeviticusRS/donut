package game

import (
	"github.com/sprinkle-it/donut/pkg/client"
)

type PlayerFactory func(*client.Client, uint16) *Player

type PlayerConfig struct {
}

func (cfg PlayerConfig) Build(client *client.Client, id uint16) *Player {
	return &Player{
		Client:   client,
		id:       id,
		position: Position{X: 3222, Z: 3222},
		sync:     NewPlayerSync(id, 2048),
	}
}

type Player struct {
	*client.Client
	id       uint16
	position Position
	sync     PlayerSync
}

func (p *Player) Id() uint16 {
	return p.id
}

func (p *Player) Updated() bool {
	return false
}

// TODO extract the relative positions based on 161 from the cache which are stored as enums (key-value store)
var fixedModeInterfaces = map[int]int{
	// Chatbox area
	7:  163, // Overlay
	29: 162, // Main chatbox interface

	// Minimap area
	9:  122, // Exp display
	28: 160, // Data Orbs

	// Sidebar tabs
	// Top row
	68: 593, // Attack Styles
	69: 320, // Skills
	70: 399, // Journey tab (quests and shit)
	71: 149, // Item Bag
	72: 387, // Worn Equipment
	73: 541, // Prayer tab
	74: 218, // Spellbook

	// Bottom row
	75: 7,   // Clan Chat
	76: 109, // Account management
	77: 429, // Friends/Ignore list
	78: 182, // Logout
	79: 261, // Settings
	80: 216, // Emotes
	81: 239, // Music
}

func (p *Player) Initialize() {
	p.Send(&Success{PlayerId: p.id})
	p.Send(&InitializeScene{Position: p.position})
	p.Send(&SetHud{Id: 161})

	for key, value := range fixedModeInterfaces {
		p.Send(&OpenChildInterface{Parent: 161<<16 | uint32(key), Id: uint16(value), Behavior: 1})
	}

	p.Flush()
}

func (p *Player) Process(w *World) {
	p.Send(p.sync.Process(w))
	p.Flush()
}
