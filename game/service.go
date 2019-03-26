package game

import (
    "github.com/sprinkle-it/donut/message"
    "github.com/sprinkle-it/donut/server"
    "go.uber.org/zap"
)

type Config struct {
    LoggerConfig     zap.Config
    SupportedVersion uint32
}

type Service struct {
    logger   *zap.Logger
    commands chan command
}

func New(config Config) (*Service, error) {
    logger, err := config.LoggerConfig.Build()
    if err != nil {
        return nil, err
    }

    return &Service{
        logger: logger,
    }, nil
}

func (s *Service) handleMail(mail server.Mail) {
    s.commands <- handleMessage{mail: mail}
}

func (s *Service) MailReceiver() server.MailReceiver {
    return server.MailReceiver{
        Handler: s.handleMail,
        Accept: []message.Config{
            handshakeConfig,
        },
    }
}

func (s *Service) Process() {
    go func() {
        for cmd := range s.commands {
            cmd.execute(s)
        }
    }()
}

type command interface {
    execute(s *Service)
}

type handleMessage struct {
    mail server.Mail
}

func (c handleMessage) execute(s *Service) {
    switch c.mail.Message.(type) {
    }
}
