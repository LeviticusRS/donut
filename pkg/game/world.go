package game

import (
    "github.com/sprinkle-it/donut/server"
    "github.com/sprinkle-it/donut/pkg/status"
    "time"
)

type WorldConfig struct {
    Delta          time.Duration
    PlayerCapacity int
    PlayerConfig   PlayerConfig
}

func (cfg WorldConfig) Build() World {
    playerIds := make(chan uint16, cfg.PlayerCapacity)
    for i := 1; i < cfg.PlayerCapacity; i++ {
        playerIds <- uint16(i)
    }

    return World{
        ticker:   time.NewTicker(cfg.Delta),
        register: make(chan CreatePlayer, cfg.PlayerCapacity),
        players:  NewPlayerList(cfg.PlayerCapacity, cfg.PlayerConfig.Build),
    }
}

type World struct {
    ticker *time.Ticker

    register chan CreatePlayer

    players PlayerList
}

type CreatePlayer struct {
    Source  *server.Client
    Profile Profile
}

func (w *World) Register(cli *server.Client, profile Profile) {
    w.register <- CreatePlayer{Source: cli, Profile: profile}
}

func (w *World) Process() {
    go func() {
        for range w.ticker.C {
            w.registerPlayers()

            for _, id := range w.players.active {
                player := w.players.Get(id)
                player.Process(w)
            }
        }
    }()
}

func (w *World) registerPlayers() {
    for i := 0; i < 50; i++ {
        select {
        case msg := <-w.register:
            player, success := w.players.New(msg.Source)
            if !success {
                _ = msg.Source.SendNow(status.Full)
                continue
            }
            player.Initialize()
        default:
            return
        }
    }
}