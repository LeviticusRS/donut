package game2

import "github.com/sprinkle-it/donut/server"

type Service struct {

}

func HandleMail(mail server.Mail) {

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
