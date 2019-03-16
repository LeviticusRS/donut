package game

import (
	"github.com/sprinkle-it/donut/pkg/client"
)

type Service struct {
	capacity int
	sessions map[uint64]*Session

	authenticator Authenticator

	commands chan command

	world World
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
    case *NewLogin:
        s.world.Register(source, Profile{})
    }
}

type unregisterSession struct {
	cli *client.Client
}

func (cmd unregisterSession) execute(s *Service) {
	delete(s.sessions, cmd.cli.Id())
	cmd.cli.Info("Unregistered game session")
}
