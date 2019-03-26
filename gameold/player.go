package gameold

import (
    "github.com/sprinkle-it/donut/server"
)

type PlayerFactory func(*server.Client, uint16) *Player

type PlayerConfig struct {
}

func (cfg PlayerConfig) Build(client *server.Client, id uint16) *Player {
    return &Player{
        Client:   client,
        id:       id,
        position: Position{X: 3222, Z: 3222,},
        sync:     NewPlayerSync(id, 2048),
    }
}

type Player struct {
    *server.Client
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

	p.Send(&DisplaySystemMessage{Text: "Welcome to Sino and Sini's Donut Camp!"})
	p.Send(&ClearPerspectiveCamera{})
	p.Send(&ClearVariables{})

	for key, value := range fixedModeInterfaces {
		p.Send(&OpenChildInterface{Parent: 161<<16 | uint32(key), Id: uint16(value), Behavior: 1})
	}

	// Removes the 'You must have a display name to chat' filter within the chatbox
	p.Send(&InvokeInterfaceScript{Id: 1105, Arguments: []ScriptArgument{1}})

	for skillId := 0; skillId < 23; skillId++ {
		p.Send(&SetSkill{Id: uint8(skillId), Level: 99, Experience: 14000000})
	}

	p.Send(&SetPlayerContextMenuOption{Slot: 0, Label: "Follow", Prioritized: true})
	p.Send(&SetPlayerContextMenuOption{Slot: 2, Label: "Follow"})
	p.Send(&SetPlayerContextMenuOption{Slot: 3, Label: "Trade with"})
	p.Send(&SetPlayerContextMenuOption{Slot: 8, Label: "Report"})

	p.Send(&SetEnergy{Percentage: 100})
	p.Send(&SetWeight{Kilograms: 180}) // hadyn's irl weight be like
	p.Send(&SetMinimapState{Id: 2})

	p.Flush()
}

func (p *Player) Process(w *World) {
	p.Send(p.sync.Process(w))
	p.Flush()
}
