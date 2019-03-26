package game

import (
    "github.com/sprinkle-it/donut/server"
)

type Service struct {
    capacity int
    commands chan command
    world    World
}

func New(config Config) (*Service, error) {
    return config.Build()
}

func (s *Service) execute(cmd command) { s.commands <- cmd }

func (s *Service) Process() {
    s.world.Process()

    go func() {
        for command := range s.commands {
            command.execute(s)
        }
    }()
}

func (s *Service) HandleMail(mail server.Mail) {
    s.execute(handleMessage{mail: mail})
}

type command interface {
    execute(s *Service)
}

type handleMessage struct {
    mail server.Mail
}

func (c handleMessage) execute(s *Service) {
    source := c.mail.Source
    switch c.mail.Message.(type) {
    case *handshake:
        _ = source.SendNow(&Ready{})
    case *Authenticate:
        s.world.Register(source, Profile{})
    }
}