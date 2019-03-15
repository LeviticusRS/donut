package game

import (
	"fmt"
	"github.com/sprinkle-it/donut/pkg/account"
	"github.com/sprinkle-it/donut/pkg/auth"
	"github.com/sprinkle-it/donut/pkg/client"
	"github.com/sprinkle-it/donut/pkg/status"
	"reflect"
)

type Service struct {
	capacity int
	sessions map[uint64]*Session

	authenticator auth.Authenticator

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
	switch msg := c.mail.Message.(type) {
	case *handshake:
		_ = source.SendNow(&Ready{})
	case *NewLogin:
		// TODO validate request
		// TODO obtain player id if available, else reject immediately

		go s.login(source, account.Email(msg.Email), account.Password(msg.Password))
	case *Reconnect:
		// TODO
	}
}

// TODO revise this. should NewLogin be the receiver type? where does the source then come from?
func (s *Service) login(source *client.Client, email account.Email, password account.Password) {
	result, err := s.authenticator.Authenticate(email, password)
	if err != nil {
		_ = source.SendNow(status.ErrorLoadingProfile)
		return
	}

	switch result.(type) {
	case auth.Success:
		// TODO load player game state
		// TODO queue player

		_ = source.SendNow(&Success{})
		_ = source.SendNow(&RebuildScene{
			InitializePlayerPositions: InitializePlayerPositions{
				LocalPosition: Position{Level: 0, X: 3200, Z: 3200},
			},
			ChunkX: 3200 >> 3,
			ChunkZ: 3200 >> 3,
		})
		_ = source.SendNow(&SetHud{
			Id: 548,
		})

	case auth.PasswordMismatch, auth.CouldNotFindAccount:
		_ = source.SendNow(status.InvalidCredentials)

	default:
		source.Fatal(fmt.Errorf("gameservice: unsupported authentication result of type %v", reflect.TypeOf(result)))
	}
}

type unregisterSession struct {
	cli *client.Client
}

func (cmd unregisterSession) execute(s *Service) {
	delete(s.sessions, cmd.cli.Id())
	cmd.cli.Info("Unregistered game session")
}
