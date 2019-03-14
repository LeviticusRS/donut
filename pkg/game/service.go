package game

import (
    "github.com/sprinkle-it/donut/pkg/client"
)

type Service struct {
    capacity int
    sessions map[uint64]*Session

    commands chan command
}

func New(config Config) (*Service, error) {
    return config.Build()
}

func (s *Service) execute(cmd command) { s.commands <- cmd }

func (s *Service) Process() {
    go func() {
        for command := range s.commands {
            command.execute(s)
        }
    }()
}

func (s *Service) HandleMail(mail client.Mail) {
    s.execute(handleMessage{mail: mail})
}

type command interface {
    execute(s *Service)
}

type handleMessage struct {
    mail client.Mail
}

func (c handleMessage) execute(s *Service) {
    source := c.mail.Source
    switch c.mail.Message.(type) {
    case *handshake:
        _ = source.SendNow(&Ready{})
    case *Authenticate:
        _ = source.SendNow(&Success{})
        _ = source.SendNow(&RebuildScene{
            InitializePlayerPositions: InitializePlayerPositions{
                LocalPosition: Position{Level: 0, X: 3200, Z: 3200,},
            },
            ChunkX: 3200 >> 3,
            ChunkZ: 3200 >> 3,
        })
        _ = source.SendNow(&SetHud{
            Id: 548,
        })
    }
}

type unregisterSession struct {
    cli *client.Client
}

func (cmd unregisterSession) execute(service *Service) {
    delete(service.sessions, cmd.cli.Id())
    cmd.cli.Info("Unregistered file session")
}
