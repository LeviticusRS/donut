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
        position: Position{X: 3222, Z: 3222,},
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

func (p *Player) Initialize() {
    p.Send(&Success{PlayerId: p.id})
    p.Send(&InitializeScene{Position: p.position})
    p.Send(&SetHud{Id: 161})
    p.Flush()
}

func (p *Player) Process(w *World) {
    p.Send(p.sync.Process(w))
    p.Flush()
}
